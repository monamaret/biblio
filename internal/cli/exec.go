package cli

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

// runWrapped runs argv in dir, forwarding combined stdout/stderr to out, and
// returns a descriptive error naming the wrapped command on non-zero exit.
// build.go and deploy.go are both thin wrappers over an external tool
// (`vite build`, `wrangler deploy`) — this is their one shared exec path so
// neither hand-rolls process plumbing (contracts/cli-commands.md `build`,
// `deploy`: "thin wrappers ... never diverge from the underlying tool's
// behavior").
func runWrapped(ctx context.Context, dir string, argv []string, out io.Writer) error {
	if len(argv) == 0 {
		return fmt.Errorf("empty command")
	}
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = dir
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", argv[0], err)
	}
	return nil
}
