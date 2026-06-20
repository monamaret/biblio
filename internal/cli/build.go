package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// BuildCommand is the wrapped command argv — `npx vite build`, which drives
// React Router v7's native `prerender` (configured in the project's
// react-router.config.ts) to pre-render every content-driven page (FR-001;
// research.md "Pre-rendering approach"). Exported so tests can substitute a
// stub command.
var BuildCommand = []string{"npx", "vite", "build"}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Pre-render every content-driven page into static output",
	Args:  cobra.NoArgs,
	RunE:  RunBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

// RunBuild implements `biblio-cli build`: it fails early with a clear
// message if the current directory isn't a scaffolded project, then runs
// `vite build` and forwards its output — including any per-file prerender
// validation failure (spec Edge Cases) — surfacing a non-zero exit on
// failure (contracts/cli-commands.md `biblio-cli build`).
func RunBuild(cmd *cobra.Command, _ []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("build: resolving working directory: %w", err)
	}

	wranglerConfig := filepath.Join(cwd, "wrangler.jsonc")
	if _, err := os.Stat(wranglerConfig); os.IsNotExist(err) {
		return fmt.Errorf("build: wrangler.jsonc not found in %s — run `biblio-cli new` first to scaffold a project", cwd)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Building (vite build, prerendering content-driven pages)…")

	if err := runWrapped(cmd.Context(), cwd, BuildCommand, cmd.OutOrStdout()); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Build complete — every content-driven page pre-rendered with no missing or broken pages.")
	return nil
}
