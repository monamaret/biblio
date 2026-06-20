package migrate

import (
	"fmt"
	"os"
	"path/filepath"
)

// Apply advances every scaffold-tracked region to its New content and
// leaves every author-customized region untouched (contracts/
// update-and-escalation.md Steps 2-3). It refuses outright if any region
// is Ambiguous — those must go through conflict resolution first
// (contracts/update-and-escalation.md Step 4); Apply never guesses.
func Apply(projectDir string, regions []Region) error {
	for _, r := range regions {
		if r.Classification == Ambiguous {
			return fmt.Errorf("migrate: region %q is ambiguous and unresolved — cannot apply", r.Path)
		}
	}

	for _, r := range regions {
		if r.Classification != ScaffoldTracked {
			continue
		}
		dest := filepath.Join(projectDir, filepath.FromSlash(r.Path))
		if r.New == nil {
			// The newer scaffold version drops this path and the author
			// never touched it (classify.go) — remove it.
			if err := os.Remove(dest); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("migrate: removing %q: %w", r.Path, err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("migrate: creating directory for %q: %w", r.Path, err)
		}
		if err := os.WriteFile(dest, r.New, 0o644); err != nil {
			return fmt.Errorf("migrate: advancing %q: %w", r.Path, err)
		}
	}

	return nil
}
