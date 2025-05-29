// Package main contains configuration management for the MOTD client.
package main

import (
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
}

// init initializes the configuration by reading from environment variables.
// Environment variables should be prefixed with "MOTD_" (e.g., MOTD_HOST, MOTD_PORT, MOTD_TIMEOUT_MS).
func init() {
	err := envconfig.Process("motd", &c)
	if err != nil {
		panic(err)
	}
}
