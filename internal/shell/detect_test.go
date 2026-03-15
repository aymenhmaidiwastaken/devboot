package shell

import (
	"os"
	"strings"
	"testing"
)

func TestDetect_ReturnsValidShellType(t *testing.T) {
	st := Detect()
	switch st {
	case Zsh, Bash, Fish:
		// valid
	default:
		t.Errorf("Detect() returned unexpected shell type: %q", st)
	}
}

func TestDetect_WithSHELLEnvVar(t *testing.T) {
	tests := []struct {
		name     string
		envVal   string
		expected ShellType
	}{
		{"zsh path", "/bin/zsh", Zsh},
		{"usr zsh", "/usr/bin/zsh", Zsh},
		{"bash path", "/bin/bash", Bash},
		{"usr bash", "/usr/bin/bash", Bash},
		{"fish path", "/usr/bin/fish", Fish},
		{"fish local", "/usr/local/bin/fish", Fish},
		{"unknown defaults to bash", "/usr/bin/sh", Bash},
		{"custom zsh path", "/opt/homebrew/bin/zsh", Zsh},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := os.Getenv("SHELL")
			t.Cleanup(func() { os.Setenv("SHELL", orig) })

			os.Setenv("SHELL", tt.envVal)
			got := Detect()
			if got != tt.expected {
				t.Errorf("Detect() with SHELL=%q = %q, want %q", tt.envVal, got, tt.expected)
			}
		})
	}
}

func TestDetect_WithEmptySHELL(t *testing.T) {
	orig := os.Getenv("SHELL")
	t.Cleanup(func() { os.Setenv("SHELL", orig) })

	os.Setenv("SHELL", "")
	st := Detect()
	// With empty SHELL, Detect falls back to LookPath("zsh") or Bash.
	// Either is valid.
	switch st {
	case Zsh, Bash:
		// valid fallback
	default:
		t.Errorf("Detect() with empty SHELL returned unexpected: %q", st)
	}
}

func TestConfigFile_ReturnsCorrectPaths(t *testing.T) {
	tests := []struct {
		shell    ShellType
		contains string
	}{
		{Zsh, ".zshrc"},
		{Bash, ".bashrc"},
		{Fish, "config.fish"},
	}

	for _, tt := range tests {
		t.Run(string(tt.shell), func(t *testing.T) {
			path := ConfigFile(tt.shell)
			if path == "" {
				t.Fatalf("ConfigFile(%s) returned empty string", tt.shell)
			}
			if !strings.Contains(path, tt.contains) {
				t.Errorf("ConfigFile(%s) = %q, want path containing %q", tt.shell, path, tt.contains)
			}
		})
	}
}

func TestConfigFile_PathsContainExpectedFilenames(t *testing.T) {
	zshPath := ConfigFile(Zsh)
	if !strings.HasSuffix(zshPath, ".zshrc") {
		t.Errorf("zsh config path %q does not end with .zshrc", zshPath)
	}

	bashPath := ConfigFile(Bash)
	if !strings.HasSuffix(bashPath, ".bashrc") {
		t.Errorf("bash config path %q does not end with .bashrc", bashPath)
	}

	fishPath := ConfigFile(Fish)
	if !strings.HasSuffix(fishPath, "config.fish") {
		t.Errorf("fish config path %q does not end with config.fish", fishPath)
	}

	// Fish config should be in a subdirectory structure
	if !strings.Contains(fishPath, "fish") {
		t.Errorf("fish config path %q does not contain 'fish' directory", fishPath)
	}
}
