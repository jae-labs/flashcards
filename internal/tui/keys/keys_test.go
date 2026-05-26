package keys

import "testing"

func TestIsQuit(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"q key", Q, true},
		{"ctrl+c key", CtrlC, true},
		{"esc key", Esc, false},
		{"enter key", Enter, false},
		{"random key", "x", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsQuit(tt.key)
			if result != tt.expected {
				t.Errorf("IsQuit(%q) = %v, want %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestIsConfirm(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"y key", Y, true},
		{"enter key", Enter, true},
		{"n key", N, false},
		{"esc key", Esc, false},
		{"random key", "x", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsConfirm(tt.key)
			if result != tt.expected {
				t.Errorf("IsConfirm(%q) = %v, want %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestIsCancel(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"n key", N, true},
		{"esc key", Esc, true},
		{"y key", Y, false},
		{"enter key", Enter, false},
		{"random key", "x", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCancel(tt.key)
			if result != tt.expected {
				t.Errorf("IsCancel(%q) = %v, want %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestIsUp(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"up arrow", Up, true},
		{"k key", K, true},
		{"down arrow", Down, false},
		{"j key", J, false},
		{"random key", "x", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUp(tt.key)
			if result != tt.expected {
				t.Errorf("IsUp(%q) = %v, want %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestIsDown(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"down arrow", Down, true},
		{"j key", J, true},
		{"up arrow", Up, false},
		{"k key", K, false},
		{"random key", "x", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDown(tt.key)
			if result != tt.expected {
				t.Errorf("IsDown(%q) = %v, want %v", tt.key, result, tt.expected)
			}
		})
	}
}
