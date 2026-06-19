package index_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/monamaret/biblio/internal/index"
)

// writeTempContent creates a minimal content/ directory for testing Generate.
// It writes site.yaml (declaring siteType) and each file in the files map
// (keys are relative paths under content/, values are raw file contents).
func writeTempContent(t *testing.T, siteType string, files map[string]string) string {
	t.Helper()
	contentDir := filepath.Join(t.TempDir(), "content")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatalf("setup: creating content dir: %v", err)
	}
	siteYAML := fmt.Sprintf("site_type: %s\nname: test\ndescription: \"\"\nmcp_enabled: true\n", siteType)
	if err := os.WriteFile(filepath.Join(contentDir, "site.yaml"), []byte(siteYAML), 0o644); err != nil {
		t.Fatalf("setup: writing site.yaml: %v", err)
	}
	for relPath, body := range files {
		dest := filepath.Join(contentDir, filepath.FromSlash(relPath))
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			t.Fatalf("setup: creating dir for %q: %v", relPath, err)
		}
		if err := os.WriteFile(dest, []byte(body), 0o644); err != nil {
			t.Fatalf("setup: writing %q: %v", relPath, err)
		}
	}
	return contentDir
}

// Generate must produce exactly one entry per .md file
// (data-model.md "Content Index"; contracts/cli-commands.md `biblio-cli index`).
func TestGenerate_OneEntryPerItem(t *testing.T) {
	contentDir := writeTempContent(t, "profile", map[string]string{
		"about.md":             "---\ntitle: \"About\"\n---\n\nBody.\n",
		"projects/my-proj.md":  "---\ntitle: \"Project\"\n---\n\nDesc.\n",
	})

	idx, errs := index.Generate(contentDir)
	if len(errs) > 0 {
		t.Fatalf("Generate() errors = %v", errs)
	}
	if len(idx.Entries) != 2 {
		t.Errorf("Generate() len(entries) = %d, want 2", len(idx.Entries))
	}
}

// Each entry must carry its correctly-derived slug and relative path.
func TestGenerate_CorrectSlugAndPath(t *testing.T) {
	contentDir := writeTempContent(t, "profile", map[string]string{
		"about.md": "---\ntitle: \"About\"\n---\n\nBody.\n",
	})

	idx, errs := index.Generate(contentDir)
	if len(errs) > 0 {
		t.Fatalf("Generate() errors = %v", errs)
	}
	if len(idx.Entries) != 1 {
		t.Fatalf("Generate() len(entries) = %d, want 1", len(idx.Entries))
	}
	e := idx.Entries[0]
	if e.Slug != "about" {
		t.Errorf("entry Slug = %q, want %q", e.Slug, "about")
	}
	if e.Path != "about.md" {
		t.Errorf("entry Path = %q, want %q", e.Path, "about.md")
	}
}

// Non-.md files (including site.yaml) must not appear in the index.
func TestGenerate_SkipsNonMarkdown(t *testing.T) {
	contentDir := writeTempContent(t, "profile", map[string]string{
		"about.md": "---\ntitle: \"About\"\n---\n\nBody.\n",
	})
	if err := os.WriteFile(filepath.Join(contentDir, "readme.txt"), []byte("ignored"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx, errs := index.Generate(contentDir)
	if len(errs) > 0 {
		t.Fatalf("Generate() errors = %v", errs)
	}
	if len(idx.Entries) != 1 {
		t.Errorf("Generate() len(entries) = %d, want 1 (non-.md files must be excluded)", len(idx.Entries))
	}
}

// When multiple files fail schema validation, Generate must report every
// failing file — not just the first — so authors can fix all errors at once
// (contracts/cli-commands.md `biblio-cli index` "Output (failure)"; spec Edge Cases).
func TestGenerate_ReportsAllValidationErrors(t *testing.T) {
	contentDir := writeTempContent(t, "profile", map[string]string{
		"a.md": "---\nrole: \"Engineer\"\n---\n\nBody.\n", // missing title
		"b.md": "---\nrole: \"Dev\"\n---\n\nBody.\n",      // missing title
	})

	_, errs := index.Generate(contentDir)
	if len(errs) < 2 {
		t.Errorf("Generate() errors = %d, want at least 2 (one per failing file)", len(errs))
	}
}

// The returned Index must be JSON-serializable with the shape contract
// (data-model.md "Content Index" — "directly consumable by prerender and MCP").
func TestGenerate_IndexIsJSONSerializable(t *testing.T) {
	contentDir := writeTempContent(t, "log", map[string]string{
		"posts/first.md": "---\ntitle: \"First\"\ndate: \"2026-06-01\"\ntags: []\n---\n\nBody.\n",
	})

	idx, errs := index.Generate(contentDir)
	if len(errs) > 0 {
		t.Fatalf("Generate() errors = %v", errs)
	}
	out, err := json.Marshal(idx)
	if err != nil {
		t.Fatalf("json.Marshal(Index) error = %v", err)
	}
	// Verify the outer shape: {"entries":[...]}
	var shape struct {
		Entries []json.RawMessage `json:"entries"`
	}
	if err := json.Unmarshal(out, &shape); err != nil {
		t.Fatalf("json.Unmarshal output shape: %v\noutput: %s", err, out)
	}
	if len(shape.Entries) != 1 {
		t.Errorf("JSON entries count = %d, want 1", len(shape.Entries))
	}
}

// Each JSON entry must include slug, path, and the frontmatter fields
// (data-model.md "Content Index" shape contract).
func TestGenerate_EntryJSONContainsRequiredKeys(t *testing.T) {
	contentDir := writeTempContent(t, "profile", map[string]string{
		"about.md": "---\ntitle: \"About Me\"\n---\n\nBody.\n",
	})

	idx, errs := index.Generate(contentDir)
	if len(errs) > 0 {
		t.Fatalf("Generate() errors = %v", errs)
	}

	out, _ := json.Marshal(idx)
	var shape struct {
		Entries []map[string]any `json:"entries"`
	}
	if err := json.Unmarshal(out, &shape); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if len(shape.Entries) == 0 {
		t.Fatal("no entries in JSON output")
	}
	e := shape.Entries[0]
	for _, key := range []string{"slug", "path", "title"} {
		if _, ok := e[key]; !ok {
			t.Errorf("entry JSON missing required key %q; got keys: %v", key, keys(e))
		}
	}
}

func keys(m map[string]any) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}
