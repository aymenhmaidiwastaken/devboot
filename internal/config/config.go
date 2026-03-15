package config

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level devboot.yaml configuration.
type Config struct {
	Tools      []string         `yaml:"tools,omitempty"`
	Shell      ShellConfig      `yaml:"shell,omitempty"`
	Git        GitConfig        `yaml:"git,omitempty"`
	VSCode     VSCodeConfig     `yaml:"vscode,omitempty"`
	Neovim     NeovimConfig     `yaml:"neovim,omitempty"`
	JetBrains  JetBrainsConfig  `yaml:"jetbrains,omitempty"`
	Dotfiles   DotfilesConfig   `yaml:"dotfiles,omitempty"`
}

type ShellConfig struct {
	Type    string            `yaml:"type,omitempty"`
	Plugins []string          `yaml:"plugins,omitempty"`
	Aliases map[string]string `yaml:"aliases,omitempty"`
	Env     map[string]string `yaml:"env,omitempty"`
}

type GitConfig struct {
	UserName          string            `yaml:"user.name,omitempty"`
	UserEmail         string            `yaml:"user.email,omitempty"`
	InitDefaultBranch string            `yaml:"init.defaultBranch,omitempty"`
	PullRebase        *bool             `yaml:"pull.rebase,omitempty"`
	Aliases           map[string]string `yaml:"aliases,omitempty"`
	SSHKey            bool              `yaml:"sshKey,omitempty"`
}

type VSCodeConfig struct {
	Extensions []string               `yaml:"extensions,omitempty"`
	Settings   map[string]interface{} `yaml:"settings,omitempty"`
}

type NeovimConfig struct {
	ConfigRepo string `yaml:"configRepo,omitempty"`
}

type JetBrainsConfig struct {
	Plugins []string `yaml:"plugins,omitempty"`
}

type DotfilesConfig struct {
	Repo      string            `yaml:"repo,omitempty"`
	Mappings  map[string]string `yaml:"mappings,omitempty"`
	Templates map[string]string `yaml:"templates,omitempty"`
}

// Load reads a config from a file path or URL.
func Load(path string) (*Config, error) {
	var data []byte
	var err error

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		data, err = fetchRemote(path)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func fetchRemote(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d fetching %s", resp.StatusCode, url)
	}

	return io.ReadAll(resp.Body)
}

// Validate checks the config for obvious errors.
func (c *Config) Validate() error {
	for _, tool := range c.Tools {
		if strings.TrimSpace(tool) == "" {
			return fmt.Errorf("empty tool name in tools list")
		}
	}

	if c.Shell.Type != "" {
		switch c.Shell.Type {
		case "bash", "zsh", "fish":
			// valid
		default:
			return fmt.Errorf("unsupported shell type: %q (supported: bash, zsh, fish)", c.Shell.Type)
		}
	}

	return nil
}

// DefaultConfigPath returns the default config file path.
func DefaultConfigPath() string {
	return "devboot.yaml"
}
