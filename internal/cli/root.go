// Package cli wires the biblio-cli Cobra command tree (new, index, dev) and
// the shared Viper configuration loaded from a project's content/site.yaml.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// version is overridden at build time via -ldflags "-X .../cli.version=...".
var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "biblio-cli",
	Short:   "Scaffold, index, and preview a themed biblio site",
	Version: version,
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
