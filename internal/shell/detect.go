package shell

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ShellType string

const (
	Zsh  ShellType = "zsh"
	Bash ShellType = "bash"
	Fish ShellType = "fish"
)

// Detect returns the user's current shell type.
func Detect() ShellType {
	shell := os.Getenv("SHELL")
	if shell == "" {
		// Fallback: check common locations
		if _, err := exec.LookPath("zsh"); err == nil {
			return Zsh
		}
		return Bash
	}

	base := filepath.Base(shell)
	switch {
	case strings.Contains(base, "zsh"):
		return Zsh
	case strings.Contains(base, "fish"):
		return Fish
	default:
		return Bash
	}
}

// ConfigFile returns the path to the shell's config file.
func ConfigFile(st ShellType) string {
	home, _ := os.UserHomeDir()
	switch st {
	case Zsh:
		return filepath.Join(home, ".zshrc")
	case Fish:
		return filepath.Join(home, ".config", "fish", "config.fish")
	default:
		return filepath.Join(home, ".bashrc")
	}
}
