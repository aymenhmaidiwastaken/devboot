package cmd

import (
	"fmt"

	"github.com/aymenhmaidiwastaken/devboot/internal/config"
	"github.com/aymenhmaidiwastaken/devboot/internal/platform"
	"github.com/aymenhmaidiwastaken/devboot/internal/shell"
	"github.com/aymenhmaidiwastaken/devboot/internal/tools"
	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff [config.yaml]",
	Short: "Show what would change if you run apply",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	configPath := config.DefaultConfigPath()
	if len(args) > 0 {
		configPath = args[0]
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	fmt.Printf("\n  devboot diff — preview changes\n")

	changes := 0

	if len(cfg.Tools) > 0 {
		ui.Section("Tools")
		for _, spec := range cfg.Tools {
			name, _ := tools.ParseTool(spec)
			if platform.HasCommand(name) {
				ui.Skip(name)
			} else {
				ui.Installing(fmt.Sprintf("would install %s", name))
				changes++
			}
		}
	}

	if len(cfg.Shell.Aliases) > 0 || len(cfg.Shell.Plugins) > 0 || len(cfg.Shell.Env) > 0 {
		ui.Section("Shell")
		shellType := shell.ShellType(cfg.Shell.Type)
		if shellType == "" {
			shellType = shell.Detect()
		}
		ui.Info(fmt.Sprintf("target: %s (%s)", shellType, shell.ConfigFile(shellType)))
		if len(cfg.Shell.Aliases) > 0 {
			ui.Info(fmt.Sprintf("%d aliases to configure", len(cfg.Shell.Aliases)))
			changes += len(cfg.Shell.Aliases)
		}
		if len(cfg.Shell.Plugins) > 0 {
			ui.Info(fmt.Sprintf("%d plugins to install", len(cfg.Shell.Plugins)))
			changes += len(cfg.Shell.Plugins)
		}
		if len(cfg.Shell.Env) > 0 {
			ui.Info(fmt.Sprintf("%d env vars to set", len(cfg.Shell.Env)))
			changes += len(cfg.Shell.Env)
		}
	}

	if cfg.Git.UserName != "" || cfg.Git.UserEmail != "" || len(cfg.Git.Aliases) > 0 {
		ui.Section("Git")
		if cfg.Git.UserName != "" {
			ui.Info(fmt.Sprintf("user.name → %s", cfg.Git.UserName))
			changes++
		}
		if cfg.Git.UserEmail != "" {
			ui.Info(fmt.Sprintf("user.email → %s", cfg.Git.UserEmail))
			changes++
		}
		if len(cfg.Git.Aliases) > 0 {
			ui.Info(fmt.Sprintf("%d git aliases to configure", len(cfg.Git.Aliases)))
			changes += len(cfg.Git.Aliases)
		}
	}

	if len(cfg.VSCode.Extensions) > 0 || len(cfg.VSCode.Settings) > 0 {
		ui.Section("VS Code")
		if len(cfg.VSCode.Extensions) > 0 {
			ui.Info(fmt.Sprintf("%d extensions to install", len(cfg.VSCode.Extensions)))
			changes += len(cfg.VSCode.Extensions)
		}
		if len(cfg.VSCode.Settings) > 0 {
			ui.Info(fmt.Sprintf("%d settings to apply", len(cfg.VSCode.Settings)))
			changes += len(cfg.VSCode.Settings)
		}
	}

	if cfg.Neovim.ConfigRepo != "" {
		ui.Section("Neovim")
		ui.Info(fmt.Sprintf("config repo: %s", cfg.Neovim.ConfigRepo))
		changes++
	}

	if len(cfg.JetBrains.Plugins) > 0 {
		ui.Section("JetBrains")
		ui.Info(fmt.Sprintf("%d plugins to install", len(cfg.JetBrains.Plugins)))
		changes += len(cfg.JetBrains.Plugins)
	}

	if cfg.Dotfiles.Repo != "" {
		ui.Section("Dotfiles")
		ui.Info(fmt.Sprintf("repo: %s", cfg.Dotfiles.Repo))
		ui.Info(fmt.Sprintf("%d file mappings", len(cfg.Dotfiles.Mappings)))
		changes += len(cfg.Dotfiles.Mappings) + 1
	}

	fmt.Println()
	if changes == 0 {
		ui.Success("everything is up to date!")
	} else {
		ui.Info(fmt.Sprintf("%d change(s) would be applied", changes))
	}
	fmt.Println()

	return nil
}
