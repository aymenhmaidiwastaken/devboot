package tools

import (
	"fmt"
	"os/exec"

	"github.com/aymenhmaidiwastaken/devboot/internal/platform"
	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
)

type BrewBackend struct{}

func (b *BrewBackend) Name() string { return "Homebrew" }

func (b *BrewBackend) IsAvailable() bool {
	return platform.HasCommand("brew")
}

func (b *BrewBackend) EnsureReady() error {
	if b.IsAvailable() {
		return nil
	}

	ui.Info("Homebrew not found, installing...")
	cmd := exec.Command("/bin/bash", "-c",
		`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Homebrew: %w", err)
	}
	ui.Success("Homebrew installed")
	return nil
}

func (b *BrewBackend) Install(pkg string) error {
	cmd := exec.Command("brew", "install", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}

func (b *BrewBackend) Uninstall(pkg string) error {
	cmd := exec.Command("brew", "uninstall", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}

func (b *BrewBackend) IsInstalled(pkg string) bool {
	cmd := exec.Command("brew", "list", pkg)
	return cmd.Run() == nil
}
