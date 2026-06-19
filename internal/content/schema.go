// Package content owns the canonical Go-side frontmatter schema
// parser/validator and slug derivation (Principle II — this logic must not be
// duplicated on the TypeScript side; data-model.md "Authored Content Item").
package content

import (
	"errors"
	"fmt"
	"strings"

	"go.yaml.in/yaml/v3"
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
}

// SiteConfig is the parsed content/site.yaml file that declares the active
// site type and top-level configuration for a Site Project.
type SiteConfig struct {
	SiteType    string `yaml:"site_type"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	MCPEnabled  bool   `yaml:"mcp_enabled"`
}

// ReadSiteConfig parses raw site.yaml bytes into a SiteConfig.
func ReadSiteConfig(data []byte) (SiteConfig, error) {
	var cfg SiteConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return SiteConfig{}, fmt.Errorf("content: parsing site.yaml: %w", err)
	}
	if cfg.SiteType == "" {
		return SiteConfig{}, errors.New("content: site.yaml missing required field 'site_type'")
	}
	return cfg, nil
}

// requiredFields maps each site type to the set of frontmatter fields that
// must be present and non-empty for any content item of that type
// (data-model.md "Authored Content Item" validation rules; spec Edge Cases).
var requiredFields = map[string][]string{
	"profile": {"title"},
	"log":     {"title", "date"},
	"docs":    {"title", "order"},
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

	return Item{
		Path:   relPath,
		Slug:   SlugFromPath(relPath),
		Fields: fields,
		Body:   body,
	}, nil
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
