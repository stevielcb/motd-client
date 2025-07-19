// Package app contains the main application logic and orchestration.
package app

import (
	"fmt"
	"log/slog"

	"github.com/stevielcb/motd-client/internal/config"
	"github.com/stevielcb/motd-client/internal/network"
	"github.com/stevielcb/motd-client/internal/terminal"
)

// App represents the main application.
type App struct {
	cfg       *config.Config
	client    network.ClientInterface
	detector  terminal.DetectorInterface
	formatter *terminal.Formatter
}

// New creates a new application instance.
func New(cfg *config.Config) *App {
	client := network.NewClient(cfg.Host, cfg.Port, cfg.Timeout())
	detector := terminal.NewDetector()

	return &App{
		cfg:      cfg,
		client:   client,
		detector: detector,
	}
}

// Run executes the main application logic.
func (a *App) Run() error {
	// Detect terminal environment
	env, err := a.detector.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect terminal environment: %w", err)
	}

	// Create formatter
	a.formatter = terminal.NewFormatter(env)

	// Connect to server
	conn, err := a.client.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	// Fetch message
	message, err := a.client.FetchMessage(conn)
	if err != nil {
		return fmt.Errorf("failed to fetch message: %w", err)
	}

	// Display message
	a.displayMessage(message)

	return nil
}

// displayMessage formats and displays the MOTD message.
func (a *App) displayMessage(message string) {
	if message == "" {
		slog.Warn("Received empty message from server")
		return
	}

	formattedMessage := a.formatter.Format(message)
	fmt.Print(formattedMessage)

	slog.Debug("Message displayed successfully", "message_length", len(message))
}
