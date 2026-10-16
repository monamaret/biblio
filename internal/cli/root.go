// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// Package cli wires the biblio-cli Cobra command tree (new, index, dev) and
// the shared Viper configuration loaded from a project's content/site.yaml.
package cli

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// version is overridden at build time via -ldflags "-X .../cli.version=...",
// which is how goreleaser-built release binaries report their tag (see
// .goreleaser.yaml). -X only works against a plain string-literal
// initializer like this one, so resolveVersion (which applies the
// build-info fallback) is deliberately kept out of this initializer and
// applied separately below. version stays "dev" for a plain `go
// build`/`go test` in this repo, and also for `go install
// .../biblio-cli@latest` — ldflags aren't applied on that path, so
// resolveVersion's fallback is what actually fixes `-v` for that install
// method (issue #193; this is also how our own deploy.yml installs
// biblio-cli in CI).
var version = "dev"

// resolveVersion returns ldflagsVersion unchanged unless it's still the
// unset default, in which case it falls back to the module version Go
// embeds automatically when the binary was built via `go install
// module@version`/`@latest` (retrievable via debug.ReadBuildInfo() — no
// ldflags required for that install path; this is also how a plain `go
// build` of this repo's own checkout reports a VCS-derived pseudo-version,
// e.g. "v0.0.0-<timestamp>-<commit>[+dirty]", which is strictly more
// useful than a flat "dev" since it pins the exact commit). The one case
// excluded is "(devel)"/empty — what `go test`'s own build reports — so
// that case still falls back to "dev" rather than surfacing that raw Go
// internal.
func resolveVersion(ldflagsVersion string) string {
	if ldflagsVersion != "dev" {
		return ldflagsVersion
	}
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return ldflagsVersion
	}
	return info.Main.Version
}

var rootCmd = &cobra.Command{
	Use:     "biblio-cli",
	Short:   "Scaffold, index, and preview a themed biblio site",
	Version: resolveVersion(version),
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig loads the current project's content/site.yaml via Viper, if present.
// A missing config is not an error here — `new` is the command that creates one.
func initConfig() {
	viper.SetConfigName("site")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./content")
	_ = viper.ReadInConfig()
}

// Execute runs the root command, exiting non-zero on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
