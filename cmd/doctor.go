package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aymenhmaidiwastaken/devboot/internal/platform"
	"github.com/aymenhmaidiwastaken/devboot/internal/shell"
	"github.com/aymenhmaidiwastaken/devboot/internal/state"
	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose common issues with your environment",
	RunE:  runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, args []string) error {
	p := platform.Detect()
	issues := 0

	fmt.Printf("\n  devboot doctor — comprehensive environment check\n")

	// System
	ui.Section("System")
	ui.Info(fmt.Sprintf("OS: %s", p.OS))
	if p.Distro.String() != "" {
		ui.Info(fmt.Sprintf("Distro: %s", p.Distro))
	}
	ui.Info(fmt.Sprintf("Arch: %s", runtime.GOARCH))
	ui.Info(fmt.Sprintf("Go: %s", runtime.Version()))

	// Shell
	ui.Section("Shell")
	shellType := shell.Detect()
	ui.Info(fmt.Sprintf("Current: %s", shellType))
	configPath := shell.ConfigFile(shellType)
	ui.Info(fmt.Sprintf("Config: %s", configPath))

	// Check if shell config exists
	if _, err := os.Stat(configPath); err != nil {
		ui.Fail(fmt.Sprintf("%s does not exist", configPath))
		ui.Info(fmt.Sprintf("  fix: touch %s", configPath))
		issues++
	} else {
		ui.Success(fmt.Sprintf("%s exists", filepath.Base(configPath)))

		// Check if shell config is actually sourced (basic check)
		data, _ := os.ReadFile(configPath)
		content := string(data)
		if strings.Contains(content, "devboot") {
			ui.Success("devboot config blocks present")
		}
	}

	// Package managers
	ui.Section("Package managers")
	hasPkgManager := false
	if checkCommandDoctor("brew", "Homebrew") {
		hasPkgManager = true
	}
	if checkCommandDoctor("apt-get", "APT") {
		hasPkgManager = true
	}
	if checkCommandDoctor("pacman", "Pacman") {
		hasPkgManager = true
	}
	if !hasPkgManager {
		ui.Warn("no supported package manager found")
		switch p.OS {
		case platform.MacOS:
			ui.Info("  fix: /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
		case platform.Linux, platform.WSL:
			ui.Info("  fix: sudo apt-get update (for Debian/Ubuntu)")
		}
		issues++
	}

	// Core tools
	ui.Section("Core tools")
	coreTools := []struct{ bin, label string }{
		{"git", "Git"},
		{"curl", "curl"},
		{"ssh-keygen", "SSH"},
	}
	for _, t := range coreTools {
		if !checkCommandDoctor(t.bin, t.label) {
			issues++
		}
	}

	// PATH analysis
	ui.Section("PATH analysis")
	pathDirs := filepath.SplitList(os.Getenv("PATH"))
	ui.Info(fmt.Sprintf("%d directories in PATH", len(pathDirs)))

	// Check for common issues
	home, _ := os.UserHomeDir()

	// Check if ~/.local/bin is in PATH
	localBin := filepath.Join(home, ".local", "bin")
	if !containsPath(pathDirs, localBin) {
		ui.Warn(fmt.Sprintf("~/.local/bin not in PATH"))
		ui.Info(fmt.Sprintf("  fix: export PATH=\"%s:$PATH\"", localBin))
		issues++
	} else {
		ui.Success("~/.local/bin in PATH")
	}

	// Check for duplicate PATH entries
	seen := make(map[string]int)
	for _, d := range pathDirs {
		seen[d]++
	}
	dupes := 0
	for _, count := range seen {
		if count > 1 {
			dupes++
		}
	}
	if dupes > 0 {
		ui.Warn(fmt.Sprintf("%d duplicate PATH entries", dupes))
	} else {
		ui.Success("no duplicate PATH entries")
	}

	// Version manager conflicts
	ui.Section("Version managers")
	nodeManagers := detectVersionManagers("node", []string{"nvm", "fnm", "volta", "mise", "asdf"})
	if len(nodeManagers) > 1 {
		ui.Warn(fmt.Sprintf("multiple Node version managers: %s — this can cause conflicts",
			strings.Join(nodeManagers, ", ")))
		issues++
	} else if len(nodeManagers) == 1 {
		ui.Success(fmt.Sprintf("Node manager: %s", nodeManagers[0]))
	}

	pyManagers := detectVersionManagers("python", []string{"pyenv", "conda", "mise", "asdf"})
	if len(pyManagers) > 1 {
		ui.Warn(fmt.Sprintf("multiple Python version managers: %s", strings.Join(pyManagers, ", ")))
		issues++
	} else if len(pyManagers) == 1 {
		ui.Success(fmt.Sprintf("Python manager: %s", pyManagers[0]))
	}

	// Git authentication
	ui.Section("Git")
	gitName := gitConfigGet("user.name")
	gitEmail := gitConfigGet("user.email")
	if gitName != "" {
		ui.Success(fmt.Sprintf("user.name: %s", gitName))
	} else {
		ui.Fail("user.name not set")
		ui.Info("  fix: git config --global user.name \"Your Name\"")
		issues++
	}
	if gitEmail != "" {
		ui.Success(fmt.Sprintf("user.email: %s", gitEmail))
	} else {
		ui.Fail("user.email not set")
		ui.Info("  fix: git config --global user.email \"you@example.com\"")
		issues++
	}

	// SSH key
	sshKeyPath := filepath.Join(home, ".ssh", "id_ed25519")
	if _, err := os.Stat(sshKeyPath); err == nil {
		ui.Success("SSH key exists (ed25519)")

		// Check ssh-agent
		if agentPID := os.Getenv("SSH_AUTH_SOCK"); agentPID != "" {
			ui.Success("ssh-agent running")

			// Check if key is added
			addedCmd := exec.Command("ssh-add", "-l")
			if out, err := addedCmd.Output(); err == nil && strings.Contains(string(out), "ed25519") {
				ui.Success("SSH key added to agent")
			} else {
				ui.Warn("SSH key not added to agent")
				ui.Info("  fix: ssh-add ~/.ssh/id_ed25519")
				issues++
			}
		} else {
			ui.Warn("ssh-agent not running")
			ui.Info("  fix: eval $(ssh-agent -s) && ssh-add")
			issues++
		}
	} else {
		rsaPath := filepath.Join(home, ".ssh", "id_rsa")
		if _, err := os.Stat(rsaPath); err == nil {
			ui.Success("SSH key exists (RSA)")
		} else {
			ui.Warn("no SSH key found")
			ui.Info("  fix: ssh-keygen -t ed25519 -C \"your@email.com\"")
			issues++
		}
	}

	// Test GitHub connectivity
	if platform.HasCommand("ssh") {
		sshCmd := exec.Command("ssh", "-T", "-o", "ConnectTimeout=5", "-o", "StrictHostKeyChecking=no", "git@github.com")
		if out, err := sshCmd.CombinedOutput(); err != nil {
			outStr := string(out)
			if strings.Contains(outStr, "successfully authenticated") {
				ui.Success("GitHub SSH authentication works")
			} else {
				ui.Warn("GitHub SSH authentication failed")
				ui.Info("  fix: add your SSH key to https://github.com/settings/keys")
				issues++
			}
		} else {
			ui.Success("GitHub SSH authentication works")
		}
	}

	// State
	ui.Section("DevBoot state")
	st, err := state.Load()
	if err != nil {
		ui.Warn(fmt.Sprintf("state file: %v", err))
	} else {
		ui.Info(fmt.Sprintf("managed tools: %d", len(st.ManagedTools)))
		ui.Info(fmt.Sprintf("recorded actions: %d", len(st.Actions)))
		if !st.LastApply.IsZero() {
			ui.Info(fmt.Sprintf("last apply: %s", st.LastApply.Format("2006-01-02 15:04")))
		}
	}

	// Config
	ui.Section("Config")
	if fileExists("devboot.yaml") {
		ui.Success("devboot.yaml found")
	} else {
		ui.Fail("devboot.yaml not found")
		ui.Info("  fix: devboot init")
		issues++
	}

	// Summary
	fmt.Println()
	if issues == 0 {
		ui.Success("all checks passed!")
	} else {
		ui.Warn(fmt.Sprintf("%d issue(s) found — see suggestions above", issues))
	}
	fmt.Println()

	return nil
}

func checkCommandDoctor(name, label string) bool {
	path, err := exec.LookPath(name)
	if err != nil {
		ui.Fail(fmt.Sprintf("%s: not found", label))
		return false
	}
	// Try to get version
	version := getToolVersion(name)
	if version != "" {
		ui.Success(fmt.Sprintf("%s: %s (%s)", label, version, path))
	} else {
		ui.Success(fmt.Sprintf("%s: %s", label, path))
	}
	return true
}

func getToolVersion(name string) string {
	var cmd *exec.Cmd
	switch name {
	case "git":
		cmd = exec.Command("git", "--version")
	case "node":
		cmd = exec.Command("node", "--version")
	case "go":
		cmd = exec.Command("go", "version")
	case "python3":
		cmd = exec.Command("python3", "--version")
	case "docker":
		cmd = exec.Command("docker", "--version")
	case "brew":
		cmd = exec.Command("brew", "--version")
	default:
		return ""
	}

	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.Split(string(out), "\n")[0])
}

func containsPath(paths []string, target string) bool {
	for _, p := range paths {
		if p == target {
			return true
		}
	}
	return false
}

func detectVersionManagers(lang string, managers []string) []string {
	var found []string
	for _, m := range managers {
		if _, err := exec.LookPath(m); err == nil {
			found = append(found, m)
		}
	}
	return found
}
