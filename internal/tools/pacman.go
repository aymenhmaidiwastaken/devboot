package tools

import (
	"fmt"
	"os/exec"
)

type PacmanBackend struct{}

func (p *PacmanBackend) Name() string { return "pacman" }

func (p *PacmanBackend) IsAvailable() bool {
	_, err := exec.LookPath("pacman")
	return err == nil
}

func (p *PacmanBackend) EnsureReady() error {
	if !p.IsAvailable() {
		return fmt.Errorf("pacman not found")
	}
	return nil
}

func (p *PacmanBackend) Install(pkg string) error {
	cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", "--needed", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}

func (p *PacmanBackend) Uninstall(pkg string) error {
	cmd := exec.Command("sudo", "pacman", "-R", "--noconfirm", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}

func (p *PacmanBackend) IsInstalled(pkg string) bool {
	cmd := exec.Command("pacman", "-Q", pkg)
	return cmd.Run() == nil
}
