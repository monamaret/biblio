// Package escalation checks a deployed Site Project's real, current scale
// against the documented D1 (storage volume / relational-query) and
// Vectorize (search-quality) escalation thresholds (FR-004, FR-005),
// without adopting either as a running dependency (Principle V).
package escalation
