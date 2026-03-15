package tools

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/aymenhmaidiwastaken/devboot/internal/deps"
	"github.com/aymenhmaidiwastaken/devboot/internal/platform"
	"github.com/aymenhmaidiwastaken/devboot/internal/state"
	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
)

// Installer handles installing tools on the current platform.
type Installer struct {
	platform platform.Platform
	backend  Backend
	State    *state.Store
}

// Backend is the interface for platform-specific package managers.
type Backend interface {
	Name() string
	IsAvailable() bool
	Install(pkg string) error
	Uninstall(pkg string) error
	IsInstalled(pkg string) bool
	EnsureReady() error
}

// NewInstaller creates a new installer for the detected platform.
func NewInstaller(p platform.Platform) (*Installer, error) {
	var backend Backend

	switch p.OS {
	case platform.MacOS:
		backend = &BrewBackend{}
	case platform.Linux, platform.WSL:
		switch p.Distro {
		case platform.DistroArch:
			backend = &PacmanBackend{}
		default:
			backend = &AptBackend{}
		}
	default:
		return nil, fmt.Errorf("unsupported platform: %s", p.OS)
	}

	st, _ := state.Load()

	return &Installer{platform: p, backend: backend, State: st}, nil
}

// ParseTool splits a tool spec like "node@22" into name and version.
func ParseTool(spec string) (name, version string) {
	parts := strings.SplitN(spec, "@", 2)
	name = parts[0]
	if len(parts) > 1 {
		version = parts[1]
	}
	return
}

// InstallAll installs a list of tools concurrently, respecting dependency order.
func (i *Installer) InstallAll(tools []string) error {
	if err := i.backend.EnsureReady(); err != nil {
		return fmt.Errorf("package manager not ready: %w", err)
	}

	// Resolve dependency order
	ordered, err := deps.Resolve(tools)
	if err != nil {
		ui.Warn(fmt.Sprintf("dependency resolution: %v (installing in original order)", err))
		ordered = tools
	}

	// Check for conflicts
	if warnings := deps.Conflicts(ordered); len(warnings) > 0 {
		for _, w := range warnings {
			ui.Warn(w)
		}
	}

	ui.Section("Installing tools via " + i.backend.Name())

	// Split into phases: dependencies first (sequential), then the rest (parallel)
	var depPhase, mainPhase []string
	requestedSet := make(map[string]bool)
	for _, t := range tools {
		name, _ := ParseTool(t)
		requestedSet[name] = true
	}
	for _, t := range ordered {
		name, _ := ParseTool(t)
		if !requestedSet[name] {
			depPhase = append(depPhase, t) // injected dependency
		} else {
			mainPhase = append(mainPhase, t)
		}
	}

	// Install injected dependencies first (sequential)
	for _, spec := range depPhase {
		i.installOne(spec)
	}

	// Install main tools (parallel)
	var wg sync.WaitGroup
	errCh := make(chan error, len(mainPhase))
	sem := make(chan struct{}, 4)

	for _, tool := range mainPhase {
		wg.Add(1)
		go func(spec string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if err := i.installOne(spec); err != nil {
				errCh <- err
			}
		}(tool)
	}

	wg.Wait()
	close(errCh)

	// Save state
	if i.State != nil {
		i.State.Save()
	}

	var errs []string
	for err := range errCh {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		return fmt.Errorf("%d tool(s) failed to install:\n  %s", len(errs), strings.Join(errs, "\n  "))
	}
	return nil
}

func (i *Installer) installOne(spec string) error {
	name, _ := ParseTool(spec)
	binName := deps.GetBinName(name)
	pkg := ResolvePackage(name, i.platform)

	// Check if already installed
	if isToolInstalled(binName) {
		ui.Skip(name)
		if i.State != nil {
			i.State.RecordTool(name, pkg, "skipped")
		}
		return nil
	}

	ui.Installing(fmt.Sprintf("installing %s...", name))
	if err := i.backend.Install(pkg); err != nil {
		ui.Fail(fmt.Sprintf("%s: %v", name, err))
		if i.State != nil {
			i.State.RecordTool(name, pkg, "failed")
		}
		return fmt.Errorf("installing %s: %w", name, err)
	}

	// Run post-install hooks
	if hooks := deps.GetPostInstall(name); len(hooks) > 0 {
		for _, hook := range hooks {
			ui.Info(fmt.Sprintf("  post-install: %s", hook))
			cmd := exec.Command("sh", "-c", hook)
			cmd.Run() // best-effort
		}
	}

	ui.Success(name)
	if i.State != nil {
		i.State.RecordTool(name, pkg, "ok")
	}
	return nil
}

// UninstallAll removes tools that were installed by devboot.
func (i *Installer) UninstallAll(tools []string) error {
	if err := i.backend.EnsureReady(); err != nil {
		return fmt.Errorf("package manager not ready: %w", err)
	}

	ui.Section("Uninstalling tools via " + i.backend.Name())

	for _, spec := range tools {
		name, _ := ParseTool(spec)
		pkg := ResolvePackage(name, i.platform)

		if !i.backend.IsInstalled(pkg) {
			ui.Skip(fmt.Sprintf("%s (not installed)", name))
			continue
		}

		ui.Installing(fmt.Sprintf("removing %s...", name))
		if err := i.backend.Uninstall(pkg); err != nil {
			ui.Fail(fmt.Sprintf("%s: %v", name, err))
			continue
		}
		ui.Success(fmt.Sprintf("%s removed", name))

		if i.State != nil {
			i.State.RemoveTool(name)
		}
	}

	if i.State != nil {
		i.State.Save()
	}

	return nil
}

func isToolInstalled(binName string) bool {
	_, err := exec.LookPath(binName)
	return err == nil
}
