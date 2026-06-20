// merge_test exercises internal/migrate's diff/merge logic: advancing
// scaffold-tracked regions and preserving author-customized regions
// byte-for-byte (contracts/update-and-escalation.md Steps 2-3).
package migrate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/monamaret/biblio/internal/migrate"
)

// Apply must overwrite a scaffold-tracked region's file with its New
// content, advancing it to the target scaffold version.
func TestApply_ScaffoldTracked_Advances(t *testing.T) {
	projectDir := newScaffoldedProject(t)
	target := filepath.Join(projectDir, "package.json")

	region := migrate.Region{
		Path:           "package.json",
		Classification: migrate.ScaffoldTracked,
		New:            []byte(`{"name":"advanced"}`),
	}

	if err := migrate.Apply(projectDir, []migrate.Region{region}); err != nil {
		t.Fatalf("Apply(): %v", err)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("reading advanced file: %v", err)
	}
	if string(got) != `{"name":"advanced"}` {
		t.Errorf("package.json content = %q, want advanced scaffold content", got)
	}
}

// Apply must leave an author-customized region's file exactly as it was —
// byte-for-byte untouched.
func TestApply_AuthorCustomized_PreservedByteForByte(t *testing.T) {
	projectDir := newScaffoldedProject(t)
	target := filepath.Join(projectDir, "package.json")

	original, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("reading original file: %v", err)
	}

	region := migrate.Region{
		Path:           "package.json",
		Classification: migrate.AuthorCustomized,
		Current:        original,
		New:            []byte(`{"name":"would-be-overwritten"}`),
	}

	if err := migrate.Apply(projectDir, []migrate.Region{region}); err != nil {
		t.Fatalf("Apply(): %v", err)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("reading file after Apply: %v", err)
	}
	if string(got) != string(original) {
		t.Errorf("package.json was modified: got %q, want untouched original %q", got, original)
	}
}

// Apply must refuse to act on an ambiguous region — that's conflict.go's
// job, never merge.go's.
func TestApply_Ambiguous_Rejected(t *testing.T) {
	projectDir := newScaffoldedProject(t)

	region := migrate.Region{
		Path:           "package.json",
		Classification: migrate.Ambiguous,
	}

	if err := migrate.Apply(projectDir, []migrate.Region{region}); err == nil {
		t.Error("Apply() with an ambiguous region: want error, got nil")
	}
}
