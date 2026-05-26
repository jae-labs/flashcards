package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.OllamaModel != "llama3.1" {
		t.Errorf("Expected default model 'llama3.1', got '%s'", cfg.OllamaModel)
	}

	if cfg.OllamaURL != "http://localhost:11434/api/generate" {
		t.Errorf("Expected default URL 'http://localhost:11434/api/generate', got '%s'", cfg.OllamaURL)
	}

	if cfg.RequestTimeout != 300 {
		t.Errorf("Expected default timeout 300, got %d", cfg.RequestTimeout)
	}

	if cfg.DataDir == "" {
		t.Error("Expected DataDir to be set")
	}

	if cfg.DatabasePath == "" {
		t.Error("Expected DatabasePath to be set")
	}
}

func TestLoadConfig(t *testing.T) {
	// Test environment variable override
	os.Setenv("FLASHCARDS_MODEL", "test-model")
	os.Setenv("FLASHCARDS_OLLAMA_URL", "http://test:1234")
	os.Setenv("FLASHCARDS_DATA_DIR", "/tmp/test-flashcards")
	defer func() {
		os.Unsetenv("FLASHCARDS_MODEL")
		os.Unsetenv("FLASHCARDS_OLLAMA_URL")
		os.Unsetenv("FLASHCARDS_DATA_DIR")
	}()

	cfg := LoadConfig()

	if cfg.OllamaModel != "test-model" {
		t.Errorf("Expected model 'test-model', got '%s'", cfg.OllamaModel)
	}

	if cfg.OllamaURL != "http://test:1234" {
		t.Errorf("Expected URL 'http://test:1234', got '%s'", cfg.OllamaURL)
	}

	if cfg.DataDir != "/tmp/test-flashcards" {
		t.Errorf("Expected DataDir '/tmp/test-flashcards', got '%s'", cfg.DataDir)
	}

	expectedDBPath := filepath.Join("/tmp/test-flashcards", "flashcards.db")
	if cfg.DatabasePath != expectedDBPath {
		t.Errorf("Expected DatabasePath '%s', got '%s'", expectedDBPath, cfg.DatabasePath)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: Config{
				OllamaURL:      "http://localhost:11434/api/generate",
				OllamaModel:    "llama3.1",
				RequestTimeout: 300,
			},
			wantErr: false,
		},
		{
			name: "empty URL",
			cfg: Config{
				OllamaURL:      "",
				OllamaModel:    "llama3.1",
				RequestTimeout: 300,
			},
			wantErr: true,
		},
		{
			name: "empty model",
			cfg: Config{
				OllamaURL:      "http://localhost:11434/api/generate",
				OllamaModel:    "",
				RequestTimeout: 300,
			},
			wantErr: true,
		},
		{
			name: "negative timeout",
			cfg: Config{
				OllamaURL:      "http://localhost:11434/api/generate",
				OllamaModel:    "llama3.1",
				RequestTimeout: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnsureDataDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	cfg := &Config{
		DataDir: filepath.Join(tempDir, "test-flashcards"),
	}

	err := cfg.EnsureDataDir()
	if err != nil {
		t.Errorf("EnsureDataDir() error = %v", err)
	}

	// Check if directory was created
	info, err := os.Stat(cfg.DataDir)
	if err != nil {
		t.Errorf("Data directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("DataDir is not a directory")
	}
}

func TestDefaultConfigErrorCase(t *testing.T) {
	// Test DefaultConfig when UserHomeDir fails
	// This is hard to test directly, but we can verify it handles the fallback
	// The function should still work even if UserHomeDir fails
	// We can't easily mock os.UserHomeDir, so we test that the function
	// doesn't panic and returns a valid config
	// Note: DefaultConfig() always returns a non-nil pointer, so we just verify it works
	cfg := DefaultConfig()
	if cfg.DataDir == "" {
		t.Error("DefaultConfig() should set DataDir even on error")
	}
}
