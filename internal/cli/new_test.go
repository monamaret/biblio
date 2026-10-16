// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/monamaret/biblio/internal/cli"
	"github.com/monamaret/biblio/internal/theme"
)

// presetIDFor returns a ready-to-use preset ID for the given site type, or
// fails the test if none is defined yet (T011 defines them).
func presetIDFor(t *testing.T, siteType string) string {
	t.Helper()
	presets := theme.PresetsForType(siteType)
	if len(presets) == 0 {
		t.Fatalf("theme.PresetsForType(%q): want at least one ready-to-use preset, got none", siteType)
	}
	return presets[0].ID
}

// An unknown --type must be rejected with an error naming the bad value
// (contracts/cli-commands.md `biblio-cli new` "Output (failure)").
func TestValidateNewFlags_UnknownType(t *testing.T) {
	destDir := t.TempDir()

	err := cli.ValidateNewFlags("my-site", "blog", "example.com", "", destDir)
	if err == nil {
		t.Fatal("ValidateNewFlags(): want error for unknown --type, got nil")
	}
	if !strings.Contains(err.Error(), "blog") {
		t.Errorf("ValidateNewFlags() error = %q, want it to name the unknown --type %q", err, "blog")
	}
}

// A --theme that doesn't match the chosen --type's presets must be rejected
// (FR-004; contracts/cli-commands.md `biblio-cli new` "Output (failure)").
func TestValidateNewFlags_ThemeNotValidForType(t *testing.T) {
	destDir := t.TempDir()
	mismatchedTheme := presetIDFor(t, "log")

	err := cli.ValidateNewFlags("my-site", "profile", "example.com", mismatchedTheme, destDir)
	if err == nil {
		t.Fatal("ValidateNewFlags(): want error for --theme not valid for --type, got nil")
	}
	if !strings.Contains(err.Error(), mismatchedTheme) {
		t.Errorf("ValidateNewFlags() error = %q, want it to name the mismatched --theme %q", err, mismatchedTheme)
	}
	if !strings.Contains(err.Error(), "profile") {
		t.Errorf("ValidateNewFlags() error = %q, want it to name the chosen --type %q", err, "profile")
	}
}

// A project-name colliding with an existing directory must be rejected
// (data-model.md "Site Project" validation rules).
func TestValidateNewFlags_CollidingProjectName(t *testing.T) {
	destDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(destDir, "my-site"), 0o755); err != nil {
		t.Fatalf("setup: creating colliding directory: %v", err)
	}

	err := cli.ValidateNewFlags("my-site", "profile", "example.com", presetIDFor(t, "profile"), destDir)
	if err == nil {
		t.Fatal("ValidateNewFlags(): want error for colliding project-name, got nil")
	}
	if !strings.Contains(err.Error(), "my-site") {
		t.Errorf("ValidateNewFlags() error = %q, want it to name the colliding project-name %q", err, "my-site")
	}
}

// A fully valid combination must pass validation cleanly.
func TestValidateNewFlags_Valid(t *testing.T) {
	destDir := t.TempDir()

	err := cli.ValidateNewFlags("my-site", "profile", "example.com", presetIDFor(t, "profile"), destDir)
	if err != nil {
		t.Errorf("ValidateNewFlags() error = %v, want nil for a valid combination", err)
	}
}
