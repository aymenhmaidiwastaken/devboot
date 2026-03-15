package shell

import (
	"fmt"
	"os"
	"strings"

	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
)

// SetAliases adds aliases to the shell config file without duplicating.
func SetAliases(shellType ShellType, aliases map[string]string) error {
	if len(aliases) == 0 {
		return nil
	}

	ui.Section("Configuring shell aliases")

	configPath := ConfigFile(shellType)
	existing, _ := os.ReadFile(configPath)
	content := string(existing)

	var added int
	var lines []string

	for name, command := range aliases {
		var aliasLine string
		if shellType == Fish {
			aliasLine = fmt.Sprintf("alias %s '%s'", name, command)
		} else {
			aliasLine = fmt.Sprintf("alias %s='%s'", name, command)
		}

		if strings.Contains(content, aliasLine) {
			ui.Skip(fmt.Sprintf("alias %s", name))
			continue
		}

		lines = append(lines, aliasLine)
		added++
	}

	if added == 0 {
		return nil
	}

	// Add a devboot marker section
	block := "\n# devboot aliases\n" + strings.Join(lines, "\n") + "\n"

	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening %s: %w", configPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(block); err != nil {
		return fmt.Errorf("writing aliases: %w", err)
	}

	for name := range aliases {
		ui.Success(fmt.Sprintf("alias %s", name))
	}

	return nil
}

// SetEnvVars adds environment variables to the shell config file.
func SetEnvVars(shellType ShellType, envVars map[string]string) error {
	if len(envVars) == 0 {
		return nil
	}

	ui.Section("Configuring environment variables")

	configPath := ConfigFile(shellType)
	existing, _ := os.ReadFile(configPath)
	content := string(existing)

	var lines []string
	var added int

	for key, value := range envVars {
		var exportLine string
		if shellType == Fish {
			exportLine = fmt.Sprintf("set -gx %s %s", key, value)
		} else {
			exportLine = fmt.Sprintf("export %s=%s", key, value)
		}

		if strings.Contains(content, exportLine) {
			ui.Skip(fmt.Sprintf("%s=%s", key, value))
			continue
		}

		lines = append(lines, exportLine)
		added++
	}

	if added == 0 {
		return nil
	}

	block := "\n# devboot env\n" + strings.Join(lines, "\n") + "\n"

	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening %s: %w", configPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(block); err != nil {
		return fmt.Errorf("writing env vars: %w", err)
	}

	for key, value := range envVars {
		ui.Success(fmt.Sprintf("%s=%s", key, value))
	}

	return nil
}
