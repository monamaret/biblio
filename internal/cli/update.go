package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/monamaret/biblio/internal/content"
	"github.com/monamaret/biblio/internal/migrate"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Bring the current project's scaffold-tracked parts up to date",
	Args:  cobra.NoArgs,
	RunE:  RunUpdate,
}

func init() {
	// Repeated --resolve path=resolution pairs (migrate.ResolveKeepAuthor /
	// migrate.ResolveUseScaffold) let the author settle a prior run's
	// Update Conflicts non-interactively.
	updateCmd.Flags().StringArray("resolve", nil, "resolve a prior Update Conflict: --resolve path=keep-author or --resolve path=use-scaffold (repeatable)")
	rootCmd.AddCommand(updateCmd)
}

// RunUpdate implements `biblio-cli update`: it classifies every region of
// the current project against the scaffold template it was generated
// from, advances scaffold-tracked regions, preserves author-customized
// regions, and surfaces any ambiguous region as an Update Conflict —
// exiting non-zero with recovery guidance until every conflict is resolved
// via --resolve (contracts/update-and-escalation.md `biblio-cli update`).
func RunUpdate(cmd *cobra.Command, _ []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("update: resolving working directory: %w", err)
	}

	siteType, err := readSiteType(cwd)
	if err != nil {
		return err
	}

	resolveFlags, err := cmd.Flags().GetStringArray("resolve")
	if err != nil {
		return fmt.Errorf("update: reading --resolve: %w", err)
	}
	resolutions, err := parseResolveFlags(resolveFlags)
	if err != nil {
		return err
	}

	run, err := migrate.RunUpdate(cwd, siteType, resolutions)
	if err != nil {
		if err == migrate.ErrUnresolvedConflicts {
			reportConflicts(cmd.ErrOrStderr(), run)
			return fmt.Errorf("update: %d conflict(s) unresolved — the project is unchanged and safe to retry; resolve each with --resolve <path>=keep-author or --resolve <path>=use-scaffold, then re-run `biblio-cli update`", len(run.Conflicts))
		}
		return fmt.Errorf("update: %w", err)
	}

	reportSuccess(cmd.OutOrStdout(), run)
	return nil
}

// readSiteType reads the declared site type from cwd/content/site.yaml,
// failing early and clearly if cwd isn't a scaffolded project.
func readSiteType(cwd string) (string, error) {
	path := filepath.Join(cwd, "content", "site.yaml")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("update: %s not found — run `biblio-cli new` first to scaffold a project", path)
	}
	if err != nil {
		return "", fmt.Errorf("update: reading %s: %w", path, err)
	}
	cfg, err := content.ReadSiteConfig(data)
	if err != nil {
		return "", fmt.Errorf("update: %w", err)
	}
	return cfg.SiteType, nil
}

// parseResolveFlags turns repeated "path=resolution" --resolve flags into
// the map[string]string migrate.RunUpdate expects.
func parseResolveFlags(flags []string) (map[string]string, error) {
	if len(flags) == 0 {
		return nil, nil
	}
	resolutions := make(map[string]string, len(flags))
	for _, f := range flags {
		path, resolution, ok := strings.Cut(f, "=")
		if !ok || path == "" || resolution == "" {
			return nil, fmt.Errorf("update: invalid --resolve %q — want <path>=keep-author or <path>=use-scaffold", f)
		}
		if resolution != migrate.ResolveKeepAuthor && resolution != migrate.ResolveUseScaffold {
			return nil, fmt.Errorf("update: invalid --resolve %q — resolution must be %q or %q", f, migrate.ResolveKeepAuthor, migrate.ResolveUseScaffold)
		}
		resolutions[path] = resolution
	}
	return resolutions, nil
}

// reportConflicts prints the contract's "Output (failure / interruption)":
// a clear statement of the site's current, consistent, recoverable state
// and how to resume — never a half-migrated site with no documented path
// forward (contracts/update-and-escalation.md).
func reportConflicts(w io.Writer, run *migrate.UpdateRun) {
	fmt.Fprintf(w, "update: %d region(s) need your decision before anything is changed:\n", len(run.Conflicts))
	for _, c := range run.Conflicts {
		if c.Resolution != "" {
			continue
		}
		fmt.Fprintf(w, "  • %s — scaffold proposes a change, but you've customized it\n", c.Region)
	}
	fmt.Fprintln(w, "Nothing has been written yet; the project is exactly as it was before this run.")
}

// reportSuccess prints the contract's "Output (success)": what was
// preserved, what was advanced, and how each conflict (if any) was
// resolved (contracts/update-and-escalation.md; data-model.md "Update
// Run").
func reportSuccess(w io.Writer, run *migrate.UpdateRun) {
	fmt.Fprintf(w, "Update complete: %d region(s) advanced, %d region(s) preserved.\n", len(run.Advanced), len(run.Preserved))
	for _, c := range run.Conflicts {
		fmt.Fprintf(w, "  • %s — resolved: %s\n", c.Region, c.Resolution)
	}
}
