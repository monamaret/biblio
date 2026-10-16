// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package update_three_sites_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/monamaret/biblio/internal/migrate"
)

// The exact same Update Run contract that holds for a profile site must
// hold, unmodified, for a project-log site (FR-003, Scenario 2).
func TestUpdateRun_ProjectLogSite_AdvancesAndPreserves(t *testing.T) {
	projectDir := scaffoldSite(t, "log", "my-project-log")

	customized := filepath.Join(projectDir, "package.json")
	customContent := []byte(`{"name":"log-author-customization"}`)
	if err := os.WriteFile(customized, customContent, 0o644); err != nil {
		t.Fatalf("simulating author customization: %v", err)
	}

	run, err := migrate.RunUpdate(projectDir, "log", nil)
	if err != nil {
		t.Fatalf("RunUpdate(): %v", err)
	}
	if len(run.Conflicts) != 0 {
		t.Errorf("RunUpdate(): got %d conflicts on a log site, want 0", len(run.Conflicts))
	}
	if len(run.Advanced) == 0 {
		t.Error("RunUpdate(): want at least one advanced region, got none")
	}

	got, err := os.ReadFile(customized)
	if err != nil {
		t.Fatalf("reading package.json after update: %v", err)
	}
	if string(got) != string(customContent) {
		t.Errorf("package.json content changed: got %q, want untouched %q", got, customContent)
	}
}
