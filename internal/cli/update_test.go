// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package cli_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/monamaret/biblio/internal/cli"
	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
)

// scaffoldedUpdateProject scaffolds a real profile project and chdirs into
// it, so `biblio-cli update` can find content/site.yaml in its working
// directory exactly as it would for a real author.
func scaffoldedUpdateProject(t *testing.T) string {
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
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: getwd: %v", err)
	}
	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("setup: chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	return projectDir
}

// updateCmdForTest builds a bare *cobra.Command carrying the same
// --resolve flag updateCmd registers, since these tests call RunUpdate
// directly rather than through updateCmd.Execute().
func updateCmdForTest(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(out)
	if errOut != nil {
		cmd.SetErr(errOut)
	}
	cmd.SetContext(context.Background())
	cmd.Flags().StringArray("resolve", nil, "")
	return cmd
}

// RunUpdate exits 0 and reports the preserved/advanced output on a clean
// run with no conflicts (contracts/update-and-escalation.md `biblio-cli
// update` "Output (success)").
func TestRunUpdate_CleanRun_ExitsZeroAndReportsOutput(t *testing.T) {
	scaffoldedUpdateProject(t)

	var out bytes.Buffer
	cmd := updateCmdForTest(&out, nil)

	if err := cli.RunUpdate(cmd, nil); err != nil {
		t.Fatalf("RunUpdate() error = %v, want nil", err)
	}
	if !strings.Contains(out.String(), "advanced") {
		t.Errorf("RunUpdate() output = %q, want it to report advanced/preserved counts", out.String())
	}
}

// RunUpdate exits non-zero (a non-nil error) and reports recovery guidance
// when a region is ambiguous and unresolved — never silently overwriting
// it (FR-002; contracts/update-and-escalation.md "Output (failure /
// interruption)").
func TestRunUpdate_UnresolvedConflict_ExitsNonZeroWithGuidance(t *testing.T) {
	projectDir := scaffoldedUpdateProject(t)

	authorContent := []byte(`{"name":"author-changed-this"}`)
	target := filepath.Join(projectDir, "package.json")
	if err := os.WriteFile(target, authorContent, 0o644); err != nil {
		t.Fatalf("simulating author customization: %v", err)
	}
	manifest, err := scaffold.ReadManifest(projectDir)
	if err != nil {
		t.Fatalf("ReadManifest(): %v", err)
	}
	manifest.Files["package.json"] = scaffold.HashBytes([]byte("a stale baseline neither side matches"))
	if err := scaffold.WriteManifest(projectDir, manifest); err != nil {
		t.Fatalf("WriteManifest(): %v", err)
	}

	var out, errOut bytes.Buffer
	cmd := updateCmdForTest(&out, &errOut)

	err = cli.RunUpdate(cmd, nil)
	if err == nil {
		t.Fatal("RunUpdate() with an unresolved conflict: want error, got nil")
	}
	if !strings.Contains(errOut.String(), "package.json") {
		t.Errorf("RunUpdate() stderr = %q, want it to name the conflicting region", errOut.String())
	}

	got, readErr := os.ReadFile(target)
	if readErr != nil {
		t.Fatalf("reading package.json: %v", readErr)
	}
	if string(got) != string(authorContent) {
		t.Errorf("package.json was modified despite an unresolved conflict: got %q, want %q", got, authorContent)
	}
}

// RunUpdate exits 0 once --resolve settles every conflict.
func TestRunUpdate_ResolveFlag_SettlesConflictAndCommits(t *testing.T) {
	projectDir := scaffoldedUpdateProject(t)

	target := filepath.Join(projectDir, "package.json")
	if err := os.WriteFile(target, []byte(`{"name":"author-changed-this"}`), 0o644); err != nil {
		t.Fatalf("simulating author customization: %v", err)
	}
	manifest, err := scaffold.ReadManifest(projectDir)
	if err != nil {
		t.Fatalf("ReadManifest(): %v", err)
	}
	manifest.Files["package.json"] = scaffold.HashBytes([]byte("a stale baseline neither side matches"))
	if err := scaffold.WriteManifest(projectDir, manifest); err != nil {
		t.Fatalf("WriteManifest(): %v", err)
	}

	var out bytes.Buffer
	cmd := updateCmdForTest(&out, nil)
	if err := cmd.Flags().Set("resolve", "package.json=keep-author"); err != nil {
		t.Fatalf("setting --resolve flag: %v", err)
	}

	if err := cli.RunUpdate(cmd, nil); err != nil {
		t.Fatalf("RunUpdate() with --resolve: error = %v, want nil", err)
	}
}
