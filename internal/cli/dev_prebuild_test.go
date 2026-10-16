// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package cli_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/monamaret/biblio/internal/cli"
)

// stubDevCommands replaces the wrapped vite/wrangler dev commands for the
// duration of a test.
func stubDevCommands(t *testing.T, vite, wrangler []string) {
	t.Helper()
	origVite, origWrangler := cli.DevViteCommand, cli.DevWranglerCommand
	cli.DevViteCommand, cli.DevWranglerCommand = vite, wrangler
	t.Cleanup(func() { cli.DevViteCommand, cli.DevWranglerCommand = origVite, origWrangler })
}

// `dev` runs the pre-build step (in Dev mode) before starting either
// server (specs/004 FR-035; contracts/cli.md `biblio-cli dev`).
func TestRunDev_RunsPrebuildBeforeServers(t *testing.T) {
	dir := scaffoldedDir(t)
	writeLogContent(t, dir, false)
	stubDevCommands(t,
		[]string{"sh", "-c", "test -f content-index.json && touch vite-saw-index"},
		// Outlives the vite stub: the supervisor stops the other process as
		// soon as either exits, so wrangler must not finish first.
		[]string{"sh", "-c", "sleep 2"},
	)

	var buf strings.Builder
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	cmd.SetContext(context.Background())

	if err := cli.RunDev(cmd, nil); err != nil {
		t.Fatalf("RunDev() error = %v\noutput:\n%s", err, buf.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "vite-saw-index")); err != nil {
		t.Error("vite dev started before content-index.json existed; the pre-build step must run first")
	}
}

// A validation error stops `dev` before either server starts.
func TestRunDev_ValidationErrorStopsBeforeServers(t *testing.T) {
	dir := scaffoldedDir(t)
	writeLogContent(t, dir, true)
	stubDevCommands(t, []string{"sh", "-c", "touch vite-ran"}, []string{"sh", "-c", "touch wrangler-ran"})

	var buf strings.Builder
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	cmd.SetContext(context.Background())

	if err := cli.RunDev(cmd, nil); err == nil {
		t.Fatal("RunDev() error = nil, want the site.yaml validation error")
	}
	if strings.Contains(buf.String(), "Starting local preview") {
		t.Error("dev announced the preview despite a validation error")
	}
	for _, marker := range []string{"vite-ran", "wrangler-ran"} {
		if _, err := os.Stat(filepath.Join(dir, marker)); err == nil {
			t.Errorf("%s exists: a server started despite a validation error", marker)
		}
	}
}
