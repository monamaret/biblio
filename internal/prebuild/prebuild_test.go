package prebuild_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/monamaret/biblio/internal/index"
	"github.com/monamaret/biblio/internal/prebuild"
	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
)

// scaffoldSite generates a fresh project of siteType with its first preset
// and returns the project directory.
func scaffoldSite(t *testing.T, siteType string) string {
	t.Helper()
	presets := theme.PresetsForType(siteType)
	if len(presets) == 0 {
		t.Fatalf("no preset for %q", siteType)
	}
	dir, err := scaffold.Generate(t.TempDir(), scaffold.Options{
		ProjectName: "site",
		SiteType:    siteType,
		Domain:      "example.com",
		Theme:       presets[0].ID,
		MCPEnabled:  true,
	})
	if err != nil {
		t.Fatalf("scaffold.Generate(%q): %v", siteType, err)
	}
	return dir
}

// expectedIndexJSON is exactly what `biblio-cli index` writes for dir's
// content/ (internal/cli/index.go formatting).
func expectedIndexJSON(t *testing.T, dir string) []byte {
	t.Helper()
	idx, errs := index.Generate(filepath.Join(dir, "content"))
	if len(errs) > 0 {
		t.Fatalf("index.Generate: %v", errs)
	}
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return data
}

// For a log site the pre-build step only validates and indexes: the index
// it writes equals what `biblio-cli index` produces, tokens.css is
// untouched, and no plates are written (spec FR-035, FR-036).
func TestRun_LogSiteValidatesAndIndexesOnly(t *testing.T) {
	dir := scaffoldSite(t, "log")
	tokensBefore, err := os.ReadFile(filepath.Join(dir, "public", "tokens.css"))
	if err != nil {
		t.Fatalf("reading tokens.css: %v", err)
	}
	if err := os.Remove(filepath.Join(dir, "content-index.json")); err != nil {
		t.Fatalf("removing shipped index: %v", err)
	}

	report, err := prebuild.Run(dir, prebuild.Build)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dir, "content-index.json"))
	if err != nil {
		t.Fatalf("content-index.json not written: %v", err)
	}
	if want := expectedIndexJSON(t, dir); string(got) != string(want) {
		t.Errorf("content-index.json differs from `biblio-cli index` output.\ngot:\n%s\nwant:\n%s", got, want)
	}

	tokensAfter, _ := os.ReadFile(filepath.Join(dir, "public", "tokens.css"))
	if string(tokensAfter) != string(tokensBefore) {
		t.Error("public/tokens.css changed for a log site; only quarterly regenerates tokens")
	}
	if _, err := os.Stat(filepath.Join(dir, "public", "plates")); !os.IsNotExist(err) {
		t.Errorf("public/plates exists for a log site (err=%v); only quarterly converts images", err)
	}

	if !slices.Equal(report.Steps, []string{"validate", "index"}) {
		t.Errorf("report.Steps = %v, want [validate index]", report.Steps)
	}
	if report.Indexed != 3 {
		t.Errorf("report.Indexed = %d, want 3 (the log template's sample posts)", report.Indexed)
	}
}

// An invalid site.yaml stops the step before anything is written.
func TestRun_InvalidSiteYAMLWritesNothing(t *testing.T) {
	dir := scaffoldSite(t, "log")
	if err := os.Remove(filepath.Join(dir, "content-index.json")); err != nil {
		t.Fatalf("removing shipped index: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "content", "site.yaml"), []byte("site_type: log\ntypesetting:\n  accent: \"#000000\"\n"), 0o644); err != nil {
		t.Fatalf("writing site.yaml: %v", err)
	}

	_, err := prebuild.Run(dir, prebuild.Build)
	if err == nil {
		t.Fatal("Run() error = nil, want a validation error")
	}
	if !strings.Contains(err.Error(), "typesetting is only supported") {
		t.Errorf("error = %q, want the site.yaml validation message", err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "content-index.json")); !os.IsNotExist(statErr) {
		t.Error("content-index.json was written despite a validation error")
	}
}

// A content validation error is reported for every failing file and also
// writes nothing.
func TestRun_ContentErrorsReportedTogether(t *testing.T) {
	dir := scaffoldSite(t, "log")
	for _, name := range []string{"bad-one.md", "bad-two.md"} {
		if err := os.WriteFile(filepath.Join(dir, "content", "posts", name), []byte("---\ntitle: x\n---\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	before, _ := os.ReadFile(filepath.Join(dir, "content-index.json"))

	_, err := prebuild.Run(dir, prebuild.Build)
	if err == nil {
		t.Fatal("Run() error = nil, want content validation errors")
	}
	for _, name := range []string{"bad-one.md", "bad-two.md"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error %q does not name %s", err, name)
		}
	}
	after, _ := os.ReadFile(filepath.Join(dir, "content-index.json"))
	if string(after) != string(before) {
		t.Error("content-index.json changed despite validation errors")
	}
}

// A directory with no content/site.yaml (not a content project) is left
// alone: the step reports it was skipped instead of failing, so `build`
// keeps its old behavior there.
func TestRun_NoSiteYAMLIsSkipped(t *testing.T) {
	dir := t.TempDir()
	report, err := prebuild.Run(dir, prebuild.Build)
	if err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}
	if !report.Skipped {
		t.Error("report.Skipped = false, want true when content/site.yaml is absent")
	}
	if _, err := os.Stat(filepath.Join(dir, "content-index.json")); !os.IsNotExist(err) {
		t.Error("content-index.json written for a directory with no content/site.yaml")
	}
}

// writeTypesetting appends a typesetting block to dir's site.yaml.
func writeTypesetting(t *testing.T, dir, block string) {
	t.Helper()
	path := filepath.Join(dir, "content", "site.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, []byte("typesetting:\n"+block)...), 0o644); err != nil {
		t.Fatal(err)
	}
}

// For quarterly, the step regenerates tokens.css from the preset plus the
// site.yaml overrides, and writes the structural switches into the index
// (specs/004 FR-024; research.md R3, R8).
func TestRun_QuarterlyRegeneratesTokensAndSwitches(t *testing.T) {
	dir := scaffoldSite(t, "quarterly")
	writeTypesetting(t, dir, "  accent: \"#1e4fd7\"\n  margin_notes: false\n  numerals: arabic\n")

	report, err := prebuild.Run(dir, prebuild.Build)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !slices.Contains(report.Steps, "tokens") {
		t.Errorf("report.Steps = %v, want a tokens step", report.Steps)
	}
	tokens, _ := os.ReadFile(filepath.Join(dir, "public", "tokens.css"))
	if !strings.Contains(string(tokens), "--color-accent: #1e4fd7;") || !strings.Contains(string(tokens), "--ts-numerals: decimal;") {
		t.Errorf("tokens.css did not apply the overrides:\n%s", tokens)
	}
	if !strings.Contains(strings.Join(report.Lines(), "\n"), "Generating public/tokens.css (quarterly-default + 3 overrides)… ok") {
		t.Errorf("report lines = %v, want the tokens line", report.Lines())
	}

	idxData, _ := os.ReadFile(filepath.Join(dir, "content-index.json"))
	var idx struct {
		Typesetting struct {
			MarginNotes bool   `json:"margin_notes"`
			Numerals    string `json:"numerals"`
		} `json:"typesetting"`
	}
	if err := json.Unmarshal(idxData, &idx); err != nil {
		t.Fatal(err)
	}
	if idx.Typesetting.MarginNotes || idx.Typesetting.Numerals != "arabic" {
		t.Errorf("index typesetting = %+v, want margin_notes false, numerals arabic", idx.Typesetting)
	}
}

// A configured font family with no matching file warns but does not fail;
// the bundled faces are found (file names compared ignoring case and
// punctuation) (FR-025; research.md R9).
func TestRun_QuarterlyFontFileWarning(t *testing.T) {
	dir := scaffoldSite(t, "quarterly")
	report, err := prebuild.Run(dir, prebuild.Build)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Warnings) != 0 {
		t.Errorf("bundled fonts produced warnings: %v", report.Warnings)
	}

	writeTypesetting(t, dir, "  font_body: \"Söhne\"\n")
	report, err = prebuild.Run(dir, prebuild.Build)
	if err != nil {
		t.Fatalf("Run() error = %v, want only a warning", err)
	}
	joined := strings.Join(report.Warnings, "\n")
	if !strings.Contains(joined, `typesetting.font_body: no file for "Söhne" in public/fonts/`) {
		t.Errorf("warnings = %q, want the missing-font warning", joined)
	}
}

// A fresh quarterly scaffold's tokens.css is exactly what its first build
// writes, so `new` and `build` never disagree (tasks.md T062/T067).
func TestRun_QuarterlyScaffoldTokensParity(t *testing.T) {
	dir := scaffoldSite(t, "quarterly")
	before, _ := os.ReadFile(filepath.Join(dir, "public", "tokens.css"))
	if _, err := prebuild.Run(dir, prebuild.Build); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(filepath.Join(dir, "public", "tokens.css"))
	if string(before) != string(after) {
		t.Errorf("scaffolded tokens.css differs from the first build's.\nscaffold:\n%s\nbuild:\n%s", before, after)
	}
}
