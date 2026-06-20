package scaffold_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
)

// presetIDFor returns a ready-to-use preset ID for the given site type, or
// fails the test if none is defined yet (T011 defines them).
func presetIDFor(t *testing.T, siteType string) string {
	t.Helper()
	presets := theme.PresetsForType(siteType)
	if len(presets) == 0 {
		t.Fatalf("theme.PresetsForType(%q): want at least one ready-to-use preset, got none", siteType)
	}
	return presets[0].ID
}

// Generate must produce the populated project directory described in
// data-model.md "Site Project": content/, content/site.yaml (mcp_enabled:
// true), wrangler.jsonc (route from --domain), and tokens.css.
func TestGenerate_CreatesPopulatedProjectDirectory(t *testing.T) {
	destDir := t.TempDir()
	opts := scaffold.Options{
		ProjectName: "my-site",
		SiteType:    "profile",
		Domain:      "example.com",
		Theme:       presetIDFor(t, "profile"),
		MCPEnabled:  true,
	}

	projectDir, err := scaffold.Generate(destDir, opts)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	wantDir := filepath.Join(destDir, opts.ProjectName)
	if projectDir != wantDir {
		t.Errorf("Generate() project dir = %q, want %q", projectDir, wantDir)
	}

	if info, err := os.Stat(filepath.Join(projectDir, "content")); err != nil || !info.IsDir() {
		t.Errorf("Generate(): content/ directory missing or not a directory: %v", err)
	}

	siteYAML, err := os.ReadFile(filepath.Join(projectDir, "content", "site.yaml"))
	if err != nil {
		t.Fatalf("Generate(): reading content/site.yaml: %v", err)
	}
	if !strings.Contains(string(siteYAML), "mcp_enabled: true") {
		t.Errorf("Generate(): content/site.yaml missing `mcp_enabled: true`:\n%s", siteYAML)
	}

	wrangler, err := os.ReadFile(filepath.Join(projectDir, "wrangler.jsonc"))
	if err != nil {
		t.Fatalf("Generate(): reading wrangler.jsonc: %v", err)
	}
	if !strings.Contains(string(wrangler), opts.Domain) {
		t.Errorf("Generate(): wrangler.jsonc missing routing for domain %q:\n%s", opts.Domain, wrangler)
	}

	if _, err := os.Stat(filepath.Join(projectDir, "tokens.css")); err != nil {
		t.Errorf("Generate(): tokens.css missing: %v", err)
	}
}

// Each of the three site types must seed content matching its shape (FR-002).
func TestGenerate_SeedsContentPerSiteType(t *testing.T) {
	for _, siteType := range scaffold.SiteTypes {
		t.Run(siteType, func(t *testing.T) {
			destDir := t.TempDir()
			opts := scaffold.Options{
				ProjectName: "site-" + siteType,
				SiteType:    siteType,
				Domain:      "example.com",
				Theme:       presetIDFor(t, siteType),
				MCPEnabled:  true,
			}

			projectDir, err := scaffold.Generate(destDir, opts)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			entries, err := os.ReadDir(filepath.Join(projectDir, "content"))
			if err != nil {
				t.Fatalf("Generate(): reading content/: %v", err)
			}
			if len(entries) == 0 {
				t.Errorf("Generate(): content/ for site-type %q has no seeded entries", siteType)
			}
		})
	}
}

// When MCPEnabled is true, wrangler.jsonc must mount the Worker `main`
// entry plus the `durable_objects` binding/`migrations` and
// `nodejs_compat` flag the in-Worker MCP server requires (v0.03 FR-001;
// data-model.md "MCP Server Mount" `mounted`; research spike 002).
func TestGenerate_MCPEnabledWiresWranglerJSONC(t *testing.T) {
	destDir := t.TempDir()
	opts := scaffold.Options{
		ProjectName: "mcp-on",
		SiteType:    "profile",
		Domain:      "example.com",
		Theme:       presetIDFor(t, "profile"),
		MCPEnabled:  true,
	}

	projectDir, err := scaffold.Generate(destDir, opts)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	wrangler, err := os.ReadFile(filepath.Join(projectDir, "wrangler.jsonc"))
	if err != nil {
		t.Fatalf("Generate(): reading wrangler.jsonc: %v", err)
	}
	got := string(wrangler)

	for _, want := range []string{
		`"main": "workers/entry.ts"`,
		`"compatibility_flags": ["nodejs_compat"]`,
		`"durable_objects"`,
		`"migrations"`,
		`"binding": "ASSETS"`,
		`"run_worker_first": true`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Generate(): wrangler.jsonc missing %q when MCPEnabled is true:\n%s", want, got)
		}
	}
}

// When MCPEnabled is false, wrangler.jsonc must contain zero MCP-related
// route, artifact, or binding — no "disabled but present" intermediate
// state (v0.03 FR-002; spec Scenario 4; data-model.md "MCP Server Mount").
func TestGenerate_MCPDisabledOmitsAllMCPWiring(t *testing.T) {
	destDir := t.TempDir()
	opts := scaffold.Options{
		ProjectName: "mcp-off",
		SiteType:    "profile",
		Domain:      "example.com",
		Theme:       presetIDFor(t, "profile"),
		MCPEnabled:  false,
	}

	projectDir, err := scaffold.Generate(destDir, opts)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	wrangler, err := os.ReadFile(filepath.Join(projectDir, "wrangler.jsonc"))
	if err != nil {
		t.Fatalf("Generate(): reading wrangler.jsonc: %v", err)
	}
	got := string(wrangler)

	for _, mustNotContain := range []string{
		`"main"`,
		`nodejs_compat`,
		`durable_objects`,
		`migrations`,
		`MyMCP`,
	} {
		if strings.Contains(got, mustNotContain) {
			t.Errorf("Generate(): wrangler.jsonc unexpectedly contains %q when MCPEnabled is false:\n%s", mustNotContain, got)
		}
	}

	siteYAML, err := os.ReadFile(filepath.Join(projectDir, "content", "site.yaml"))
	if err != nil {
		t.Fatalf("Generate(): reading content/site.yaml: %v", err)
	}
	if !strings.Contains(string(siteYAML), "mcp_enabled: false") {
		t.Errorf("Generate(): content/site.yaml missing `mcp_enabled: false`:\n%s", siteYAML)
	}
}

// A project-name colliding with an existing directory must be rejected, and
// must never leave a partially-written project behind (contracts/cli-commands.md).
func TestGenerate_CollisionDetection(t *testing.T) {
	destDir := t.TempDir()
	opts := scaffold.Options{
		ProjectName: "taken",
		SiteType:    "profile",
		Domain:      "example.com",
		Theme:       presetIDFor(t, "profile"),
		MCPEnabled:  true,
	}

	existing := filepath.Join(destDir, opts.ProjectName)
	if err := os.Mkdir(existing, 0o755); err != nil {
		t.Fatalf("setup: creating colliding directory: %v", err)
	}
	preexisting, err := os.ReadDir(existing)
	if err != nil {
		t.Fatalf("setup: reading colliding directory: %v", err)
	}

	if _, err := scaffold.Generate(destDir, opts); err == nil {
		t.Fatal("Generate(): want error for colliding project-name, got nil")
	}

	after, err := os.ReadDir(existing)
	if err != nil {
		t.Fatalf("Generate(): colliding directory vanished: %v", err)
	}
	if len(after) != len(preexisting) {
		t.Errorf("Generate(): colliding directory was modified on collision (want untouched): had %d entries, now %d", len(preexisting), len(after))
	}
}
