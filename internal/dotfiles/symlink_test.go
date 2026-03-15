package dotfiles

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateSymlink(t *testing.T) {
	dir := t.TempDir()

	// Create source file
	src := filepath.Join(dir, "source.txt")
	os.WriteFile(src, []byte("hello"), 0644)

	// Create symlink
	dst := filepath.Join(dir, "link.txt")
	if err := createSymlink(src, dst); err != nil {
		if os.IsPermission(err) || strings.Contains(err.Error(), "privilege") {
			t.Skip("symlinks require elevated privileges on Windows")
		}
		t.Fatalf("createSymlink failed: %v", err)
	}

	// Verify symlink exists and points to source
	target, err := os.Readlink(dst)
	if err != nil {
		t.Fatalf("Readlink failed: %v", err)
	}
	absSrc, _ := filepath.Abs(src)
	absTarget, _ := filepath.Abs(target)
	if absSrc != absTarget {
		t.Errorf("symlink target = %s, want %s", absTarget, absSrc)
	}

	// Run again — should be idempotent
	if err := createSymlink(src, dst); err != nil {
		t.Fatalf("second createSymlink failed: %v", err)
	}
}

func TestCreateSymlinkBackup(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "source.txt")
	os.WriteFile(src, []byte("new"), 0644)

	// Create existing file at destination
	dst := filepath.Join(dir, "existing.txt")
	os.WriteFile(dst, []byte("old"), 0644)

	if err := createSymlink(src, dst); err != nil {
		if os.IsPermission(err) || strings.Contains(err.Error(), "privilege") {
			t.Skip("symlinks require elevated privileges on Windows")
		}
		t.Fatalf("createSymlink failed: %v", err)
	}

	// Verify backup was created
	matches, _ := filepath.Glob(filepath.Join(dir, "existing.txt.devboot-backup.*"))
	if len(matches) == 0 {
		t.Error("expected backup file to be created")
	}
}

func TestCreateSymlinkMissingSource(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "nonexistent")
	dst := filepath.Join(dir, "link")

	if err := createSymlink(src, dst); err == nil {
		t.Error("expected error for missing source")
	}
}
