// Package main contains configuration management for the MOTD client.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	c Config // Global configuration instance
)

// Config holds the application configuration settings.
// It can be configured through environment variables with the "MOTD_" prefix.
type Config struct {
	Host      string `default:"localhost"`              // Server hostname
	Port      int    `default:"4200"`                   // Server port
	TimeoutMs int    `default:"100" split_words:"true"` // Connection timeout in milliseconds
	LogLevel  string `default:"info"`                   // Log level (debug, info, warn, error)
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}
	if c.TimeoutMs <= 0 {
		return fmt.Errorf("timeout must be positive, got %d", c.TimeoutMs)
	}
	return nil
}

// Timeout returns the timeout as a time.Duration.
func (c *Config) Timeout() time.Duration {
	return time.Duration(c.TimeoutMs) * time.Millisecond
}

// init initializes the configuration by reading from environment variables.
// Environment variables should be prefixed with "MOTD_" (e.g., MOTD_HOST, MOTD_PORT, MOTD_TIMEOUT_MS).
func init() {
	err := envconfig.Process("motd", &c)
	if err != nil {
		slog.Error("Failed to process environment configuration", "error", err)
		os.Exit(1)
	}

	if err := c.Validate(); err != nil {
		slog.Error("Invalid configuration", "error", err)
		os.Exit(1)
	}

	// Setup structured logging
	setupLogging()
}

// setupLogging configures structured logging based on the log level.
func setupLogging() {
	var level slog.Level
	switch c.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: false, // Disable source to reduce noise
	}

	handler := slog.NewTextHandler(os.Stderr, opts)
	slog.SetDefault(slog.New(handler))

	slog.Debug("MOTD client initialized",
		"host", c.Host,
		"port", c.Port,
		"timeout", c.Timeout(),
		"log_level", c.LogLevel)
}
