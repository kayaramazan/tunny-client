package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Tunny configuration",
	Long: `Initialize creates a configuration file at ~/.tunny/config.json with default or custom settings.

Note: Server URL is fixed - you'll connect to the hosted tunnel service.

Examples:
  # Create config with default values
  tunny init

  # Create config with custom token
  tunny init --token your-token --subdomain myapp
`,
	RunE: runInit,
}

var (
	initToken     string
	initSubdomain string
)

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&initToken, "token", "", "Authentication token")
	initCmd.Flags().StringVar(&initSubdomain, "subdomain", "", "Preferred subdomain")
}

func runInit(cmd *cobra.Command, args []string) error {
	cfg := &Config{
		Token:     initToken,
		Subdomain: initSubdomain,
	}

	if err := SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println()
	fmt.Println("✅ Configuration saved!")
	fmt.Println()
	fmt.Println("📄 Location: ~/.tunny/config.json")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Printf("  Server URL:  %s (fixed)\n", DefaultServerURL)
	fmt.Printf("  Token:       %s\n", cfg.Token)
	fmt.Printf("  Subdomain:   %s\n", cfg.Subdomain)
	fmt.Println()
	fmt.Println("💡 You can override token/subdomain with:")
	fmt.Println("   - CLI flags: tunny connect --token <token> localhost:3000")
	fmt.Println("   - Environment variables: TUNNY_TOKEN, TUNNY_SUBDOMAIN")
	fmt.Println()

	return nil
}
