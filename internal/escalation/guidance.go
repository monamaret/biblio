// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// Package escalation checks a deployed Site Project's real, current scale
// against the documented D1 (storage volume / relational-query) and
// Vectorize (search-quality) escalation thresholds (FR-004, FR-005),
// without adopting either as a running dependency (Principle V).
package escalation

import (
	"encoding/json"
	"fmt"
	"os"
)

// D1EntryCountThreshold and VectorizeEntryCountThreshold are the
// documented, checkable entry-count thresholds for the two conditional
// escalation paths (docs/storage-escalation-guide.md "Thresholds"). They
// are content-index entry-count proxies, not the qualitative judgments
// Principle V actually names: D1 escalates on real volume/relational-query
// growth, and Vectorize escalates only when search *quality* demonstrably
// fails on conceptual queries — a thing volume alone cannot prove. Crossing
// VectorizeEntryCountThreshold is a prompt to go check real query quality
// by hand, not a substitute for that judgment (see the guide).
const (
	D1EntryCountThreshold        = 500
	VectorizeEntryCountThreshold = 2000
)

// ThresholdStatus is whether a site has crossed a documented escalation
// threshold (data-model.md "Storage-Escalation Guidance Check").
type ThresholdStatus string

const (
	NotCrossed ThresholdStatus = "not-crossed"
	Crossed    ThresholdStatus = "crossed"
)

// GuidanceCheck is the documented-threshold check run against one real
// site's actual current scale (data-model.md "Storage-Escalation Guidance
// Check"; FR-004, FR-005, Scenario 3).
type GuidanceCheck struct {
	Site            string
	EntryCount      int
	D1Status        ThresholdStatus
	VectorizeStatus ThresholdStatus
}

// contentIndex is the minimal shape of content-index.json this package
// needs — entry count, not full entry contents.
type contentIndex struct {
	Entries []json.RawMessage `json:"entries"`
}

// CheckSite reads the content-index.json at contentIndexPath and reports
// its threshold-status for both conditional escalation paths, correctly
// reflecting the site's real, current scale (FR-004) — never a synthetic
// projection.
func CheckSite(contentIndexPath, site string) (GuidanceCheck, error) {
	data, err := os.ReadFile(contentIndexPath)
	if err != nil {
		return GuidanceCheck{}, fmt.Errorf("escalation: reading %s: %w", contentIndexPath, err)
	}
	var idx contentIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return GuidanceCheck{}, fmt.Errorf("escalation: decoding %s: %w", contentIndexPath, err)
	}

	count := len(idx.Entries)
	check := GuidanceCheck{Site: site, EntryCount: count}

	if count > D1EntryCountThreshold {
		check.D1Status = Crossed
	} else {
		check.D1Status = NotCrossed
	}
	if count > VectorizeEntryCountThreshold {
		check.VectorizeStatus = Crossed
	} else {
		check.VectorizeStatus = NotCrossed
	}

	return check, nil
}
