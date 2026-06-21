package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// ensureAssetsDirectory must create wrangler.jsonc's configured
// assets.directory, defaulting to "dist/client" (the scaffolded default)
// when the field is missing — `wrangler dev` requires this directory to
// exist at startup, before any build has ever run, but `biblio-cli new`
// never created it (issue #194).
func TestEnsureAssetsDirectory_CreatesDefaultWhenAssetsFieldMissing(t *testing.T) {
	cwd := t.TempDir()
	wranglerConfig := filepath.Join(cwd, "wrangler.jsonc")
	if err := os.WriteFile(wranglerConfig, []byte(`{"name":"site"}`), 0o644); err != nil {
		t.Fatalf("setup: writing wrangler.jsonc: %v", err)
	}

	if err := ensureAssetsDirectory(cwd, wranglerConfig); err != nil {
		t.Fatalf("ensureAssetsDirectory() error = %v", err)
	}

	want := filepath.Join(cwd, "dist", "client")
	if info, err := os.Stat(want); err != nil || !info.IsDir() {
		t.Errorf("ensureAssetsDirectory(): %s missing or not a directory: %v", want, err)
	}
}

// A project that customizes wrangler.jsonc's assets.directory must get
// that directory created instead of the default "dist/client".
func TestEnsureAssetsDirectory_UsesConfiguredDirectory(t *testing.T) {
	cwd := t.TempDir()
	wranglerConfig := filepath.Join(cwd, "wrangler.jsonc")
	config := `{"assets":{"directory":"./custom/out"}}`
	if err := os.WriteFile(wranglerConfig, []byte(config), 0o644); err != nil {
		t.Fatalf("setup: writing wrangler.jsonc: %v", err)
	}

	if err := ensureAssetsDirectory(cwd, wranglerConfig); err != nil {
		t.Fatalf("ensureAssetsDirectory() error = %v", err)
	}

	want := filepath.Join(cwd, "custom", "out")
	if info, err := os.Stat(want); err != nil || !info.IsDir() {
		t.Errorf("ensureAssetsDirectory(): configured directory %s missing or not a directory: %v", want, err)
	}

	unwanted := filepath.Join(cwd, "dist", "client")
	if _, err := os.Stat(unwanted); err == nil {
		t.Errorf("ensureAssetsDirectory(): also created the default %s when a custom directory was configured", unwanted)
	}
}

// Running against an already-existing assets directory (e.g. after a
// prior build) must succeed, not error on the directory already being
// there.
func TestEnsureAssetsDirectory_IdempotentWhenAlreadyExists(t *testing.T) {
	cwd := t.TempDir()
	wranglerConfig := filepath.Join(cwd, "wrangler.jsonc")
	if err := os.WriteFile(wranglerConfig, []byte(`{}`), 0o644); err != nil {
		t.Fatalf("setup: writing wrangler.jsonc: %v", err)
	}
	existing := filepath.Join(cwd, "dist", "client")
	if err := os.MkdirAll(existing, 0o755); err != nil {
		t.Fatalf("setup: pre-creating dist/client: %v", err)
	}
	marker := filepath.Join(existing, "already-built.html")
	if err := os.WriteFile(marker, []byte("<html></html>"), 0o644); err != nil {
		t.Fatalf("setup: writing marker file: %v", err)
	}

	if err := ensureAssetsDirectory(cwd, wranglerConfig); err != nil {
		t.Fatalf("ensureAssetsDirectory() error = %v, want nil when the directory already exists", err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Errorf("ensureAssetsDirectory(): pre-existing build output was disturbed: %v", err)
	}
}
