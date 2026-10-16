// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// guidance_test exercises internal/escalation's Storage-Escalation
// Guidance Check: threshold-status computation against a real
// content-index.json's derived scale (data-model.md "Storage-Escalation
// Guidance Check"; FR-004, FR-005).
package escalation_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/monamaret/biblio/internal/escalation"
)

// writeContentIndex writes a content-index.json with entryCount synthetic
// entries into a temp dir and returns its path.
func writeContentIndex(t *testing.T, entryCount int) string {
	t.Helper()
	entries := make([]map[string]any, entryCount)
	for i := range entries {
		entries[i] = map[string]any{"slug": "entry", "path": "entry.md", "body": "hello"}
	}
	data, err := json.Marshal(map[string]any{"entries": entries})
	if err != nil {
		t.Fatalf("marshalling synthetic content-index.json: %v", err)
	}
	path := filepath.Join(t.TempDir(), "content-index.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("writing synthetic content-index.json: %v", err)
	}
	return path
}

// A site well under both documented thresholds must report not-crossed for
// both conditional paths.
func TestCheckSite_BelowBothThresholds_NotCrossed(t *testing.T) {
	path := writeContentIndex(t, 5)

	check, err := escalation.CheckSite(path, "my-profile-site")
	if err != nil {
		t.Fatalf("CheckSite(): %v", err)
	}
	if check.D1Status != escalation.NotCrossed {
		t.Errorf("D1Status = %q, want %q", check.D1Status, escalation.NotCrossed)
	}
	if check.VectorizeStatus != escalation.NotCrossed {
		t.Errorf("VectorizeStatus = %q, want %q", check.VectorizeStatus, escalation.NotCrossed)
	}
}

// A site past the documented D1 entry-count threshold, but below the
// (higher) Vectorize threshold, must report D1 crossed and Vectorize not.
func TestCheckSite_PastD1Threshold_OnlyD1Crossed(t *testing.T) {
	path := writeContentIndex(t, escalation.D1EntryCountThreshold+1)

	check, err := escalation.CheckSite(path, "my-project-log")
	if err != nil {
		t.Fatalf("CheckSite(): %v", err)
	}
	if check.D1Status != escalation.Crossed {
		t.Errorf("D1Status = %q, want %q", check.D1Status, escalation.Crossed)
	}
	if check.VectorizeStatus != escalation.NotCrossed {
		t.Errorf("VectorizeStatus = %q, want %q", check.VectorizeStatus, escalation.NotCrossed)
	}
}

// A site past the documented Vectorize entry-count threshold must report
// both paths crossed (Vectorize's threshold is documented as the higher of
// the two — data-model.md "Storage-Escalation Guidance Check").
func TestCheckSite_PastVectorizeThreshold_BothCrossed(t *testing.T) {
	path := writeContentIndex(t, escalation.VectorizeEntryCountThreshold+1)

	check, err := escalation.CheckSite(path, "my-docs-site")
	if err != nil {
		t.Fatalf("CheckSite(): %v", err)
	}
	if check.D1Status != escalation.Crossed {
		t.Errorf("D1Status = %q, want %q", check.D1Status, escalation.Crossed)
	}
	if check.VectorizeStatus != escalation.Crossed {
		t.Errorf("VectorizeStatus = %q, want %q", check.VectorizeStatus, escalation.Crossed)
	}
}

// CheckSite must reflect the real entry count, not a synthetic projection
// (FR-004).
func TestCheckSite_ReportsRealEntryCount(t *testing.T) {
	path := writeContentIndex(t, 42)

	check, err := escalation.CheckSite(path, "my-profile-site")
	if err != nil {
		t.Fatalf("CheckSite(): %v", err)
	}
	if check.EntryCount != 42 {
		t.Errorf("EntryCount = %d, want 42", check.EntryCount)
	}
	if check.Site != "my-profile-site" {
		t.Errorf("Site = %q, want %q", check.Site, "my-profile-site")
	}
}
