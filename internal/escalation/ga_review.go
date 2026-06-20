package escalation

// DeliverableStatus is the closing status a roadmap deliverable is checked
// against (data-model.md "GA Readiness Review").
type DeliverableStatus string

const (
	// CompleteAndDemonstrated is the bar every ordinary roadmap
	// deliverable must meet.
	CompleteAndDemonstrated DeliverableStatus = "complete-and-demonstrated"
	// ProvenExecutableOnDemand is the bar the two explicitly-conditional
	// storage-escalation deliverables (§10) meet instead — proving the
	// path is walkable, never requiring it to actually be carried out
	// (FR-006, spec Assumptions "What GA means here").
	ProvenExecutableOnDemand DeliverableStatus = "proven-executable-on-demand"
)

// Deliverable is one roadmap §1-§10 item being checked against this
// release's acceptance bar.
type Deliverable struct {
	ID     string
	Name   string
	Status DeliverableStatus
	// Reference points at the already-closed evidence demonstrating
	// Status — a prior release's closing artifact, an integration test
	// package, or a guide — never a re-run of that evidence (SC-003).
	Reference string
}

// meetsBar reports whether d's Status meets the bar for its deliverable —
// every ordinary deliverable must be CompleteAndDemonstrated; only the bar
// itself differs for conditional deliverables, which may instead be
// ProvenExecutableOnDemand.
func (d Deliverable) meetsBar() bool {
	return d.Status == CompleteAndDemonstrated || d.Status == ProvenExecutableOnDemand
}

// Evidence is one aggregated suite's pass/fail result and a pointer to
// where it lives, for the GA Readiness Review to reference rather than
// re-run (SC-003).
type Evidence struct {
	Reference string
	Passed    bool
}

// GAReadinessReview is the closing, cross-release artifact that marks the
// platform Generally Available (data-model.md "GA Readiness Review";
// FR-006, Scenario 4, SC-002).
type GAReadinessReview struct {
	Deliverables       []Deliverable
	UpdateEvidence     string
	EscalationEvidence string
	GADetermination    bool
}

// Evaluate cross-references every deliverable against its bar and
// aggregates the update and escalation evidence by reference —
// ga-determination is true only when every deliverable meets its bar and
// both evidence aggregates passed (data-model.md "GA Readiness Review"
// validation rules; this is the literal, checkable definition of GA this
// release establishes).
func Evaluate(deliverables []Deliverable, updateEvidence, escalationEvidence Evidence) GAReadinessReview {
	determination := updateEvidence.Passed && escalationEvidence.Passed
	for _, d := range deliverables {
		if !d.meetsBar() {
			determination = false
			break
		}
	}

	return GAReadinessReview{
		Deliverables:       deliverables,
		UpdateEvidence:     updateEvidence.Reference,
		EscalationEvidence: escalationEvidence.Reference,
		GADetermination:    determination,
	}
}
