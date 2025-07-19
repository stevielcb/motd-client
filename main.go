// Package main implements a client for fetching and displaying Message of the Day (MOTD)
// from a TCP server. It handles various terminal environments and formatting requirements.
package main

import (
	"log/slog"
	"os"

	"github.com/stevielcb/motd-client/internal/app"
	"github.com/stevielcb/motd-client/internal/config"
	"github.com/stevielcb/motd-client/internal/logger"
)

// main is the entry point of the application.
func main() {
	if err := run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}

// run contains the main application logic with proper error handling.
func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Setup logging
	logger.Setup(cfg.LogLevel)

	slog.Debug("MOTD client initialized",
		"host", cfg.Host,
		"port", cfg.Port,
		"timeout", cfg.Timeout(),
		"log_level", cfg.LogLevel)

	// Create and run application
	application := app.New(cfg)
	return application.Run()
}
