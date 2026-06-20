package migrate

import (
	"errors"
	"fmt"
)

// ErrUnresolvedConflicts is returned by RunUpdate when one or more regions
// classified as Ambiguous have no caller-supplied resolution. No write
// occurs in this case — classification always happens, and is reported,
// before any write (data-model.md "Update Conflict" validation rules).
var ErrUnresolvedConflicts = errors.New("migrate: update run has unresolved conflicts")

// UpdateRun is a single execution of `biblio-cli update`/`upgrade` against
// one deployed Site Project (data-model.md "Update Run"; FR-001,
// Scenario 1).
type UpdateRun struct {
	ProjectDir string
	SiteType   string
	// Preserved lists author-customized regions carried through untouched.
	Preserved []string
	// Advanced lists scaffold-tracked regions brought up to date.
	Advanced []string
	// Conflicts lists every region that was ambiguous, with its Resolution
	// recorded if one was supplied (empty if RunUpdate returned before
	// committing).
	Conflicts []UpdateConflict
}

// RunUpdate classifies every region of the project at projectDir against
// the embedded scaffold template for siteType (contracts/
// update-and-escalation.md Step 1), then advances scaffold-tracked
// regions, preserves author-customized regions, and surfaces any ambiguous
// region as an Update Conflict (Steps 2-4). resolutions maps a region path
// to a caller-supplied Resolution (ResolveKeepAuthor/ResolveUseScaffold);
// it may be nil.
//
// If any region remains ambiguous after resolutions are applied, RunUpdate
// returns ErrUnresolvedConflicts and performs no write at all — the
// returned UpdateRun still reports every conflict (with whatever
// resolutions were supplied) so the caller can present them, but Preserved
// and Advanced are empty since nothing was committed (Step 5: only after
// all conflicts are resolved does the flow apply changes).
func RunUpdate(projectDir, siteType string, resolutions map[string]string) (*UpdateRun, error) {
	regions, err := ClassifySite(projectDir, siteType)
	if err != nil {
		return nil, fmt.Errorf("migrate: update run: %w", err)
	}

	conflicts := Conflicts(regions)
	for i, c := range conflicts {
		if r, ok := resolutions[c.Region]; ok {
			conflicts[i].Resolution = r
		}
	}

	run := &UpdateRun{ProjectDir: projectDir, SiteType: siteType, Conflicts: conflicts}

	for _, c := range conflicts {
		if c.Resolution == "" {
			return run, ErrUnresolvedConflicts
		}
	}

	if err := ApplyWithResolutions(projectDir, regions, conflicts); err != nil {
		return run, fmt.Errorf("migrate: update run: %w", err)
	}

	resolutionByPath := make(map[string]string, len(conflicts))
	for _, c := range conflicts {
		resolutionByPath[c.Region] = c.Resolution
	}

	for _, r := range regions {
		switch {
		case r.Classification == ScaffoldTracked:
			run.Advanced = append(run.Advanced, r.Path)
		case r.Classification == AuthorCustomized:
			run.Preserved = append(run.Preserved, r.Path)
		case r.Classification == Ambiguous && resolutionByPath[r.Path] == ResolveUseScaffold:
			run.Advanced = append(run.Advanced, r.Path)
		case r.Classification == Ambiguous && resolutionByPath[r.Path] == ResolveKeepAuthor:
			run.Preserved = append(run.Preserved, r.Path)
		}
	}

	return run, nil
}
