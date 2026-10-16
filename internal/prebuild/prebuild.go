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

	"github.com/monamaret/biblio/internal/content"
	"github.com/monamaret/biblio/internal/index"
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
	if _, err := content.ReadSiteConfig(siteYAML); err != nil {
		return Report{}, err
	}
	report.Steps = append(report.Steps, "validate")

	// 4. Index (all types).
	idx, errs := index.GenerateWithOptions(contentDir, index.Options{IncludeDrafts: mode == Dev})
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
