package tools

import (
	"fmt"
	"os/exec"
	"sync"
)

type AptBackend struct {
	updated bool
	mu      sync.Mutex
}

func (a *AptBackend) Name() string { return "apt" }

func (a *AptBackend) IsAvailable() bool {
	_, err := exec.LookPath("apt-get")
	return err == nil
}

func (a *AptBackend) EnsureReady() error {
	if !a.IsAvailable() {
		return fmt.Errorf("apt-get not found")
	}
	return a.update()
}

func (a *AptBackend) update() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.updated {
		return nil
	}

	cmd := exec.Command("sudo", "apt-get", "update", "-qq")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("apt-get update failed: %w", err)
	}
	a.updated = true
	return nil
}

func (a *AptBackend) Install(pkg string) error {
	cmd := exec.Command("sudo", "apt-get", "install", "-y", "-qq", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}

func (a *AptBackend) Uninstall(pkg string) error {
	cmd := exec.Command("sudo", "apt-get", "remove", "-y", "-qq", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}

func (a *AptBackend) IsInstalled(pkg string) bool {
	cmd := exec.Command("dpkg", "-s", pkg)
	return cmd.Run() == nil
}
