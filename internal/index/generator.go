// Package index owns the canonical content-index.json generation step
// (FR-003; data-model.md "Content Index"; Principle II — biblio-cli is the
// sole producer; the generated JSON is consumed unchanged by the React Router
// prerender step and the in-Worker MCP query layer).
package index

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/monamaret/biblio/internal/content"
)

// Entry is one element of content-index.json (data-model.md "Content Index").
// Its JSON representation is a flat object that merges Slug, Path, Body, and
// all frontmatter Fields — consumable directly by both the prerender step
// (v0.02) and the MCP query layer (v0.03) without further transformation.
type Entry struct {
	Slug   string
	Path   string
	Body   string
	Fields map[string]any
}

// MarshalJSON produces a flat JSON object: Slug, Path, Body, and every
// frontmatter field at the top level (data-model.md "Content Index" shape).
func (e Entry) MarshalJSON() ([]byte, error) {
	merged := make(map[string]any, len(e.Fields)+3)
	for k, v := range e.Fields {
		merged[k] = v
	}
	merged["slug"] = e.Slug
	merged["path"] = e.Path
	merged["body"] = e.Body
	return json.Marshal(merged)
}

// Index is the content-index.json artifact produced by Generate. Every
// type carries Site (the page <title> comes from it, #215) and Entries;
// profile, log, and docs indexes carry nothing else, since the
// quarterly-only fields are omitted from JSON when empty
// (specs/004-quarterly-site-type data-model.md "content-index.json shape").
type Index struct {
	Site         *SiteMeta         `json:"site,omitempty"`
	Entries      []Entry           `json:"entries"`
	Issues       []IssueSummary    `json:"issues,omitempty"`
	Categories   []CategorySummary `json:"categories,omitempty"`
	CurrentIssue *string           `json:"current_issue,omitempty"`
	Typesetting  *IndexTypesetting `json:"typesetting,omitempty"`
	// Warnings collects non-fatal problems found while indexing (e.g. a
	// footnote reference with no definition). Never serialized.
	Warnings []string `json:"-"`
	// DraftIssuesExcluded counts draft issues left out of a publish-mode
	// quarterly index, for the CLI's report line. Never serialized.
	DraftIssuesExcluded int `json:"-"`
}

// IssueSummary is one element of a quarterly index's `issues` array,
// newest first.
type IssueSummary struct {
	Issue        string `json:"issue"`
	Slug         string `json:"slug"`
	Number       int    `json:"number"`
	Title        string `json:"title"`
	Date         string `json:"date"`
	Season       string `json:"season,omitempty"`
	CoverPlate   string `json:"cover_plate,omitempty"`
	CoverAlt     string `json:"cover_alt,omitempty"`
	Draft        bool   `json:"draft"`
	ArticleCount int    `json:"article_count"`
}

// CategorySummary is one element of a quarterly index's `categories`
// array.
type CategorySummary struct {
	Key          string `json:"key"`
	Label        string `json:"label"`
	ArticleCount int    `json:"article_count"`
}

// IndexTypesetting carries the structural typesetting switches components
// read at prerender time (research.md R3); everything visual lives in
// tokens.css instead.
type IndexTypesetting struct {
	MarginNotes bool   `json:"margin_notes"`
	Numerals    string `json:"numerals"`
}

// Options tunes Generate for the pre-build step. The zero value is what
// `biblio-cli index` uses: drafts excluded, no plates.
type Options struct {
	// IncludeDrafts keeps draft issues (and their articles) in the index,
	// marked draft: true — used by `biblio-cli dev` only.
	IncludeDrafts bool
	// Plates maps a content-relative image path to its converted plate URL
	// (/plates/<hash>.<ext>). Nil means no image rewriting.
	Plates map[string]string
	// Typesetting supplies the structural switches written to a quarterly
	// index; nil means the quarterly-default values.
	Typesetting *IndexTypesetting
}

// Generate reads every .md file under contentDir, parses and validates each
// against the site type declared in contentDir/site.yaml, and returns the
// completed Index together with any per-file validation errors.
//
// All validation errors are collected before returning — the caller receives
// the full list so authors can fix every failing file at once
// (contracts/cli-commands.md `biblio-cli index` "Output (failure)";
// spec Edge Cases). When errors are non-empty, the returned Index is partial
// and should not be written to disk.
func Generate(contentDir string) (Index, []error) {
	return GenerateWithOptions(contentDir, Options{})
}

// GenerateWithOptions is Generate with the pre-build step's options. For
// profile, log, and docs sites, opts has no effect.
func GenerateWithOptions(contentDir string, opts Options) (Index, []error) {
	cfgData, err := os.ReadFile(filepath.Join(contentDir, "site.yaml"))
	if err != nil {
		return Index{}, []error{fmt.Errorf("index: reading site.yaml: %w", err)}
	}
	cfg, err := content.ReadSiteConfig(cfgData)
	if err != nil {
		return Index{}, []error{err}
	}
	if cfg.SiteType == "quarterly" {
		return generateQuarterly(contentDir, cfg, opts)
	}

	var entries []Entry
	var errs []error

	walkErr := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		relPath, err := filepath.Rel(contentDir, path)
		if err != nil {
			errs = append(errs, fmt.Errorf("index: resolving path %q: %w", path, err))
			return nil
		}
		relPath = filepath.ToSlash(relPath)

		data, err := os.ReadFile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("index: reading %q: %w", relPath, err))
			return nil
		}

		item, err := content.ParseItem(relPath, data, cfg.SiteType)
		if err != nil {
			errs = append(errs, err)
			return nil
		}

		entries = append(entries, Entry{
			Slug:   item.Slug,
			Path:   item.Path,
			Body:   item.Body,
			Fields: item.Fields,
		})
		return nil
	})
	if walkErr != nil {
		errs = append(errs, fmt.Errorf("index: walking content directory: %w", walkErr))
	}

	return Index{Site: &SiteMeta{Name: cfg.Name, Description: cfg.Description}, Entries: entries}, errs
}

// WriteFile writes idx to path exactly as `biblio-cli index` always has
// (two-space indented JSON, no trailing newline), so every writer of
// content-index.json — the index command and the pre-build step — produces
// byte-identical output.
func WriteFile(path string, idx Index) error {
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("index: serialising content-index.json: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("index: writing content-index.json: %w", err)
	}
	return nil
}
