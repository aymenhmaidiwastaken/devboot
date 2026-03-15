package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setTempHome overrides HOME (and USERPROFILE on Windows) so ConfigFile()
// returns paths inside t.TempDir(). Returns the temp home directory.
func setTempHome(t *testing.T) string {
	t.Helper()
	tmpHome := t.TempDir()

	origHome := os.Getenv("HOME")
	origUserProfile := os.Getenv("USERPROFILE")
	t.Cleanup(func() {
		os.Setenv("HOME", origHome)
		os.Setenv("USERPROFILE", origUserProfile)
	})

	os.Setenv("HOME", tmpHome)
	os.Setenv("USERPROFILE", tmpHome)
	return tmpHome
}

// createConfigFile creates the shell config file inside tmpHome and any
// required parent directories, optionally seeding it with initial content.
func createConfigFile(t *testing.T, tmpHome string, shellType ShellType, initialContent string) string {
	t.Helper()
	configPath := ConfigFile(shellType)
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", dir, err)
	}
	if err := os.WriteFile(configPath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("WriteFile(%s): %v", configPath, err)
	}
	return configPath
}

// --- SetAliases tests ---

func TestSetAliases_ZshFormat(t *testing.T) {
	tmpHome := setTempHome(t)
	configPath := createConfigFile(t, tmpHome, Zsh, "# existing config\n")

	aliases := map[string]string{
		"ll": "ls -la",
		"gs": "git status",
	}

	if err := SetAliases(Zsh, aliases); err != nil {
		t.Fatalf("SetAliases() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "alias ll='ls -la'") {
		t.Error("expected zsh-format alias ll='ls -la' in config")
	}
	if !strings.Contains(text, "alias gs='git status'") {
		t.Error("expected zsh-format alias gs='git status' in config")
	}
}

func TestSetAliases_FishFormat(t *testing.T) {
	tmpHome := setTempHome(t)
	configPath := createConfigFile(t, tmpHome, Fish, "")

	aliases := map[string]string{
		"ll": "ls -la",
	}

	if err := SetAliases(Fish, aliases); err != nil {
		t.Fatalf("SetAliases() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "alias ll 'ls -la'") {
		t.Errorf("expected fish-format alias, got:\n%s", text)
	}
	// Fish format should NOT use = sign
	if strings.Contains(text, "alias ll=") {
		t.Error("fish alias should use space, not '='")
	}
}

func TestSetAliases_SkipsExistingAliases(t *testing.T) {
	tmpHome := setTempHome(t)
	existing := "# config\nalias ll='ls -la'\n"
	configPath := createConfigFile(t, tmpHome, Zsh, existing)

	aliases := map[string]string{
		"ll": "ls -la",
	}

	if err := SetAliases(Zsh, aliases); err != nil {
		t.Fatalf("SetAliases() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	text := string(content)
	count := strings.Count(text, "alias ll='ls -la'")
	if count != 1 {
		t.Errorf("expected alias to appear exactly once, found %d times", count)
	}
}

func TestSetAliases_EmptyMap(t *testing.T) {
	tmpHome := setTempHome(t)
	initial := "# original content\n"
	configPath := createConfigFile(t, tmpHome, Zsh, initial)

	if err := SetAliases(Zsh, map[string]string{}); err != nil {
		t.Fatalf("SetAliases() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if string(content) != initial {
		t.Errorf("empty aliases should not modify config, got:\n%s", string(content))
	}
}

func TestSetAliases_WritesDevbootMarker(t *testing.T) {
	tmpHome := setTempHome(t)
	configPath := createConfigFile(t, tmpHome, Bash, "")

	aliases := map[string]string{
		"k": "kubectl",
	}

	if err := SetAliases(Bash, aliases); err != nil {
		t.Fatalf("SetAliases() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if !strings.Contains(string(content), "# devboot aliases") {
		t.Error("expected '# devboot aliases' marker comment in config")
	}
}

// --- SetEnvVars tests ---

func TestSetEnvVars_BashZshExportFormat(t *testing.T) {
	tmpHome := setTempHome(t)
	configPath := createConfigFile(t, tmpHome, Bash, "")

	envVars := map[string]string{
		"EDITOR": "vim",
		"GOPATH": "$HOME/go",
	}

	if err := SetEnvVars(Bash, envVars); err != nil {
		t.Fatalf("SetEnvVars() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "export EDITOR=vim") {
		t.Error("expected 'export EDITOR=vim' in config")
	}
	if !strings.Contains(text, "export GOPATH=$HOME/go") {
		t.Error("expected 'export GOPATH=$HOME/go' in config")
	}
}

func TestSetEnvVars_FishSetFormat(t *testing.T) {
	tmpHome := setTempHome(t)
	configPath := createConfigFile(t, tmpHome, Fish, "")

	envVars := map[string]string{
		"EDITOR": "vim",
	}

	if err := SetEnvVars(Fish, envVars); err != nil {
		t.Fatalf("SetEnvVars() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "set -gx EDITOR vim") {
		t.Errorf("expected fish 'set -gx' format, got:\n%s", text)
	}
	// Should NOT contain export
	if strings.Contains(text, "export") {
		t.Error("fish env vars should use 'set -gx', not 'export'")
	}
}

func TestSetEnvVars_SkipsExistingVars(t *testing.T) {
	tmpHome := setTempHome(t)
	existing := "export EDITOR=vim\n"
	configPath := createConfigFile(t, tmpHome, Bash, existing)

	envVars := map[string]string{
		"EDITOR": "vim",
	}

	if err := SetEnvVars(Bash, envVars); err != nil {
		t.Fatalf("SetEnvVars() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	count := strings.Count(string(content), "export EDITOR=vim")
	if count != 1 {
		t.Errorf("expected env var to appear exactly once, found %d times", count)
	}
}

func TestSetEnvVars_EmptyMap(t *testing.T) {
	tmpHome := setTempHome(t)
	initial := "# original\n"
	configPath := createConfigFile(t, tmpHome, Zsh, initial)

	if err := SetEnvVars(Zsh, map[string]string{}); err != nil {
		t.Fatalf("SetEnvVars() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if string(content) != initial {
		t.Errorf("empty env vars should not modify config, got:\n%s", string(content))
	}
}

func TestSetEnvVars_WritesDevbootMarker(t *testing.T) {
	tmpHome := setTempHome(t)
	configPath := createConfigFile(t, tmpHome, Zsh, "")

	envVars := map[string]string{
		"FOO": "bar",
	}

	if err := SetEnvVars(Zsh, envVars); err != nil {
		t.Fatalf("SetEnvVars() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if !strings.Contains(string(content), "# devboot env") {
		t.Error("expected '# devboot env' marker comment in config")
	}
}
