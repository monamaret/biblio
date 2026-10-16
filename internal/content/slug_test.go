// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package content_test

import (
	"testing"

	"github.com/monamaret/biblio/internal/content"
)

// SlugFromPath derives the canonical slug from a relative content path.
// Slugs must be stable and collision-free within a site
// (data-model.md "Authored Content Item"; Principle VI — test-first).

func TestSlugFromPath_TopLevel(t *testing.T) {
	got := content.SlugFromPath("about.md")
	if got != "about" {
		t.Errorf("SlugFromPath(%q) = %q, want %q", "about.md", got, "about")
	}
}

func TestSlugFromPath_Nested(t *testing.T) {
	got := content.SlugFromPath("projects/sample-project.md")
	if got != "projects/sample-project" {
		t.Errorf("SlugFromPath(%q) = %q, want %q", "projects/sample-project.md", got, "projects/sample-project")
	}
}

func TestSlugFromPath_DatePrefix(t *testing.T) {
	got := content.SlugFromPath("posts/2026-06-01-hello-world.md")
	if got != "posts/2026-06-01-hello-world" {
		t.Errorf("SlugFromPath(%q) = %q, want %q", "posts/2026-06-01-hello-world.md", got, "posts/2026-06-01-hello-world")
	}
}

// Two files with the same name in different directories must produce
// different slugs (collision-free invariant from data-model.md).
func TestSlugFromPath_NoCollision(t *testing.T) {
	slug1 := content.SlugFromPath("a/page.md")
	slug2 := content.SlugFromPath("b/page.md")
	if slug1 == slug2 {
		t.Errorf("SlugFromPath: different paths must not produce the same slug; both returned %q", slug1)
	}
}

// The same path must always return the same slug (stability invariant).
func TestSlugFromPath_Stable(t *testing.T) {
	path := "posts/2026-06-01-hello-world.md"
	first := content.SlugFromPath(path)
	second := content.SlugFromPath(path)
	if first != second {
		t.Errorf("SlugFromPath(%q) is not stable: got %q then %q", path, first, second)
	}
}

// index.md at the top level produces slug "index" (not an empty string).
func TestSlugFromPath_IndexMd(t *testing.T) {
	got := content.SlugFromPath("index.md")
	if got != "index" {
		t.Errorf("SlugFromPath(%q) = %q, want %q", "index.md", got, "index")
	}
}
