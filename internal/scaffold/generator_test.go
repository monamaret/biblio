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

// Generate must also record a scaffold manifest (.biblio/scaffold-manifest.json)
// capturing the site type and a content hash for every embedded template
// file outside content/ (content/ is the author's domain from the moment
// it's scaffolded — Principle II — and is never part of the diffable
// baseline). internal/migrate's update flow uses this manifest as its
// 3-way-diff baseline (data-model.md "Scaffold Region Classification"
// classification source).
func TestGenerate_WritesScaffoldManifest(t *testing.T) {
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

	manifest, err := scaffold.ReadManifest(projectDir)
	if err != nil {
		t.Fatalf("ReadManifest() error = %v", err)
	}

	if manifest.SiteType != "profile" {
		t.Errorf("manifest.SiteType = %q, want %q", manifest.SiteType, "profile")
	}
	if len(manifest.Files) == 0 {
		t.Fatal("manifest.Files: want at least one entry, got none")
	}
	for path := range manifest.Files {
		if strings.HasPrefix(path, "content/") {
			t.Errorf("manifest.Files contains %q: content/ paths must never be part of the baseline", path)
		}
	}
	if _, ok := manifest.Files["package.json"]; !ok {
		t.Error(`manifest.Files["package.json"]: want an entry, got none`)
	}

	// The hash must match the file actually written to disk, so a later
	// diff against a freshly-embedded template is meaningful.
	data, err := os.ReadFile(filepath.Join(projectDir, "package.json"))
	if err != nil {
		t.Fatalf("reading written package.json: %v", err)
	}
	gotHash := scaffold.HashBytes(data)
	if manifest.Files["package.json"] != gotHash {
		t.Errorf("manifest.Files[%q] = %q, want %q (hash of the file actually written)", "package.json", manifest.Files["package.json"], gotHash)
	}
}

// WorkerName must produce a valid Cloudflare Workers script name (lowercase
// letters, digits, hyphens only) from any project name, including the
// domain-shaped names people naturally pick (issue #195).
func TestWorkerName(t *testing.T) {
	validName := func(t *testing.T, name string) {
		t.Helper()
		if name == "" {
			t.Error("WorkerName(): want non-empty result")
		}
		for _, r := range name {
			isLower := r >= 'a' && r <= 'z'
			isDigit := r >= '0' && r <= '9'
			if !isLower && !isDigit && r != '-' {
				t.Errorf("WorkerName() = %q: contains invalid character %q", name, r)
			}
		}
		if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
			t.Errorf("WorkerName() = %q: must not start or end with a hyphen", name)
		}
	}

	tests := []struct {
		projectName string
		want        string
	}{
		// The real-world case this issue was filed against: a project
		// named after the domain it backs.
		{"notebook.monamaret.com", "notebook-monamaret-com"},
		{"notebook-monamaret.com", "notebook-monamaret-com"},
		{"My Cool Site", "my-cool-site"},
		{"already-valid-name", "already-valid-name"},
	}

	for _, tt := range tests {
		t.Run(tt.projectName, func(t *testing.T) {
			got := scaffold.WorkerName(tt.projectName)
			validName(t, got)
			if got != tt.want {
				t.Errorf("WorkerName(%q) = %q, want %q", tt.projectName, got, tt.want)
			}
		})
	}
}

// A project name containing characters invalid for a Workers script name
// (e.g. dots, as in a domain-shaped name) must still scaffold successfully,
// with wrangler.jsonc's "name" sanitized but the directory name and
// content/site.yaml's `name` left exactly as given — neither of those has
// the Workers naming constraint (issue #195).
func TestGenerate_SanitizesWranglerNameButNotDirectoryOrSiteYAML(t *testing.T) {
	destDir := t.TempDir()
	opts := scaffold.Options{
		ProjectName: "notebook.monamaret.com",
		SiteType:    "log",
		Domain:      "notebook.monamaret.com",
		Theme:       presetIDFor(t, "log"),
		MCPEnabled:  true,
	}

	projectDir, err := scaffold.Generate(destDir, opts)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	wantDir := filepath.Join(destDir, "notebook.monamaret.com")
	if projectDir != wantDir {
		t.Errorf("Generate() project dir = %q, want %q (directory name must be unaffected by Workers-name sanitization)", projectDir, wantDir)
	}

	siteYAML, err := os.ReadFile(filepath.Join(projectDir, "content", "site.yaml"))
	if err != nil {
		t.Fatalf("Generate(): reading content/site.yaml: %v", err)
	}
	if !strings.Contains(string(siteYAML), "name: notebook.monamaret.com") {
		t.Errorf("Generate(): content/site.yaml name must keep the literal project name:\n%s", siteYAML)
	}

	wrangler, err := os.ReadFile(filepath.Join(projectDir, "wrangler.jsonc"))
	if err != nil {
		t.Fatalf("Generate(): reading wrangler.jsonc: %v", err)
	}
	got := string(wrangler)
	if !strings.Contains(got, `"name": "notebook-monamaret-com"`) {
		t.Errorf("Generate(): wrangler.jsonc \"name\" must be the sanitized Workers-safe slug:\n%s", got)
	}
	if strings.Contains(got, `"name": "notebook.monamaret.com"`) {
		t.Errorf("Generate(): wrangler.jsonc \"name\" must not be the unsanitized project name:\n%s", got)
	}
}
