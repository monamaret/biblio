// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// conflict_test exercises internal/migrate's conflict surfacing: every
// ambiguous region becomes an Update Conflict that is never auto-resolved
// (data-model.md "Update Conflict"; FR-002), and Apply only ever commits
// once every conflict carries a recorded resolution (contracts/
// update-and-escalation.md Step 4-5).
package migrate_test

import (
	"testing"

	"github.com/monamaret/biblio/internal/migrate"
)

// Conflicts must extract exactly one Update Conflict per ambiguous region,
// carrying both sides of the diff, and skip non-ambiguous regions entirely.
func TestConflicts_OnePerAmbiguousRegion(t *testing.T) {
	regions := []migrate.Region{
		{Path: "package.json", Classification: migrate.ScaffoldTracked, New: []byte("new")},
		{
			Path:           "workers/entry.ts",
			Classification: migrate.Ambiguous,
			Current:        []byte("author version"),
			New:            []byte("scaffold version"),
		},
	}

	conflicts := migrate.Conflicts(regions)
	if len(conflicts) != 1 {
		t.Fatalf("Conflicts(): got %d conflicts, want 1", len(conflicts))
	}

	c := conflicts[0]
	if c.Region != "workers/entry.ts" {
		t.Errorf("conflict.Region = %q, want %q", c.Region, "workers/entry.ts")
	}
	if string(c.AuthorCurrent) != "author version" {
		t.Errorf("conflict.AuthorCurrent = %q, want %q", c.AuthorCurrent, "author version")
	}
	if string(c.ScaffoldProposed) != "scaffold version" {
		t.Errorf("conflict.ScaffoldProposed = %q, want %q", c.ScaffoldProposed, "scaffold version")
	}
	if c.Resolution != "" {
		t.Errorf("conflict.Resolution = %q, want empty (unresolved) until the author decides", c.Resolution)
	}
}

// Apply must refuse to commit any region while a conflict for that region
// remains unresolved — the flow waits for an explicit resolution and never
// auto-resolves (FR-002, data-model.md "Update Conflict" validation rules).
func TestApply_UnresolvedConflict_NeverReachesCommit(t *testing.T) {
	projectDir := newScaffoldedProject(t)

	regions := []migrate.Region{
		{Path: "package.json", Classification: migrate.Ambiguous, Current: []byte("author"), New: []byte("scaffold")},
	}
	conflicts := migrate.Conflicts(regions)
	if len(conflicts) != 1 {
		t.Fatalf("Conflicts(): got %d, want 1", len(conflicts))
	}

	if err := migrate.ApplyWithResolutions(projectDir, regions, conflicts); err == nil {
		t.Error("ApplyWithResolutions() with an unresolved conflict: want error, got nil")
	}
}

// Once a conflict carries a resolution, ApplyWithResolutions must commit
// using that resolution's content for the resolved region.
func TestApply_ResolvedConflict_CommitsResolution(t *testing.T) {
	projectDir := newScaffoldedProject(t)

	regions := []migrate.Region{
		{Path: "package.json", Classification: migrate.Ambiguous, Current: []byte("author"), New: []byte("scaffold")},
	}
	conflicts := migrate.Conflicts(regions)
	conflicts[0].Resolution = "keep-author"

	if err := migrate.ApplyWithResolutions(projectDir, regions, conflicts); err != nil {
		t.Fatalf("ApplyWithResolutions() with a resolved conflict: %v", err)
	}
}
