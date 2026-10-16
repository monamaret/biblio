// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/monamaret/biblio/internal/devserver"
	"github.com/monamaret/biblio/internal/prebuild"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start local preview: vite dev + wrangler dev as one supervised command",
	Args:  cobra.NoArgs,
	RunE:  RunDev,
}

// DevViteCommand and DevWranglerCommand are the two supervised dev-server
// argvs. Exported so tests can substitute stub commands, as BuildCommand
// is for `build`.
var (
	DevViteCommand     = []string{"npx", "vite", "dev"}
	DevWranglerCommand = []string{"npx", "wrangler", "dev"}
)

func init() {
	rootCmd.AddCommand(devCmd)
}

// RunDev implements `biblio-cli dev`: it resolves the wrangler.jsonc in the
// current project, runs the pre-build step once in Dev mode (drafts
// included; specs/004-quarterly-site-type FR-012, FR-035), launches vite dev and wrangler dev as a supervised pair,
// surfaces both servers' addresses, and exits cleanly on interrupt
// (contracts/cli-commands.md `biblio-cli dev`; FR-005).
func RunDev(cmd *cobra.Command, _ []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("dev: resolving working directory: %w", err)
	}

	// Fail early with a clear message if wrangler.jsonc is missing — matching
	// the contract's documented failure mode for this command.
	wranglerConfig := filepath.Join(cwd, "wrangler.jsonc")
	if _, err := os.Stat(wranglerConfig); os.IsNotExist(err) {
		return fmt.Errorf("dev: wrangler.jsonc not found in %s — run `biblio-cli new` first to scaffold a project", cwd)
	}

	// wrangler dev requires its assets directory (wrangler.jsonc
	// `assets.directory`) to exist at startup, before any build has ever
	// run — `biblio-cli new` doesn't create it (dist/ is generated, not
	// scaffolded), and it's also gitignored, so a fresh clone of an
	// already-scaffolded project hits the same gap on a clean checkout,
	// not just right after `new`. Ensuring it here, right before
	// launching, covers both cases at the one point that actually needs
	// it to exist (issue #194).
	if err := ensureAssetsDirectory(cwd, wranglerConfig); err != nil {
		return fmt.Errorf("dev: %w", err)
	}

	if err := runPrebuild(cmd, cwd, prebuild.Dev); err != nil {
		return fmt.Errorf("dev: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Starting local preview (vite dev + wrangler dev)…")
	fmt.Fprintln(cmd.OutOrStdout(), "Press Ctrl+C to stop.")

	s := devserver.NewSupervisor(devserver.Config{
		ViteCmd:     DevViteCommand,
		WranglerCmd: DevWranglerCommand,
		Output:      cmd.OutOrStdout(),
	})

	return s.Run(cmd.Context())
}

// ensureAssetsDirectory creates wrangler.jsonc's configured
// `assets.directory` (relative to cwd) if it doesn't already exist,
// defaulting to "dist/client" — the scaffolded default — when the field is
// missing or unparseable, rather than failing dev over a config-reading
// detail unrelated to its actual purpose.
func ensureAssetsDirectory(cwd, wranglerConfig string) error {
	dir := filepath.Join(cwd, "dist", "client")

	if data, err := os.ReadFile(wranglerConfig); err == nil {
		var cfg struct {
			Assets struct {
				Directory string `json:"directory"`
			} `json:"assets"`
		}
		if json.Unmarshal(data, &cfg) == nil && cfg.Assets.Directory != "" {
			dir = filepath.Join(cwd, filepath.FromSlash(cfg.Assets.Directory))
		}
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating assets directory %s: %w", dir, err)
	}
	return nil
}
