package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aymenhmaidiwastaken/devboot/internal/config"
	"github.com/aymenhmaidiwastaken/devboot/internal/deps"
	"github.com/aymenhmaidiwastaken/devboot/internal/editor"
	"github.com/aymenhmaidiwastaken/devboot/internal/platform"
	"github.com/aymenhmaidiwastaken/devboot/internal/state"
	"github.com/aymenhmaidiwastaken/devboot/internal/tools"
	"github.com/aymenhmaidiwastaken/devboot/internal/tui"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [config.yaml]",
	Short: "Show what's installed vs what's configured",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	configPath := config.DefaultConfigPath()
	if len(args) > 0 {
		configPath = args[0]
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	p := platform.Detect()
	st, _ := state.Load()

	var sections []tui.StatusSection

	// Platform info
	sections = append(sections, tui.StatusSection{
		Name: fmt.Sprintf("System — %s (%s)", p.OS, p.Arch),
		Items: []tui.StatusItem{
			{Label: fmt.Sprintf("managed tools: %d", len(st.ManagedTools)), Status: "info"},
			{Label: fmt.Sprintf("recorded actions: %d", len(st.Actions)), Status: "info"},
		},
	})

	// Tools
	if len(cfg.Tools) > 0 {
		var items []tui.StatusItem
		installed, missing := 0, 0
		for _, spec := range cfg.Tools {
			name, _ := tools.ParseTool(spec)
			binName := deps.GetBinName(name)
			if platform.HasCommand(binName) {
				ver := getToolVersion(binName)
				label := name
				if ver != "" {
					label = fmt.Sprintf("%s (%s)", name, ver)
				}
				if st.IsManagedTool(name) {
					label += " [managed]"
				}
				items = append(items, tui.StatusItem{Label: label, Status: "ok"})
				installed++
			} else {
				items = append(items, tui.StatusItem{Label: fmt.Sprintf("%s (missing)", name), Status: "missing"})
				missing++
			}
		}
		sectionName := fmt.Sprintf("Tools — %d installed, %d missing", installed, missing)
		sections = append(sections, tui.StatusSection{Name: sectionName, Items: items})
	}

	// Shell
	if cfg.Shell.Type != "" || len(cfg.Shell.Aliases) > 0 {
		items := []tui.StatusItem{
			{Label: fmt.Sprintf("type: %s", cfg.Shell.Type), Status: "info"},
			{Label: fmt.Sprintf("aliases: %d", len(cfg.Shell.Aliases)), Status: "info"},
			{Label: fmt.Sprintf("plugins: %d", len(cfg.Shell.Plugins)), Status: "info"},
			{Label: fmt.Sprintf("env vars: %d", len(cfg.Shell.Env)), Status: "info"},
		}
		sections = append(sections, tui.StatusSection{Name: "Shell", Items: items})
	}

	// Git
	if cfg.Git.UserName != "" || cfg.Git.UserEmail != "" {
		var items []tui.StatusItem
		if cfg.Git.UserName != "" {
			current := gitConfigGet("user.name")
			if current == cfg.Git.UserName {
				items = append(items, tui.StatusItem{Label: fmt.Sprintf("user.name: %s", current), Status: "ok"})
			} else {
				items = append(items, tui.StatusItem{Label: fmt.Sprintf("user.name: want %q, have %q", cfg.Git.UserName, current), Status: "warn"})
			}
		}
		if cfg.Git.UserEmail != "" {
			current := gitConfigGet("user.email")
			if current == cfg.Git.UserEmail {
				items = append(items, tui.StatusItem{Label: fmt.Sprintf("user.email: %s", current), Status: "ok"})
			} else {
				items = append(items, tui.StatusItem{Label: fmt.Sprintf("user.email: want %q, have %q", cfg.Git.UserEmail, current), Status: "warn"})
			}
		}
		items = append(items, tui.StatusItem{Label: fmt.Sprintf("aliases: %d configured", len(cfg.Git.Aliases)), Status: "info"})
		sections = append(sections, tui.StatusSection{Name: "Git", Items: items})
	}

	// VS Code
	if len(cfg.VSCode.Extensions) > 0 || len(cfg.VSCode.Settings) > 0 {
		var items []tui.StatusItem
		installed, codePath := editor.VSCodeStatus()
		if codePath != "" {
			items = append(items, tui.StatusItem{Label: fmt.Sprintf("CLI: %s", codePath), Status: "ok"})
			items = append(items, tui.StatusItem{Label: fmt.Sprintf("extensions installed: %d", len(installed)), Status: "info"})
		} else {
			items = append(items, tui.StatusItem{Label: "VS Code CLI not found", Status: "missing"})
		}
		items = append(items, tui.StatusItem{Label: fmt.Sprintf("extensions configured: %d", len(cfg.VSCode.Extensions)), Status: "info"})
		sections = append(sections, tui.StatusSection{Name: "VS Code", Items: items})
	}

	// Neovim
	if cfg.Neovim.ConfigRepo != "" {
		var items []tui.StatusItem
		if platform.HasCommand("nvim") {
			items = append(items, tui.StatusItem{Label: "nvim installed", Status: "ok"})
		} else {
			items = append(items, tui.StatusItem{Label: "nvim not found", Status: "missing"})
		}
		home, _ := os.UserHomeDir()
		nvimDir := filepath.Join(home, ".config", "nvim")
		if fileExists(nvimDir) {
			items = append(items, tui.StatusItem{Label: "config present", Status: "ok"})
		} else {
			items = append(items, tui.StatusItem{Label: "config missing", Status: "missing"})
		}
		sections = append(sections, tui.StatusSection{Name: "Neovim", Items: items})
	}

	// Dotfiles
	if cfg.Dotfiles.Repo != "" {
		items := []tui.StatusItem{
			{Label: fmt.Sprintf("repo: %s", cfg.Dotfiles.Repo), Status: "info"},
			{Label: fmt.Sprintf("mappings: %d", len(cfg.Dotfiles.Mappings)), Status: "info"},
		}
		sections = append(sections, tui.StatusSection{Name: "Dotfiles", Items: items})
	}

	fmt.Println(tui.RenderStatusDashboard(sections))
	return nil
}
