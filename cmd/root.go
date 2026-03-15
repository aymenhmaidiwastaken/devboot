package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devboot",
	Short: "Dev environment bootstrapper",
	Long:  `DevBoot sets up a complete development environment from a single config file. Fresh machine to productive in one command.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
