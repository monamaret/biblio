package content_test

import (
	"strings"
	"testing"

	"github.com/monamaret/biblio/internal/content"
)

// ParseItemKind picks the quarterly schema from the file's location:
// issues/<dir>/issue.md is an issue, any other .md in an issue folder is
// an article, colophon.md is the colophon (specs/004-quarterly-site-type
// FR-005–FR-008; data-model.md "On-disk layout").
func TestParseItemKind_KindFromPath(t *testing.T) {
	cases := []struct {
		path, front, wantKind string
	}{
		{"issues/01-spring-2026/issue.md", "title: T\nnumber: 1\ndate: 2026-04-01\n", "issue"},
		{"issues/01-spring-2026/on-margins.md", "title: T\ncategory: Essays\n", "article"},
		{"colophon.md", "title: Colophon\n", "colophon"},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			item, err := content.ParseItemKind(tc.path, []byte("---\n"+tc.front+"---\n\nBody.\n"), "quarterly")
			if err != nil {
				t.Fatalf("ParseItemKind() error = %v", err)
			}
			if item.Kind != tc.wantKind {
				t.Errorf("Kind = %q, want %q", item.Kind, tc.wantKind)
			}
			if item.Slug != strings.TrimSuffix(tc.path, ".md") {
				t.Errorf("Slug = %q, want %q", item.Slug, strings.TrimSuffix(tc.path, ".md"))
			}
		})
	}
}

// Unquoted YAML dates (`date: 2026-04-01`) are accepted and normalized to
// the YYYY-MM-DD string the index and template expect.
func TestParseItemKind_UnquotedDateNormalized(t *testing.T) {
	item, err := content.ParseItemKind("issues/01-a/issue.md", []byte("---\ntitle: T\nnumber: 1\ndate: 2026-04-01\n---\n"), "quarterly")
	if err != nil {
		t.Fatalf("ParseItemKind() error = %v", err)
	}
	if item.Fields["date"] != "2026-04-01" {
		t.Errorf("Fields[date] = %#v, want the string \"2026-04-01\"", item.Fields["date"])
	}
}

// Every per-kind rule fails with a message naming the file and the field.
func TestParseItemKind_ValidationErrors(t *testing.T) {
	cases := []struct {
		name, path, front string
		want              []string
	}{
		{"issue missing title", "issues/01-a/issue.md", "number: 1\ndate: 2026-04-01\n", []string{"issues/01-a/issue.md", `"title"`}},
		{"issue missing number", "issues/01-a/issue.md", "title: T\ndate: 2026-04-01\n", []string{"issues/01-a/issue.md", `"number"`}},
		{"issue missing date", "issues/01-a/issue.md", "title: T\nnumber: 1\n", []string{"issues/01-a/issue.md", `"date"`}},
		{"issue number not integer", "issues/01-a/issue.md", "title: T\nnumber: one\ndate: 2026-04-01\n", []string{`"number"`, "integer"}},
		{"issue number zero", "issues/01-a/issue.md", "title: T\nnumber: 0\ndate: 2026-04-01\n", []string{`"number"`, "1 or more"}},
		{"issue bad date", "issues/01-a/issue.md", "title: T\nnumber: 1\ndate: April\n", []string{`"date"`, "YYYY-MM-DD"}},
		{"issue empty season", "issues/01-a/issue.md", "title: T\nnumber: 1\ndate: 2026-04-01\nseason: \"\"\n", []string{`"season"`}},
		{"issue cover without alt", "issues/01-a/issue.md", "title: T\nnumber: 1\ndate: 2026-04-01\ncover: c.jpg\n", []string{`"cover_alt"`, `"cover"`}},
		{"issue draft not bool", "issues/01-a/issue.md", "title: T\nnumber: 1\ndate: 2026-04-01\ndraft: yes please\n", []string{`"draft"`, "true or false"}},
		{"article missing title", "issues/01-a/x.md", "category: Essays\n", []string{"issues/01-a/x.md", `"title"`}},
		{"article missing category", "issues/01-a/x.md", "title: T\n", []string{"issues/01-a/x.md", `"category"`}},
		{"article blank category", "issues/01-a/x.md", "title: T\ncategory: \"   \"\n", []string{`"category"`}},
		{"article image without alt", "issues/01-a/x.md", "title: T\ncategory: E\nimage: i.jpg\n", []string{`"image_alt"`, `"image"`}},
		{"article order not integer", "issues/01-a/x.md", "title: T\ncategory: E\norder: first\n", []string{`"order"`, "integer"}},
		{"article tags not list", "issues/01-a/x.md", "title: T\ncategory: E\ntags: grid\n", []string{`"tags"`, "list"}},
		{"colophon missing title", "colophon.md", "season: x\n", []string{"colophon.md", `"title"`}},
		{"md outside issues", "about.md", "title: T\n", []string{"about.md", "issues/<folder>/"}},
		{"nested too deep", "issues/a/b/x.md", "title: T\ncategory: E\n", []string{"issues/a/b/x.md", "issues/<folder>/"}},
		{"md directly in issues", "issues/x.md", "title: T\ncategory: E\n", []string{"issues/x.md", "issues/<folder>/"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := content.ParseItemKind(tc.path, []byte("---\n"+tc.front+"---\n"), "quarterly")
			if err == nil {
				t.Fatal("ParseItemKind() error = nil, want a validation error")
			}
			for _, sub := range tc.want {
				if !strings.Contains(err.Error(), sub) {
					t.Errorf("error %q does not contain %q", err, sub)
				}
			}
		})
	}
}

// For the other site types ParseItemKind is exactly ParseItem.
func TestParseItemKind_DelegatesForOtherTypes(t *testing.T) {
	data := []byte("---\ntitle: \"Hello\"\ndate: \"2026-06-01\"\n---\n\nHi.\n")
	want, wantErr := content.ParseItem("posts/2026-06-01-hello.md", data, "log")
	got, gotErr := content.ParseItemKind("posts/2026-06-01-hello.md", data, "log")
	if (wantErr == nil) != (gotErr == nil) || got.Slug != want.Slug || got.Kind != "" {
		t.Errorf("ParseItemKind(log) = %+v, %v; want %+v, %v", got, gotErr, want, wantErr)
	}
}

// Category labels that differ only in case or spacing share one key; the
// key is URL-safe ASCII (FR-011; research.md R10).
func TestCategoryKey(t *testing.T) {
	cases := map[string]string{
		"Essays":       "essays",
		" essays ":     "essays",
		"ESSAYS":       "essays",
		"Field  Notes": "field-notes",
		"Ukázky":       "ukazky",
		"Łódź & Brno":  "lodz-brno",
	}
	for label, want := range cases {
		got, err := content.CategoryKey(label)
		if err != nil {
			t.Errorf("CategoryKey(%q) error = %v", label, err)
			continue
		}
		if got != want {
			t.Errorf("CategoryKey(%q) = %q, want %q", label, got, want)
		}
	}
	if _, err := content.CategoryKey("!!!"); err == nil {
		t.Error(`CategoryKey("!!!") error = nil, want an error for a label with no letters or digits`)
	}
}
