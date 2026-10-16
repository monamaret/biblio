// Package pilotsitequarterly proves a freshly scaffolded quarterly site's
// own sample content goes through the pre-build step cleanly, and keeps the
// template's shipped content-index.json honest (specs/004-quarterly-site-type
// US1, FR-002, SC-002; tasks.md T024).
package pilotsitequarterly_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/monamaret/biblio/internal/index"
	"github.com/monamaret/biblio/internal/prebuild"
	"github.com/monamaret/biblio/internal/scaffold"
)

const templateDir = "../../../templates/site/quarterly"

func scaffoldQuarterly(t *testing.T) string {
	t.Helper()
	dir, err := scaffold.Generate(t.TempDir(), scaffold.Options{
		ProjectName: "field-notes",
		SiteType:    "quarterly",
		Domain:      "example.com",
		Theme:       "quarterly-default",
		MCPEnabled:  true,
	})
	if err != nil {
		t.Fatalf("scaffold.Generate(quarterly): %v", err)
	}
	return dir
}

// The scaffolded sample content indexes with no errors and has the shape
// FR-002 promises: at least two issues, two categories, and a colophon.
func TestPilotQuarterly_SampleContentBuilds(t *testing.T) {
	dir := scaffoldQuarterly(t)
	report, err := prebuild.Run(dir, prebuild.Build)
	if err != nil {
		t.Fatalf("prebuild.Run: %v", err)
	}
	if len(report.Warnings) != 0 {
		t.Errorf("sample content produced warnings: %v", report.Warnings)
	}

	data, err := os.ReadFile(filepath.Join(dir, "content-index.json"))
	if err != nil {
		t.Fatalf("reading content-index.json: %v", err)
	}
	var idx struct {
		Entries []struct {
			Kind string `json:"kind"`
		} `json:"entries"`
		Issues     []json.RawMessage `json:"issues"`
		Categories []json.RawMessage `json:"categories"`
	}
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("parsing content-index.json: %v", err)
	}
	if len(idx.Issues) < 2 {
		t.Errorf("issues = %d, want at least 2", len(idx.Issues))
	}
	if len(idx.Categories) < 2 {
		t.Errorf("categories = %d, want at least 2", len(idx.Categories))
	}
	colophon := false
	for _, e := range idx.Entries {
		colophon = colophon || e.Kind == "colophon"
	}
	if !colophon {
		t.Error("no colophon entry in the sample content")
	}
}

// The template's shipped content-index.json equals what `biblio-cli index`
// writes for the template's own content/ — no plates, drafts excluded —
// which is what a fresh scaffold should carry before its first build
// (tasks.md T024, T045).
func TestPilotQuarterly_ShippedIndexMatchesIndexCommand(t *testing.T) {
	idx, errs := index.GenerateWithOptions(filepath.Join(templateDir, "content"), index.Options{})
	if len(errs) > 0 {
		t.Fatalf("index.GenerateWithOptions: %v", errs)
	}
	want := filepath.Join(t.TempDir(), "content-index.json")
	if err := index.WriteFile(want, idx); err != nil {
		t.Fatal(err)
	}
	wantData, _ := os.ReadFile(want)
	got, err := os.ReadFile(filepath.Join(templateDir, "content-index.json"))
	if err != nil {
		t.Fatalf("reading shipped content-index.json: %v", err)
	}
	if string(got) != string(wantData) {
		t.Error("templates/site/quarterly/content-index.json is stale; regenerate it with `go run ./cmd/biblio-cli index templates/site/quarterly/content`")
	}
}
