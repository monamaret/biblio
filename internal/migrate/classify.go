// Package migrate carries an existing deployed Site Project forward to a
// newer scaffold-template version: classifying each region as
// scaffold-tracked, author-customized, or ambiguous, advancing or
// preserving it accordingly, and surfacing any ambiguous region as an
// Update Conflict rather than silently resolving it (FR-001, FR-002).
package migrate
