package platform

import (
	"runtime"
	"testing"
)

// --- Detect ---

func TestDetect(t *testing.T) {
	p := Detect()

	if p.Arch != runtime.GOARCH {
		t.Errorf("expected arch %s, got %s", runtime.GOARCH, p.Arch)
	}

	switch p.OS {
	case MacOS, Linux, WSL, Windows:
		// valid
	default:
		t.Errorf("unexpected OS: %v", p.OS)
	}
}

func TestDetectArchIsSet(t *testing.T) {
	p := Detect()
	if p.Arch == "" {
		t.Error("expected Arch to be set")
	}
}

func TestDetectOSMatchesRuntime(t *testing.T) {
	p := Detect()
	switch runtime.GOOS {
	case "darwin":
		if p.OS != MacOS {
			t.Errorf("expected MacOS on darwin, got %v", p.OS)
		}
	case "windows":
		if p.OS != Windows {
			t.Errorf("expected Windows on windows, got %v", p.OS)
		}
	case "linux":
		if p.OS != Linux && p.OS != WSL {
			t.Errorf("expected Linux or WSL on linux, got %v", p.OS)
		}
	}
}

// --- OS.String ---

func TestOSString(t *testing.T) {
	tests := []struct {
		os       OS
		expected string
	}{
		{MacOS, "macOS"},
		{Linux, "Linux"},
		{WSL, "WSL"},
		{Windows, "Windows"},
		{Unknown, "Unknown"},
	}
	for _, tt := range tests {
		if got := tt.os.String(); got != tt.expected {
			t.Errorf("OS(%d).String() = %q, want %q", tt.os, got, tt.expected)
		}
	}
}

func TestOSStringCoversAllValues(t *testing.T) {
	// Ensure no OS value returns empty string
	for _, os := range []OS{MacOS, Linux, WSL, Windows, Unknown} {
		s := os.String()
		if s == "" {
			t.Errorf("OS(%d).String() returned empty string", os)
		}
	}
}

// --- Distro.String ---

func TestDistroString(t *testing.T) {
	tests := []struct {
		distro   Distro
		expected string
	}{
		{DistroDebian, "Debian"},
		{DistroUbuntu, "Ubuntu"},
		{DistroArch, "Arch"},
		{DistroFedora, "Fedora"},
		{DistroNone, ""},
		{DistroOther, ""},
	}
	for _, tt := range tests {
		if got := tt.distro.String(); got != tt.expected {
			t.Errorf("Distro(%d).String() = %q, want %q", tt.distro, got, tt.expected)
		}
	}
}

// --- HasCommand ---

func TestHasCommandExists(t *testing.T) {
	if !HasCommand("go") {
		t.Error("expected 'go' to be found")
	}
}

func TestHasCommandNotExists(t *testing.T) {
	if HasCommand("nonexistent-binary-12345-xyz") {
		t.Error("expected nonexistent binary to not be found")
	}
}

func TestHasCommandEmptyString(t *testing.T) {
	if HasCommand("") {
		t.Error("expected empty string to not be found")
	}
}

func TestHasCommandCommonTools(t *testing.T) {
	// "git" is typically available in CI/test environments
	if !HasCommand("git") {
		t.Log("git not found (may be expected in minimal environments)")
	}
}

// --- OS type values ---

func TestOSTypeValues(t *testing.T) {
	// Ensure enum values are distinct
	values := []OS{MacOS, Linux, WSL, Windows, Unknown}
	seen := make(map[OS]bool)
	for _, v := range values {
		if seen[v] {
			t.Errorf("duplicate OS value: %d", v)
		}
		seen[v] = true
	}
}

func TestDistroTypeValues(t *testing.T) {
	values := []Distro{DistroNone, DistroDebian, DistroUbuntu, DistroArch, DistroFedora, DistroOther}
	seen := make(map[Distro]bool)
	for _, v := range values {
		if seen[v] {
			t.Errorf("duplicate Distro value: %d", v)
		}
		seen[v] = true
	}
}
