package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aymenhmaidiwastaken/devboot/internal/config"
	"github.com/aymenhmaidiwastaken/devboot/internal/dotfiles"
	"github.com/aymenhmaidiwastaken/devboot/internal/editor"
	gitcfg "github.com/aymenhmaidiwastaken/devboot/internal/git"
	"github.com/aymenhmaidiwastaken/devboot/internal/platform"
	"github.com/aymenhmaidiwastaken/devboot/internal/shell"
	"github.com/aymenhmaidiwastaken/devboot/internal/state"
	"github.com/aymenhmaidiwastaken/devboot/internal/tools"
	"github.com/aymenhmaidiwastaken/devboot/internal/tui"
	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
	"github.com/spf13/cobra"
)

var (
	onlySection string
	noTUI       bool
)

var applyCmd = &cobra.Command{
	Use:   "apply [config.yaml]",
	Short: "Apply configuration to set up your dev environment",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runApply,
}

func init() {
	applyCmd.Flags().StringVar(&onlySection, "only", "", "apply only a specific section (tools, shell, git, vscode, neovim, jetbrains, dotfiles)")
	applyCmd.Flags().BoolVar(&noTUI, "no-tui", false, "disable interactive TUI (use plain output)")
	rootCmd.AddCommand(applyCmd)
}

func runApply(cmd *cobra.Command, args []string) error {
	configPath := config.DefaultConfigPath()
	if len(args) > 0 {
		configPath = args[0]
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	p := platform.Detect()

	// Load state
	st, _ := state.Load()
	st.MarkApply()
	defer st.Save()

	if !noTUI {
		return runApplyWithTUI(cfg, p, st)
	}

	return runApplyPlain(cfg, p, st)
}

func runApplyWithTUI(cfg *config.Config, p platform.Platform, st *state.Store) error {
	return tui.RunApplyTUI(func(m *tui.ApplyModel) error {
		sections := buildSections(cfg, p, st)

		if onlySection != "" {
			fn, ok := sections[onlySection]
			if !ok {
				return fmt.Errorf("unknown section: %q", onlySection)
			}
			return fn(m)
		}

		order := []string{"tools", "shell", "git", "vscode", "neovim", "jetbrains", "dotfiles"}
		for _, name := range order {
			if fn, ok := sections[name]; ok {
				if err := fn(m); err != nil {
					// Continue with other sections
					_ = err
				}
			}
		}
		return nil
	})
}

func runApplyPlain(cfg *config.Config, p platform.Platform, st *state.Store) error {
	fmt.Printf("\n  devboot — %s (%s)\n", p.OS, p.Arch)

	sections := buildSectionsPlain(cfg, p, st)

	if onlySection != "" {
		fn, ok := sections[onlySection]
		if !ok {
			return fmt.Errorf("unknown section: %q (available: tools, shell, git, vscode, neovim, jetbrains, dotfiles)", onlySection)
		}
		return fn()
	}

	order := []string{"tools", "shell", "git", "vscode", "neovim", "jetbrains", "dotfiles"}
	for _, name := range order {
		if fn, ok := sections[name]; ok {
			if err := fn(); err != nil {
				ui.Warn(fmt.Sprintf("%s section had errors: %v", name, err))
			}
		}
	}

	fmt.Println()
	ui.Success("devboot apply complete!")
	fmt.Println()
	return nil
}

func buildSections(cfg *config.Config, p platform.Platform, st *state.Store) map[string]func(m *tui.ApplyModel) error {
	sections := make(map[string]func(m *tui.ApplyModel) error)

	if len(cfg.Tools) > 0 {
		sections["tools"] = func(m *tui.ApplyModel) error {
			return applyToolsTUI(cfg, p, st, m)
		}
	}
	if cfg.Shell.Type != "" || len(cfg.Shell.Aliases) > 0 || len(cfg.Shell.Env) > 0 || len(cfg.Shell.Plugins) > 0 {
		sections["shell"] = func(m *tui.ApplyModel) error {
			return applyShellTUI(cfg, st, m)
		}
	}
	if cfg.Git.UserName != "" || cfg.Git.UserEmail != "" || len(cfg.Git.Aliases) > 0 || cfg.Git.SSHKey {
		sections["git"] = func(m *tui.ApplyModel) error {
			return applyGitTUI(cfg, st, m)
		}
	}
	if len(cfg.VSCode.Extensions) > 0 || len(cfg.VSCode.Settings) > 0 {
		sections["vscode"] = func(m *tui.ApplyModel) error {
			return applyVSCodePlainWrap(cfg)
		}
	}
	if cfg.Neovim.ConfigRepo != "" {
		sections["neovim"] = func(m *tui.ApplyModel) error {
			return editor.NeovimApply(cfg.Neovim.ConfigRepo)
		}
	}
	if len(cfg.JetBrains.Plugins) > 0 {
		sections["jetbrains"] = func(m *tui.ApplyModel) error {
			return editor.JetBrainsApply(cfg.JetBrains.Plugins)
		}
	}
	if cfg.Dotfiles.Repo != "" {
		sections["dotfiles"] = func(m *tui.ApplyModel) error {
			return applyDotfilesAll(cfg, st)
		}
	}

	return sections
}

func buildSectionsPlain(cfg *config.Config, p platform.Platform, st *state.Store) map[string]func() error {
	sections := make(map[string]func() error)

	if len(cfg.Tools) > 0 {
		sections["tools"] = func() error { return applyTools(cfg, p) }
	}
	if cfg.Shell.Type != "" || len(cfg.Shell.Aliases) > 0 || len(cfg.Shell.Env) > 0 || len(cfg.Shell.Plugins) > 0 {
		sections["shell"] = func() error { return applyShell(cfg) }
	}
	if cfg.Git.UserName != "" || cfg.Git.UserEmail != "" || len(cfg.Git.Aliases) > 0 || cfg.Git.PullRebase != nil || cfg.Git.SSHKey {
		sections["git"] = func() error { return applyGit(cfg) }
	}
	if len(cfg.VSCode.Extensions) > 0 || len(cfg.VSCode.Settings) > 0 {
		sections["vscode"] = func() error { return applyVSCodePlainWrap(cfg) }
	}
	if cfg.Neovim.ConfigRepo != "" {
		sections["neovim"] = func() error { return editor.NeovimApply(cfg.Neovim.ConfigRepo) }
	}
	if len(cfg.JetBrains.Plugins) > 0 {
		sections["jetbrains"] = func() error { return editor.JetBrainsApply(cfg.JetBrains.Plugins) }
	}
	if cfg.Dotfiles.Repo != "" {
		sections["dotfiles"] = func() error { return applyDotfilesAll(cfg, nil) }
	}

	return sections
}

func applyToolsTUI(cfg *config.Config, p platform.Platform, st *state.Store, m *tui.ApplyModel) error {
	installer, err := tools.NewInstaller(p)
	if err != nil {
		return err
	}
	installer.State = st
	return installer.InstallAll(cfg.Tools)
}

func applyShellTUI(cfg *config.Config, st *state.Store, m *tui.ApplyModel) error {
	return applyShell(cfg)
}

func applyGitTUI(cfg *config.Config, st *state.Store, m *tui.ApplyModel) error {
	return applyGit(cfg)
}

func applyTools(cfg *config.Config, p platform.Platform) error {
	if len(cfg.Tools) == 0 {
		return nil
	}
	installer, err := tools.NewInstaller(p)
	if err != nil {
		return err
	}
	return installer.InstallAll(cfg.Tools)
}

func applyShell(cfg *config.Config) error {
	shellType := shell.ShellType(cfg.Shell.Type)
	if shellType == "" {
		shellType = shell.Detect()
	}

	if err := shell.InstallPlugins(shellType, cfg.Shell.Plugins); err != nil {
		ui.Warn(fmt.Sprintf("plugin installation: %v", err))
	}
	if err := shell.SetAliases(shellType, cfg.Shell.Aliases); err != nil {
		return err
	}
	if err := shell.SetEnvVars(shellType, cfg.Shell.Env); err != nil {
		return err
	}
	return nil
}

func applyGit(cfg *config.Config) error {
	if err := gitcfg.Configure(cfg.Git); err != nil {
		return err
	}
	if cfg.Git.SSHKey {
		if err := gitcfg.EnsureSSHKey(cfg.Git.UserEmail); err != nil {
			ui.Warn(fmt.Sprintf("SSH key: %v", err))
		}
	}
	return nil
}

func applyVSCodePlainWrap(cfg *config.Config) error {
	return editor.VSCodeApply(cfg.VSCode.Extensions, cfg.VSCode.Settings)
}

func applyDotfilesAll(cfg *config.Config, st *state.Store) error {
	if err := dotfiles.Sync(cfg.Dotfiles.Repo, cfg.Dotfiles.Mappings); err != nil {
		return err
	}

	if len(cfg.Dotfiles.Templates) > 0 {
		vars := dotfiles.DefaultTemplateVars(cfg.Git.UserName, cfg.Git.UserEmail, platform.Detect().OS.String())
		for src, dst := range cfg.Dotfiles.Templates {
			home, _ := os.UserHomeDir()
			srcPath := src
			if !filepath.IsAbs(srcPath) {
				srcPath = filepath.Join(home, dotfiles.DotfilesDirPath(), src)
			}
			dstPath := dotfiles.ExpandHome(dst)
			if err := dotfiles.RenderFile(srcPath, dstPath, vars); err != nil {
				ui.Warn(fmt.Sprintf("template %s: %v", src, err))
			}
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
