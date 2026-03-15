package cmd

import (
	"fmt"

	"github.com/aymenhmaidiwastaken/devboot/internal/config"
	"github.com/aymenhmaidiwastaken/devboot/internal/platform"
	"github.com/aymenhmaidiwastaken/devboot/internal/profile"
	"github.com/aymenhmaidiwastaken/devboot/internal/tools"
	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage curated tool profiles",
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available profiles",
	RunE:  runProfileList,
}

var profileShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show details of a profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileShow,
}

var profileApplyCmd = &cobra.Command{
	Use:   "apply <name>",
	Short: "Apply a profile (install tools + configure)",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileApply,
}

var profileExportCmd = &cobra.Command{
	Use:   "export <name>",
	Short: "Export a profile as devboot.yaml",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileExport,
}

var profileSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search profiles by keyword",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileSearch,
}

func init() {
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileShowCmd)
	profileCmd.AddCommand(profileApplyCmd)
	profileCmd.AddCommand(profileExportCmd)
	profileCmd.AddCommand(profileSearchCmd)
	rootCmd.AddCommand(profileCmd)
}

func runProfileList(cmd *cobra.Command, args []string) error {
	fmt.Printf("\n  Available profiles:\n\n")

	for _, name := range profile.List() {
		p, _ := profile.Get(name)
		fmt.Printf("  %-12s %s (%d tools)\n", p.Name, p.Description, len(p.Tools))
	}

	fmt.Printf("\n  Use: devboot profile show <name>   for details\n")
	fmt.Printf("       devboot profile apply <name>  to apply\n\n")
	return nil
}

func runProfileShow(cmd *cobra.Command, args []string) error {
	p, ok := profile.Get(args[0])
	if !ok {
		return fmt.Errorf("unknown profile: %q\nAvailable: %v", args[0], profile.List())
	}

	fmt.Println()
	fmt.Println(profile.Describe(p))
	return nil
}

func runProfileApply(cmd *cobra.Command, args []string) error {
	p, ok := profile.Get(args[0])
	if !ok {
		return fmt.Errorf("unknown profile: %q", args[0])
	}

	fmt.Printf("\n  Applying profile: %s\n", p.Name)
	fmt.Printf("  %s\n", p.Description)

	// Install tools
	if len(p.Tools) > 0 {
		plat := platform.Detect()
		installer, err := tools.NewInstaller(plat)
		if err != nil {
			return err
		}
		if err := installer.InstallAll(p.Tools); err != nil {
			ui.Warn(fmt.Sprintf("some tools failed: %v", err))
		}
	}

	// Apply shell config
	if len(p.Shell.Aliases) > 0 || len(p.Shell.Plugins) > 0 || len(p.Shell.Env) > 0 {
		cfg := &config.Config{Shell: p.Shell}
		if err := applyShell(cfg); err != nil {
			ui.Warn(fmt.Sprintf("shell config: %v", err))
		}
	}

	// Apply vscode
	if len(p.VSCode.Extensions) > 0 {
		cfg := &config.Config{VSCode: p.VSCode}
		if err := applyVSCodePlainWrap(cfg); err != nil {
			ui.Warn(fmt.Sprintf("vscode: %v", err))
		}
	}

	fmt.Println()
	ui.Success(fmt.Sprintf("profile %q applied!", p.Name))
	fmt.Println()
	return nil
}

func runProfileExport(cmd *cobra.Command, args []string) error {
	p, ok := profile.Get(args[0])
	if !ok {
		return fmt.Errorf("unknown profile: %q", args[0])
	}

	cfg := profile.ToConfig(p)
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	fmt.Printf("# devboot.yaml — %s profile\n%s", p.Name, string(data))
	return nil
}

func runProfileSearch(cmd *cobra.Command, args []string) error {
	results := profile.Search(args[0])
	if len(results) == 0 {
		fmt.Printf("  No profiles matching %q\n", args[0])
		return nil
	}

	fmt.Printf("\n  Profiles matching %q:\n\n", args[0])
	for _, p := range results {
		fmt.Printf("  %-12s %s\n", p.Name, p.Description)
	}
	fmt.Println()
	return nil
}
