package editor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
)

type NeovimConfig struct {
	ConfigRepo string `yaml:"configRepo,omitempty"`
}

// NeovimApply sets up Neovim by cloning a config repo.
func NeovimApply(configRepo string) error {
	if configRepo == "" {
		return nil
	}

	ui.Section("Configuring Neovim")

	if _, err := exec.LookPath("nvim"); err != nil {
		ui.Warn("nvim not found — install neovim first")
		return nil
	}

	home, _ := os.UserHomeDir()
	nvimConfigDir := filepath.Join(home, ".config", "nvim")

	// Check if config already exists
	if _, err := os.Stat(filepath.Join(nvimConfigDir, "init.lua")); err == nil {
		ui.Skip("Neovim config already exists")
		return nil
	}
	if _, err := os.Stat(filepath.Join(nvimConfigDir, "init.vim")); err == nil {
		ui.Skip("Neovim config already exists")
		return nil
	}

	// Backup existing config if present
	if info, err := os.Stat(nvimConfigDir); err == nil && info.IsDir() {
		backupDir := nvimConfigDir + ".devboot-backup"
		ui.Info(fmt.Sprintf("backing up existing config to %s", backupDir))
		os.Rename(nvimConfigDir, backupDir)
	}

	ui.Installing(fmt.Sprintf("cloning %s...", configRepo))
	cmd := exec.Command("git", "clone", "--depth=1", configRepo, nvimConfigDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("cloning neovim config: %s: %w", string(output), err)
	}

	ui.Success("Neovim config installed")
	ui.Info("Run nvim to complete plugin installation")
	return nil
}
