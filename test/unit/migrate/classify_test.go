// Package migrate_test exercises internal/migrate's Scaffold Region
// Classification: the per-path judgment the update flow makes before
// touching anything (data-model.md "Scaffold Region Classification";
// FR-001, FR-002).
package migrate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/monamaret/biblio/internal/migrate"
	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
)

// newScaffoldedProject scaffolds a fresh profile project into a temp dir,
// the same real-generator pattern test/integration/pilot-site uses, so
// classification is tested against the actual embedded template tree.
func newScaffoldedProject(t *testing.T) string {
	t.Helper()
	presets := theme.PresetsForType("profile")
	if len(presets) == 0 {
		t.Fatal(`theme.PresetsForType("profile"): want at least one preset, got none`)
	}
	projectDir, err := scaffold.Generate(t.TempDir(), scaffold.Options{
		ProjectName: "my-site",
		SiteType:    "profile",
		Domain:      "example.com",
		Theme:       presets[0].ID,
		MCPEnabled:  true,
	})
	if err != nil {
		t.Fatalf("scaffold.Generate(): %v", err)
	}
	return projectDir
}

func regionFor(t *testing.T, regions []migrate.Region, path string) migrate.Region {
	t.Helper()
	for _, r := range regions {
		if r.Path == path {
			return r
		}
	}
	t.Fatalf("ClassifySite(): no region for path %q", path)
	return migrate.Region{}
}

// An untouched file, identical to the manifest baseline, must classify as
// scaffold-tracked: safe for update to advance.
func TestClassifySite_UntouchedFile_ScaffoldTracked(t *testing.T) {
	projectDir := newScaffoldedProject(t)

	regions, err := migrate.ClassifySite(projectDir, "profile")
	if err != nil {
		t.Fatalf("ClassifySite(): %v", err)
	}

	got := regionFor(t, regions, "package.json")
	if got.Classification != migrate.ScaffoldTracked {
		t.Errorf("package.json classification = %q, want %q", got.Classification, migrate.ScaffoldTracked)
	}
}

// A file the author has edited, where the scaffold template hasn't changed
// that same file since scaffolding, must classify as author-customized:
// update must leave it untouched.
func TestClassifySite_AuthorEditedUnchangedScaffold_AuthorCustomized(t *testing.T) {
	projectDir := newScaffoldedProject(t)

	target := filepath.Join(projectDir, "package.json")
	if err := os.WriteFile(target, []byte(`{"name":"author-customized"}`), 0o644); err != nil {
		t.Fatalf("writing author customization: %v", err)
	}

	regions, err := migrate.ClassifySite(projectDir, "profile")
	if err != nil {
		t.Fatalf("ClassifySite(): %v", err)
	}

	got := regionFor(t, regions, "package.json")
	if got.Classification != migrate.AuthorCustomized {
		t.Errorf("package.json classification = %q, want %q", got.Classification, migrate.AuthorCustomized)
	}
}

// content/ paths must never be classified at all — they are the author's
// domain from the moment of scaffolding (Principle II) and are excluded
// from the diffable region set entirely.
func TestClassifySite_ContentPaths_NeverClassified(t *testing.T) {
	projectDir := newScaffoldedProject(t)

	regions, err := migrate.ClassifySite(projectDir, "profile")
	if err != nil {
		t.Fatalf("ClassifySite(): %v", err)
	}

	for _, r := range regions {
		if filepath.Dir(r.Path) == "content" || r.Path == "content" {
			t.Errorf("ClassifySite(): region for content/ path %q must not exist", r.Path)
		}
	}
}

// Every region must be classified before any write occurs — ClassifySite
// itself performs no writes, only reads and reports.
func TestClassifySite_NoRegionsLeftUnclassified(t *testing.T) {
	projectDir := newScaffoldedProject(t)

	regions, err := migrate.ClassifySite(projectDir, "profile")
	if err != nil {
		t.Fatalf("ClassifySite(): %v", err)
	}
	if len(regions) == 0 {
		t.Fatal("ClassifySite(): want at least one region, got none")
	}
	for _, r := range regions {
		switch r.Classification {
		case migrate.ScaffoldTracked, migrate.AuthorCustomized, migrate.Ambiguous:
			// valid
		default:
			t.Errorf("region %q has unrecognized classification %q", r.Path, r.Classification)
		}
	}
}
