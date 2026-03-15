package cmd

import (
	"fmt"

	"github.com/aymenhmaidiwastaken/devboot/internal/platform"
	"github.com/aymenhmaidiwastaken/devboot/internal/state"
	"github.com/aymenhmaidiwastaken/devboot/internal/tools"
	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
	"github.com/spf13/cobra"
)

var uninstallAll bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [tool...]",
	Short: "Uninstall tools that were installed by devboot",
	RunE:  runUninstall,
}

func init() {
	uninstallCmd.Flags().BoolVar(&uninstallAll, "all", false, "uninstall all devboot-managed tools")
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	p := platform.Detect()
	installer, err := tools.NewInstaller(p)
	if err != nil {
		return err
	}
	installer.State = st

	var toRemove []string

	if uninstallAll {
		if len(st.ManagedTools) == 0 {
			ui.Info("no tools managed by devboot")
			return nil
		}
		toRemove = st.ManagedTools
		ui.Warn(fmt.Sprintf("this will uninstall %d tools: %v", len(toRemove), toRemove))
	} else if len(args) > 0 {
		for _, arg := range args {
			if st.IsManagedTool(arg) {
				toRemove = append(toRemove, arg)
			} else {
				ui.Warn(fmt.Sprintf("%s was not installed by devboot (use your package manager directly)", arg))
			}
		}
	} else {
		return fmt.Errorf("specify tools to uninstall or use --all")
	}

	if len(toRemove) == 0 {
		return nil
	}

	return installer.UninstallAll(toRemove)
}
