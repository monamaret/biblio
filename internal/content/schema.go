// Package content owns the canonical Go-side frontmatter schema
// parser/validator and slug derivation (Principle II — this logic must not be
// duplicated on the TypeScript side; data-model.md "Authored Content Item").
package content

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"

	"github.com/monamaret/biblio/internal/theme"
)

// Item is a parsed, validated Authored Content Item (data-model.md).
type Item struct {
	// Path is the relative path under content/ (e.g. "posts/2026-06-01-hello.md").
	Path string
	// Slug is the canonical, structurally-derived identifier (e.g. "posts/2026-06-01-hello").
	Slug string
	// Fields holds the parsed YAML frontmatter key/value pairs.
	Fields map[string]any
	// Body holds the raw Markdown body following the closing frontmatter delimiter.
	Body string
	// Kind is the quarterly file kind ("issue", "article", or "colophon");
	// empty for the other site types (specs/004-quarterly-site-type FR-008).
	Kind string
}

// SiteConfig is the parsed content/site.yaml file that declares the active
// site type and top-level configuration for a Site Project.
type SiteConfig struct {
	SiteType    string `yaml:"site_type"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	MCPEnabled  bool   `yaml:"mcp_enabled"`
	// Typesetting holds the validated quarterly `typesetting:` overrides,
	// or nil when the block is absent (specs/004-quarterly-site-type
	// FR-022). Only quarterly sites may set it.
	Typesetting *theme.Typesetting `yaml:"-"`
}

// ReadSiteConfig parses raw site.yaml bytes into a SiteConfig. The rest of
// site.yaml stays leniently parsed, as before; a `typesetting:` block is
// validated key by key, and every problem is reported together, each line
// prefixed "site.yaml:" (contracts/config-and-tokens.md "Error format").
func ReadSiteConfig(data []byte) (SiteConfig, error) {
	var raw struct {
		SiteConfig  `yaml:",inline"`
		Typesetting map[string]any `yaml:"typesetting"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return SiteConfig{}, fmt.Errorf("content: parsing site.yaml: %w", err)
	}
	cfg := raw.SiteConfig
	if cfg.SiteType == "" {
		return SiteConfig{}, errors.New("content: site.yaml missing required field 'site_type'")
	}
	if raw.Typesetting != nil {
		if cfg.SiteType != "quarterly" {
			return SiteConfig{}, fmt.Errorf("site.yaml: typesetting is only supported for site_type quarterly (this site is %q)", cfg.SiteType)
		}
		ts, errs := theme.ParseTypesetting(raw.Typesetting)
		if len(errs) > 0 {
			lines := make([]string, len(errs))
			for i, e := range errs {
				lines[i] = "site.yaml: " + e.Error()
			}
			return SiteConfig{}, errors.New(strings.Join(lines, "\n"))
		}
		cfg.Typesetting = &ts
	}
	return cfg, nil
}

// requiredFields maps each site type to the set of frontmatter fields that
// must be present and non-empty for any content item of that type
// (data-model.md "Authored Content Item" validation rules; spec Edge Cases).
var requiredFields = map[string][]string{
	"profile": {"title"},
	"log":     {"title", "date"},
	"docs":    {"title", "section", "order"},
}

// ParseItem parses and validates a single Markdown+frontmatter file.
// relPath is the file's relative path under content/ (used to derive the slug
// and set Item.Path). data is the raw file content. siteType selects which
// required-field set to validate against.
//
// Returns a non-nil error for any of: malformed YAML frontmatter, missing or
// empty required fields, or an unclosed frontmatter delimiter — never
// silently returns an incomplete Item (spec Edge Cases).
func ParseItem(relPath string, data []byte, siteType string) (Item, error) {
	frontRaw, body, err := splitFrontmatter(data)
	if err != nil {
		return Item{}, fmt.Errorf("content: %s: %w", relPath, err)
	}

	var fields map[string]any
	if len(frontRaw) > 0 {
		if err := yaml.Unmarshal(frontRaw, &fields); err != nil {
			return Item{}, fmt.Errorf("content: %s: parsing frontmatter: %w", relPath, err)
		}
	}
	if fields == nil {
		fields = map[string]any{}
	}

	required, known := requiredFields[siteType]
	if !known {
		return Item{}, fmt.Errorf("content: unknown site type %q", siteType)
	}
	for _, field := range required {
		v, ok := fields[field]
		if !ok || v == nil || fmt.Sprintf("%v", v) == "" {
			return Item{}, fmt.Errorf("content: %s: missing required frontmatter field %q (required for site type %q)", relPath, field, siteType)
		}
	}

	if siteType == "log" {
		if err := validateLogFields(relPath, fields); err != nil {
			return Item{}, err
		}
	}

	if siteType == "docs" {
		if err := validateDocsFields(relPath, fields); err != nil {
			return Item{}, err
		}
	}

	return Item{
		Path:   relPath,
		Slug:   SlugFromPath(relPath),
		Fields: fields,
		Body:   body,
	}, nil
}

// validateLogFields checks the structural shape of the project-log-specific
// fields beyond mere presence (data-model.md "Project-Log Entry" validation
// rules; spec Edge Cases): `date` must parse as YYYY-MM-DD, `tags` (if
// present) must be a list of strings, and `project` (if present) must be a
// string. `tags` and `project` are both optional — their absence is valid.
func validateLogFields(relPath string, fields map[string]any) error {
	if date, ok := fields["date"]; ok {
		s, isString := date.(string)
		if !isString {
			return fmt.Errorf("content: %s: field %q must be a date string (YYYY-MM-DD)", relPath, "date")
		}
		if _, err := time.Parse("2006-01-02", s); err != nil {
			return fmt.Errorf("content: %s: field %q is not a valid date (want YYYY-MM-DD): %w", relPath, "date", err)
		}
	}

	if tags, ok := fields["tags"]; ok {
		list, isList := tags.([]any)
		if !isList {
			return fmt.Errorf("content: %s: field %q must be a list of strings", relPath, "tags")
		}
		for _, tag := range list {
			if _, isString := tag.(string); !isString {
				return fmt.Errorf("content: %s: field %q contains a non-string element %v", relPath, "tags", tag)
			}
		}
	}

	if project, ok := fields["project"]; ok {
		if _, isString := project.(string); !isString {
			return fmt.Errorf("content: %s: field %q must be a string", relPath, "project")
		}
	}

	return nil
}

// validateDocsFields checks the structural shape of the documentation-type-
// specific fields beyond mere presence (data-model.md "Documentation Page"
// validation rules; spec Edge Cases): `section` must be a non-empty string,
// and `order` must parse as an integer. Both are already required to be
// present by requiredFields["docs"] — this only validates their shape once
// present.
func validateDocsFields(relPath string, fields map[string]any) error {
	if section, ok := fields["section"]; ok {
		s, isString := section.(string)
		if !isString || s == "" {
			return fmt.Errorf("content: %s: field %q must be a non-empty string", relPath, "section")
		}
	}

	if order, ok := fields["order"]; ok {
		switch v := order.(type) {
		case int, int64:
			// valid
		case float64:
			if v != float64(int(v)) {
				return fmt.Errorf("content: %s: field %q must be an integer", relPath, "order")
			}
		default:
			return fmt.Errorf("content: %s: field %q must be an integer", relPath, "order")
		}
	}

	return nil
}

// splitFrontmatter splits a Markdown file's raw bytes into YAML frontmatter
// and body. The frontmatter is enclosed by leading and closing "---" lines.
// Returns (nil, data, nil) when no frontmatter delimiter is found.
func splitFrontmatter(data []byte) (front []byte, body string, err error) {
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return nil, string(data), nil
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			front = []byte(strings.Join(lines[1:i], "\n"))
			bodyLines := lines[i+1:]
			// Strip one leading blank line that conventionally follows the closer.
			if len(bodyLines) > 0 && strings.TrimSpace(bodyLines[0]) == "" {
				bodyLines = bodyLines[1:]
			}
			return front, strings.Join(bodyLines, "\n"), nil
		}
	}
	return nil, "", errors.New("unclosed frontmatter: missing closing ---")
}
