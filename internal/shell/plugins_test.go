package shell

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallPlugins_EmptyList(t *testing.T) {
	// Should be a no-op and return nil
	err := InstallPlugins(Zsh, []string{})
	if err != nil {
		t.Fatalf("InstallPlugins with empty list should return nil, got: %v", err)
	}

	err = InstallPlugins(Zsh, nil)
	if err != nil {
		t.Fatalf("InstallPlugins with nil list should return nil, got: %v", err)
	}
}

func TestInstallPlugins_UnknownPluginName(t *testing.T) {
	tmpHome := setTempHome(t)
	_ = createConfigFile(t, tmpHome, Zsh, "")

	// A name that is not in pluginRepos and has no "/" (so it's not a GitHub shorthand)
	// This should call ui.Fail and continue, not return an error.
	err := InstallPlugins(Zsh, []string{"nonexistent-plugin-xyz"})
	if err != nil {
		t.Fatalf("InstallPlugins should not return error for unknown plugin, got: %v", err)
	}
}

func TestInstallPlugins_PluginAlreadyExists(t *testing.T) {
	tmpHome := setTempHome(t)
	_ = createConfigFile(t, tmpHome, Zsh, "")

	// Pre-create the plugin directory so InstallPlugins sees it exists
	pluginDir := filepath.Join(tmpHome, ".local", "share", "devboot", "plugins", "zsh-autosuggestions")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	// Create a plugin file so findPluginFile can find it
	pluginFile := filepath.Join(pluginDir, "zsh-autosuggestions.zsh")
	if err := os.WriteFile(pluginFile, []byte("# plugin"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	err := InstallPlugins(Zsh, []string{"zsh-autosuggestions"})
	if err != nil {
		t.Fatalf("InstallPlugins should skip existing plugin, got error: %v", err)
	}
}

func TestFindPluginFile_FindsPluginZshFile(t *testing.T) {
	dir := t.TempDir()
	pluginName := "my-plugin"

	// Create a .plugin.zsh file
	pluginFile := filepath.Join(dir, pluginName+".plugin.zsh")
	if err := os.WriteFile(pluginFile, []byte("# plugin code"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got := findPluginFile(dir, pluginName)
	if got != pluginFile {
		t.Errorf("findPluginFile() = %q, want %q", got, pluginFile)
	}
}

func TestFindPluginFile_FindsZshFile(t *testing.T) {
	dir := t.TempDir()
	pluginName := "my-plugin"

	// Create only a .zsh file (no .plugin.zsh)
	zshFile := filepath.Join(dir, pluginName+".zsh")
	if err := os.WriteFile(zshFile, []byte("# zsh code"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got := findPluginFile(dir, pluginName)
	if got != zshFile {
		t.Errorf("findPluginFile() = %q, want %q", got, zshFile)
	}
}

func TestFindPluginFile_ReturnsEmptyForEmptyDir(t *testing.T) {
	dir := t.TempDir()

	got := findPluginFile(dir, "anything")
	if got != "" {
		t.Errorf("findPluginFile() on empty dir = %q, want empty string", got)
	}
}

func TestFindPluginFile_FallbackFindsAnyPluginZsh(t *testing.T) {
	dir := t.TempDir()

	// Create a .plugin.zsh file with a different name than what we search for
	otherFile := filepath.Join(dir, "other-name.plugin.zsh")
	if err := os.WriteFile(otherFile, []byte("# fallback"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got := findPluginFile(dir, "nonmatching-name")
	if got != otherFile {
		t.Errorf("findPluginFile() fallback = %q, want %q", got, otherFile)
	}
}
