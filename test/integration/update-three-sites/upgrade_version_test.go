package update_three_sites_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/monamaret/biblio/internal/cli"
)

// upgradeCmdForTest builds a bare *cobra.Command carrying the same
// --resolve flag the real upgradeCmd registers, chdir'd into projectDir.
func upgradeCmdForTest(t *testing.T, projectDir string, out *bytes.Buffer) *cobra.Command {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: getwd: %v", err)
	}
	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("setup: chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	cmd := &cobra.Command{}
	cmd.SetOut(out)
	cmd.SetContext(context.Background())
	cmd.Flags().StringArray("resolve", nil, "")
	return cmd
}

// `biblio-cli upgrade <version>` must succeed, with the same advance/
// preserve contract as `update`, across all three real site types when the
// target version matches this binary's own embedded scaffold version
// (FR-003, Scenario 2) — this binary embeds exactly one scaffold-template
// version, so "upgrading" to it is the only version this command can
// actually carry out.
func TestRunUpgrade_MatchingVersion_SucceedsAcrossAllThreeSiteTypes(t *testing.T) {
	for _, siteType := range []string{"profile", "log", "docs"} {
		t.Run(siteType, func(t *testing.T) {
			projectDir := scaffoldSite(t, siteType, "my-"+siteType+"-site")

			var out bytes.Buffer
			cmd := upgradeCmdForTest(t, projectDir, &out)

			if err := cli.RunUpgrade(cmd, []string{cli.Version()}); err != nil {
				t.Fatalf("RunUpgrade(%q): %v", cli.Version(), err)
			}
		})
	}
}

// Targeting a version this binary doesn't actually embed must fail clearly
// rather than silently pretending to upgrade to a version it can't
// produce — this binary has no historical multi-version template storage.
func TestRunUpgrade_MismatchedVersion_FailsClearly(t *testing.T) {
	projectDir := scaffoldSite(t, "profile", "my-profile-site")

	var out bytes.Buffer
	cmd := upgradeCmdForTest(t, projectDir, &out)

	err := cli.RunUpgrade(cmd, []string{"v99.99.99-does-not-exist"})
	if err == nil {
		t.Fatal("RunUpgrade() with a version this binary doesn't embed: want error, got nil")
	}
}
