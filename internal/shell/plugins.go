package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
)

// Plugin repos for common zsh plugins.
var pluginRepos = map[string]string{
	"zsh-autosuggestions":     "https://github.com/zsh-users/zsh-autosuggestions.git",
	"zsh-syntax-highlighting": "https://github.com/zsh-users/zsh-syntax-highlighting.git",
	"fzf-tab":                 "https://github.com/Aloxaf/fzf-tab.git",
	"zsh-completions":         "https://github.com/zsh-users/zsh-completions.git",
	"zsh-history-substring-search": "https://github.com/zsh-users/zsh-history-substring-search.git",
}

// InstallPlugins installs shell plugins via git clone.
func InstallPlugins(shellType ShellType, plugins []string) error {
	if len(plugins) == 0 {
		return nil
	}

	ui.Section("Installing shell plugins")

	home, _ := os.UserHomeDir()
	pluginDir := filepath.Join(home, ".local", "share", "devboot", "plugins")
	os.MkdirAll(pluginDir, 0755)

	configPath := ConfigFile(shellType)
	existing, _ := os.ReadFile(configPath)
	content := string(existing)

	var sourceLines []string

	for _, plugin := range plugins {
		dest := filepath.Join(pluginDir, plugin)

		if _, err := os.Stat(dest); err == nil {
			ui.Skip(plugin)
		} else {
			repo, ok := pluginRepos[plugin]
			if !ok {
				// Assume it's a GitHub shorthand or full URL
				if strings.Contains(plugin, "/") {
					repo = "https://github.com/" + plugin + ".git"
				} else {
					ui.Fail(fmt.Sprintf("unknown plugin: %s", plugin))
					continue
				}
			}

			ui.Installing(fmt.Sprintf("cloning %s...", plugin))
			cmd := exec.Command("git", "clone", "--depth=1", repo, dest)
			if output, err := cmd.CombinedOutput(); err != nil {
				ui.Fail(fmt.Sprintf("%s: %s", plugin, string(output)))
				continue
			}
			ui.Success(plugin)
		}

		// Build source line for shell config
		if shellType == Zsh {
			// Find the main .zsh file
			mainFile := findPluginFile(dest, plugin)
			if mainFile != "" {
				sourceLine := fmt.Sprintf("source %s", mainFile)
				if !strings.Contains(content, sourceLine) {
					sourceLines = append(sourceLines, sourceLine)
				}
			}
		}
	}

	// Add source lines to shell config
	if len(sourceLines) > 0 {
		block := "\n# devboot plugins\n" + strings.Join(sourceLines, "\n") + "\n"
		f, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("updating shell config: %w", err)
		}
		defer f.Close()
		f.WriteString(block)
	}

	return nil
}

func findPluginFile(dir, name string) string {
	// Try common patterns
	candidates := []string{
		filepath.Join(dir, name+".zsh"),
		filepath.Join(dir, name+".plugin.zsh"),
		filepath.Join(dir, name+".zsh-theme"),
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}

	// Fallback: look for any .plugin.zsh file
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".plugin.zsh") {
			return filepath.Join(dir, e.Name())
		}
	}

	return ""
}
