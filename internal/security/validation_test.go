package security

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateFilePath(t *testing.T) {
	// Create a temporary file for testing
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.md")
	err := os.WriteFile(testFile, []byte("# Test"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid file",
			path:    testFile,
			wantErr: false,
		},
		{
			name:    "valid directory",
			path:    tempDir,
			wantErr: false,
		},
		{
			name:    "path traversal",
			path:    filepath.Join(tempDir, "..", "etc", "passwd"),
			wantErr: true,
		},
		{
			name:    "non-existent file",
			path:    filepath.Join(tempDir, "nonexistent.md"),
			wantErr: true,
		},
		{
			name:    "relative path with traversal",
			path:    "../test.md",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFilePathSpecialCases(t *testing.T) {
	tempDir := t.TempDir()

	// Create a symlink (if supported)
	symlinkPath := filepath.Join(tempDir, "symlink")
	testFile := filepath.Join(tempDir, "target.md")
	err := os.WriteFile(testFile, []byte("# Test"), 0600)
	if err == nil {
		// Try to create symlink (may fail on Windows)
		_ = os.Symlink(testFile, symlinkPath)
		// Test symlink - should work if symlink was created
		if _, err := os.Lstat(symlinkPath); err == nil {
			err := ValidateFilePath(symlinkPath)
			// Symlink to a file should be valid (it's a regular file)
			if err != nil {
				t.Logf("Symlink validation: %v (may be expected)", err)
			}
		}
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr bool
	}{
		{
			name:    "valid localhost HTTP",
			rawURL:  "http://localhost:11434/api/generate",
			wantErr: false,
		},
		{
			name:    "valid 127.0.0.1 HTTP",
			rawURL:  "http://127.0.0.1:11434/api/generate",
			wantErr: false,
		},
		{
			name:    "valid localhost HTTPS",
			rawURL:  "https://localhost:11434/api/generate",
			wantErr: false,
		},
		{
			name:    "invalid scheme",
			rawURL:  "ftp://localhost:11434/api/generate",
			wantErr: true,
		},
		{
			name:    "external host",
			rawURL:  "http://example.com:11434/api/generate",
			wantErr: true,
		},
		{
			name:    "malformed URL",
			rawURL:  "not-a-url",
			wantErr: true,
		},
		{
			name:    "missing scheme",
			rawURL:  "localhost:11434/api/generate",
			wantErr: true,
		},
		{
			name:    "empty URL",
			rawURL:  "",
			wantErr: true,
		},
		{
			name:    "localhost without port",
			rawURL:  "http://localhost/api/generate",
			wantErr: false,
		},
		{
			name:    "127.0.0.1 without port",
			rawURL:  "http://127.0.0.1/api/generate",
			wantErr: false,
		},
		{
			name:    "invalid hostname",
			rawURL:  "http://192.168.1.1:11434/api/generate",
			wantErr: true,
		},
		{
			name:    "localhost with subdomain",
			rawURL:  "http://sub.localhost:11434/api/generate",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "text with null bytes",
			input:    "Hello\x00World",
			expected: "HelloWorld",
		},
		{
			name:     "text with carriage return",
			input:    "Hello\rWorld",
			expected: "Hello World",
		},
		{
			name:     "text with newline",
			input:    "Hello\nWorld",
			expected: "Hello World",
		},
		{
			name:     "text with extra whitespace",
			input:    "  Hello World  ",
			expected: "Hello World",
		},
		{
			name:     "mixed control characters",
			input:    "  Hello\x00\r\nWorld  ",
			expected: "Hello  World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeInput(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeInput() = %q, expected %q", result, tt.expected)
			}
		})
	}
}
