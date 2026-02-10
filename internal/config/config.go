package config

import (
	"os"
	"path/filepath"
)

// Config for the application
type Config struct {
	MaxMemoryClips int    `json:"max_memory_clips"` // Default: 50
	DBPath         string `json:"db_path"`          // SQLite path
	SocketPath     string `json:"socket_path"`      // Unix socket path
}

// Default configuration values
const (
	DefaultMaxMemoryClips = 50
	DefaultSocketPath     = "/tmp/clipnest.sock"
)

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	homeDir, _ := os.UserHomeDir()
	return Config{
		MaxMemoryClips: DefaultMaxMemoryClips,
		DBPath:         filepath.Join(homeDir, "Library", "Application Support", "ClipNest", "clipnest.db"),
		SocketPath:     DefaultSocketPath,
	}
}

// Load loads configuration (expand environment variables)
func Load(path string) (Config, error) {
	// For now, use defaults
	// TODO: Load from JSON/YAML file
	return DefaultConfig(), nil
}

// GetDBPath returns the database path
func GetDBPath() string {
	cfg := DefaultConfig()
	return cfg.DBPath
}

// GetSocketPath returns the socket path
func GetSocketPath() string {
	cfg := DefaultConfig()
	return cfg.SocketPath
}

// EnsureDirectories creates necessary directories
func EnsureDirectories() error {
	cfg := DefaultConfig()
	return os.MkdirAll(filepath.Dir(cfg.DBPath), 0755)
}
