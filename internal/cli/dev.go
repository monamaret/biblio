package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/monamaret/biblio/internal/devserver"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start local preview: vite dev + wrangler dev as one supervised command",
	Args:  cobra.NoArgs,
	RunE:  runDev,
}

func init() {
	rootCmd.AddCommand(devCmd)
}

// runDev implements `biblio-cli dev`: it resolves the wrangler.jsonc in the
// current project, launches vite dev and wrangler dev as a supervised pair,
// surfaces both servers' addresses, and exits cleanly on interrupt
// (contracts/cli-commands.md `biblio-cli dev`; FR-005).
func runDev(cmd *cobra.Command, _ []string) error {
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

	fmt.Fprintln(cmd.OutOrStdout(), "Starting local preview (vite dev + wrangler dev)…")
	fmt.Fprintln(cmd.OutOrStdout(), "Press Ctrl+C to stop.")

	s := devserver.NewSupervisor(devserver.Config{
		ViteCmd:     []string{"npx", "vite", "dev"},
		WranglerCmd: []string{"npx", "wrangler", "dev"},
		Output:      cmd.OutOrStdout(),
	})

	return s.Run(cmd.Context())
}
