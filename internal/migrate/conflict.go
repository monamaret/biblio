package migrate

import "fmt"

// Resolution values an author can record against an UpdateConflict.
const (
	ResolveKeepAuthor  = "keep-author"
	ResolveUseScaffold = "use-scaffold"
)

// UpdateConflict is a surfaced case where advancing the scaffold would
// otherwise overwrite an author's intentional customization (data-model.md
// "Update Conflict"; FR-002). The flow waits for and applies the author's
// explicit Resolution — it is never auto-resolved.
type UpdateConflict struct {
	Region           string
	ScaffoldProposed []byte
	AuthorCurrent    []byte
	// Resolution is empty until the author decides; ApplyWithResolutions
	// refuses to commit any region whose conflict has no Resolution.
	Resolution string
}

// Conflicts extracts one UpdateConflict per Ambiguous region, carrying
// both sides of the diff for the author to review.
func Conflicts(regions []Region) []UpdateConflict {
	var conflicts []UpdateConflict
	for _, r := range regions {
		if r.Classification != Ambiguous {
			continue
		}
		conflicts = append(conflicts, UpdateConflict{
			Region:           r.Path,
			ScaffoldProposed: r.New,
			AuthorCurrent:    r.Current,
		})
	}
	return conflicts
}

// ApplyWithResolutions commits regions exactly as Apply does, except every
// Ambiguous region in regions must have a matching, resolved UpdateConflict
// in conflicts — an Update Run with any non-empty, unresolved conflict MUST
// NOT proceed to overwrite those regions (data-model.md "Update Conflict"
// validation rules). A resolved conflict is converted into a concrete write
// (ResolveUseScaffold) or no-op (ResolveKeepAuthor) before committing, so
// the interruption-recovery guarantees Apply already provides still hold.
func ApplyWithResolutions(projectDir string, regions []Region, conflicts []UpdateConflict) error {
	resolutionByPath := make(map[string]string, len(conflicts))
	for _, c := range conflicts {
		resolutionByPath[c.Region] = c.Resolution
	}

	resolved := make([]Region, len(regions))
	for i, r := range regions {
		resolved[i] = r
		if r.Classification != Ambiguous {
			continue
		}
		resolution, ok := resolutionByPath[r.Path]
		if !ok || resolution == "" {
			return fmt.Errorf("migrate: region %q has an unresolved conflict — cannot commit", r.Path)
		}
		switch resolution {
		case ResolveUseScaffold:
			resolved[i].Classification = ScaffoldTracked
		case ResolveKeepAuthor:
			resolved[i].Classification = AuthorCustomized
		default:
			return fmt.Errorf("migrate: region %q has unrecognized resolution %q", r.Path, resolution)
		}
	}

	return Apply(projectDir, resolved)
}
