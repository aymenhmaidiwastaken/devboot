package cmd

import (
	"fmt"
	"os"

	"github.com/aymenhmaidiwastaken/devboot/internal/config"
	"github.com/aymenhmaidiwastaken/devboot/internal/platform"
	"github.com/aymenhmaidiwastaken/devboot/internal/tools"
	"github.com/aymenhmaidiwastaken/devboot/internal/tui"
	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var addCmd = &cobra.Command{
	Use:   "add [tool...]",
	Short: "Add and install tools (interactive if no args)",
	RunE:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	var toolsToAdd []string

	if len(args) > 0 {
		toolsToAdd = args
	} else {
		// Interactive mode
		selected, err := tui.RunAddTool()
		if err != nil {
			return err
		}
		if len(selected) == 0 {
			fmt.Println("  No tools selected.")
			return nil
		}
		toolsToAdd = selected
	}

	p := platform.Detect()
	installer, err := tools.NewInstaller(p)
	if err != nil {
		return err
	}

	// Install the tools
	if err := installer.InstallAll(toolsToAdd); err != nil {
		return err
	}

	// Update devboot.yaml if it exists
	configPath := config.DefaultConfigPath()
	if fileExists(configPath) {
		if err := addToolsToConfig(configPath, toolsToAdd); err != nil {
			ui.Warn(fmt.Sprintf("could not update %s: %v", configPath, err))
		} else {
			ui.Success(fmt.Sprintf("updated %s", configPath))
		}
	}

	return nil
}

func addToolsToConfig(path string, newTools []string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Get existing tools
	existing := make(map[string]bool)
	if toolsRaw, ok := raw["tools"]; ok {
		if toolsList, ok := toolsRaw.([]interface{}); ok {
			for _, t := range toolsList {
				if s, ok := t.(string); ok {
					existing[s] = true
				}
			}
		}
	}

	// Add new tools
	var toolsList []string
	if toolsRaw, ok := raw["tools"]; ok {
		if tl, ok := toolsRaw.([]interface{}); ok {
			for _, t := range tl {
				if s, ok := t.(string); ok {
					toolsList = append(toolsList, s)
				}
			}
		}
	}

	for _, t := range newTools {
		if !existing[t] {
			toolsList = append(toolsList, t)
		}
	}
	raw["tools"] = toolsList

	out, err := yaml.Marshal(raw)
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0644)
}
