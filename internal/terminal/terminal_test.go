package terminal

import (
	"os"
	"testing"
)

func TestDetector_Detect(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		wantErr     bool
		checkResult func(*Environment) bool
	}{
		{
			name: "no TERM variable",
			envVars: map[string]string{
				"TERM": "", // Empty TERM variable
			},
			wantErr: true,
		},
		{
			name: "iTerm2 terminal",
			envVars: map[string]string{
				"TERM":         "xterm-256color",
				"TERM_PROGRAM": "iTerm.app",
			},
			wantErr: false,
			checkResult: func(env *Environment) bool {
				return env.IsITerm2 && !env.IsSSH && !env.IsTmux &&
					env.StartSeq == "\033]" && env.EndSeq == "\a"
			},
		},
		{
			name: "SSH client",
			envVars: map[string]string{
				"TERM":       "xterm-256color",
				"SSH_CLIENT": "192.168.1.1 12345 22",
			},
			wantErr: false,
			checkResult: func(env *Environment) bool {
				return !env.IsITerm2 && env.IsSSH && !env.IsTmux &&
					env.StartSeq == "\033]" && env.EndSeq == "\a"
			},
		},
		{
			name: "tmux session",
			envVars: map[string]string{
				"TERM": "screen-256color",
			},
			wantErr: false,
			checkResult: func(env *Environment) bool {
				return !env.IsITerm2 && !env.IsSSH && env.IsTmux &&
					env.StartSeq == "\033Ptmux;\033\033]" && env.EndSeq == "\a\033\\"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables for test
			for key, value := range tt.envVars {
				if value == "" {
					os.Unsetenv(key)
				} else {
					os.Setenv(key, value)
				}
			}
			defer func() {
				// Clean up environment variables
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			detector := NewDetector()
			env, err := detector.Detect()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Detect() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Detect() unexpected error: %v", err)
				return
			}

			if tt.checkResult != nil && !tt.checkResult(env) {
				t.Errorf("Detect() result validation failed")
			}
		})
	}
}

func TestFormatter_Format(t *testing.T) {
	tests := []struct {
		name     string
		env      *Environment
		message  string
		expected string
	}{
		{
			name: "empty message",
			env: &Environment{
				StartSeq: "\033]",
				EndSeq:   "\a",
			},
			message:  "",
			expected: "",
		},
		{
			name: "standard terminal",
			env: &Environment{
				StartSeq: "\033]",
				EndSeq:   "\a",
			},
			message:  "Hello, World!",
			expected: "\033]Hello, World!\a",
		},
		{
			name: "tmux terminal",
			env: &Environment{
				StartSeq: "\033Ptmux;\033\033]",
				EndSeq:   "\a\033\\",
			},
			message:  "Test message",
			expected: "\033Ptmux;\033\033]Test message\a\033\\",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(tt.env)
			result := formatter.Format(tt.message)

			if result != tt.expected {
				t.Errorf("Format() = %q, want %q", result, tt.expected)
			}
		})
	}
}
