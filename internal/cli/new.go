package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
	"github.com/monamaret/biblio/internal/tui"
)

var (
	newType       string
	newDomain     string
	newTheme      string
	newMCPEnabled bool
)

var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Scaffold a complete, runnable, themed combined site project",
	Args:  cobra.ExactArgs(1),
	RunE:  runNew,
}

func init() {
	newCmd.Flags().StringVar(&newType, "type", "", "site type: profile, log, or docs (prompted interactively when omitted)")
	newCmd.Flags().StringVar(&newDomain, "domain", "", "the site's domain, used to pre-fill wrangler.jsonc routing (prompted interactively when omitted)")
	newCmd.Flags().StringVar(&newTheme, "theme", "", "theme preset id, must match --type (prompted interactively when omitted)")
	newCmd.Flags().BoolVar(&newMCPEnabled, "mcp-enabled", true, "whether the agent-facing MCP server is mounted on this site (v0.03 FR-002); the sole MCP-related configuration surface")
	rootCmd.AddCommand(newCmd)
}

// runNew implements `biblio-cli new <project-name>`: it fills any omitted
// flags via the guided Bubble Tea flow, validates the resulting inputs, and
// invokes the scaffold generator — surfacing the contract's success/failure
// output and exit codes (contracts/cli-commands.md `biblio-cli new`).
func runNew(cmd *cobra.Command, args []string) error {
	projectName := args[0]
	siteType, domain, themeID := newType, newDomain, newTheme

	if siteType == "" || domain == "" || themeID == "" {
		result, err := tui.RunNewFlow()
		if err != nil {
			return err
		}
		if siteType == "" {
			siteType = result.SiteType
		}
		if domain == "" {
			domain = result.Domain
		}
		if themeID == "" {
			themeID = result.Theme
		}
	}

	destDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("new: resolving working directory: %w", err)
	}

	if err := ValidateNewFlags(projectName, siteType, domain, themeID, destDir); err != nil {
		return err
	}

	projectDir, err := scaffold.Generate(destDir, scaffold.Options{
		ProjectName: projectName,
		SiteType:    siteType,
		Domain:      domain,
		Theme:       themeID,
		MCPEnabled:  newMCPEnabled,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Scaffolded %q (%s, theme %s) at %s\n", projectName, siteType, themeID, projectDir)
	return nil
}

// ValidateNewFlags checks `new`'s inputs against the contract's documented
// failure modes — an unknown --type, a --theme that doesn't match the chosen
// --type's presets, or a project-name colliding with an existing directory —
// each surfaced as a clear, actionable error naming the bad input
// (contracts/cli-commands.md `biblio-cli new` "Output (failure)";
// data-model.md "Site Project" validation rules).
func ValidateNewFlags(projectName, siteType, domain, themeID, destDir string) error {
	if !isValidSiteType(siteType) {
		return fmt.Errorf("new: unknown --type %q (must be one of %v)", siteType, scaffold.SiteTypes)
	}

	validTheme := false
	for _, p := range theme.PresetsForType(siteType) {
		if p.ID == themeID {
			validTheme = true
			break
		}
	}
	if !validTheme {
		return fmt.Errorf("new: --theme %q is not a valid preset for --type %q", themeID, siteType)
	}

	if _, err := os.Stat(filepath.Join(destDir, projectName)); err == nil {
		return fmt.Errorf("new: project directory %q already exists", projectName)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("new: checking project directory %q: %w", projectName, err)
	}

	return nil
}

func isValidSiteType(siteType string) bool {
	for _, t := range scaffold.SiteTypes {
		if t == siteType {
			return true
		}
	}
	return false
}
