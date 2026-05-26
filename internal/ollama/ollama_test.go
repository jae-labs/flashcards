package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestParseFlashcards(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []map[string]string
	}{
		{
			name:  "single Q&A",
			input: "Q: What is Go?\nA: A programming language",
			expected: []map[string]string{
				{"question": "What is Go?", "answer": "A programming language"},
			},
		},
		{
			name:  "multiple Q&A",
			input: "Q: What is Go?\nA: A programming language\nQ: What is Python?\nA: Another language",
			expected: []map[string]string{
				{"question": "What is Go?", "answer": "A programming language"},
				{"question": "What is Python?", "answer": "Another language"},
			},
		},
		{
			name:     "no Q&A",
			input:    "Some text",
			expected: []map[string]string{},
		},
		{
			name:     "incomplete Q",
			input:    "Q: What is Go?",
			expected: []map[string]string{},
		},
		{
			name:  "with extra spaces",
			input: "Q:  What is Go?  \nA:  A programming language  ",
			expected: []map[string]string{
				{"question": "What is Go?", "answer": "A programming language"},
			},
		},
		{
			name:  "multiple Q without A",
			input: "Q: Question 1\nQ: Question 2\nA: Answer 2",
			expected: []map[string]string{
				{"question": "Question 2", "answer": "Answer 2"},
			},
		},
		{
			name:  "empty lines",
			input: "Q: What is Go?\n\n\nA: A language\n\n",
			expected: []map[string]string{
				{"question": "What is Go?", "answer": "A language"},
			},
		},
		{
			name:  "lines too short",
			input: "Q:\nA:\nQ: Q\nA: A",
			expected: []map[string]string{
				{"question": "Q", "answer": "A"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseFlashcards(tt.input)
			if err != nil {
				t.Errorf("ParseFlashcards() error = %v", err)
				return
			}
			if len(result) != len(tt.expected) {
				t.Errorf("ParseFlashcards() len = %d, expected %d", len(result), len(tt.expected))
				return
			}
			for i, qa := range result {
				if qa["question"] != tt.expected[i]["question"] || qa["answer"] != tt.expected[i]["answer"] {
					t.Errorf("ParseFlashcards() = %v, expected %v", qa, tt.expected[i])
				}
			}
		})
	}
}

func TestGenerateQA(t *testing.T) {
	// Create a test server that mocks Ollama API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		// Verify content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		// Parse request body
		var req OllamaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		// Verify request fields
		if req.Model == "" {
			t.Error("Model should not be empty")
		}
		if req.Prompt == "" {
			t.Error("Prompt should not be empty")
		}

		// Send streaming response (Ollama format)
		w.Header().Set("Content-Type", "application/json")

		// Send multiple chunks
		chunks := []map[string]interface{}{
			{"response": "Q: What is Go?\nA: ", "done": false},
			{"response": "A programming language", "done": false},
			{"response": "\nQ: What is Python?\nA: ", "done": false},
			{"response": "Another language", "done": true},
		}

		for _, chunk := range chunks {
			if err := json.NewEncoder(w).Encode(chunk); err != nil {
				t.Errorf("Failed to encode response: %v", err)
				return
			}
			w.(http.Flusher).Flush()
		}
	}))
	defer server.Close()

	ctx := context.Background()
	model := "test-model"
	url := server.URL
	prompt := "Generate flashcards"

	result, err := GenerateQA(ctx, model, url, prompt)
	if err != nil {
		t.Fatalf("GenerateQA() error = %v", err)
	}

	expected := "Q: What is Go?\nA: A programming language\nQ: What is Python?\nA: Another language"
	if result != expected {
		t.Errorf("GenerateQA() = %q, expected %q", result, expected)
	}
}

func TestGenerateQAErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "invalid JSON response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("invalid json"))
			},
			wantErr: false, // Should handle gracefully
		},
		{
			name: "empty response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(""))
			},
			wantErr: false, // Should return empty string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			ctx := context.Background()
			result, err := GenerateQA(ctx, "test-model", server.URL, "test prompt")

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateQA() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result == "" && tt.name != "empty response" {
				t.Logf("GenerateQA() returned empty result for %s", tt.name)
			}
		})
	}
}

func TestGenerateQATimeout(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Delay longer than context timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := GenerateQA(ctx, "test-model", server.URL, "test prompt")
	if err == nil {
		t.Error("GenerateQA() should return error on timeout")
	}
}

func TestGenerateQAInvalidURL(t *testing.T) {
	ctx := context.Background()
	_, err := GenerateQA(ctx, "test-model", "http://invalid-url-that-does-not-exist:12345", "test prompt")
	if err == nil {
		t.Error("GenerateQA() should return error for invalid URL")
	}
}
