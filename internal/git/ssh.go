package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
)

// EnsureSSHKey generates an SSH key if one doesn't exist.
func EnsureSSHKey(email string) error {
	home, _ := os.UserHomeDir()
	keyPath := filepath.Join(home, ".ssh", "id_ed25519")

	if _, err := os.Stat(keyPath); err == nil {
		ui.Skip("SSH key already exists")
		pubKey, _ := os.ReadFile(keyPath + ".pub")
		if len(pubKey) > 0 {
			ui.Info(fmt.Sprintf("Public key: %s", string(pubKey)))
		}
		return nil
	}

	ui.Installing("generating SSH key...")

	// Ensure .ssh directory exists
	sshDir := filepath.Join(home, ".ssh")
	os.MkdirAll(sshDir, 0700)

	comment := email
	if comment == "" {
		comment = "devboot"
	}

	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-C", comment, "-f", keyPath, "-N", "")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ssh-keygen failed: %s: %w", string(output), err)
	}

	ui.Success("SSH key generated")

	pubKey, _ := os.ReadFile(keyPath + ".pub")
	if len(pubKey) > 0 {
		ui.Info(fmt.Sprintf("Public key:\n%s", string(pubKey)))
	}

	return nil
}
