// Package config manages application configuration for Flashcards
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds all configuration for the Flashcards application
type Config struct {
	// Database settings
	DatabasePath string

	// Ollama settings
	OllamaURL      string
	OllamaModel    string
	RequestTimeout int // seconds

	// Application settings
	DataDir string
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "." // fallback to current directory
	}

	dataDir := filepath.Join(homeDir, ".flashcards")

	return &Config{
		DatabasePath:   filepath.Join(dataDir, "flashcards.db"),
		OllamaURL:      "http://localhost:11434/api/generate",
		OllamaModel:    "llama3.1",
		RequestTimeout: 300, // 5 minutes
		DataDir:        dataDir,
	}
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	cfg := DefaultConfig()

	// Override with environment variables if present
	if model := os.Getenv("FLASHCARDS_MODEL"); model != "" {
		cfg.OllamaModel = model
	}

	if url := os.Getenv("FLASHCARDS_OLLAMA_URL"); url != "" {
		cfg.OllamaURL = url
	}

	if dataDir := os.Getenv("FLASHCARDS_DATA_DIR"); dataDir != "" {
		cfg.DataDir = dataDir
		cfg.DatabasePath = filepath.Join(dataDir, "flashcards.db")
	}

	return cfg
}

// EnsureDataDir creates the data directory if it doesn't exist
func (c *Config) EnsureDataDir() error {
	return os.MkdirAll(c.DataDir, 0700)
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.OllamaURL == "" {
		return fmt.Errorf("ollama URL cannot be empty")
	}
	if c.OllamaModel == "" {
		return fmt.Errorf("ollama model cannot be empty")
	}
	if c.RequestTimeout <= 0 {
		return fmt.Errorf("request timeout must be positive")
	}
	return nil
}
