// Package config handles application configuration management.
package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
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

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	var cfg Config

	err := envconfig.Process("motd", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to process environment configuration: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}
