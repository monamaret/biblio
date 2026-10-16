// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package update_three_sites_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/monamaret/biblio/internal/migrate"
)

// An interrupted update (simulated here as a panic mid-run, the same way a
// killed process would leave things) must never leave a half-written
// file — Apply writes each region's full content in one os.WriteFile call,
// so any region untouched when the run stops is exactly as it was before
// (spec Edge Cases "does the site remain in a consistent, recoverable
// state").
func TestUpdateRun_Interruption_NoHalfWrittenFiles(t *testing.T) {
	projectDir := scaffoldSite(t, "profile", "my-profile-site")

	customized := filepath.Join(projectDir, "package.json")
	customContent := []byte(`{"name":"author-customization-survives-interruption"}`)
	if err := os.WriteFile(customized, customContent, 0o644); err != nil {
		t.Fatalf("simulating author customization: %v", err)
	}

	regions, err := migrate.ClassifySite(projectDir, "profile")
	if err != nil {
		t.Fatalf("ClassifySite(): %v", err)
	}

	// Simulate an interrupted Apply: only commit the first region, as if
	// the process had been killed right after it (deliberately exercising
	// a partial run rather than the full RunUpdate path).
	if len(regions) == 0 {
		t.Fatal("ClassifySite(): want at least one region, got none")
	}
	if err := migrate.Apply(projectDir, regions[:1]); err != nil {
		t.Fatalf("Apply() (simulated partial run): %v", err)
	}

	// Whatever happened to region[0], every other region's file on disk
	// must be byte-identical to what it was before the interrupted run —
	// nothing half-migrated.
	for _, r := range regions[1:] {
		path := filepath.Join(projectDir, filepath.FromSlash(r.Path))
		got, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) && r.Current == nil {
				continue // region never existed on disk; consistent
			}
			t.Fatalf("reading %q after interruption: %v", r.Path, err)
		}
		if string(got) != string(r.Current) {
			t.Errorf("region %q changed despite not being committed: got %q, want untouched %q", r.Path, got, r.Current)
		}
	}

	// The author's customization specifically must have survived,
	// regardless of where the interruption landed.
	got, err := os.ReadFile(customized)
	if err != nil {
		t.Fatalf("reading package.json after interruption: %v", err)
	}
	if string(got) != string(customContent) {
		t.Errorf("author customization lost after interruption: got %q, want %q", got, customContent)
	}

	// Re-running the update after the interruption must succeed cleanly —
	// the site is recoverable, not stuck (spec Edge Cases).
	run, err := migrate.RunUpdate(projectDir, "profile", nil)
	if err != nil {
		t.Fatalf("RunUpdate() resuming after interruption: %v", err)
	}
	if len(run.Conflicts) != 0 {
		t.Errorf("RunUpdate() resuming after interruption: got %d conflicts, want 0", len(run.Conflicts))
	}
}

// A run with any unresolved Update Conflict must never reach the commit
// step at all — confirmed end to end here against a real scaffolded site
// (contracts/update-and-escalation.md Step 5; unit-level coverage already
// exists in test/unit/migrate/conflict_test.go).
func TestUpdateRun_UnresolvedConflict_NeverCommits(t *testing.T) {
	projectDir := scaffoldSite(t, "profile", "my-profile-site")

	before, err := os.ReadFile(filepath.Join(projectDir, "package.json"))
	if err != nil {
		t.Fatalf("reading package.json: %v", err)
	}

	regions := []migrate.Region{
		{Path: "package.json", Classification: migrate.Ambiguous, Current: before, New: []byte("would-be-scaffold-content")},
	}
	if err := migrate.ApplyWithResolutions(projectDir, regions, migrate.Conflicts(regions)); err == nil {
		t.Fatal("ApplyWithResolutions() with an unresolved conflict: want error, got nil")
	}

	after, err := os.ReadFile(filepath.Join(projectDir, "package.json"))
	if err != nil {
		t.Fatalf("reading package.json after a rejected commit: %v", err)
	}
	if string(after) != string(before) {
		t.Errorf("package.json was modified despite a rejected commit: got %q, want untouched %q", after, before)
	}
}
