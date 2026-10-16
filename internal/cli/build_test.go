package cli_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/monamaret/biblio/internal/cli"
)

// scaffoldedDir creates a temp dir containing a minimal wrangler.jsonc so
// build/deploy preconditions pass, and chdirs into it for the duration of
// the test.
func scaffoldedDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "wrangler.jsonc"), []byte("{}"), 0o644); err != nil {
		t.Fatalf("setup: writing wrangler.jsonc: %v", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("setup: chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	return dir
}

// Build succeeds with zero missing/broken pages when the wrapped `vite
// build` exits cleanly (FR-001; spec Scenario 1; contracts/cli-commands.md
// `build` "Output (success)").
func TestRunBuild_Success(t *testing.T) {
	scaffoldedDir(t)
	original := cli.BuildCommand
	cli.BuildCommand = []string{"sh", "-c", "echo prerendered 4 pages && exit 0"}
	defer func() { cli.BuildCommand = original }()

	var buf strings.Builder
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	cmd.SetContext(context.Background())

	if err := cli.RunBuild(cmd, nil); err != nil {
		t.Fatalf("RunBuild() error = %v, want nil", err)
	}
	if !strings.Contains(buf.String(), "Build complete") {
		t.Errorf("RunBuild() output = %q, want it to report a complete build", buf.String())
	}
}

// Build fails with a clear, per-file message rather than silently skipping
// the page when the wrapped `vite build` surfaces a prerender validation
// failure (spec Edge Cases; contracts/cli-commands.md `build` "Output
// (failure)").
func TestRunBuild_FailsWithClearMessageOnMalformedContent(t *testing.T) {
	scaffoldedDir(t)
	original := cli.BuildCommand
	cli.BuildCommand = []string{"sh", "-c", "echo 'content/about.md: malformed frontmatter' >&2 && exit 1"}
	defer func() { cli.BuildCommand = original }()

	var buf strings.Builder
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	cmd.SetContext(context.Background())

	err := cli.RunBuild(cmd, nil)
	if err == nil {
		t.Fatal("RunBuild() error = nil, want non-nil on prerender validation failure")
	}
	if !strings.Contains(buf.String(), "content/about.md: malformed frontmatter") {
		t.Errorf("RunBuild() output = %q, want it to surface the per-file failure message", buf.String())
	}
}

// Build fails early with a clear message when run outside a scaffolded
// project — never attempting to invoke vite build without a wrangler.jsonc.
func TestRunBuild_FailsWithoutScaffoldedProject(t *testing.T) {
	dir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("setup: chdir: %v", err)
	}
	defer func() { _ = os.Chdir(cwd) }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err = cli.RunBuild(cmd, nil)
	if err == nil {
		t.Fatal("RunBuild() error = nil, want non-nil outside a scaffolded project")
	}
	if !strings.Contains(err.Error(), "wrangler.jsonc") {
		t.Errorf("RunBuild() error = %q, want it to name the missing wrangler.jsonc", err)
	}
}

// writeLogContent turns dir into a minimal log content project: a valid
// site.yaml plus one entry, or an invalid site.yaml when invalid is true.
func writeLogContent(t *testing.T, dir string, invalid bool) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, "content", "posts"), 0o755); err != nil {
		t.Fatal(err)
	}
	siteYAML := "site_type: log\nname: test\n"
	if invalid {
		siteYAML += "typesetting:\n  accent: \"#000000\"\n"
	}
	if err := os.WriteFile(filepath.Join(dir, "content", "site.yaml"), []byte(siteYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	post := "---\ntitle: \"Hello\"\ndate: \"2026-06-01\"\n---\n\nHi.\n"
	if err := os.WriteFile(filepath.Join(dir, "content", "posts", "2026-06-01-hello.md"), []byte(post), 0o644); err != nil {
		t.Fatal(err)
	}
}

// The pre-build step (validate + index) runs before the wrapped `vite
// build`, so the build always sees a fresh content-index.json (specs/004
// FR-035; contracts/cli.md `biblio-cli build`).
func TestRunBuild_RunsPrebuildBeforeVite(t *testing.T) {
	dir := scaffoldedDir(t)
	writeLogContent(t, dir, false)
	original := cli.BuildCommand
	cli.BuildCommand = []string{"sh", "-c", "test -f content-index.json && touch saw-index"}
	defer func() { cli.BuildCommand = original }()

	var buf strings.Builder
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	cmd.SetContext(context.Background())

	if err := cli.RunBuild(cmd, nil); err != nil {
		t.Fatalf("RunBuild() error = %v\noutput:\n%s", err, buf.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "saw-index")); err != nil {
		t.Error("vite build ran before content-index.json existed; the pre-build step must run first")
	}
	for _, want := range []string{"Validating content/site.yaml", "Indexing: 1 item(s)"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("output = %q, want it to contain %q", buf.String(), want)
		}
	}
}

// A settings validation error stops `build` before `vite build` starts.
func TestRunBuild_ValidationErrorStopsBeforeVite(t *testing.T) {
	dir := scaffoldedDir(t)
	writeLogContent(t, dir, true)
	original := cli.BuildCommand
	cli.BuildCommand = []string{"sh", "-c", "touch vite-ran"}
	defer func() { cli.BuildCommand = original }()

	cmd := &cobra.Command{}
	cmd.SetOut(&strings.Builder{})
	cmd.SetContext(context.Background())

	err := cli.RunBuild(cmd, nil)
	if err == nil {
		t.Fatal("RunBuild() error = nil, want the site.yaml validation error")
	}
	if !strings.Contains(err.Error(), "typesetting is only supported") {
		t.Errorf("error = %q, want the validation message", err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "vite-ran")); statErr == nil {
		t.Error("vite build ran despite a validation error")
	}
}
