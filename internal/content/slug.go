// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package content

import (
	"path/filepath"
	"strings"
)

// SlugFromPath derives the canonical, stable slug from a relative content
// path (data-model.md "Authored Content Item" — slug field). The slug is the
// relative path with the .md extension removed and path separators normalised
// to forward slashes, making it stable and collision-free within a site.
//
// Examples:
//
//	"about.md"                         → "about"
//	"projects/my-project.md"           → "projects/my-project"
//	"posts/2026-06-01-hello-world.md"  → "posts/2026-06-01-hello-world"
func SlugFromPath(relPath string) string {
	noExt := strings.TrimSuffix(relPath, filepath.Ext(relPath))
	return filepath.ToSlash(noExt)
}
