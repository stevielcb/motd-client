// Package terminal handles terminal environment detection and message formatting.
package terminal

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Error definitions
var (
	ErrTerminalNotSet = errors.New("TERM environment variable not set")
)

// DetectorInterface defines the interface for terminal detection
type DetectorInterface interface {
	Detect() (*Environment, error)
}

// Environment represents the detected terminal environment.
type Environment struct {
	IsITerm2 bool
	IsSSH    bool
	IsTmux   bool
	StartSeq string
	EndSeq   string
}

// Detector handles terminal environment detection.
type Detector struct{}

// NewDetector creates a new terminal detector.
func NewDetector() *Detector {
	return &Detector{}
}

// Detect determines the terminal environment and returns appropriate formatting.
func (d *Detector) Detect() (*Environment, error) {
	// Check if TERM environment variable exists
	_, ok := os.LookupEnv("TERM")
	if !ok {
		return nil, ErrTerminalNotSet
	}

	env := &Environment{}

	// Check for iTerm2/VSCode terminal
	termProgram, ok := os.LookupEnv("TERM_PROGRAM")
	if ok {
		if strings.HasPrefix(termProgram, "iTerm") || strings.HasPrefix(termProgram, "vscode") {
			env.IsITerm2 = true
		}
	}

	// Check for SSH client
	_, ok = os.LookupEnv("SSH_CLIENT")
	if ok {
		env.IsSSH = true
	}

	// Check for tmux
	term := os.Getenv("TERM")
	if !env.IsITerm2 && !env.IsSSH && strings.HasPrefix(term, "screen") {
		env.IsTmux = true
		env.StartSeq = "\033Ptmux;\033\033]"
		env.EndSeq = "\a\033\\"
	} else {
		env.StartSeq = "\033]"
		env.EndSeq = "\a"
	}

	return env, nil
}

// Formatter handles message formatting for different terminal environments.
type Formatter struct {
	env *Environment
}

// NewFormatter creates a new message formatter.
func NewFormatter(env *Environment) *Formatter {
	return &Formatter{env: env}
}

// Format formats a message for the detected terminal environment.
func (f *Formatter) Format(message string) string {
	if message == "" {
		return ""
	}
	return fmt.Sprintf("%s%s%s", f.env.StartSeq, message, f.env.EndSeq)
}
