package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List active tunnels",
	Long: `List shows all currently active tunnels on the server.

Examples:
  # List tunnels in table format
  tunny list

  # List tunnels in JSON format
  tunny list --json
`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
}

type TunnelInfo struct {
	TunnelID   string    `json:"tunnel_id"`
	Subdomain  string    `json:"subdomain"`
	Connected  bool      `json:"connected"`
	NumStreams int       `json:"num_streams"`
	AttachedAt time.Time `json:"attached_at"`
}

type TunnelsResponse struct {
	Count   int          `json:"count"`
	Tunnels []TunnelInfo `json:"tunnels"`
}

func runList(cmd *cobra.Command, args []string) error {
	// Use the constant server URL
	wsURL := DefaultServerURL

	// Convert WebSocket URL to HTTP URL
	httpURL := strings.Replace(wsURL, "ws://", "http://", 1)
	httpURL = strings.Replace(httpURL, "wss://", "https://", 1)
	serverURL := strings.TrimSuffix(httpURL, "/ws")

	// Fetch tunnels from server
	resp, err := http.Get(serverURL + "/tunnels")
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned error: %s - %s", resp.Status, string(body))
	}

	var tunnelsResp TunnelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tunnelsResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Output
	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(tunnelsResp)
	}

	printTunnelsTable(tunnelsResp, serverURL)
	return nil
}

func printTunnelsTable(resp TunnelsResponse, serverURL string) {
	if resp.Count == 0 {
		fmt.Println("\n❌ No active tunnels found.\n")
		return
	}

	fmt.Printf("\n🌐 Active Tunnels (%d)\n\n", resp.Count)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "TUNNEL ID\tSUBDOMAIN\tSTATUS\tSTREAMS\tUPTIME\tPUBLIC URL")
	fmt.Fprintln(w, "─────────────────\t──────────\t────────\t────────\t────────\t─────────────────────────────")

	for _, tunnel := range resp.Tunnels {
		status := "🟢 online"
		if !tunnel.Connected {
			status = "🔴 offline"
		}

		uptime := time.Since(tunnel.AttachedAt).Round(time.Second)
		publicURL := fmt.Sprintf("%s/%s", serverURL, tunnel.TunnelID)

		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n",
			tunnel.TunnelID,
			tunnel.Subdomain,
			status,
			tunnel.NumStreams,
			formatDuration(uptime),
			publicURL,
		)
	}

	w.Flush()
	fmt.Println()
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}
