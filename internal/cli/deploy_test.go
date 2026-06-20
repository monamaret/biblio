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

// builtDir creates a temp dir with a scaffolded wrangler.jsonc (route
// "my-site.example.com") and a built dist/client output, and chdirs into it
// for the duration of the test.
func builtDir(t *testing.T) string {
	t.Helper()
	dir := scaffoldedDirWithRoute(t, "my-site.example.com")
	if err := os.MkdirAll(filepath.Join(dir, "dist", "client"), 0o755); err != nil {
		t.Fatalf("setup: creating dist/client: %v", err)
	}
	return dir
}

// scaffoldedDirWithRoute is like scaffoldedDir but writes a wrangler.jsonc
// with a real route pattern so deploy can surface the configured public
// address (data-model.md "Deployment" `public-address`).
func scaffoldedDirWithRoute(t *testing.T, pattern string) string {
	t.Helper()
	dir := t.TempDir()
	wrangler := `{"name":"my-site","routes":[{"pattern":"` + pattern + `","custom_domain":true}]}`
	if err := os.WriteFile(filepath.Join(dir, "wrangler.jsonc"), []byte(wrangler), 0o644); err != nil {
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

// Deploy succeeds against a built project and surfaces the configured
// public address (FR-002; contracts/cli-commands.md `deploy` "Output
// (success)"; data-model.md "Deployment" `public-address`).
func TestRunDeploy_SucceedsAndSurfacesPublicAddress(t *testing.T) {
	builtDir(t)
	original := cli.DeployCommand
	cli.DeployCommand = []string{"sh", "-c", "exit 0"}
	defer func() { cli.DeployCommand = original }()

	var buf strings.Builder
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	cmd.SetContext(context.Background())

	if err := cli.RunDeploy(cmd, nil); err != nil {
		t.Fatalf("RunDeploy() error = %v, want nil", err)
	}
	if !strings.Contains(buf.String(), "my-site.example.com") {
		t.Errorf("RunDeploy() output = %q, want it to surface the configured public address", buf.String())
	}
}

// Deploy surfaces the underlying wrangler/CI failure with diagnosable
// context rather than a silent partial publish (contracts/cli-commands.md
// `deploy` "Output (failure)").
func TestRunDeploy_SurfacesUnderlyingFailure(t *testing.T) {
	builtDir(t)
	original := cli.DeployCommand
	cli.DeployCommand = []string{"sh", "-c", "echo 'Authentication error: not logged in' >&2 && exit 1"}
	defer func() { cli.DeployCommand = original }()

	var buf strings.Builder
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	cmd.SetContext(context.Background())

	err := cli.RunDeploy(cmd, nil)
	if err == nil {
		t.Fatal("RunDeploy() error = nil, want non-nil on wrangler failure")
	}
	if !strings.Contains(buf.String(), "Authentication error") {
		t.Errorf("RunDeploy() output = %q, want it to surface the underlying failure", buf.String())
	}
}

// Deploy fails early with a clear message when no build artifact exists —
// never attempting to invoke wrangler deploy without one.
func TestRunDeploy_FailsWithoutBuildArtifact(t *testing.T) {
	scaffoldedDirWithRoute(t, "my-site.example.com")

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := cli.RunDeploy(cmd, nil)
	if err == nil {
		t.Fatal("RunDeploy() error = nil, want non-nil without a build artifact")
	}
	if !strings.Contains(err.Error(), "biblio-cli build") {
		t.Errorf("RunDeploy() error = %q, want it to direct the author to run `biblio-cli build`", err)
	}
}

// Re-running deploy against an already-live site (a second deploy with the
// same wrangler.jsonc and build output) updates cleanly in place — it must
// succeed again rather than erroring on "already deployed" (spec Edge
// Cases; data-model.md "Deployment" re-deploy-semantics invariant).
func TestRunDeploy_ReDeployAgainstLiveSiteUpdatesCleanly(t *testing.T) {
	builtDir(t)
	original := cli.DeployCommand
	cli.DeployCommand = []string{"sh", "-c", "exit 0"}
	defer func() { cli.DeployCommand = original }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	if err := cli.RunDeploy(cmd, nil); err != nil {
		t.Fatalf("first RunDeploy() error = %v, want nil", err)
	}

	var buf strings.Builder
	cmd2 := &cobra.Command{}
	cmd2.SetOut(&buf)
	cmd2.SetContext(context.Background())
	if err := cli.RunDeploy(cmd2, nil); err != nil {
		t.Fatalf("re-deploy RunDeploy() error = %v, want a clean in-place update with no error", err)
	}
	if !strings.Contains(buf.String(), "Deploy complete") {
		t.Errorf("re-deploy output = %q, want a clean success report on re-deploy", buf.String())
	}
}
