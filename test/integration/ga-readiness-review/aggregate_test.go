// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// Package ga_readiness_review exercises the closing GA Readiness Review
// artifact: cross-referencing every roadmap deliverable §1-§10 against
// this release's acceptance bar, aggregating the Update Run suite
// (US1/US2) and Storage-Escalation Guidance Checks (US3) by reference —
// never re-running them (FR-006, Scenario 4, SC-002, SC-003).
package ga_readiness_review_test

import (
	"testing"

	"github.com/monamaret/biblio/internal/escalation"
)

var (
	passingUpdateEvidence     = escalation.Evidence{Reference: "test/integration/update-three-sites", Passed: true}
	passingEscalationEvidence = escalation.Evidence{Reference: "test/integration/escalation-guidance", Passed: true}
)

// Evaluate must require every deliverable to meet its bar before
// determining GA — one incomplete, non-conditional deliverable is enough
// to withhold ga-determination.
func TestEvaluate_OneIncompleteDeliverable_WithholdsGA(t *testing.T) {
	deliverables := []escalation.Deliverable{
		{ID: "§1", Name: "biblio-cli foundations", Status: escalation.CompleteAndDemonstrated, Reference: "test/integration/pilot-site"},
		{ID: "§9", Name: "Replicate to remaining site types", Status: ""},
	}

	review := escalation.Evaluate(deliverables, passingUpdateEvidence, passingEscalationEvidence)
	if review.GADetermination {
		t.Error("Evaluate(): GADetermination = true with an incomplete deliverable, want false")
	}
}

// Evaluate must accept proven-executable-on-demand as meeting the bar for
// the two explicitly-conditional storage-escalation deliverables
// specifically (FR-006, spec Assumptions "What GA means here") — it must
// NOT require them to be complete-and-demonstrated like every other
// deliverable.
func TestEvaluate_ConditionalDeliverableProvenExecutable_MeetsBar(t *testing.T) {
	deliverables := []escalation.Deliverable{
		{ID: "§1", Name: "biblio-cli foundations", Status: escalation.CompleteAndDemonstrated, Reference: "test/integration/pilot-site"},
		{ID: "§10", Name: "Storage escalation (conditional)", Status: escalation.ProvenExecutableOnDemand, Reference: "docs/storage-escalation-guide.md"},
	}

	review := escalation.Evaluate(deliverables, passingUpdateEvidence, passingEscalationEvidence)
	if !review.GADetermination {
		t.Error("Evaluate(): GADetermination = false with all deliverables meeting their bar, want true")
	}
}

// Evaluate must withhold GA if the update-evidence or escalation-evidence
// aggregate itself failed, even when every deliverable status looks fine —
// the evidence is the thing actually being aggregated, not just the labels.
func TestEvaluate_FailedEvidence_WithholdsGA(t *testing.T) {
	deliverables := []escalation.Deliverable{
		{ID: "§1", Name: "biblio-cli foundations", Status: escalation.CompleteAndDemonstrated, Reference: "test/integration/pilot-site"},
	}

	failedUpdateEvidence := escalation.Evidence{Reference: "test/integration/update-three-sites", Passed: false}
	review := escalation.Evaluate(deliverables, failedUpdateEvidence, passingEscalationEvidence)
	if review.GADetermination {
		t.Error("Evaluate(): GADetermination = true with failed update-evidence, want false")
	}
}

// Evaluate must record the update/escalation evidence references it
// aggregated by, not re-derive them — SC-003 requires this release's
// evidence to stand on its own without re-proving prior releases' work,
// and the reverse: this review references rather than re-runs THIS
// release's own US1/US2/US3 suites too.
func TestEvaluate_RecordsEvidenceByReference(t *testing.T) {
	deliverables := []escalation.Deliverable{
		{ID: "§9", Name: "Replicate to remaining site types", Status: escalation.CompleteAndDemonstrated, Reference: "test/integration/update-three-sites"},
	}

	review := escalation.Evaluate(deliverables, passingUpdateEvidence, passingEscalationEvidence)
	if review.UpdateEvidence != passingUpdateEvidence.Reference {
		t.Errorf("Evaluate(): UpdateEvidence = %q, want %q", review.UpdateEvidence, passingUpdateEvidence.Reference)
	}
	if review.EscalationEvidence != passingEscalationEvidence.Reference {
		t.Errorf("Evaluate(): EscalationEvidence = %q, want %q", review.EscalationEvidence, passingEscalationEvidence.Reference)
	}
}
