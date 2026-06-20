// Package update_three_sites exercises `biblio-cli update`'s Update Run
// against real scaffolded sites end to end — first a single site (FR-001,
// FR-002, Scenario 1; this file), then all three site types (FR-003,
// Scenario 2; cross_type_test.go and friends).
package update_three_sites_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/monamaret/biblio/internal/migrate"
	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
)

// scaffoldSite scaffolds a real project of the given site type into a temp
// dir, the same real-generator pattern test/integration/pilot-site uses.
func scaffoldSite(t *testing.T, siteType, name string) string {
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
	return projectDir
}

// An Update Run against an untouched real site must advance every
// scaffold-tracked region, report nothing as a conflict, and leave the
// author's (here: nonexistent) customizations alone (data-model.md
// "Update Run" success triple).
func TestUpdateRun_SingleSite_AdvancesUntouchedScaffold(t *testing.T) {
	projectDir := scaffoldSite(t, "profile", "my-profile-site")

	run, err := migrate.RunUpdate(projectDir, "profile", nil)
	if err != nil {
		t.Fatalf("RunUpdate(): %v", err)
	}
	if len(run.Conflicts) != 0 {
		t.Errorf("RunUpdate(): got %d conflicts on an untouched site, want 0", len(run.Conflicts))
	}
	if len(run.Advanced) == 0 {
		t.Error("RunUpdate(): want at least one advanced region, got none")
	}
}

// An Update Run must preserve an author's customization untouched and
// report it as preserved, not silently overwrite it (FR-001, FR-002).
func TestUpdateRun_SingleSite_PreservesAuthorCustomization(t *testing.T) {
	projectDir := scaffoldSite(t, "profile", "my-profile-site")

	customized := filepath.Join(projectDir, "package.json")
	customContent := []byte(`{"name":"my-real-author-customization"}`)
	if err := os.WriteFile(customized, customContent, 0o644); err != nil {
		t.Fatalf("simulating author customization: %v", err)
	}

	run, err := migrate.RunUpdate(projectDir, "profile", nil)
	if err != nil {
		t.Fatalf("RunUpdate(): %v", err)
	}

	found := false
	for _, p := range run.Preserved {
		if p == "package.json" {
			found = true
		}
	}
	if !found {
		t.Errorf("RunUpdate(): package.json not reported as preserved; Preserved=%v", run.Preserved)
	}

	got, err := os.ReadFile(customized)
	if err != nil {
		t.Fatalf("reading package.json after update: %v", err)
	}
	if string(got) != string(customContent) {
		t.Errorf("package.json content changed: got %q, want untouched %q", got, customContent)
	}
}

// When a region has genuinely changed on both sides (simulated here by
// hand-editing the manifest's recorded baseline, since this binary embeds
// only one scaffold version and can't naturally diverge from itself), the
// Update Run must surface it as a conflict and refuse to commit until the
// author resolves it — never auto-resolved (FR-002).
func TestUpdateRun_SingleSite_SurfacesGenuineConflict(t *testing.T) {
	projectDir := scaffoldSite(t, "profile", "my-profile-site")

	authorContent := []byte(`{"name":"author-changed-this"}`)
	target := filepath.Join(projectDir, "package.json")
	if err := os.WriteFile(target, authorContent, 0o644); err != nil {
		t.Fatalf("simulating author customization: %v", err)
	}

	manifest, err := scaffold.ReadManifest(projectDir)
	if err != nil {
		t.Fatalf("ReadManifest(): %v", err)
	}
	// Simulate "the scaffold also changed this path between versions" by
	// recording a baseline hash that matches neither the author's edit nor
	// the (single, current) embedded template.
	manifest.Files["package.json"] = scaffold.HashBytes([]byte("a stale baseline neither side matches"))
	if err := scaffold.WriteManifest(projectDir, manifest); err != nil {
		t.Fatalf("WriteManifest(): %v", err)
	}

	run, err := migrate.RunUpdate(projectDir, "profile", nil)
	if err != migrate.ErrUnresolvedConflicts {
		t.Fatalf("RunUpdate() error = %v, want ErrUnresolvedConflicts", err)
	}
	if len(run.Conflicts) != 1 || run.Conflicts[0].Region != "package.json" {
		t.Fatalf("RunUpdate(): Conflicts = %+v, want exactly one for package.json", run.Conflicts)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("reading package.json after a halted run: %v", err)
	}
	if string(got) != string(authorContent) {
		t.Errorf("package.json was modified despite an unresolved conflict: got %q, want %q", got, authorContent)
	}

	// Resolving the conflict must commit using the resolution.
	resolved, err := migrate.RunUpdate(projectDir, "profile", map[string]string{
		"package.json": migrate.ResolveKeepAuthor,
	})
	if err != nil {
		t.Fatalf("RunUpdate() with resolution: %v", err)
	}
	if len(resolved.Conflicts) != 1 || resolved.Conflicts[0].Resolution != migrate.ResolveKeepAuthor {
		t.Errorf("resolved.Conflicts = %+v, want one resolved with %q", resolved.Conflicts, migrate.ResolveKeepAuthor)
	}
}
