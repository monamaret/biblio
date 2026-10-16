package index_test

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/monamaret/biblio/internal/index"
)

// fixtureDir is the shared quarterly fixture site (repo testdata/quarterly):
// issue 01 (2 articles, cover), issue 02 (3 articles incl. dense notes and
// plates), draft issue 03 (1 article), and a colophon.
const fixtureDir = "../../testdata/quarterly"

// copyFixture copies the fixture into a temp dir so a test can add files.
func copyFixture(t *testing.T) string {
	t.Helper()
	dst := t.TempDir()
	err := filepath.WalkDir(fixtureDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(fixtureDir, p)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		t.Fatalf("copying fixture: %v", err)
	}
	return dst
}

func generate(t *testing.T, dir string, opts index.Options) index.Index {
	t.Helper()
	idx, errs := index.GenerateWithOptions(dir, opts)
	if len(errs) > 0 {
		t.Fatalf("GenerateWithOptions errors: %v", errs)
	}
	return idx
}

func slugs(idx index.Index) []string {
	var s []string
	for _, e := range idx.Entries {
		s = append(s, e.Slug)
	}
	return s
}

func entryBySlug(t *testing.T, idx index.Index, slug string) index.Entry {
	t.Helper()
	for _, e := range idx.Entries {
		if e.Slug == slug {
			return e
		}
	}
	t.Fatalf("no entry %q in %v", slug, slugs(idx))
	return index.Entry{}
}

// Entries are ordered: issues newest first; within an issue, the issue
// entry then its articles in contents order (grouped by category, the
// categories in order of their first article); the colophon last
// (data-model.md "content-index.json shape", "Ordering rules").
func TestQuarterly_EntryOrder(t *testing.T) {
	idx := generate(t, fixtureDir, index.Options{})
	want := []string{
		"issues/02-summer-2026",
		"issues/02-summer-2026/dense-notes",
		"issues/02-summer-2026/unordered",
		"issues/02-summer-2026/plates",
		"issues/01-spring-2026",
		"issues/01-spring-2026/first-essay",
		"issues/01-spring-2026/grid-systems",
		"colophon",
	}
	if got := slugs(idx); !slices.Equal(got, want) {
		t.Errorf("entry order =\n%v\nwant\n%v", got, want)
	}
}

// Articles carry their issue, category key and contents position, derived
// from the folder — never from frontmatter.
func TestQuarterly_DerivedArticleFields(t *testing.T) {
	idx := generate(t, fixtureDir, index.Options{})
	cases := map[string]struct {
		key      string
		position int
	}{
		"issues/02-summer-2026/dense-notes": {"essays", 1},
		"issues/02-summer-2026/unordered":   {"essays", 2},
		"issues/02-summer-2026/plates":      {"design", 3},
	}
	for slug, want := range cases {
		e := entryBySlug(t, idx, slug)
		f := e.Fields
		if f["kind"] != "article" || f["issue"] != "02-summer-2026" || f["issue_number"] != 2 || f["issue_date"] != "2026-07-01" {
			t.Errorf("%s: kind/issue fields = %v/%v/%v/%v", slug, f["kind"], f["issue"], f["issue_number"], f["issue_date"])
		}
		if f["category_key"] != want.key {
			t.Errorf("%s: category_key = %v, want %s", slug, f["category_key"], want.key)
		}
		if f["position"] != want.position {
			t.Errorf("%s: position = %v, want %d", slug, f["position"], want.position)
		}
	}
	issue := entryBySlug(t, idx, "issues/02-summer-2026")
	if issue.Fields["kind"] != "issue" || issue.Fields["issue"] != "02-summer-2026" || issue.Fields["article_count"] != 3 {
		t.Errorf("issue entry fields = %v", issue.Fields)
	}
	if entryBySlug(t, idx, "colophon").Fields["kind"] != "colophon" {
		t.Error("colophon entry kind != colophon")
	}
}

// Site-level arrays: issues newest first, categories by first appearance
// with their first-seen label, the current issue, and the default
// structural typesetting switches.
func TestQuarterly_SiteLevelArrays(t *testing.T) {
	idx := generate(t, fixtureDir, index.Options{})

	if len(idx.Issues) != 2 || idx.Issues[0].Issue != "02-summer-2026" || idx.Issues[1].Issue != "01-spring-2026" {
		t.Fatalf("issues = %+v, want 02 then 01", idx.Issues)
	}
	first := idx.Issues[1]
	if first.Number != 1 || first.Title != "Margins and Grids" || first.Date != "2026-04-01" || first.Season != "Spring 2026" ||
		first.CoverAlt == "" || first.ArticleCount != 2 || first.Draft || first.Slug != "issues/01-spring-2026" {
		t.Errorf("issue 01 summary = %+v", first)
	}

	wantCats := []index.CategorySummary{
		{Key: "essays", Label: "essays", ArticleCount: 3},
		{Key: "design", Label: "Design", ArticleCount: 2},
	}
	if !slices.Equal(idx.Categories, wantCats) {
		t.Errorf("categories = %+v, want %+v", idx.Categories, wantCats)
	}

	if idx.CurrentIssue == nil || *idx.CurrentIssue != "02-summer-2026" {
		t.Errorf("current_issue = %v, want 02-summer-2026", idx.CurrentIssue)
	}
	if idx.Typesetting == nil || !idx.Typesetting.MarginNotes || idx.Typesetting.Numerals != "roman" {
		t.Errorf("typesetting = %+v, want margin_notes true, numerals roman", idx.Typesetting)
	}
}

// Drafts: excluded (and counted) in publish mode; included and marked in
// dev mode; never the current issue (FR-012).
func TestQuarterly_Drafts(t *testing.T) {
	published := generate(t, fixtureDir, index.Options{})
	for _, s := range slugs(published) {
		if strings.HasPrefix(s, "issues/03-autumn-2026") {
			t.Errorf("draft entry %q present in a publish-mode index", s)
		}
	}
	if published.DraftIssuesExcluded != 1 {
		t.Errorf("DraftIssuesExcluded = %d, want 1", published.DraftIssuesExcluded)
	}

	dev := generate(t, fixtureDir, index.Options{IncludeDrafts: true})
	if len(dev.Issues) != 3 || dev.Issues[0].Issue != "03-autumn-2026" || !dev.Issues[0].Draft {
		t.Fatalf("dev issues = %+v, want draft 03 first", dev.Issues)
	}
	if dev.CurrentIssue == nil || *dev.CurrentIssue != "02-summer-2026" {
		t.Errorf("dev current_issue = %v, want 02-summer-2026 (a draft is never current)", dev.CurrentIssue)
	}
	draftArticle := entryBySlug(t, dev, "issues/03-autumn-2026/draft-piece")
	if draftArticle.Fields["draft"] != true {
		t.Errorf("draft article draft field = %v, want true", draftArticle.Fields["draft"])
	}
	if dev.DraftIssuesExcluded != 0 {
		t.Errorf("dev DraftIssuesExcluded = %d, want 0", dev.DraftIssuesExcluded)
	}
}

// Two issues with the same number fail, naming both folders; drafts count.
func TestQuarterly_DuplicateIssueNumber(t *testing.T) {
	dir := copyFixture(t)
	if err := os.MkdirAll(filepath.Join(dir, "issues", "04-dupe"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "issues", "04-dupe", "issue.md"), []byte("---\ntitle: Dupe\nnumber: 3\ndate: 2026-12-01\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, errs := index.GenerateWithOptions(dir, index.Options{})
	joined := ""
	for _, e := range errs {
		joined += e.Error() + "\n"
	}
	if !strings.Contains(joined, "03-autumn-2026") || !strings.Contains(joined, "04-dupe") {
		t.Errorf("errors = %q, want both folders sharing number 3 named", joined)
	}
}

// An issue folder without issue.md is an error; an issue with no articles
// is fine and counts zero.
func TestQuarterly_IssueFolderRules(t *testing.T) {
	dir := copyFixture(t)
	if err := os.MkdirAll(filepath.Join(dir, "issues", "05-empty"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "issues", "05-empty", "issue.md"), []byte("---\ntitle: Empty\nnumber: 5\ndate: 2027-01-01\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	idx := generate(t, dir, index.Options{})
	if idx.Issues[0].Issue != "05-empty" || idx.Issues[0].ArticleCount != 0 {
		t.Errorf("issues[0] = %+v, want 05-empty with 0 articles", idx.Issues[0])
	}
	if *idx.CurrentIssue != "05-empty" {
		t.Errorf("current_issue = %s, want 05-empty", *idx.CurrentIssue)
	}

	if err := os.MkdirAll(filepath.Join(dir, "issues", "06-orphan"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "issues", "06-orphan", "x.md"), []byte("---\ntitle: X\ncategory: E\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, errs := index.GenerateWithOptions(dir, index.Options{})
	if len(errs) == 0 || !strings.Contains(errs[0].Error(), "06-orphan") || !strings.Contains(errs[0].Error(), "issue.md") {
		t.Errorf("errors = %v, want one naming 06-orphan's missing issue.md", errs)
	}
}

// No published issue at all: current_issue is null, not missing.
func TestQuarterly_NoPublishedIssue(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "site.yaml"), []byte("site_type: quarterly\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	idx := generate(t, dir, index.Options{})
	data, _ := json.Marshal(idx)
	if !strings.Contains(string(data), `"current_issue":null`) {
		t.Errorf("JSON = %s, want \"current_issue\":null", data)
	}
	if !strings.Contains(string(data), `"issues":[]`) || !strings.Contains(string(data), `"categories":[]`) {
		t.Errorf("JSON = %s, want empty issues and categories arrays", data)
	}
}

// Profile/log/docs JSON never gains the quarterly keys.
func TestQuarterly_OtherTypesJSONUnchanged(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "site.yaml"), []byte("site_type: profile\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "about.md"), []byte("---\ntitle: About\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	idx := generate(t, dir, index.Options{IncludeDrafts: true})
	data, _ := json.Marshal(idx)
	for _, key := range []string{"issues", "categories", "current_issue", "typesetting"} {
		if strings.Contains(string(data), `"`+key+`"`) {
			t.Errorf("profile index JSON contains %q: %s", key, data)
		}
	}
}
