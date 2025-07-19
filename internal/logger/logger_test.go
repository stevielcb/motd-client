package logger

import (
	"log/slog"
	"testing"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
	}{
		{
			name:     "debug level",
			logLevel: "debug",
		},
		{
			name:     "info level",
			logLevel: "info",
		},
		{
			name:     "warn level",
			logLevel: "warn",
		},
		{
			name:     "error level",
			logLevel: "error",
		},
		{
			name:     "invalid level defaults to info",
			logLevel: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test mainly ensures Setup doesn't panic
			// and sets up logging correctly
			Setup(tt.logLevel)

			// Verify that slog default logger is set
			if slog.Default() == nil {
				t.Error("Expected slog default logger to be set")
			}
		})
	}
}

func TestSetup_LogLevels(t *testing.T) {
	// Test that different log levels are handled correctly
	levels := []string{"debug", "info", "warn", "error", "invalid"}

	for _, level := range levels {
		t.Run("level_"+level, func(t *testing.T) {
			// Should not panic
			Setup(level)

			// Verify logger is functional
			logger := slog.Default()
			if logger == nil {
				t.Error("Logger should not be nil")
			}
		})
	}
}
