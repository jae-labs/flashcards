// Package keys provides centralized key binding constants and helper functions
// for consistent keyboard interaction across all TUI screens.
package keys

// Common key constants - these are the string representations of keys
// as returned by bubbletea's KeyMsg.String() method.
const (
	Esc      = "esc"
	Enter    = "enter"
	Tab      = "tab"
	Space    = " "
	Up       = "up"
	Down     = "down"
	Left     = "left"
	Right    = "right"
	K        = "k"
	J        = "j"
	Q        = "q"
	C        = "c"
	E        = "e"
	D        = "d"
	I        = "i"
	Y        = "y"
	N        = "n"
	R        = "r"
	B        = "b"
	CtrlC    = "ctrl+c"
	PageUp   = "pgup"
	PageDown = "pgdown"

	// Number keys for shortcuts
	One   = "1"
	Three = "3"
	Seven = "7"
	Nine  = "9"
)

// IsQuit checks if the key is a quit key (q or ctrl+c).
func IsQuit(key string) bool {
	return key == Q || key == CtrlC
}

// IsConfirm checks if the key is a confirmation key (y or enter).
func IsConfirm(key string) bool {
	return key == Y || key == Enter
}

// IsCancel checks if the key is a cancel key (n or esc).
func IsCancel(key string) bool {
	return key == N || key == Esc
}

// IsUp checks if the key is an upward navigation key (up or k).
func IsUp(key string) bool {
	return key == Up || key == K
}

// IsDown checks if the key is a downward navigation key (down or j).
func IsDown(key string) bool {
	return key == Down || key == J
}
