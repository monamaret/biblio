// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// Package pilot_site_mcp is the release-closing integration check for
// v0.03 (FR-001..FR-004; SC-001, SC-002, SC-003; constitution Release
// Process release-closing test) — it exercises the Go-side scaffold
// wiring this release added (the conditional `wrangler.jsonc` Worker
// `main`/`durable_objects`/`run_worker_first`/`nodejs_compat` block that
// mounts/omits the in-Worker MCP server), the same way
// test/integration/pilot-site/ does for v0.02's build/deploy wrapper.
//
// What this test automates: scaffold (mcp_enabled true and false) ->
// index -> the resulting wrangler.jsonc's MCP wiring (or lack thereof).
//
// What this test does NOT automate, by necessity, and was instead
// verified by hand in this session against a real Worker:
//   - A real MCP client (the official @modelcontextprotocol/inspector
//     CLI) completing a genuine connect -> list-tools -> call-tool
//     exchange against the real Workers runtime (`wrangler dev`),
//     receiving all 5 tools (list_pages, get_page, search_content,
//     get_profile, list_projects) and a correctly structured
//     text+resource result from search_content (SC-001, Scenario 2).
//   - A terminal/curl negotiation check: `curl -H "Accept:
//     text/markdown"` returning clean Markdown, and with a curl
//     User-Agent additionally returning real ANSI escape codes via
//     md4x's renderToAnsi(), both carrying `Vary: Accept, User-Agent`
//     (SC-002).
//   - The mcp_enabled:false "zero reachable trace" check: with no
//     Worker `main` deployed at all, `curl -I /mcp` against the real
//     Workers runtime returns a plain 404 with no MCP-flavored
//     response whatsoever (SC-003).
//
// Re-running that exact verification requires Node/npm, network access
// to install dependencies, and a real `wrangler dev` process — none of
// which `go test ./...` assumes elsewhere in this repo (matching
// test/integration/pilot-site/'s own precedent of stubbing the
// build/deploy wrappers rather than literally invoking npm/vite/wrangler
// in CI).
package pilot_site_mcp_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/monamaret/biblio/internal/index"
	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
)

// TestPilotSiteMCP_EnabledWiresFullMCPSurface scaffolds a real pilot
// profile project with mcp_enabled true and confirms wrangler.jsonc mounts
// every piece the in-Worker MCP server needs (FR-001; SC-001's
// precondition).
func TestPilotSiteMCP_EnabledWiresFullMCPSurface(t *testing.T) {
	presets := theme.PresetsForType("profile")
	if len(presets) == 0 {
		t.Fatal(`theme.PresetsForType("profile"): want at least one preset, got none`)
	}

	destDir := t.TempDir()
	projectDir, err := scaffold.Generate(destDir, scaffold.Options{
		ProjectName: "pilot-site-mcp-on",
		SiteType:    "profile",
		Domain:      "pilot.example.com",
		Theme:       presets[0].ID,
		MCPEnabled:  true,
	})
	if err != nil {
		t.Fatalf("scaffold.Generate(): %v", err)
	}

	contentDir := filepath.Join(projectDir, "content")
	idx, errs := index.Generate(contentDir)
	if len(errs) > 0 {
		t.Fatalf("index.Generate(): %v", errs)
	}
	if len(idx.Entries) == 0 {
		t.Fatal("index.Generate(): want at least one authored entry, got zero")
	}

	wrangler, err := os.ReadFile(filepath.Join(projectDir, "wrangler.jsonc"))
	if err != nil {
		t.Fatalf("reading wrangler.jsonc: %v", err)
	}
	got := string(wrangler)

	for _, want := range []string{
		`"main": "workers/entry.ts"`,
		`"compatibility_flags": ["nodejs_compat"]`,
		`"run_worker_first": true`,
		`"durable_objects"`,
		`"migrations"`,
		"pilot.example.com",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("wrangler.jsonc missing %q for an mcp_enabled site:\n%s", want, got)
		}
	}
}

// TestPilotSiteMCP_DisabledLeavesZeroTrace scaffolds the same pilot
// profile project with mcp_enabled false and confirms wrangler.jsonc
// contains absolutely no MCP-related route, binding, or artifact — the
// build-time guarantee SC-003's live "curl -I /mcp -> 404" check (verified
// by hand against a real Workers runtime) depends on.
func TestPilotSiteMCP_DisabledLeavesZeroTrace(t *testing.T) {
	presets := theme.PresetsForType("profile")
	if len(presets) == 0 {
		t.Fatal(`theme.PresetsForType("profile"): want at least one preset, got none`)
	}

	destDir := t.TempDir()
	projectDir, err := scaffold.Generate(destDir, scaffold.Options{
		ProjectName: "pilot-site-mcp-off",
		SiteType:    "profile",
		Domain:      "pilot.example.com",
		Theme:       presets[0].ID,
		MCPEnabled:  false,
	})
	if err != nil {
		t.Fatalf("scaffold.Generate(): %v", err)
	}

	wrangler, err := os.ReadFile(filepath.Join(projectDir, "wrangler.jsonc"))
	if err != nil {
		t.Fatalf("reading wrangler.jsonc: %v", err)
	}
	got := string(wrangler)

	for _, mustNotContain := range []string{`"main"`, "nodejs_compat", "durable_objects", "migrations", "run_worker_first", "MyMCP"} {
		if strings.Contains(got, mustNotContain) {
			t.Errorf("wrangler.jsonc unexpectedly contains %q for an mcp_enabled:false site:\n%s", mustNotContain, got)
		}
	}

	var parsed map[string]any
	if err := json.Unmarshal(wrangler, &parsed); err != nil {
		t.Fatalf("wrangler.jsonc is not valid JSON: %v", err)
	}
	if _, hasAssets := parsed["assets"]; !hasAssets {
		t.Error("wrangler.jsonc: want the assets-only config preserved even with MCP disabled")
	}
}
