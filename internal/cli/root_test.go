// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package cli

import "testing"

// resolveVersion must never override an already-set (non-"dev") ldflags
// version — the build-info fallback only kicks in for the unset default
// (issue #193).
func TestResolveVersion_PassesThroughNonDevValue(t *testing.T) {
	got := resolveVersion("v1.2.3")
	if got != "v1.2.3" {
		t.Errorf("resolveVersion(%q) = %q, want %q", "v1.2.3", got, "v1.2.3")
	}
}

// Under `go test`, Go's build info reports "(devel)"/empty rather than a
// real module version — resolveVersion must leave the "dev" default alone
// in that case rather than surfacing that raw Go internal. (A plain `go
// build` of the actual CLI binary behaves differently — see root.go's
// resolveVersion doc comment — but isn't observable from within a test
// binary.)
func TestResolveVersion_LocalBuildStaysDev(t *testing.T) {
	got := resolveVersion("dev")
	if got != "dev" {
		t.Errorf(`resolveVersion("dev") = %q, want "dev" (go test's own build info shouldn't override it)`, got)
	}
}
