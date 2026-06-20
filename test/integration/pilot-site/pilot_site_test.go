// Package pilot_site exercises the full content -> build -> deploy path
// (FR-004; data-model.md "Pilot Site Validation Record"; constitution
// Release Process release-closing test) against a scaffolded profile
// project with real authored content.
//
// What this test automates: scaffold -> index -> build-precondition ->
// deploy-precondition, using the real Go generators and the real `build`/
// `deploy` CLI wrappers (with stub external commands standing in for
// `vite build`/`wrangler deploy`, matching internal/cli's own test
// convention so this suite runs without Node or a Cloudflare account).
//
// What this test does NOT automate, by necessity: publishing to a real
// public address and a human visitor confirming the live site in a browser
// (spec SC-001, SC-003). That leg requires the author's own Cloudflare
// account and one real, fully authored site, and is the actual closing
// evidence for this release — see quickstart.md Scenarios 1-3 for the
// manual steps once a real domain and Cloudflare credentials are in hand.
package pilot_site_test

import (
	"context"
	"encoding/json"
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

// TestPilotSite_ContentToBuildToDeploy runs the full content -> build ->
// deploy path end to end against a scaffolded profile project, confirming
// it completes with no manual intervention beyond the documented commands
// (SC-002) and that the deploy step surfaces the configured public address
// (FR-002, FR-003).
func TestPilotSite_ContentToBuildToDeploy(t *testing.T) {
	// 1. Author content: scaffold a real profile project (FR-004 — the
	// pilot site is a personal-profile/CV project, not a synthetic fixture).
	presets := theme.PresetsForType("profile")
	if len(presets) == 0 {
		t.Fatal(`theme.PresetsForType("profile"): want at least one preset, got none`)
	}

	destDir := t.TempDir()
	projectDir, err := scaffold.Generate(destDir, scaffold.Options{
		ProjectName: "pilot-site",
		SiteType:    "profile",
		Domain:      "pilot.example.com",
		Theme:       presets[0].ID,
		MCPEnabled:  true,
	})
	if err != nil {
		t.Fatalf("scaffold.Generate(): %v", err)
	}

	authorAdditionalContent(t, projectDir)

	// 2. Index: biblio-cli index's actual generator, against every authored
	// page (FR-001's source data).
	contentDir := filepath.Join(projectDir, "content")
	idx, errs := index.Generate(contentDir)
	if len(errs) > 0 {
		t.Fatalf("index.Generate(): %v", errs)
	}
	if len(idx.Entries) < 3 {
		t.Fatalf("index.Generate(): want at least 3 authored entries (seeded + added), got %d", len(idx.Entries))
	}
	writeContentIndex(t, projectDir, idx)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(cwd) }()

	// 3. Build: the real `biblio-cli build` wrapper, with a stub standing in
	// for `vite build` (the actual prerender mechanics are proven separately
	// against real tooling in templates/site/profile's own vitest suite and
	// a manual `vite build` run — this verifies the CLI wrapper's own
	// precondition/success-reporting contract against a real scaffolded
	// project, not vite/prerender correctness itself).
	originalBuildCmd := cli.BuildCommand
	cli.BuildCommand = []string{"sh", "-c", "mkdir -p dist/client && exit 0"}
	defer func() { cli.BuildCommand = originalBuildCmd }()

	buildCmd := &cobra.Command{}
	buildCmd.SetContext(context.Background())
	var buildOut strings.Builder
	buildCmd.SetOut(&buildOut)
	if err := cli.RunBuild(buildCmd, nil); err != nil {
		t.Fatalf("RunBuild(): %v, output:\n%s", err, buildOut.String())
	}
	if !strings.Contains(buildOut.String(), "Build complete") {
		t.Errorf("RunBuild() output = %q, want a complete-build report", buildOut.String())
	}

	// 4. Deploy: the real `biblio-cli deploy` wrapper, with a stub standing
	// in for `wrangler deploy`.
	originalDeployCmd := cli.DeployCommand
	cli.DeployCommand = []string{"sh", "-c", "exit 0"}
	defer func() { cli.DeployCommand = originalDeployCmd }()

	deployCmd := &cobra.Command{}
	deployCmd.SetContext(context.Background())
	var deployOut strings.Builder
	deployCmd.SetOut(&deployOut)
	if err := cli.RunDeploy(deployCmd, nil); err != nil {
		t.Fatalf("RunDeploy(): %v, output:\n%s", err, deployOut.String())
	}

	// FR-002/FR-003: deploy surfaces the public address pre-filled from
	// --domain at scaffold time.
	if !strings.Contains(deployOut.String(), "pilot.example.com") {
		t.Errorf("RunDeploy() output = %q, want it to surface the configured public address", deployOut.String())
	}
}

// authorAdditionalContent adds real-looking authored pages beyond the
// scaffold's seed content, so this test exercises more than the bare
// minimum a fresh `biblio-cli new` produces.
func authorAdditionalContent(t *testing.T, projectDir string) {
	t.Helper()
	experience := filepath.Join(projectDir, "content", "projects", "second-project.md")
	body := []byte(`---
title: "Second Project"
summary: "Another project worth showcasing."
url: "https://example.com/second"
---

Details about the second project go here.
`)
	if err := os.WriteFile(experience, body, 0o644); err != nil {
		t.Fatalf("authoring additional content: %v", err)
	}
}

// writeContentIndex mirrors `biblio-cli index`'s own write step (the real
// generator already validated the content above; this writes its output to
// disk exactly where `react-router.config.ts`'s prerender step reads it
// from).
func writeContentIndex(t *testing.T, projectDir string, idx index.Index) {
	t.Helper()
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		t.Fatalf("marshalling content-index.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "content-index.json"), data, 0o644); err != nil {
		t.Fatalf("writing content-index.json: %v", err)
	}
}
