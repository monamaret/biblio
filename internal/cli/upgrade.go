// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/monamaret/biblio/internal/migrate"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade <version>",
	Short: "Bring the current project's scaffold-tracked parts up to a specific scaffold version",
	Args:  cobra.ExactArgs(1),
	RunE:  RunUpgrade,
}

func init() {
	// Shares the same --resolve <path>=keep-author|use-scaffold flag as
	// `update` — both commands share the same internal/migrate machinery.
	upgradeCmd.Flags().StringArray("resolve", nil, "resolve a prior Update Conflict: --resolve path=keep-author or --resolve path=use-scaffold (repeatable)")
	rootCmd.AddCommand(upgradeCmd)
}

// Version returns the biblio-cli version this binary was built with
// (root.go's version, overridden at build time via -ldflags), exported so
// `upgrade`'s version check and its tests have a single source of truth
// for "what version does this binary actually embed."
func Version() string {
	return version
}

// RunUpgrade implements `biblio-cli upgrade <version>`: it shares update's
// classify/advance/preserve/surface machinery (internal/migrate), but
// first confirms the requested target version is the one this binary
// actually has embedded. This binary embeds exactly one scaffold-template
// version at a time — there is no historical multi-version template
// store — so "upgrading" to any other version isn't something this binary
// can carry out (contracts/update-and-escalation.md `biblio-cli upgrade`).
func RunUpgrade(cmd *cobra.Command, args []string) error {
	target := args[0]
	if target != version {
		return fmt.Errorf("upgrade: this biblio-cli binary embeds scaffold-template version %q, not %q — install biblio-cli %s to upgrade to it", version, target, target)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("upgrade: resolving working directory: %w", err)
	}

	siteType, err := readSiteType(cwd)
	if err != nil {
		return err
	}

	resolveFlags, err := cmd.Flags().GetStringArray("resolve")
	if err != nil {
		return fmt.Errorf("upgrade: reading --resolve: %w", err)
	}
	resolutions, err := parseResolveFlags(resolveFlags)
	if err != nil {
		return err
	}

	run, err := migrate.RunUpdate(cwd, siteType, resolutions)
	if err != nil {
		if err == migrate.ErrUnresolvedConflicts {
			reportConflicts(cmd.ErrOrStderr(), run)
			return fmt.Errorf("upgrade: %d conflict(s) unresolved — the project is unchanged and safe to retry; resolve each with --resolve <path>=keep-author or --resolve <path>=use-scaffold, then re-run `biblio-cli upgrade %s`", len(run.Conflicts), target)
		}
		return fmt.Errorf("upgrade: %w", err)
	}

	reportSuccess(cmd.OutOrStdout(), run)
	return nil
}
