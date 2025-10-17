package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// DefaultServerURL is the production tunnel server URL
// This is hardcoded - users always connect to YOUR hosted server
// Updated to the new Cloud Run URL
const DefaultServerURL = "wss://tunny-server-251376301627.us-central1.run.app/ws"

// Config represents the Tunny configuration
// Note: ServerURL is no longer in config - it's a constant
type Config struct {
	Token     string `json:"token"`
	Subdomain string `json:"subdomain"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Token:     "",
		Subdomain: "",
	}
}

// LoadConfig loads configuration from file, environment variables, and defaults
// Priority: CLI flags > Environment variables > Config file > Defaults
func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	// Try to load from config file
	if err := loadConfigFile(cfg); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// Override with environment variables
	loadFromEnv(cfg)

	return cfg, nil
}

// loadConfigFile loads configuration from ~/.tunny/config.json
func loadConfigFile(cfg *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".tunny", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, cfg)
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(cfg *Config) {
	// ServerURL is no longer configurable - it's a constant
	if token := os.Getenv("TUNNY_TOKEN"); token != "" {
		cfg.Token = token
	}
	if subdomain := os.Getenv("TUNNY_SUBDOMAIN"); subdomain != "" {
		cfg.Subdomain = subdomain
	}
}

// SaveConfig saves the configuration to ~/.tunny/config.json
func SaveConfig(cfg *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".tunny")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
