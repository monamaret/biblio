package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// DeployCommand is the wrapped command argv — `npx wrangler deploy`, which
// publishes the built static output to the route/custom-domain binding
// already present in wrangler.jsonc (FR-002; research.md "Build/deploy
// command surface"). Exported so tests can substitute a stub command.
var DeployCommand = []string{"npx", "wrangler", "deploy"}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Publish a built site to its configured public address",
	Args:  cobra.NoArgs,
	RunE:  RunDeploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

// RunDeploy implements `biblio-cli deploy`: it fails early with a clear
// message if the current directory isn't a scaffolded project or hasn't
// been built yet, then runs `wrangler deploy` and forwards its output —
// surfacing the underlying failure (auth, route-binding, build-artifact) on
// a non-zero exit. Re-running against an already-live site is `wrangler
// deploy`'s own in-place update behavior, so no special-casing is needed
// here (contracts/cli-commands.md `biblio-cli deploy`; spec Edge Cases).
func RunDeploy(cmd *cobra.Command, _ []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("deploy: resolving working directory: %w", err)
	}

	wranglerConfig := filepath.Join(cwd, "wrangler.jsonc")
	if _, err := os.Stat(wranglerConfig); os.IsNotExist(err) {
		return fmt.Errorf("deploy: wrangler.jsonc not found in %s — run `biblio-cli new` first to scaffold a project", cwd)
	}

	buildOutput := filepath.Join(cwd, "dist", "client")
	if _, err := os.Stat(buildOutput); os.IsNotExist(err) {
		return fmt.Errorf("deploy: no build output found in %s — run `biblio-cli build` first", buildOutput)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Deploying (wrangler deploy)…")

	if err := runWrapped(cmd.Context(), cwd, DeployCommand, cmd.OutOrStdout()); err != nil {
		return fmt.Errorf("deploy: %w", err)
	}

	address, err := publicAddress(wranglerConfig)
	if err != nil {
		// wrangler itself succeeded; a missing/malformed route in our own
		// generated config is not a deploy failure — just skip the address line.
		fmt.Fprintln(cmd.OutOrStdout(), "Deploy complete — the site is live at its configured public address.")
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Deploy complete — the site is live at %s\n", address)
	return nil
}

// publicAddress reads the route pattern configured in wrangler.jsonc — the
// public address `biblio-cli new` pre-filled from --domain (data-model.md
// "Deployment" `public-address`).
func publicAddress(wranglerConfig string) (string, error) {
	data, err := os.ReadFile(wranglerConfig)
	if err != nil {
		return "", err
	}
	var cfg struct {
		Routes []struct {
			Pattern string `json:"pattern"`
		} `json:"routes"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}
	if len(cfg.Routes) == 0 || cfg.Routes[0].Pattern == "" {
		return "", fmt.Errorf("deploy: no route configured in %s", wranglerConfig)
	}
	return cfg.Routes[0].Pattern, nil
}
