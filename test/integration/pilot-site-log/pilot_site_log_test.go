// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// Package pilot_site_log proves the project-log site type's central
// generalization claim: adding a new day's entry uses the exact same
// authoring flow as the previous one, with zero one-off structural or
// config adjustments, repeatable across several consecutive entries
// (FR-004; spec Scenario 3, SC-002).
//
// What this test automates: scaffold once, then for three separate days —
// author one new dated entry, index, build (stubbed `vite build`, matching
// test/integration/pilot-site's own convention so this suite runs without
// Node) — asserting wrangler.jsonc, content/site.yaml, and package.json
// are byte-for-byte unchanged between cycles. If any cycle required a
// structural or config change to publish, that diff would show up here.
package pilot_site_log_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/monamaret/biblio/internal/cli"
	"github.com/monamaret/biblio/internal/index"
	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
)

// snapshot captures the bytes of the files this release's authoring-rhythm
// claim says must never change between consecutive entries.
type snapshot struct {
	wranglerJSONC string
	siteYAML      string
	packageJSON   string
}

func takeSnapshot(t *testing.T, projectDir string) snapshot {
	t.Helper()
	read := func(rel string) string {
		data, err := os.ReadFile(filepath.Join(projectDir, rel))
		if err != nil {
			t.Fatalf("reading %q: %v", rel, err)
		}
		return string(data)
	}
	return snapshot{
		wranglerJSONC: read("wrangler.jsonc"),
		siteYAML:      read(filepath.Join("content", "site.yaml")),
		packageJSON:   read("package.json"),
	}
}

// TestAuthoringRhythm_ThreeConsecutiveEntries scaffolds a project-log site
// once, then repeats "author a new dated entry -> index -> build" for
// three separate days, asserting zero structural/config drift between
// cycles (FR-004, SC-002; quickstart.md Scenario 3).
func TestAuthoringRhythm_ThreeConsecutiveEntries(t *testing.T) {
	presets := theme.PresetsForType("log")
	if len(presets) == 0 {
		t.Fatal(`theme.PresetsForType("log"): want at least one preset, got none`)
	}

	destDir := t.TempDir()
	projectDir, err := scaffold.Generate(destDir, scaffold.Options{
		ProjectName: "pilot-log",
		SiteType:    "log",
		Domain:      "pilot-log.example.com",
		Theme:       presets[0].ID,
		MCPEnabled:  true,
	})
	if err != nil {
		t.Fatalf("scaffold.Generate(): %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(cwd) }()

	originalBuildCmd := cli.BuildCommand
	cli.BuildCommand = []string{"sh", "-c", "mkdir -p dist/client && exit 0"}
	defer func() { cli.BuildCommand = originalBuildCmd }()

	var baseline *snapshot
	dates := []string{"2026-06-15", "2026-06-16", "2026-06-17"}

	for i, date := range dates {
		day := i + 1
		entryPath := filepath.Join(projectDir, "content", "posts", fmt.Sprintf("%s-day-%d-update.md", date, day))
		body := fmt.Sprintf("---\ntitle: \"Day %d update\"\ndate: \"%s\"\ntags:\n  - daily\n---\n\nProgress notes for day %d.\n", day, date, day)
		if err := os.WriteFile(entryPath, []byte(body), 0o644); err != nil {
			t.Fatalf("day %d: authoring entry: %v", day, err)
		}

		idx, errs := index.Generate(filepath.Join(projectDir, "content"))
		if len(errs) > 0 {
			t.Fatalf("day %d: index.Generate(): %v", day, errs)
		}
		const seededEntries = 3 // templates/site/log/content/posts/ ships 3 seed entries (T020)
		want := seededEntries + day
		if len(idx.Entries) != want {
			t.Fatalf("day %d: index.Generate(): want %d entries (seeded + %d authored), got %d", day, want, day, len(idx.Entries))
		}

		buildCmd := &cobra.Command{}
		buildCmd.SetContext(context.Background())
		var buildOut strings.Builder
		buildCmd.SetOut(&buildOut)
		if err := cli.RunBuild(buildCmd, nil); err != nil {
			t.Fatalf("day %d: RunBuild(): %v, output:\n%s", day, err, buildOut.String())
		}
		if !strings.Contains(buildOut.String(), "Build complete") {
			t.Errorf("day %d: RunBuild() output = %q, want a complete-build report", day, buildOut.String())
		}

		current := takeSnapshot(t, projectDir)
		if baseline == nil {
			baseline = &current
			continue
		}
		if current.wranglerJSONC != baseline.wranglerJSONC {
			t.Errorf("day %d: wrangler.jsonc changed between cycles — authoring a new entry must require zero structural adjustments (FR-004)", day)
		}
		if current.siteYAML != baseline.siteYAML {
			t.Errorf("day %d: content/site.yaml changed between cycles (FR-004)", day)
		}
		if current.packageJSON != baseline.packageJSON {
			t.Errorf("day %d: package.json changed between cycles (FR-004)", day)
		}
	}
}
