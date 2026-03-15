package dotfiles

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
)

const dotfilesDir = ".local/share/devboot/dotfiles"

// DotfilesDirPath returns the relative path to the dotfiles directory.
func DotfilesDirPath() string {
	return dotfilesDir
}

// ExpandHome expands ~ to the user's home directory (exported wrapper).
func ExpandHome(path string) string {
	return expandHome(path)
}

// Sync clones a dotfiles repo and symlinks files according to the mappings.
func Sync(repo string, mappings map[string]string) error {
	if repo == "" {
		return nil
	}

	ui.Section("Syncing dotfiles")

	home, _ := os.UserHomeDir()
	dest := filepath.Join(home, dotfilesDir)

	// Clone or pull
	if _, err := os.Stat(filepath.Join(dest, ".git")); err == nil {
		ui.Installing("pulling latest dotfiles...")
		cmd := exec.Command("git", "-C", dest, "pull", "--ff-only")
		if output, err := cmd.CombinedOutput(); err != nil {
			ui.Warn(fmt.Sprintf("git pull: %s", string(output)))
		} else {
			ui.Success("dotfiles updated")
		}
	} else {
		ui.Installing(fmt.Sprintf("cloning %s...", repo))
		os.MkdirAll(filepath.Dir(dest), 0755)
		cmd := exec.Command("git", "clone", "--depth=1", repo, dest)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("cloning dotfiles: %s: %w", string(output), err)
		}
		ui.Success("dotfiles cloned")
	}

	// Apply mappings
	if len(mappings) > 0 {
		for src, dst := range mappings {
			srcPath := filepath.Join(dest, src)
			dstPath := expandHome(dst)

			if err := createSymlink(srcPath, dstPath); err != nil {
				ui.Fail(fmt.Sprintf("%s → %s: %v", src, dst, err))
			} else {
				ui.Success(fmt.Sprintf("%s → %s", src, dst))
			}
		}
	}

	return nil
}

// Push commits and pushes changes from the dotfiles repo.
func Push(message string) error {
	home, _ := os.UserHomeDir()
	dest := filepath.Join(home, dotfilesDir)

	if _, err := os.Stat(filepath.Join(dest, ".git")); err != nil {
		return fmt.Errorf("dotfiles repo not found at %s — run devboot apply first", dest)
	}

	ui.Section("Pushing dotfiles")

	// Check for changes
	cmd := exec.Command("git", "-C", dest, "status", "--porcelain")
	out, _ := cmd.Output()
	if len(out) == 0 {
		ui.Skip("no changes to push")
		return nil
	}

	if message == "" {
		message = "devboot: update dotfiles"
	}

	// Add, commit, push
	commands := [][]string{
		{"git", "-C", dest, "add", "-A"},
		{"git", "-C", dest, "commit", "-m", message},
		{"git", "-C", dest, "push"},
	}

	for _, args := range commands {
		c := exec.Command(args[0], args[1:]...)
		if output, err := c.CombinedOutput(); err != nil {
			return fmt.Errorf("%s: %s: %w", args[0], string(output), err)
		}
	}

	ui.Success("dotfiles pushed")
	return nil
}

func expandHome(path string) string {
	if len(path) > 1 && path[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
