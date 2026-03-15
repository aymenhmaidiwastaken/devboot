package editor

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
)

type JetBrainsConfig struct {
	Plugins []string `yaml:"plugins,omitempty"`
}

// JetBrainsApply installs JetBrains IDE plugins via the command-line plugin manager.
func JetBrainsApply(plugins []string) error {
	if len(plugins) == 0 {
		return nil
	}

	ui.Section("Configuring JetBrains IDE")

	// Try to find any JetBrains IDE CLI
	ideBin := findJetBrainsBin()
	if ideBin == "" {
		ui.Warn("no JetBrains IDE CLI found — skipping plugins")
		return nil
	}

	for _, plugin := range plugins {
		ui.Installing(fmt.Sprintf("plugin %s...", plugin))
		cmd := exec.Command(ideBin, "installPlugins", plugin)
		if output, err := cmd.CombinedOutput(); err != nil {
			ui.Fail(fmt.Sprintf("plugin %s: %s", plugin, strings.TrimSpace(string(output))))
		} else {
			ui.Success(fmt.Sprintf("plugin %s", plugin))
		}
	}

	return nil
}

func findJetBrainsBin() string {
	candidates := []string{
		"idea", "goland", "pycharm", "webstorm", "phpstorm",
		"clion", "rubymine", "rider", "datagrip",
	}
	for _, name := range candidates {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}
	return ""
}
