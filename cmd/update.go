package cmd

import (
	"fmt"

	"github.com/aymenhmaidiwastaken/devboot/internal/platform"
	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all installed tools",
	RunE:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	p := platform.Detect()
	fmt.Printf("\n  devboot update — %s\n", p.OS)

	ui.Section("Updating package manager")

	switch p.OS {
	case platform.MacOS:
		if platform.HasCommand("brew") {
			ui.Installing("brew update && brew upgrade...")
			if err := runShell("brew update && brew upgrade"); err != nil {
				ui.Warn(fmt.Sprintf("brew upgrade: %v", err))
			} else {
				ui.Success("Homebrew packages updated")
			}
		}
	case platform.Linux, platform.WSL:
		if platform.HasCommand("apt-get") {
			ui.Installing("apt update && apt upgrade...")
			if err := runShell("sudo apt-get update -qq && sudo apt-get upgrade -y -qq"); err != nil {
				ui.Warn(fmt.Sprintf("apt upgrade: %v", err))
			} else {
				ui.Success("APT packages updated")
			}
		}
		if platform.HasCommand("pacman") {
			ui.Installing("pacman -Syu...")
			if err := runShell("sudo pacman -Syu --noconfirm"); err != nil {
				ui.Warn(fmt.Sprintf("pacman upgrade: %v", err))
			} else {
				ui.Success("Pacman packages updated")
			}
		}
	}

	fmt.Println()
	ui.Success("update complete!")
	fmt.Println()
	return nil
}
