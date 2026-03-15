package cmd

import (
	"fmt"
	"os"

	"github.com/aymenhmaidiwastaken/devboot/internal/tui"
	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
	"github.com/spf13/cobra"
)

var initPlain bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a starter devboot.yaml (interactive wizard)",
	RunE:  runInit,
}

func init() {
	initCmd.Flags().BoolVar(&initPlain, "plain", false, "generate a plain starter config without wizard")
	rootCmd.AddCommand(initCmd)
}

const starterConfig = `# devboot.yaml — your dev environment in one file
# Docs: https://github.com/aymenhmaidiwastaken/devboot
# Run:  devboot apply

tools:
  - git
  - node
  - python
  - ripgrep
  - fzf
  - jq

shell:
  type: zsh
  plugins:
    - zsh-autosuggestions
    - zsh-syntax-highlighting
  aliases:
    g: git
    ll: ls -la
    dc: docker compose
  env:
    EDITOR: vim

git:
  user.name: "Your Name"
  user.email: "you@example.com"
  init.defaultBranch: main
  pull.rebase: true
  aliases:
    co: checkout
    br: branch
    st: status
    lg: "log --oneline --graph --decorate"

# vscode:
#   extensions:
#     - ms-python.python
#     - golang.go
#     - esbenp.prettier-vscode
#   settings:
#     editor.fontSize: 14
#     editor.tabSize: 2
#     editor.formatOnSave: true

# neovim:
#   configRepo: https://github.com/aymenhmaidiwastaken/nvim-config.git

# dotfiles:
#   repo: https://github.com/aymenhmaidiwastaken/dotfiles.git
#   mappings:
#     .vimrc: ~/.vimrc
#     .tmux.conf: ~/.tmux.conf
#     starship.toml: ~/.config/starship.toml
`

func runInit(cmd *cobra.Command, args []string) error {
	if initPlain {
		return runInitPlain()
	}
	return tui.RunInitWizard()
}

func runInitPlain() error {
	path := "devboot.yaml"

	if fileExists(path) {
		return fmt.Errorf("%s already exists — remove it first or edit it directly", path)
	}

	if err := os.WriteFile(path, []byte(starterConfig), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	ui.Success(fmt.Sprintf("Created %s", path))
	ui.Info("Edit the file, then run: devboot apply")
	return nil
}
