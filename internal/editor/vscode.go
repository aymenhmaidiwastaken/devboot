package editor

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
)

// VSCodeApply installs extensions and applies settings.
func VSCodeApply(extensions []string, settings map[string]interface{}) error {
	if len(extensions) == 0 && len(settings) == 0 {
		return nil
	}

	ui.Section("Configuring VS Code")

	codeBin := findVSCodeBin()
	if codeBin == "" {
		ui.Warn("VS Code CLI (code) not found — skipping extensions")
	}

	if codeBin != "" && len(extensions) > 0 {
		installed := listInstalledExtensions(codeBin)
		for _, ext := range extensions {
			if installed[strings.ToLower(ext)] {
				ui.Skip(fmt.Sprintf("extension %s", ext))
				continue
			}
			ui.Installing(fmt.Sprintf("extension %s...", ext))
			cmd := exec.Command(codeBin, "--install-extension", ext, "--force")
			if output, err := cmd.CombinedOutput(); err != nil {
				ui.Fail(fmt.Sprintf("extension %s: %s", ext, strings.TrimSpace(string(output))))
			} else {
				ui.Success(fmt.Sprintf("extension %s", ext))
			}
		}
	}

	if len(settings) > 0 {
		if err := applyVSCodeSettings(settings); err != nil {
			ui.Warn(fmt.Sprintf("settings: %v", err))
		}
	}

	return nil
}

// VSCodeStatus returns installed extensions.
func VSCodeStatus() (installed []string, codePath string) {
	codeBin := findVSCodeBin()
	if codeBin == "" {
		return nil, ""
	}
	codePath = codeBin
	exts := listInstalledExtensions(codeBin)
	for ext := range exts {
		installed = append(installed, ext)
	}
	return
}

func findVSCodeBin() string {
	for _, name := range []string{"code", "code-insiders", "codium"} {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}
	return ""
}

func listInstalledExtensions(codeBin string) map[string]bool {
	cmd := exec.Command(codeBin, "--list-extensions")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	result := make(map[string]bool)
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			result[strings.ToLower(line)] = true
		}
	}
	return result
}

func vsCodeSettingsPath() string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Code", "User", "settings.json")
	case "linux":
		return filepath.Join(home, ".config", "Code", "User", "settings.json")
	default: // windows
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			appdata = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appdata, "Code", "User", "settings.json")
	}
}

func applyVSCodeSettings(settings map[string]interface{}) error {
	settingsPath := vsCodeSettingsPath()

	// Read existing settings
	existing := make(map[string]interface{})
	if data, err := os.ReadFile(settingsPath); err == nil {
		json.Unmarshal(data, &existing)
	}

	// Merge settings (new values override)
	changed := false
	for key, val := range settings {
		if existingVal, ok := existing[key]; ok {
			if fmt.Sprintf("%v", existingVal) == fmt.Sprintf("%v", val) {
				ui.Skip(fmt.Sprintf("setting %s", key))
				continue
			}
		}
		existing[key] = val
		changed = true
		ui.Success(fmt.Sprintf("setting %s = %v", key, val))
	}

	if !changed {
		return nil
	}

	// Write back
	os.MkdirAll(filepath.Dir(settingsPath), 0755)
	data, err := json.MarshalIndent(existing, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(settingsPath, data, 0644)
}
