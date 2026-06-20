// Package integration tests the full scaffold → index path for each of the
// three site types (SC-001, SC-002; quickstart.md Scenarios 1 and 2).
// These tests call the scaffold and index generators directly rather than
// shelling out to the compiled binary so that CI does not require a
// pre-built artifact — the binary path (T028/main.go) is validated separately.
package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/monamaret/biblio/internal/index"
	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
)

// scaffoldAndIndex runs the full scaffold → index pipeline for one site type
// and returns the parsed content-index.json output.
func scaffoldAndIndex(t *testing.T, siteType string) (map[string]any, string) {
	t.Helper()
	presets := theme.PresetsForType(siteType)
	if len(presets) == 0 {
		t.Fatalf("theme.PresetsForType(%q): no presets defined", siteType)
	}

	destDir := t.TempDir()
	projectDir, err := scaffold.Generate(destDir, scaffold.Options{
		ProjectName: "test-" + siteType,
		SiteType:    siteType,
		Domain:      "example.com",
		Theme:       presets[0].ID,
		MCPEnabled:  true,
	})
	if err != nil {
		t.Fatalf("scaffold.Generate(%q): %v", siteType, err)
	}

	contentDir := filepath.Join(projectDir, "content")
	idx, errs := index.Generate(contentDir)
	if len(errs) > 0 {
		t.Fatalf("index.Generate(%q): %v", siteType, errs)
	}
	if len(idx.Entries) == 0 {
		t.Fatalf("index.Generate(%q): want at least one entry, got zero", siteType)
	}

	data, err := json.Marshal(idx)
	if err != nil {
		t.Fatalf("json.Marshal(idx): %v", err)
	}
	var shape map[string]any
	if err := json.Unmarshal(data, &shape); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	return shape, projectDir
}

// TestScaffoldThenIndex_Profile validates the full scaffold → index path for
// the profile site type (SC-002).
func TestScaffoldThenIndex_Profile(t *testing.T) {
	shape, _ := scaffoldAndIndex(t, "profile")
	entries, ok := shape["entries"].([]any)
	if !ok || len(entries) == 0 {
		t.Fatalf("entries missing or empty in content-index.json for profile")
	}
	for _, raw := range entries {
		e, ok := raw.(map[string]any)
		if !ok {
			t.Fatal("entry is not an object")
		}
		if e["slug"] == "" || e["slug"] == nil {
			t.Errorf("entry missing slug: %v", e)
		}
		if e["path"] == "" || e["path"] == nil {
			t.Errorf("entry missing path: %v", e)
		}
		if e["title"] == "" || e["title"] == nil {
			t.Errorf("entry missing title: %v", e)
		}
	}
}

// TestScaffoldThenIndex_Log validates the full scaffold → index path for the
// log site type including the date field required by that site type (SC-002).
func TestScaffoldThenIndex_Log(t *testing.T) {
	shape, _ := scaffoldAndIndex(t, "log")
	entries, _ := shape["entries"].([]any)
	if len(entries) == 0 {
		t.Fatal("entries missing or empty in content-index.json for log")
	}
	for _, raw := range entries {
		e := raw.(map[string]any)
		if e["date"] == nil || e["date"] == "" {
			t.Errorf("log entry missing required 'date' field: %v", e)
		}
	}
}

// TestScaffoldThenIndex_Docs validates the full scaffold → index path for the
// docs site type including the order field required by that site type (SC-002).
func TestScaffoldThenIndex_Docs(t *testing.T) {
	shape, _ := scaffoldAndIndex(t, "docs")
	entries, _ := shape["entries"].([]any)
	if len(entries) == 0 {
		t.Fatal("entries missing or empty in content-index.json for docs")
	}
	for _, raw := range entries {
		e := raw.(map[string]any)
		if e["order"] == nil {
			t.Errorf("docs entry missing required 'order' field: %v", e)
		}
	}
}

// TestScaffoldedProject_RequiredFiles verifies that every scaffolded project
// directory contains the full artifact set (SC-001 / quickstart.md Scenario 1).
func TestScaffoldedProject_RequiredFiles(t *testing.T) {
	for _, siteType := range scaffold.SiteTypes {
		t.Run(siteType, func(t *testing.T) {
			_, projectDir := scaffoldAndIndex(t, siteType)
			for _, rel := range []string{
				filepath.Join("content", "site.yaml"),
				"wrangler.jsonc",
				"tokens.css",
			} {
				path := filepath.Join(projectDir, rel)
				if _, err := os.Stat(path); err != nil {
					t.Errorf("scaffold(%q): required file %q missing: %v", siteType, rel, err)
				}
			}

			siteYAML, _ := os.ReadFile(filepath.Join(projectDir, "content", "site.yaml"))
			if !strings.Contains(string(siteYAML), "mcp_enabled: true") {
				t.Errorf("scaffold(%q): content/site.yaml missing mcp_enabled: true", siteType)
			}
		})
	}
}

// TestScaffoldedProject_TokensCSS verifies that tokens.css contains a :root
// block with CSS custom properties from the chosen preset (FR-004).
func TestScaffoldedProject_TokensCSS(t *testing.T) {
	for _, siteType := range scaffold.SiteTypes {
		t.Run(siteType, func(t *testing.T) {
			_, projectDir := scaffoldAndIndex(t, siteType)
			tokensCSS, err := os.ReadFile(filepath.Join(projectDir, "tokens.css"))
			if err != nil {
				t.Fatalf("scaffold(%q): tokens.css: %v", siteType, err)
			}
			if !strings.Contains(string(tokensCSS), ":root") {
				t.Errorf("scaffold(%q): tokens.css missing :root block", siteType)
			}
			if !strings.Contains(string(tokensCSS), "--color-bg") {
				t.Errorf("scaffold(%q): tokens.css missing --color-bg custom property", siteType)
			}
		})
	}
}
