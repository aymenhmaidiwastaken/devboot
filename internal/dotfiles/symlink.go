package dotfiles

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// createSymlink creates a symlink from src to dst, backing up any existing file.
func createSymlink(src, dst string) error {
	// Verify source exists
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("source not found: %s", src)
	}

	// Check if symlink already points to correct target
	if target, err := os.Readlink(dst); err == nil {
		absSrc, _ := filepath.Abs(src)
		absTarget, _ := filepath.Abs(target)
		if absSrc == absTarget {
			return nil // already correct
		}
	}

	// Backup existing file if it's not a symlink
	if info, err := os.Lstat(dst); err == nil {
		if info.Mode()&os.ModeSymlink == 0 {
			backupPath := dst + ".devboot-backup." + time.Now().Format("20060102-150405")
			if err := os.Rename(dst, backupPath); err != nil {
				return fmt.Errorf("backing up %s: %w", dst, err)
			}
		} else {
			os.Remove(dst)
		}
	}

	// Ensure parent directory exists
	os.MkdirAll(filepath.Dir(dst), 0755)

	return os.Symlink(src, dst)
}
