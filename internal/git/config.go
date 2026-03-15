package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/aymenhmaidiwastaken/devboot/internal/config"
	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
)

// Configure applies git configuration from the config.
func Configure(cfg config.GitConfig) error {
	ui.Section("Configuring Git")

	settings := map[string]string{}

	if cfg.UserName != "" {
		settings["user.name"] = cfg.UserName
	}
	if cfg.UserEmail != "" {
		settings["user.email"] = cfg.UserEmail
	}
	if cfg.InitDefaultBranch != "" {
		settings["init.defaultBranch"] = cfg.InitDefaultBranch
	}
	if cfg.PullRebase != nil && *cfg.PullRebase {
		settings["pull.rebase"] = "true"
	}

	for key, value := range settings {
		current := getGitConfig(key)
		if current == value {
			ui.Skip(fmt.Sprintf("%s = %s", key, value))
			continue
		}

		if err := setGitConfig(key, value); err != nil {
			ui.Fail(fmt.Sprintf("%s: %v", key, err))
			return err
		}
		ui.Success(fmt.Sprintf("%s = %s", key, value))
	}

	// Git aliases
	if len(cfg.Aliases) > 0 {
		for name, command := range cfg.Aliases {
			key := "alias." + name
			current := getGitConfig(key)
			if current == command {
				ui.Skip(fmt.Sprintf("git %s", name))
				continue
			}
			if err := setGitConfig(key, command); err != nil {
				ui.Fail(fmt.Sprintf("git alias %s: %v", name, err))
				continue
			}
			ui.Success(fmt.Sprintf("git %s → %s", name, command))
		}
	}

	return nil
}

func getGitConfig(key string) string {
	cmd := exec.Command("git", "config", "--global", "--get", key)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func setGitConfig(key, value string) error {
	cmd := exec.Command("git", "config", "--global", key, value)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git config --global %s %s: %w", key, value, err)
	}
	return nil
}
