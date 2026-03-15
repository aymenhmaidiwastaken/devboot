package cmd

import (
	"github.com/aymenhmaidiwastaken/devboot/internal/dotfiles"
	"github.com/spf13/cobra"
)

var dotfilesCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "Manage dotfiles",
}

var dotfilesPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Commit and push changes to your dotfiles repo",
	RunE: func(cmd *cobra.Command, args []string) error {
		msg, _ := cmd.Flags().GetString("message")
		return dotfiles.Push(msg)
	},
}

func init() {
	dotfilesPushCmd.Flags().StringP("message", "m", "", "commit message")
	dotfilesCmd.AddCommand(dotfilesPushCmd)
	rootCmd.AddCommand(dotfilesCmd)
}
