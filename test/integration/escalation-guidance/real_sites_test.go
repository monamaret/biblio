// Package escalation_guidance runs the documented storage-escalation
// thresholds against real scaffolded sites' actual content-index.json
// scale (FR-004, FR-005, Scenario 3).
package escalation_guidance_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/monamaret/biblio/internal/escalation"
	"github.com/monamaret/biblio/internal/index"
	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
)

// scaffoldAndIndex scaffolds a real project of siteType, generates its
// content-index.json from the seed content `new` writes, and returns the
// index path — the same real scaffold -> index pipeline `biblio-cli`
// itself runs.
func scaffoldAndIndex(t *testing.T, siteType, name string) string {
	t.Helper()
	presets := theme.PresetsForType(siteType)
	if len(presets) == 0 {
		t.Fatalf("theme.PresetsForType(%q): want at least one preset, got none", siteType)
	}
	projectDir, err := scaffold.Generate(t.TempDir(), scaffold.Options{
		ProjectName: name,
		SiteType:    siteType,
		Domain:      name + ".example.com",
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
	data, err := json.Marshal(idx)
	if err != nil {
		t.Fatalf("marshalling content-index.json: %v", err)
	}
	indexPath := filepath.Join(projectDir, "content-index.json")
	if err := os.WriteFile(indexPath, data, 0o644); err != nil {
		t.Fatalf("writing content-index.json: %v", err)
	}
	return indexPath
}

// The guidance check must correctly reflect each real site's actual,
// current scale — a freshly scaffolded site with only seed content is
// nowhere near either documented threshold (FR-004, Scenario 3).
func TestCheckSite_AllThreeRealSites_ReportNotCrossed(t *testing.T) {
	sites := map[string]string{
		"profile": "my-profile-site",
		"log":     "my-project-log",
		"docs":    "my-docs-site",
	}

	for siteType, name := range sites {
		t.Run(siteType, func(t *testing.T) {
			indexPath := scaffoldAndIndex(t, siteType, name)

			check, err := escalation.CheckSite(indexPath, name)
			if err != nil {
				t.Fatalf("CheckSite(): %v", err)
			}
			if check.D1Status != escalation.NotCrossed {
				t.Errorf("D1Status = %q for a freshly scaffolded site, want %q", check.D1Status, escalation.NotCrossed)
			}
			if check.VectorizeStatus != escalation.NotCrossed {
				t.Errorf("VectorizeStatus = %q for a freshly scaffolded site, want %q", check.VectorizeStatus, escalation.NotCrossed)
			}
			if check.EntryCount <= 0 {
				t.Errorf("EntryCount = %d, want at least the seed content entries", check.EntryCount)
			}
		})
	}
}
