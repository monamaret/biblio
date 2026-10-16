// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// Package crosstypereview — terminal-rendering half of the Cross-Type
// Consistency Record (FR-004, Scenario 3, SC-002). The TS-side negotiation
// logic itself is exercised per type by each template's own
// app/negotiation/accept.test.ts (profile, log, docs all have one) — this
// Go-side review instead confirms the two structural guarantees a
// single-process Go test can check quickly and deterministically, without
// needing to npm-install and execute three separate Node projects just to
// re-prove what each template's own vitest suite already proves:
//
//  1. The negotiation source itself (accept.ts/ansi.ts/markdown.ts/wasm.d.ts)
//     is byte-identical across all three site types — there is one negotiation
//     owner, never a per-type fork (Principle II) — so "each site renders a
//     clean Markdown/ANSI view" can never silently diverge by code drift.
//  2. Each site type's actual seed content carries the field shape
//     appropriate to its own content model (profile's role/summary, log's
//     date/tags, docs's section/order) — confirming the *content* differs
//     appropriately per type even though the *rendering code* deliberately
//     does not (FR-004 — "appropriate to its own content shape").
package crosstypereview_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/monamaret/biblio/internal/index"
	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
)

func scaffoldProject(t *testing.T, siteType string) string {
	t.Helper()
	presets := theme.PresetsForType(siteType)
	if len(presets) == 0 {
		t.Fatalf("theme.PresetsForType(%q): no presets defined", siteType)
	}
	destDir := t.TempDir()
	projectDir, err := scaffold.Generate(destDir, scaffold.Options{
		ProjectName: "review-" + siteType,
		SiteType:    siteType,
		Domain:      "example.com",
		Theme:       presets[0].ID,
		MCPEnabled:  true,
	})
	if err != nil {
		t.Fatalf("scaffold.Generate(%q): %v", siteType, err)
	}
	return projectDir
}

// TestCrossTypeTerminalRendering_OneSharedNegotiationOwner confirms the
// negotiation logic is never forked per type (Principle II) — a precondition
// for "each site returns a clean, correct, type-appropriate rendering" ever
// being a property of the shared code rather than something that could
// silently diverge between types.
func TestCrossTypeTerminalRendering_OneSharedNegotiationOwner(t *testing.T) {
	siteTypes := scaffold.SiteTypes
	files := []string{
		filepath.Join("app", "negotiation", "accept.ts"),
		filepath.Join("app", "negotiation", "ansi.ts"),
		filepath.Join("app", "negotiation", "wasm.d.ts"),
	}

	projectDirs := map[string]string{}
	for _, siteType := range siteTypes {
		projectDirs[siteType] = scaffoldProject(t, siteType)
	}

	for _, file := range files {
		var reference []byte
		var referenceType string
		for _, siteType := range siteTypes {
			data, err := os.ReadFile(filepath.Join(projectDirs[siteType], file))
			if err != nil {
				t.Fatalf("reading %s for %q: %v", file, siteType, err)
			}
			if reference == nil {
				reference, referenceType = data, siteType
				continue
			}
			if !bytes.Equal(reference, data) {
				t.Errorf("%s differs between %q and %q — negotiation logic must have exactly one owner (Principle II), never a per-type fork", file, referenceType, siteType)
			}
		}
	}
}

// TestCrossTypeTerminalRendering_ContentShapeIsAppropriatePerType confirms
// each site type's seed content carries the field shape appropriate to its
// own content model — the "type-appropriate" half of FR-004/Scenario 3 that
// the (deliberately type-agnostic) negotiation/rendering code itself cannot
// express; this is the contract that markdown.ts/accept.ts render correctly
// *because* each type's content-index.json already carries the right shape.
func TestCrossTypeTerminalRendering_ContentShapeIsAppropriatePerType(t *testing.T) {
	cases := []struct {
		siteType     string
		wantFields   []string
		forbidFields []string
	}{
		{siteType: "profile", wantFields: []string{"role", "summary"}, forbidFields: []string{"date", "section", "order"}},
		{siteType: "log", wantFields: []string{"date", "tags"}, forbidFields: []string{"section", "order"}},
		{siteType: "docs", wantFields: []string{"section", "order"}, forbidFields: []string{"date", "tags", "project"}},
		{siteType: "quarterly", wantFields: []string{"kind", "issue", "category", "footnotes", "margin_notes", "markdown_body"}, forbidFields: []string{"section", "project", "role"}},
	}
	if len(cases) != len(scaffold.SiteTypes) {
		t.Fatalf("content-shape cases cover %d site types, scaffold.SiteTypes has %d — add the new type's expected shape", len(cases), len(scaffold.SiteTypes))
	}

	for _, c := range cases {
		c := c
		t.Run(c.siteType, func(t *testing.T) {
			projectDir := scaffoldProject(t, c.siteType)
			contentDir := filepath.Join(projectDir, "content")
			idx, errs := index.Generate(contentDir)
			if len(errs) > 0 {
				t.Fatalf("index.Generate(%q): %v", c.siteType, errs)
			}
			if len(idx.Entries) == 0 {
				t.Fatalf("no entries generated for %q", c.siteType)
			}

			foundAny := map[string]bool{}
			for _, entry := range idx.Entries {
				for _, want := range c.wantFields {
					if _, ok := entry.Fields[want]; ok {
						foundAny[want] = true
					}
				}
				for _, forbid := range c.forbidFields {
					if _, ok := entry.Fields[forbid]; ok {
						t.Errorf("%q entry %q unexpectedly has field %q — not part of this site type's content shape", c.siteType, entry.Slug, forbid)
					}
				}
			}
			for _, want := range c.wantFields {
				if !foundAny[want] {
					t.Errorf("%q seed content never carries field %q — expected at least one entry with it", c.siteType, want)
				}
			}
		})
	}
}
