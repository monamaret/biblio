// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package content_test

import (
	"strings"
	"testing"

	"github.com/monamaret/biblio/internal/content"
)

// ParseItem must parse valid frontmatter and body for each site type
// (data-model.md "Authored Content Item"; Principle VI — test-first).

func TestParseItem_ValidProfile(t *testing.T) {
	data := []byte("---\ntitle: \"About Me\"\nrole: \"Software Engineer\"\nsummary: \"A bio.\"\n---\n\nWelcome!\n")
	item, err := content.ParseItem("about.md", data, "profile")
	if err != nil {
		t.Fatalf("ParseItem() error = %v", err)
	}
	if item.Slug != "about" {
		t.Errorf("ParseItem() slug = %q, want %q", item.Slug, "about")
	}
	if item.Fields["title"] != "About Me" {
		t.Errorf("ParseItem() Fields[title] = %v, want %q", item.Fields["title"], "About Me")
	}
	if !strings.Contains(item.Body, "Welcome!") {
		t.Errorf("ParseItem() body = %q, want to contain %q", item.Body, "Welcome!")
	}
}

func TestParseItem_ValidLog(t *testing.T) {
	data := []byte("---\ntitle: \"Hello, World\"\ndate: \"2026-06-01\"\ntags:\n  - meta\n---\n\nFirst entry.\n")
	item, err := content.ParseItem("posts/2026-06-01-hello-world.md", data, "log")
	if err != nil {
		t.Fatalf("ParseItem() error = %v", err)
	}
	if item.Slug != "posts/2026-06-01-hello-world" {
		t.Errorf("ParseItem() slug = %q, want %q", item.Slug, "posts/2026-06-01-hello-world")
	}
	if item.Fields["date"] != "2026-06-01" {
		t.Errorf("ParseItem() Fields[date] = %v, want %q", item.Fields["date"], "2026-06-01")
	}
}

func TestParseItem_ValidDocs(t *testing.T) {
	data := []byte("---\ntitle: \"Getting Started\"\nsection: \"Overview\"\norder: 2\n---\n\nFirst steps.\n")
	item, err := content.ParseItem("getting-started.md", data, "docs")
	if err != nil {
		t.Fatalf("ParseItem() error = %v", err)
	}
	if item.Slug != "getting-started" {
		t.Errorf("ParseItem() slug = %q, want %q", item.Slug, "getting-started")
	}
}

// Missing 'title' (required for all site types) must be a clear error.
func TestParseItem_MissingTitle(t *testing.T) {
	data := []byte("---\nrole: \"Engineer\"\n---\n\nBody.\n")
	if _, err := content.ParseItem("about.md", data, "profile"); err == nil {
		t.Fatal("ParseItem(): want error for missing required field 'title', got nil")
	}
}

// Missing 'date' (required for log) must be a clear error.
func TestParseItem_MissingDate_Log(t *testing.T) {
	data := []byte("---\ntitle: \"Post\"\ntags:\n  - meta\n---\n\nBody.\n")
	if _, err := content.ParseItem("posts/post.md", data, "log"); err == nil {
		t.Fatal("ParseItem(): want error for missing required field 'date' (log), got nil")
	}
}

// Missing 'order' (required for docs) must be a clear error.
func TestParseItem_MissingOrder_Docs(t *testing.T) {
	data := []byte("---\ntitle: \"Page\"\n---\n\nBody.\n")
	if _, err := content.ParseItem("page.md", data, "docs"); err == nil {
		t.Fatal("ParseItem(): want error for missing required field 'order' (docs), got nil")
	}
}

// Missing 'section' (required for docs, alongside 'order') must be a clear
// error (data-model.md "Documentation Page" validation rules; spec Edge
// Cases).
func TestParseItem_MissingSection_Docs(t *testing.T) {
	data := []byte("---\ntitle: \"Page\"\norder: 1\n---\n\nBody.\n")
	if _, err := content.ParseItem("page.md", data, "docs"); err == nil {
		t.Fatal("ParseItem(): want error for missing required field 'section' (docs), got nil")
	}
}

// A docs page's 'order' must parse as an integer — a non-numeric value is a
// clear, per-file error, not silently accepted (data-model.md "Documentation
// Page" validation rules; spec Edge Cases).
func TestParseItem_DocsOrderNotInteger(t *testing.T) {
	data := []byte("---\ntitle: \"Page\"\nsection: \"Overview\"\norder: \"first\"\n---\n\nBody.\n")
	if _, err := content.ParseItem("page.md", data, "docs"); err == nil {
		t.Fatal("ParseItem(): want error for non-integer 'order' (docs), got nil")
	}
}

// A docs page's 'section' must be a non-empty string — a non-string value
// is a clear, per-file error, not silently accepted.
func TestParseItem_DocsSectionNotString(t *testing.T) {
	data := []byte("---\ntitle: \"Page\"\nsection:\n  - not-a-string\norder: 1\n---\n\nBody.\n")
	if _, err := content.ParseItem("page.md", data, "docs"); err == nil {
		t.Fatal("ParseItem(): want error for 'section' that isn't a string (docs), got nil")
	}
}

// Malformed YAML frontmatter must return a clear error — never silently
// produce an incomplete or zero-value Item (spec Edge Cases).
func TestParseItem_MalformedYAML(t *testing.T) {
	data := []byte("---\ntitle: : bad yaml here\n---\n\nBody.\n")
	if _, err := content.ParseItem("about.md", data, "profile"); err == nil {
		t.Fatal("ParseItem(): want error for malformed YAML frontmatter, got nil")
	}
}

// The Markdown body must be preserved exactly after the closing frontmatter
// delimiter (data-model.md "Authored Content Item" body field).
func TestParseItem_BodyPreserved(t *testing.T) {
	body := "This is the **body** content."
	data := []byte("---\ntitle: \"Page\"\n---\n\n" + body)
	item, err := content.ParseItem("page.md", data, "profile")
	if err != nil {
		t.Fatalf("ParseItem() error = %v", err)
	}
	if !strings.Contains(item.Body, body) {
		t.Errorf("ParseItem() body = %q, want to contain %q", item.Body, body)
	}
}

// A malformed 'date' on a log entry (not a parseable YYYY-MM-DD date) must
// be a clear, per-file error rather than silently accepted (data-model.md
// "Project-Log Entry" validation rules; spec Edge Cases).
func TestParseItem_LogMalformedDate(t *testing.T) {
	data := []byte("---\ntitle: \"Post\"\ndate: \"not-a-date\"\n---\n\nBody.\n")
	if _, err := content.ParseItem("posts/post.md", data, "log"); err == nil {
		t.Fatal("ParseItem(): want error for malformed 'date' (log), got nil")
	}
}

// A 'tags' field on a log entry that isn't a list of strings (e.g. a single
// scalar) must be a clear error, not silently coerced or accepted.
func TestParseItem_LogTagsNotAList(t *testing.T) {
	data := []byte("---\ntitle: \"Post\"\ndate: \"2026-06-01\"\ntags: \"meta\"\n---\n\nBody.\n")
	if _, err := content.ParseItem("posts/post.md", data, "log"); err == nil {
		t.Fatal("ParseItem(): want error for 'tags' that isn't a list (log), got nil")
	}
}

// A 'tags' list containing a non-string element must be a clear error.
func TestParseItem_LogTagsElementNotString(t *testing.T) {
	data := []byte("---\ntitle: \"Post\"\ndate: \"2026-06-01\"\ntags:\n  - 5\n---\n\nBody.\n")
	if _, err := content.ParseItem("posts/post.md", data, "log"); err == nil {
		t.Fatal("ParseItem(): want error for non-string 'tags' element (log), got nil")
	}
}

// A 'project' field on a log entry that isn't a string must be a clear
// error. 'project' itself is optional — its absence is valid (data-model.md
// "Project-Log Entry" — project optional).
func TestParseItem_LogProjectNotString(t *testing.T) {
	data := []byte("---\ntitle: \"Post\"\ndate: \"2026-06-01\"\nproject:\n  - not-a-string\n---\n\nBody.\n")
	if _, err := content.ParseItem("posts/post.md", data, "log"); err == nil {
		t.Fatal("ParseItem(): want error for 'project' that isn't a string (log), got nil")
	}
}

// A log entry with no 'tags' or 'project' at all is valid — both are
// optional (data-model.md "Project-Log Entry").
func TestParseItem_LogTagsAndProjectOptional(t *testing.T) {
	data := []byte("---\ntitle: \"Post\"\ndate: \"2026-06-01\"\n---\n\nBody.\n")
	if _, err := content.ParseItem("posts/post.md", data, "log"); err != nil {
		t.Fatalf("ParseItem() error = %v, want nil (tags/project are optional)", err)
	}
}

// Path and Slug on the returned Item must reflect the relPath argument.
func TestParseItem_PathAndSlug(t *testing.T) {
	data := []byte("---\ntitle: \"Project\"\n---\n\nDesc.\n")
	item, err := content.ParseItem("projects/my-project.md", data, "profile")
	if err != nil {
		t.Fatalf("ParseItem() error = %v", err)
	}
	if item.Path != "projects/my-project.md" {
		t.Errorf("ParseItem() Path = %q, want %q", item.Path, "projects/my-project.md")
	}
	if item.Slug != "projects/my-project" {
		t.Errorf("ParseItem() Slug = %q, want %q", item.Slug, "projects/my-project")
	}
}

// A quarterly site.yaml's `typesetting:` block is parsed into the typed,
// validated settings (specs/004-quarterly-site-type FR-022, FR-023).
func TestReadSiteConfig_QuarterlyTypesetting(t *testing.T) {
	data := []byte("site_type: quarterly\nname: q\ntypesetting:\n  accent: \"#1e4fd7\"\n  numerals: arabic\n")
	cfg, err := content.ReadSiteConfig(data)
	if err != nil {
		t.Fatalf("ReadSiteConfig() error = %v", err)
	}
	if cfg.Typesetting == nil {
		t.Fatal("ReadSiteConfig(): Typesetting = nil, want parsed settings")
	}
	if cfg.Typesetting.Accent == nil || *cfg.Typesetting.Accent != "#1e4fd7" {
		t.Errorf("Typesetting.Accent = %v, want #1e4fd7", cfg.Typesetting.Accent)
	}
	if cfg.Typesetting.Numerals == nil || *cfg.Typesetting.Numerals != "arabic" {
		t.Errorf("Typesetting.Numerals = %v, want arabic", cfg.Typesetting.Numerals)
	}
}

// No `typesetting:` block leaves Typesetting nil (preset defaults apply).
func TestReadSiteConfig_QuarterlyWithoutTypesetting(t *testing.T) {
	cfg, err := content.ReadSiteConfig([]byte("site_type: quarterly\nname: q\n"))
	if err != nil {
		t.Fatalf("ReadSiteConfig() error = %v", err)
	}
	if cfg.Typesetting != nil {
		t.Errorf("Typesetting = %+v, want nil", cfg.Typesetting)
	}
}

// Every invalid typesetting value is reported, each prefixed with
// "site.yaml:" (contracts/config-and-tokens.md "Error format").
func TestReadSiteConfig_InvalidTypesettingReportsEveryKey(t *testing.T) {
	_, err := content.ReadSiteConfig([]byte("site_type: quarterly\ntypesetting:\n  scale: 2\n  align: center\n"))
	if err == nil {
		t.Fatal("ReadSiteConfig() error = nil, want validation errors")
	}
	for _, want := range []string{"site.yaml: typesetting.scale", "site.yaml: typesetting.align"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not contain %q", err, want)
		}
	}
}

// A `typesetting:` block on any other site type is an error naming the
// site type, so a setting is never silently ignored (research.md R8).
func TestReadSiteConfig_TypesettingRejectedOnOtherTypes(t *testing.T) {
	_, err := content.ReadSiteConfig([]byte("site_type: log\ntypesetting:\n  accent: \"#000000\"\n"))
	if err == nil {
		t.Fatal("ReadSiteConfig() error = nil, want an error for typesetting on a log site")
	}
	if !strings.Contains(err.Error(), `typesetting is only supported for site_type quarterly (this site is "log")`) {
		t.Errorf("error = %q, want it to name the site type", err)
	}
}
