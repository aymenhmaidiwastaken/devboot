package tools

import (
	"fmt"
	"testing"

	"github.com/aymenhmaidiwastaken/devboot/internal/platform"
)

// --- ParseTool ---

func TestParseTool(t *testing.T) {
	tests := []struct {
		spec    string
		name    string
		version string
	}{
		{"git", "git", ""},
		{"node@22", "node", "22"},
		{"python@3.12", "python", "3.12"},
		{"rust", "rust", ""},
		{"go@1.22.1", "go", "1.22.1"},
	}

	for _, tt := range tests {
		name, version := ParseTool(tt.spec)
		if name != tt.name {
			t.Errorf("ParseTool(%q) name = %q, want %q", tt.spec, name, tt.name)
		}
		if version != tt.version {
			t.Errorf("ParseTool(%q) version = %q, want %q", tt.spec, version, tt.version)
		}
	}
}

func TestParseToolEdgeCases(t *testing.T) {
	// Only name, no @
	name, version := ParseTool("docker")
	if name != "docker" || version != "" {
		t.Errorf("expected docker/'', got %q/%q", name, version)
	}

	// Multiple @
	name, version = ParseTool("tool@1@2")
	if name != "tool" || version != "1@2" {
		t.Errorf("expected tool/1@2, got %q/%q", name, version)
	}

	// Empty string
	name, version = ParseTool("")
	if name != "" || version != "" {
		t.Errorf("expected empty/empty, got %q/%q", name, version)
	}

	// Just @
	name, version = ParseTool("@")
	if name != "" || version != "" {
		t.Errorf("expected empty/empty for @, got %q/%q", name, version)
	}

	// Version only
	name, version = ParseTool("@22")
	if name != "" || version != "22" {
		t.Errorf("expected empty/22 for @22, got %q/%q", name, version)
	}
}

// --- ResolvePackage ---

func TestResolvePackage(t *testing.T) {
	tests := []struct {
		name     string
		os       platform.OS
		expected string
	}{
		{"node", platform.MacOS, "node"},
		{"node", platform.Linux, "nodejs"},
		{"node", platform.WSL, "nodejs"},
		{"nodejs", platform.MacOS, "node"},
		{"python", platform.MacOS, "python@3"},
		{"python", platform.Linux, "python3"},
		{"python", platform.WSL, "python3"},
		{"go", platform.MacOS, "go"},
		{"go", platform.Linux, "golang"},
		{"golang", platform.Linux, "golang"},
		{"rust", platform.MacOS, "rustup"},
		{"rust", platform.Linux, "rustup"},
		{"docker", platform.MacOS, "docker"},
		{"docker", platform.Linux, "docker.io"},
		{"docker", platform.WSL, "docker.io"},
		{"kubectl", platform.MacOS, "kubectl"},
		{"ripgrep", platform.Linux, "ripgrep"},
		{"fzf", platform.MacOS, "fzf"},
		{"neovim", platform.MacOS, "neovim"},
		{"nvim", platform.Linux, "neovim"},
		{"tmux", platform.MacOS, "tmux"},
		{"jq", platform.Linux, "jq"},
		{"curl", platform.WSL, "curl"},
		{"wget", platform.MacOS, "wget"},
	}

	for _, tt := range tests {
		p := platform.Platform{OS: tt.os}
		got := ResolvePackage(tt.name, p)
		if got != tt.expected {
			t.Errorf("ResolvePackage(%q, %v) = %q, want %q", tt.name, tt.os, got, tt.expected)
		}
	}
}

func TestResolvePackageUnknown(t *testing.T) {
	// Unknown tools should fall through to name as-is
	unknowns := []string{"unknowntool", "my-custom-pkg", "zzz-tool"}
	for _, name := range unknowns {
		for _, os := range []platform.OS{platform.MacOS, platform.Linux, platform.WSL} {
			p := platform.Platform{OS: os}
			got := ResolvePackage(name, p)
			if got != name {
				t.Errorf("ResolvePackage(%q, %v) = %q, want %q (fallthrough)", name, os, got, name)
			}
		}
	}
}

func TestResolvePackageAllMappingsHaveThreeOS(t *testing.T) {
	// Every tool in the mapping should have entries for MacOS, Linux, and WSL
	for name, mapping := range toolMapping {
		for _, os := range []platform.OS{platform.MacOS, platform.Linux, platform.WSL} {
			if _, ok := mapping[os]; !ok {
				t.Errorf("toolMapping[%q] missing entry for OS %v", name, os)
			}
		}
	}
}

// --- isToolInstalled ---

func TestIsToolInstalled(t *testing.T) {
	// "go" should be installed in test environment
	if !isToolInstalled("go") {
		t.Error("expected 'go' to be installed")
	}

	// non-existent binary
	if isToolInstalled("nonexistent-binary-xyz-12345") {
		t.Error("expected nonexistent binary to not be installed")
	}
}

// --- NewInstaller ---

func TestNewInstallerMacOS(t *testing.T) {
	p := platform.Platform{OS: platform.MacOS}
	inst, err := NewInstaller(p)
	if err != nil {
		t.Fatalf("NewInstaller failed: %v", err)
	}
	if inst.backend.Name() != "Homebrew" {
		t.Errorf("expected Homebrew backend, got %s", inst.backend.Name())
	}
}

func TestNewInstallerLinuxDefault(t *testing.T) {
	p := platform.Platform{OS: platform.Linux, Distro: platform.DistroUbuntu}
	inst, err := NewInstaller(p)
	if err != nil {
		t.Fatalf("NewInstaller failed: %v", err)
	}
	if inst.backend.Name() != "apt" {
		t.Errorf("expected apt backend, got %s", inst.backend.Name())
	}
}

func TestNewInstallerArch(t *testing.T) {
	p := platform.Platform{OS: platform.Linux, Distro: platform.DistroArch}
	inst, err := NewInstaller(p)
	if err != nil {
		t.Fatalf("NewInstaller failed: %v", err)
	}
	if inst.backend.Name() != "pacman" {
		t.Errorf("expected pacman backend, got %s", inst.backend.Name())
	}
}

func TestNewInstallerWSL(t *testing.T) {
	p := platform.Platform{OS: platform.WSL, Distro: platform.DistroDebian}
	inst, err := NewInstaller(p)
	if err != nil {
		t.Fatalf("NewInstaller failed: %v", err)
	}
	if inst.backend.Name() != "apt" {
		t.Errorf("expected apt backend for WSL, got %s", inst.backend.Name())
	}
}

func TestNewInstallerUnsupported(t *testing.T) {
	p := platform.Platform{OS: platform.Windows}
	_, err := NewInstaller(p)
	if err == nil {
		t.Error("expected error for unsupported platform")
	}
}

func TestNewInstallerHasState(t *testing.T) {
	p := platform.Platform{OS: platform.MacOS}
	inst, _ := NewInstaller(p)
	if inst.State == nil {
		t.Error("expected state to be initialized")
	}
}

// --- Mock backend for InstallAll/UninstallAll tests ---

type mockBackend struct {
	name       string
	installed  map[string]bool
	installErr map[string]error
}

func newMockBackend() *mockBackend {
	return &mockBackend{
		name:       "mock",
		installed:  make(map[string]bool),
		installErr: make(map[string]error),
	}
}

func (m *mockBackend) Name() string        { return m.name }
func (m *mockBackend) IsAvailable() bool    { return true }
func (m *mockBackend) EnsureReady() error   { return nil }
func (m *mockBackend) IsInstalled(pkg string) bool { return m.installed[pkg] }
func (m *mockBackend) Install(pkg string) error {
	if err, ok := m.installErr[pkg]; ok {
		return err
	}
	m.installed[pkg] = true
	return nil
}
func (m *mockBackend) Uninstall(pkg string) error {
	delete(m.installed, pkg)
	return nil
}

func TestMockBackendInstallUninstall(t *testing.T) {
	mock := newMockBackend()

	if mock.IsInstalled("git") {
		t.Error("git should not be installed initially")
	}
	if err := mock.Install("git"); err != nil {
		t.Fatalf("Install failed: %v", err)
	}
	if !mock.IsInstalled("git") {
		t.Error("git should be installed after Install()")
	}
	if err := mock.Uninstall("git"); err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}
	if mock.IsInstalled("git") {
		t.Error("git should not be installed after Uninstall()")
	}
}

func TestMockBackendInstallError(t *testing.T) {
	mock := newMockBackend()
	mock.installErr["badpkg"] = fmt.Errorf("download failed")

	err := mock.Install("badpkg")
	if err == nil {
		t.Error("expected error for badpkg")
	}
	if mock.IsInstalled("badpkg") {
		t.Error("badpkg should not be marked installed on error")
	}
}

func TestMockBackendProperties(t *testing.T) {
	mock := newMockBackend()
	if mock.Name() != "mock" {
		t.Errorf("expected name 'mock', got %q", mock.Name())
	}
	if !mock.IsAvailable() {
		t.Error("mock should always be available")
	}
	if err := mock.EnsureReady(); err != nil {
		t.Errorf("EnsureReady should not fail: %v", err)
	}
}
