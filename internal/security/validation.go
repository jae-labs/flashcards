// Package security provides validation and security utilities for Flashcards
package security

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ValidateFilePath checks if a file path is safe to access
func ValidateFilePath(path string) error {
	// Clean the path to resolve any relative components
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected: %s", path)
	}

	// Ensure the file exists and is accessible
	info, err := os.Stat(cleanPath)
	if err != nil {
		return fmt.Errorf("cannot access file: %w", err)
	}

	// Check if it's a regular file or directory
	if !info.Mode().IsRegular() && !info.IsDir() {
		return fmt.Errorf("path is not a regular file or directory: %s", cleanPath)
	}

	return nil
}

// ValidateURL checks if a URL is safe and properly formatted
func ValidateURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Only allow http and https schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("only HTTP/HTTPS URLs are allowed, got: %s", parsedURL.Scheme)
	}

	// Restrict to localhost for Ollama API
	host := parsedURL.Hostname()
	if host != "localhost" && host != "127.0.0.1" {
		return fmt.Errorf("only localhost URLs are allowed for security, got: %s", host)
	}

	return nil
}

// SanitizeInput removes potentially dangerous characters from user input
func SanitizeInput(input string) string {
	// Remove null bytes and other control characters
	sanitized := strings.ReplaceAll(input, "\x00", "")
	sanitized = strings.ReplaceAll(sanitized, "\r", " ")
	sanitized = strings.ReplaceAll(sanitized, "\n", " ")

	// Trim whitespace
	return strings.TrimSpace(sanitized)
}
