package update_three_sites_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/monamaret/biblio/internal/migrate"
	"github.com/monamaret/biblio/internal/prebuild"
	"github.com/monamaret/biblio/internal/scaffold"
)

// snapshotTree maps every file under dir (slash-separated, relative) to
// its contents; missing dirs give an empty map.
func snapshotTree(t *testing.T, dir string) map[string]string {
	t.Helper()
	files := map[string]string{}
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		rel, _ := filepath.Rel(dir, p)
		files[filepath.ToSlash(rel)] = string(data)
		return nil
	})
	return files
}

// The Update Run contract holds for a quarterly site with authored content
// and typesetting that has been built: template files advance, content/
// and site.yaml are untouched, and the files the build generates —
// content-index.json, public/tokens.css, public/plates/ — are never
// reported as conflicts or overwritten, even when the template's shipped
// content-index.json changed upstream (specs/004-quarterly-site-type FR-003,
// US6; tasks.md T088).
func TestUpdateRun_QuarterlySite_AdvancesAndPreservesGenerated(t *testing.T) {
	projectDir := scaffoldSite(t, "quarterly", "my-quarterly")

	article := "---\ntitle: A New Piece\ncategory: Essays\nauthor: Me\n---\n\nFresh.[^1]\n\n[^1]: A note.\n"
	if err := os.WriteFile(filepath.Join(projectDir, "content", "issues", "02-autumn-2026", "new-piece.md"), []byte(article), 0o644); err != nil {
		t.Fatal(err)
	}
	siteYAML := filepath.Join(projectDir, "content", "site.yaml")
	data, _ := os.ReadFile(siteYAML)
	if err := os.WriteFile(siteYAML, append(data, []byte("typesetting:\n  numerals: arabic\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := prebuild.Run(projectDir, prebuild.Build); err != nil {
		t.Fatalf("prebuild.Run: %v", err)
	}

	// Simulate a newer template whose shipped content-index.json differs
	// from the one this site was scaffolded with.
	manifest, err := scaffold.ReadManifest(projectDir)
	if err != nil {
		t.Fatal(err)
	}
	if _, tracked := manifest.Files["content-index.json"]; tracked {
		manifest.Files["content-index.json"] = scaffold.HashBytes([]byte("an older shipped index"))
		if err := scaffold.WriteManifest(projectDir, manifest); err != nil {
			t.Fatal(err)
		}
	}

	contentBefore := snapshotTree(t, filepath.Join(projectDir, "content"))
	generatedBefore := map[string]string{}
	for _, rel := range []string{"content-index.json", "public/tokens.css"} {
		b, _ := os.ReadFile(filepath.Join(projectDir, rel))
		generatedBefore[rel] = string(b)
	}
	platesBefore := snapshotTree(t, filepath.Join(projectDir, "public", "plates"))
	if len(platesBefore) == 0 {
		t.Fatal("the build produced no plates to protect")
	}

	run, err := migrate.RunUpdate(projectDir, "quarterly", nil)
	if err != nil {
		t.Fatalf("RunUpdate(): %v (conflicts: %+v)", err, run.Conflicts)
	}
	if len(run.Conflicts) != 0 {
		t.Errorf("RunUpdate(): conflicts %+v, want none — generated files are not template regions", run.Conflicts)
	}
	if len(run.Advanced) == 0 {
		t.Error("RunUpdate(): want template regions advanced, got none")
	}

	if got := snapshotTree(t, filepath.Join(projectDir, "content")); len(got) != len(contentBefore) {
		t.Errorf("content/ file count changed: %d → %d", len(contentBefore), len(got))
	} else {
		for k, v := range contentBefore {
			if got[k] != v {
				t.Errorf("content/%s changed during update", k)
			}
		}
	}
	for rel, want := range generatedBefore {
		got, _ := os.ReadFile(filepath.Join(projectDir, rel))
		if string(got) != want {
			t.Errorf("%s was overwritten by update", rel)
		}
	}
	if got := snapshotTree(t, filepath.Join(projectDir, "public", "plates")); len(got) != len(platesBefore) {
		t.Errorf("public/plates/ changed during update: %d → %d files", len(platesBefore), len(got))
	}
}
