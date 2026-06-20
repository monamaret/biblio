package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/monamaret/biblio/internal/index"
)

var indexCmd = &cobra.Command{
	Use:   "index [path]",
	Short: "Generate the canonical content-index.json from a content/ directory",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runIndex,
}

func init() {
	rootCmd.AddCommand(indexCmd)
}

// runIndex implements `biblio-cli index [path]`: it reads every .md file
// under the content directory (defaulting to ./content in the working
// directory), validates each item against the site type declared in
// site.yaml, writes content-index.json on success, and reports every
// failing file on validation error — never just the first
// (contracts/cli-commands.md `biblio-cli index`).
func runIndex(cmd *cobra.Command, args []string) error {
	contentDir, err := resolveContentDir(args)
	if err != nil {
		return err
	}

	idx, errs := index.Generate(contentDir)
	if len(errs) > 0 {
		var sb strings.Builder
		fmt.Fprintf(&sb, "index: %d file(s) failed schema validation:\n", len(errs))
		for _, e := range errs {
			fmt.Fprintf(&sb, "  • %v\n", e)
		}
		return fmt.Errorf("%s", sb.String())
	}

	outPath := filepath.Join(filepath.Dir(contentDir), "content-index.json")
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("index: serialising content-index.json: %w", err)
	}
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		return fmt.Errorf("index: writing content-index.json: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Indexed %d item(s) → %s\n", len(idx.Entries), outPath)
	return nil
}

// resolveContentDir returns the absolute path to the content/ directory.
// If an explicit path was given as an argument it is used; otherwise the
// working directory's ./content subdirectory is the default
// (contracts/cli-commands.md `biblio-cli index` argument behaviour).
func resolveContentDir(args []string) (string, error) {
	if len(args) == 1 {
		abs, err := filepath.Abs(args[0])
		if err != nil {
			return "", fmt.Errorf("index: resolving path %q: %w", args[0], err)
		}
		return abs, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("index: resolving working directory: %w", err)
	}
	return filepath.Join(cwd, "content"), nil
}
