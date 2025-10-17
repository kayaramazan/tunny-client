package cli

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hashicorp/yamux"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"nhooyr.io/websocket"

	"github.com/kayaramazan/tunny/internal/common"
)

var (
	token     string
	subdomain string
	tunnelID  string
	devMode   bool
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect [target]",
	Short: "Create a tunnel to your local server",
	Long: `Connect creates a secure tunnel from the Tunny server to your local server.

The target should be in the format host:port (e.g., localhost:3000, 127.0.0.1:8080)

Examples:
  # Connect to localhost:3000 with auto-generated tunnel ID
  tunny connect localhost:3000

  # Connect with a custom tunnel ID
  tunny connect --id my-api localhost:8080

  # Use a specific token and subdomain
  tunny connect --token mytoken --subdomain myapp localhost:3000
`,
	Args: cobra.ExactArgs(1),
	RunE: runConnect,
}

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().StringVarP(&token, "token", "t", "", "Authentication token (overrides config/env)")
	connectCmd.Flags().StringVar(&subdomain, "subdomain", "", "Subdomain for the tunnel (overrides config/env)")
	connectCmd.Flags().StringVarP(&tunnelID, "id", "i", "", "Custom tunnel ID (auto-generated if not provided)")
	connectCmd.Flags().BoolVarP(&devMode, "dev", "d", false, "Enable development mode with verbose logging")
}

func runConnect(cmd *cobra.Command, args []string) error {
	target := args[0]

	// Validate target format
	if !strings.Contains(target, ":") {
		return fmt.Errorf("target must be in format host:port (e.g., localhost:3000)")
	}

	// Load configuration (config file + env vars + defaults)
	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Server URL is hardcoded - users always connect to YOUR server
	serverURL := DefaultServerURL

	// CLI flags override config/env (only if explicitly set)
	if token == "" {
		token = cfg.Token
	}
	if subdomain == "" {
		subdomain = cfg.Subdomain
	}

	log := common.NewLogger(devMode)
	defer log.Sync()

	// Generate tunnel ID if not provided
	if tunnelID == "" {
		tunnelID = generateTunnelID()
	}

	// Parse and prepare server URL
	u, err := url.Parse(serverURL)
	if err != nil {
		return fmt.Errorf("invalid server URL: %w", err)
	}

	q := u.Query()
	q.Set("token", token)
	q.Set("subdomain", subdomain)
	q.Set("tunnel_id", tunnelID)
	u.RawQuery = q.Encode()

	// Connect to server
	ctx := context.Background()
	c, _, err := websocket.Dial(ctx, u.String(), &websocket.DialOptions{
		CompressionMode: websocket.CompressionDisabled,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer c.Close(websocket.StatusNormalClosure, "bye")

	nc := websocket.NetConn(ctx, c, websocket.MessageBinary)

	sess, err := yamux.Client(nc, nil)
	if err != nil {
		return fmt.Errorf("failed to create tunnel session: %w", err)
	}
	defer sess.Close()

	// Display tunnel information
	serverHost := extractServerHost(u)
	publicURL := fmt.Sprintf("http://%s/%s", serverHost, tunnelID)

	printTunnelInfo(publicURL, tunnelID, target)

	log.Info("tunnel established",
		zap.String("tunnel_id", tunnelID),
		zap.String("subdomain", subdomain),
		zap.String("target", target))

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Accept incoming streams
	errCh := make(chan error, 1)
	go func() {
		for {
			stream, err := sess.AcceptStream()
			if err != nil {
				errCh <- err
				return
			}
			log.Info("stream accepted, forwarding to target")
			go handleStream(log, stream, target)
		}
	}()

	// Wait for shutdown signal or error
	select {
	case <-sigCh:
		fmt.Println("\n\n🛑 Shutting down tunnel...")
		return nil
	case err := <-errCh:
		if err != nil && err != io.EOF {
			log.Error("tunnel error", zap.Error(err))
			return err
		}
		return nil
	}
}

func printTunnelInfo(publicURL, tunnelID, target string) {
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║              ✨ Tunny - Tunnel Established ✨             ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("  🌐 Public URL:    \033[1;36m%s\033[0m\n", publicURL)
	fmt.Printf("  🔑 Tunnel ID:     \033[1;33m%s\033[0m\n", tunnelID)
	fmt.Printf("  📡 Forwarding:    %s → %s\n", publicURL, target)
	fmt.Println()
	fmt.Println("─────────────────────────────────────────────────────────────")
	fmt.Printf("  Try it:  \033[1;32mcurl %s/hello\033[0m\n", publicURL)
	fmt.Println("─────────────────────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("  Press Ctrl+C to stop the tunnel")
	fmt.Println()
}

func extractServerHost(u *url.URL) string {
	serverHost := u.Host
	if u.Scheme == "ws" || u.Scheme == "wss" {
		// Remove port if default
		if strings.HasSuffix(serverHost, ":80") && u.Scheme == "ws" {
			serverHost = strings.TrimSuffix(serverHost, ":80")
		} else if strings.HasSuffix(serverHost, ":443") && u.Scheme == "wss" {
			serverHost = strings.TrimSuffix(serverHost, ":443")
		}
	}
	return serverHost
}

func generateTunnelID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func handleStream(log *zap.Logger, stream net.Conn, target string) {
	defer stream.Close()
	log.Info("handleStream started", zap.String("target", target))

	// Read HTTP request from stream (sent by server's proxy)
	headReader, err := common.ReadHTTPRequestHead(stream)
	if err != nil {
		log.Warn("read head failed", zap.Error(err))
		return
	}
	req, err := http.ReadRequest(headReader)
	if err != nil {
		log.Warn("parse request failed", zap.Error(err))
		return
	}
	log.Info("request parsed", zap.String("method", req.Method), zap.String("path", req.URL.Path))

	// IMPORTANT: Read body from stream if present
	// We must read the entire body before forwarding because the stream
	// won't send EOF until the response is read
	if req.ContentLength > 0 || req.TransferEncoding != nil {
		var buf strings.Builder
		if req.ContentLength > 0 {
			_, err := io.CopyN(&buf, stream, req.ContentLength)
			if err != nil {
				log.Warn("read body failed", zap.Error(err))
				return
			}
		} else {
			_, err := io.Copy(&buf, io.LimitReader(stream, 10*1024*1024)) // 10MB max
			if err != nil {
				log.Warn("read chunked body failed", zap.Error(err))
				return
			}
		}
		req.Body = io.NopCloser(strings.NewReader(buf.String()))
		log.Info("body read", zap.Int("size", buf.Len()))
	} else {
		req.Body = http.NoBody
	}

	// Connect to local target
	up, err := net.Dial("tcp", target)
	if err != nil {
		log.Warn("dial local target failed", zap.Error(err), zap.String("target", target))
		resp := &http.Response{
			StatusCode: http.StatusBadGateway,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader("local target unavailable")),
		}
		resp.Header.Set("Content-Type", "text/plain")
		_ = resp.Write(stream)
		return
	}
	defer up.Close()
	log.Info("connected to local target", zap.String("target", target))

	// Forward request to local server
	// req.Write will write headers + body from req.Body (which reads from stream)
	if err := req.Write(up); err != nil {
		log.Warn("forward request failed", zap.Error(err))
		return
	}
	log.Info("request forwarded to local target")

	// Read response from local server and write back to stream
	br := bufio.NewReader(up)
	resp, err := http.ReadResponse(br, req)
	if err != nil {
		log.Warn("read response failed", zap.Error(err))
		return
	}
	defer resp.Body.Close()
	log.Info("response received from local target", zap.Int("status", resp.StatusCode))

	// Write response back to server stream
	if err := resp.Write(stream); err != nil {
		log.Warn("write response failed", zap.Error(err))
		return
	}
	log.Info("response written back to stream")
}
