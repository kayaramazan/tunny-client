// Package cli provides the CLI interface for Tunny.
package cli

import (
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tunny",
	Short: "Tunny - Expose your local server to the internet",
	Long: `Tunny is a simple ngrok-like reverse tunnel service.
Expose your local HTTP server to the internet via a secure tunnel.

Example:
  tunny connect localhost:3000
  tunny connect --id my-api localhost:8080
  tunny list
`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetVersionTemplate(`Tunny version {{.Version}}
`)
}
