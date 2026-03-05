package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// OllamaRequest represents the request to the Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// GenerateQA sends a prompt to the Ollama API and returns the response
func GenerateQA(ctx context.Context, model, url, prompt string) (string, error) {
	body, err := json.Marshal(OllamaRequest{Model: model, Prompt: prompt})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response bytes.Buffer
	dec := json.NewDecoder(resp.Body)
	for {
		var chunk map[string]interface{}
		if err := dec.Decode(&chunk); err != nil {
			break
		}
		if r, ok := chunk["response"].(string); ok {
			response.WriteString(r)
		}
		if done, ok := chunk["done"].(bool); ok && done {
			break
		}
	}
	return response.String(), nil
}

// ParseFlashcards parses the Ollama response and returns a list of questions and answers
func ParseFlashcards(response string) ([]map[string]string, error) {
	var qas []map[string]string
	lines := strings.Split(response, "\n")
	var q, a string
	for _, line := range lines {
		l := strings.TrimSpace(line)
		l = strings.Trim(l, "*: ")
		if len(l) < 2 {
			continue
		}
		switch l[:2] {
		case "Q:":
			q = strings.TrimSpace(strings.Trim(l[2:], "*: "))
		case "A:":
			a = strings.TrimSpace(strings.Trim(l[2:], "*: "))
			if q != "" && a != "" {
				qas = append(qas, map[string]string{"question": q, "answer": a})
				q = "" // Reset for next Q/A pair
			}
		}
	}
	return qas, nil
}
