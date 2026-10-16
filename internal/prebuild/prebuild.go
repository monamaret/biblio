// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// Package prebuild owns the step `biblio-cli build` and `biblio-cli dev`
// run before handing off to Vite: validate content/site.yaml, regenerate
// public/tokens.css (quarterly), convert images to plates (quarterly), and
// regenerate content-index.json — in that order, in one place
// (specs/004-quarterly-site-type FR-035; research.md R8; contracts/cli.md).
// Profile, log, and docs sites run only the validate and index steps, so
// their pages, styles, and assets are unchanged (FR-036).
package prebuild

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/monamaret/biblio/internal/content"
	"github.com/monamaret/biblio/internal/index"
	"github.com/monamaret/biblio/internal/plates"
	"github.com/monamaret/biblio/internal/theme"
)

// Mode selects publish semantics (Build) or preview semantics (Dev, which
// keeps draft issues in the index).
type Mode int

const (
	// Build is used by `biblio-cli build`: drafts excluded.
	Build Mode = iota
	// Dev is used by `biblio-cli dev`: drafts included and marked.
	Dev
)

// Report describes what one Run did, for the CLI to print.
type Report struct {
	// Skipped is true when the directory has no content/site.yaml, so there
	// was nothing to validate or index.
	Skipped bool
	// Steps lists the steps that ran, in order: "validate", "tokens",
	// "plates", "index".
	Steps []string
	// Preset and Overrides describe the tokens step: the quarterly preset
	// and how many typesetting keys site.yaml overrides.
	Preset    string
	Overrides int
	// Plates counts the images step's work.
	Plates plates.Stats
	// Indexed is the number of entries written to content-index.json.
	Indexed int
	// DraftIssuesExcluded counts draft issues left out (Build mode only).
	DraftIssuesExcluded int
	// Warnings are non-fatal problems; the build continues.
	Warnings []string
}

// Lines renders the report in the format of contracts/cli.md
// "biblio-cli build".
func (r Report) Lines() []string {
	if r.Skipped {
		return []string{"No content/site.yaml found — skipping the pre-build step."}
	}
	var lines []string
	for _, step := range r.Steps {
		switch step {
		case "validate":
			lines = append(lines, "Validating content/site.yaml… ok")
		case "tokens":
			lines = append(lines, fmt.Sprintf("Generating public/tokens.css (%s + %d overrides)… ok", r.Preset, r.Overrides))
		case "plates":
			lines = append(lines, fmt.Sprintf("Converting images: %d converted, %d skipped, %d removed → public/plates/", r.Plates.Converted, r.Plates.Skipped, r.Plates.Removed))
		case "index":
			line := fmt.Sprintf("Indexing: %d item(s)", r.Indexed)
			if r.DraftIssuesExcluded > 0 {
				line += fmt.Sprintf(", %d draft issue(s) excluded", r.DraftIssuesExcluded)
			}
			lines = append(lines, line+" → content-index.json")
		}
	}
	for _, w := range r.Warnings {
		lines = append(lines, "warning: "+w)
	}
	return lines
}

// Run executes the pre-build step for the project at projectDir. Any
// validation error stops it before anything is written, and every failing
// file or key is reported together.
func Run(projectDir string, mode Mode) (Report, error) {
	contentDir := filepath.Join(projectDir, "content")
	siteYAML, err := os.ReadFile(filepath.Join(contentDir, "site.yaml"))
	if errors.Is(err, os.ErrNotExist) {
		return Report{Skipped: true}, nil
	}
	if err != nil {
		return Report{}, fmt.Errorf("prebuild: reading content/site.yaml: %w", err)
	}

	var report Report

	// 1. Validate settings (all types).
	cfg, err := content.ReadSiteConfig(siteYAML)
	if err != nil {
		return Report{}, err
	}
	report.Steps = append(report.Steps, "validate")
	opts := index.Options{IncludeDrafts: mode == Dev}

	if cfg.SiteType == "quarterly" {
		// 2. Regenerate tokens.css from the preset plus the overrides.
		eff, preset, err := writeQuarterlyTokens(projectDir, cfg)
		if err != nil {
			return Report{}, err
		}
		report.Steps = append(report.Steps, "tokens")
		report.Preset, report.Overrides = preset, len(eff.Overridden)
		report.Warnings = append(report.Warnings, missingFontWarnings(projectDir, eff)...)
		opts.Typesetting = &index.IndexTypesetting{MarginNotes: eff.MarginNotes, Numerals: eff.Numerals}

		// 3. Convert images to plates, and give the index their URLs.
		settings := plates.DefaultSettings
		settings.Grayscale = eff.GrayscaleImages
		urls, stats, err := plates.Convert(contentDir, filepath.Join(projectDir, "public", "plates"), settings)
		if err != nil {
			return Report{}, err
		}
		report.Steps = append(report.Steps, "plates")
		report.Plates = stats
		opts.Plates = urls
	}

	// 4. Index (all types).
	idx, errs := index.GenerateWithOptions(contentDir, opts)
	if len(errs) > 0 {
		var sb strings.Builder
		fmt.Fprintf(&sb, "index: %d file(s) failed schema validation:", len(errs))
		for _, e := range errs {
			fmt.Fprintf(&sb, "\n  • %v", e)
		}
		return Report{}, errors.New(sb.String())
	}
	if err := index.WriteFile(filepath.Join(projectDir, "content-index.json"), idx); err != nil {
		return Report{}, err
	}
	report.Steps = append(report.Steps, "index")
	report.Indexed = len(idx.Entries)
	report.DraftIssuesExcluded = idx.DraftIssuesExcluded
	report.Warnings = append(report.Warnings, idx.Warnings...)

	return report, nil
}

// writeQuarterlyTokens resolves the quarterly typesetting (preset plus the
// validated site.yaml overrides) and writes public/tokens.css from it
// (specs/004-quarterly-site-type FR-024; research.md R8).
func writeQuarterlyTokens(projectDir string, cfg content.SiteConfig) (theme.Effective, string, error) {
	presets := theme.PresetsForType("quarterly")
	if len(presets) == 0 {
		return theme.Effective{}, "", errors.New("prebuild: no quarterly theme preset")
	}
	preset := presets[0]
	var overrides theme.Typesetting
	if cfg.Typesetting != nil {
		overrides = *cfg.Typesetting
	}
	eff, err := theme.Resolve(preset, overrides)
	if err != nil {
		return theme.Effective{}, "", err
	}
	css, err := theme.GenerateQuarterlyTokensCSS(preset, eff)
	if err != nil {
		return theme.Effective{}, "", err
	}
	publicDir := filepath.Join(projectDir, "public")
	if err := os.MkdirAll(publicDir, 0o755); err != nil {
		return theme.Effective{}, "", fmt.Errorf("prebuild: creating public/: %w", err)
	}
	if err := os.WriteFile(filepath.Join(publicDir, "tokens.css"), []byte(css), 0o644); err != nil {
		return theme.Effective{}, "", fmt.Errorf("prebuild: writing public/tokens.css: %w", err)
	}
	return eff, preset.ID, nil
}

// missingFontWarnings warns for every configured font family with no
// matching file in public/fonts/; the page then falls back to the generic
// family in the stack (FR-025). File names and family names are compared
// lowercased with everything but letters and digits removed, so
// "IBM Plex Mono" matches ibm-plex-mono-latin-400-normal.woff2.
func missingFontWarnings(projectDir string, eff theme.Effective) []string {
	var files []string
	entries, _ := os.ReadDir(filepath.Join(projectDir, "public", "fonts"))
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, fontKey(e.Name()))
		}
	}
	checks := []struct{ key, family, fallback string }{
		{"font_body", eff.FontBody, "sans-serif"},
		{"font_display", eff.FontDisplay, "sans-serif"},
		{"font_mono", eff.FontMono, "monospace"},
	}
	var warnings []string
	seen := map[string]bool{}
	for _, c := range checks {
		want := fontKey(c.family)
		found := false
		for _, f := range files {
			if strings.HasPrefix(f, want) {
				found = true
				break
			}
		}
		if !found && !seen[c.key] {
			seen[c.key] = true
			warnings = append(warnings, fmt.Sprintf("typesetting.%s: no file for %q in public/fonts/ — pages will use the fallback %s", c.key, c.family, c.fallback))
		}
	}
	return warnings
}

// fontKey lowercases s and keeps only its letters and digits.
func fontKey(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}
