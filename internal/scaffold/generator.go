package scaffold

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/monamaret/biblio/internal/theme"
)

// Options describes a single `biblio-cli new` invocation — the inputs that
// together determine a Site Project's generated shape (data-model.md
// "Site Project").
type Options struct {
	// ProjectName is the new project's directory name.
	ProjectName string
	// SiteType selects the scaffold template, theme presets, and seed
	// content shape — one of SiteTypes.
	SiteType string
	// Domain pre-fills the route/custom-domain binding in wrangler.jsonc.
	Domain string
	// Theme is the selected preset's ID; it must belong to SiteType.
	Theme string
	// MCPEnabled is the sole MCP-related configuration surface
	// (Security & Operational Posture): when true, wrangler.jsonc mounts
	// the Worker `main` entry, its `durable_objects` binding/`migrations`
	// for the in-Worker MCP server, and `nodejs_compat` (data-model.md
	// "MCP Server Mount" `mounted`; research spike 002). When false, none
	// of that is present at all — a build-time decision, not a runtime
	// toggle (v0.03 FR-002).
	MCPEnabled bool
}

// Generate scaffolds a new Site Project as destDir/<ProjectName>: it seeds
// content/ from the embedded site-type template, then writes
// content/site.yaml (with mcp_enabled: true), wrangler.jsonc (routing from
// Domain), and tokens.css (from the selected Theme preset) — the artifact
// set the `new` contract promises (contracts/cli-commands.md `biblio-cli
// new`). On any failure the project directory is removed, leaving destDir
// exactly as it was — never a partially-written project.
func Generate(destDir string, opts Options) (string, error) {
	if !isSiteType(opts.SiteType) {
		return "", fmt.Errorf("scaffold: unknown site type %q (want one of %v)", opts.SiteType, SiteTypes)
	}

	preset, ok := theme.PresetByID(opts.Theme)
	if !ok || preset.SiteType != opts.SiteType {
		return "", fmt.Errorf("scaffold: theme %q is not a valid preset for site type %q", opts.Theme, opts.SiteType)
	}

	projectDir := filepath.Join(destDir, opts.ProjectName)
	if _, err := os.Stat(projectDir); err == nil {
		return "", fmt.Errorf("scaffold: project directory %q already exists", opts.ProjectName)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("scaffold: checking project directory %q: %w", opts.ProjectName, err)
	}

	tokensCSS, err := theme.GenerateTokensCSS(preset)
	if err != nil {
		return "", fmt.Errorf("scaffold: generating tokens.css: %w", err)
	}

	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		return "", fmt.Errorf("scaffold: creating project directory %q: %w", opts.ProjectName, err)
	}
	wrote := false
	defer func() {
		if !wrote {
			os.RemoveAll(projectDir)
		}
	}()

	if err := copyTemplateTree(projectDir, opts.SiteType); err != nil {
		return "", err
	}

	contentDir := filepath.Join(projectDir, "content")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		return "", fmt.Errorf("scaffold: creating content directory: %w", err)
	}
	if err := os.WriteFile(filepath.Join(contentDir, "site.yaml"), []byte(renderSiteYAML(opts)), 0o644); err != nil {
		return "", fmt.Errorf("scaffold: writing content/site.yaml: %w", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "wrangler.jsonc"), []byte(renderWranglerJSONC(opts)), 0o644); err != nil {
		return "", fmt.Errorf("scaffold: writing wrangler.jsonc: %w", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "tokens.css"), []byte(tokensCSS), 0o644); err != nil {
		return "", fmt.Errorf("scaffold: writing tokens.css: %w", err)
	}

	manifest, err := BuildManifest(opts.SiteType)
	if err != nil {
		return "", err
	}
	if err := WriteManifest(projectDir, manifest); err != nil {
		return "", err
	}

	wrote = true
	return projectDir, nil
}

func isSiteType(siteType string) bool {
	for _, t := range SiteTypes {
		if t == siteType {
			return true
		}
	}
	return false
}

// copyTemplateTree copies every embedded template file for siteType into
// dir, preserving relative paths (Walk skips the .gitkeep VCS placeholders).
func copyTemplateTree(dir, siteType string) error {
	assets, err := Assets(siteType)
	if err != nil {
		return fmt.Errorf("scaffold: loading template assets for %q: %w", siteType, err)
	}
	return Walk(siteType, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(assets, path)
		if err != nil {
			return fmt.Errorf("scaffold: reading template %q: %w", path, err)
		}
		dest := filepath.Join(dir, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("scaffold: creating directory for %q: %w", path, err)
		}
		if err := os.WriteFile(dest, data, 0o644); err != nil {
			return fmt.Errorf("scaffold: writing %q: %w", path, err)
		}
		return nil
	})
}

// renderSiteYAML produces content/site.yaml: the site's declared type,
// name/description, and the sole MCP config surface, mcp_enabled (default
// true; data-model.md "Site Project"; Security & Operational Posture).
func renderSiteYAML(opts Options) string {
	return fmt.Sprintf(`site_type: %s
name: %s
description: ""
mcp_enabled: %t
`, opts.SiteType, opts.ProjectName, opts.MCPEnabled)
}

// renderWranglerJSONC produces wrangler.jsonc with its route/custom-domain
// binding pre-filled from --domain (data-model.md "Site Project"). When
// opts.MCPEnabled is true, it also mounts the Worker `main` entry plus the
// `durable_objects` binding/`migrations` and `nodejs_compat` flag the
// in-Worker MCP server (`agents`/`McpAgent`) requires — confirmed via
// research spike 002 — and sets `assets.run_worker_first: true`, without
// which Workers Static Assets serves any matching prerendered page
// directly and never invokes workers/entry.ts at all, silently breaking
// content negotiation for every existing page (verified against a real
// `wrangler dev` run). When MCPEnabled is false, none of this is present
// at all: no "disabled but present" intermediate state (v0.03 FR-002;
// data-model.md "MCP Server Mount").
func renderWranglerJSONC(opts Options) string {
	if !opts.MCPEnabled {
		return fmt.Sprintf(`{
  "name": %q,
  "compatibility_date": "2026-01-01",
  "assets": {
    "directory": "./dist/client"
  },
  "routes": [
    { "pattern": %q, "custom_domain": true }
  ]
}
`, opts.ProjectName, opts.Domain)
	}

	return fmt.Sprintf(`{
  "name": %q,
  "main": "workers/entry.ts",
  "compatibility_date": "2026-01-01",
  "compatibility_flags": ["nodejs_compat"],
  "assets": {
    "directory": "./dist/client",
    "binding": "ASSETS",
    "run_worker_first": true
  },
  "durable_objects": {
    "bindings": [{ "name": "MyMCP", "class_name": "MyMCP" }]
  },
  "migrations": [{ "tag": "v1", "new_sqlite_classes": ["MyMCP"] }],
  "routes": [
    { "pattern": %q, "custom_domain": true }
  ]
}
`, opts.ProjectName, opts.Domain)
}
