package config

import (
	"os"
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Host:      "localhost",
				Port:      8080,
				TimeoutMs: 100,
				LogLevel:  "info",
			},
			wantErr: false,
		},
		{
			name: "empty host",
			config: Config{
				Host:      "",
				Port:      8080,
				TimeoutMs: 100,
				LogLevel:  "info",
			},
			wantErr: true,
		},
		{
			name: "invalid port - zero",
			config: Config{
				Host:      "localhost",
				Port:      0,
				TimeoutMs: 100,
				LogLevel:  "info",
			},
			wantErr: true,
		},
		{
			name: "invalid port - too high",
			config: Config{
				Host:      "localhost",
				Port:      70000,
				TimeoutMs: 100,
				LogLevel:  "info",
			},
			wantErr: true,
		},
		{
			name: "invalid timeout - zero",
			config: Config{
				Host:      "localhost",
				Port:      8080,
				TimeoutMs: 0,
				LogLevel:  "info",
			},
			wantErr: true,
		},
		{
			name: "invalid timeout - negative",
			config: Config{
				Host:      "localhost",
				Port:      8080,
				TimeoutMs: -100,
				LogLevel:  "info",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_Timeout(t *testing.T) {
	config := Config{
		TimeoutMs: 1500,
	}

	expected := 1500 * time.Millisecond
	result := config.Timeout()

	if result != expected {
		t.Errorf("Config.Timeout() = %v, want %v", result, expected)
	}
}

func TestLoad(t *testing.T) {
	// Save original environment variables
	originalEnv := make(map[string]string)
	envVars := []string{"MOTD_HOST", "MOTD_PORT", "MOTD_TIMEOUT_MS", "MOTD_LOG_LEVEL"}

	for _, key := range envVars {
		if val := os.Getenv(key); val != "" {
			originalEnv[key] = val
		}
	}

	// Clean up after test
	defer func() {
		for key := range originalEnv {
			os.Setenv(key, originalEnv[key])
		}
		for _, key := range envVars {
			if _, exists := originalEnv[key]; !exists {
				os.Unsetenv(key)
			}
		}
	}()

	tests := []struct {
		name        string
		envVars     map[string]string
		wantErr     bool
		checkResult func(*Config) bool
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			wantErr: false,
			checkResult: func(cfg *Config) bool {
				return cfg.Host == "localhost" && cfg.Port == 4200 &&
					cfg.TimeoutMs == 100 && cfg.LogLevel == "info"
			},
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"MOTD_HOST":       "example.com",
				"MOTD_PORT":       "8080",
				"MOTD_TIMEOUT_MS": "5000",
				"MOTD_LOGLEVEL":   "debug",
			},
			wantErr: false,
			checkResult: func(cfg *Config) bool {
				return cfg.Host == "example.com" && cfg.Port == 8080 &&
					cfg.TimeoutMs == 5000 && cfg.LogLevel == "debug"
			},
		},
		{
			name: "invalid port",
			envVars: map[string]string{
				"MOTD_PORT": "99999",
			},
			wantErr: true,
		},
		{
			name: "invalid timeout",
			envVars: map[string]string{
				"MOTD_TIMEOUT_MS": "-100",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables for test
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Load() unexpected error: %v", err)
				return
			}

			if tt.checkResult != nil && !tt.checkResult(cfg) {
				t.Errorf("Load() result validation failed")
			}
		})
	}
}
