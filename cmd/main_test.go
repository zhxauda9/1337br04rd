package main

import (
	"flag"
	"os"
	"testing"
)

func TestMainFlags(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantPort    int
		wantHelp    bool
		shouldError bool
	}{
		{"default values", []string{}, 8080, false, false},
		{"custom port", []string{"--port", "9090"}, 9090, false, false},
		{"help flag", []string{"--help"}, 8080, true, false},
		{"invalid port", []string{"--port", "8080"}, 0, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags before each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			port = 0
			showHelp = false

			// Save original args and restore when done
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			// Set test args
			os.Args = append([]string{"cmd"}, tt.args...)

			// Capture any panic from flag.Parse()
			defer func() {
				if r := recover(); r != nil {
					if !tt.shouldError {
						t.Errorf("Unexpected error: %v", r)
					}
				}
			}()

			initFlags()

			if tt.shouldError {
				// If we expected an error, we shouldn't check the values
				return
			}

			if port != tt.wantPort {
				t.Errorf("Expected port %d, got %d", tt.wantPort, port)
			}
			if showHelp != tt.wantHelp {
				t.Errorf("Expected help %v, got %v", tt.wantHelp, showHelp)
			}
		})
	}
}
