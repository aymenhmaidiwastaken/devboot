package profile

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aymenhmaidiwastaken/devboot/internal/config"
)

// Profile is a named, composable set of tools and config.
type Profile struct {
	Name        string
	Description string
	Tools       []string
	Shell       config.ShellConfig
	VSCode      config.VSCodeConfig
}

// BuiltinProfiles contains all built-in profiles.
var BuiltinProfiles = map[string]Profile{
	"frontend": {
		Name:        "frontend",
		Description: "Frontend web development (Node, Bun, linters, formatters)",
		Tools:       []string{"node", "bun", "git", "ripgrep", "fzf"},
		Shell: config.ShellConfig{
			Aliases: map[string]string{
				"nr":  "npm run",
				"ni":  "npm install",
				"dev": "npm run dev",
			},
		},
		VSCode: config.VSCodeConfig{
			Extensions: []string{
				"esbenp.prettier-vscode",
				"dbaeumer.vscode-eslint",
				"bradlc.vscode-tailwindcss",
				"dsznajder.es7-react-js-snippets",
				"formulahendry.auto-rename-tag",
			},
		},
	},
	"backend": {
		Name:        "backend",
		Description: "Backend development (Go, Docker, Kubernetes, databases)",
		Tools:       []string{"go", "docker", "kubectl", "jq", "curl", "git", "ripgrep"},
		Shell: config.ShellConfig{
			Aliases: map[string]string{
				"k":  "kubectl",
				"dc": "docker compose",
				"dps": "docker ps",
			},
		},
		VSCode: config.VSCodeConfig{
			Extensions: []string{
				"golang.go",
				"ms-azuretools.vscode-docker",
				"redhat.vscode-yaml",
			},
		},
	},
	"devops": {
		Name:        "devops",
		Description: "DevOps & SRE toolkit (Docker, K8s, Terraform, monitoring)",
		Tools:       []string{"docker", "kubectl", "helm", "terraform", "jq", "curl", "git", "gh"},
		Shell: config.ShellConfig{
			Aliases: map[string]string{
				"k":   "kubectl",
				"tf":  "terraform",
				"dc":  "docker compose",
				"dps": "docker ps",
			},
		},
	},
	"data": {
		Name:        "data",
		Description: "Data science & ML (Python, Jupyter, data tools)",
		Tools:       []string{"python", "git", "jq", "curl"},
		Shell: config.ShellConfig{
			Aliases: map[string]string{
				"py":  "python3",
				"pip": "python3 -m pip",
				"jnb": "jupyter notebook",
			},
		},
		VSCode: config.VSCodeConfig{
			Extensions: []string{
				"ms-python.python",
				"ms-toolsai.jupyter",
				"ms-python.vscode-pylance",
			},
		},
	},
	"rust": {
		Name:        "rust",
		Description: "Rust development ecosystem",
		Tools:       []string{"rust", "git", "ripgrep"},
		Shell: config.ShellConfig{
			Env: map[string]string{
				"CARGO_HOME": "~/.cargo",
			},
		},
		VSCode: config.VSCodeConfig{
			Extensions: []string{
				"rust-lang.rust-analyzer",
				"serayuzgur.crates",
				"tamasfe.even-better-toml",
			},
		},
	},
	"terminal": {
		Name:        "terminal",
		Description: "Power user terminal setup (modern CLI tools, shell enhancements)",
		Tools:       []string{"ripgrep", "fzf", "bat", "eza", "fd", "zoxide", "jq", "htop", "lazygit", "neovim", "tmux", "starship"},
		Shell: config.ShellConfig{
			Plugins: []string{
				"zsh-autosuggestions",
				"zsh-syntax-highlighting",
				"fzf-tab",
			},
			Aliases: map[string]string{
				"ls":  "eza",
				"ll":  "eza -la",
				"cat": "bat",
				"cd":  "z",
				"g":   "git",
				"lg":  "lazygit",
				"vim": "nvim",
			},
		},
	},
	"fullstack": {
		Name:        "fullstack",
		Description: "Full-stack web development (Node, Go/Python, Docker, DBs)",
		Tools:       []string{"node", "go", "python", "docker", "kubectl", "git", "gh", "jq", "curl", "ripgrep", "fzf"},
		Shell: config.ShellConfig{
			Aliases: map[string]string{
				"g":   "git",
				"k":   "kubectl",
				"dc":  "docker compose",
				"nr":  "npm run",
				"dev": "npm run dev",
			},
		},
		VSCode: config.VSCodeConfig{
			Extensions: []string{
				"esbenp.prettier-vscode",
				"dbaeumer.vscode-eslint",
				"golang.go",
				"ms-python.python",
				"ms-azuretools.vscode-docker",
			},
		},
	},
}

// Get returns a profile by name.
func Get(name string) (Profile, bool) {
	p, ok := BuiltinProfiles[strings.ToLower(name)]
	return p, ok
}

// List returns all available profile names sorted.
func List() []string {
	var names []string
	for name := range BuiltinProfiles {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Search finds profiles matching a keyword.
func Search(query string) []Profile {
	query = strings.ToLower(query)
	var results []Profile
	for _, p := range BuiltinProfiles {
		if strings.Contains(strings.ToLower(p.Name), query) ||
			strings.Contains(strings.ToLower(p.Description), query) {
			results = append(results, p)
		}
		// Also search tool lists
		for _, t := range p.Tools {
			if strings.Contains(t, query) {
				results = append(results, p)
				break
			}
		}
	}
	// Deduplicate
	seen := make(map[string]bool)
	var unique []Profile
	for _, p := range results {
		if !seen[p.Name] {
			seen[p.Name] = true
			unique = append(unique, p)
		}
	}
	return unique
}

// MergeIntoConfig merges a profile into an existing config.
func MergeIntoConfig(cfg *config.Config, p Profile) {
	// Merge tools (deduplicate)
	existing := make(map[string]bool)
	for _, t := range cfg.Tools {
		existing[t] = true
	}
	for _, t := range p.Tools {
		if !existing[t] {
			cfg.Tools = append(cfg.Tools, t)
		}
	}

	// Merge aliases
	if len(p.Shell.Aliases) > 0 {
		if cfg.Shell.Aliases == nil {
			cfg.Shell.Aliases = make(map[string]string)
		}
		for k, v := range p.Shell.Aliases {
			if _, exists := cfg.Shell.Aliases[k]; !exists {
				cfg.Shell.Aliases[k] = v
			}
		}
	}

	// Merge env
	if len(p.Shell.Env) > 0 {
		if cfg.Shell.Env == nil {
			cfg.Shell.Env = make(map[string]string)
		}
		for k, v := range p.Shell.Env {
			if _, exists := cfg.Shell.Env[k]; !exists {
				cfg.Shell.Env[k] = v
			}
		}
	}

	// Merge plugins (deduplicate)
	if len(p.Shell.Plugins) > 0 {
		pluginSet := make(map[string]bool)
		for _, pl := range cfg.Shell.Plugins {
			pluginSet[pl] = true
		}
		for _, pl := range p.Shell.Plugins {
			if !pluginSet[pl] {
				cfg.Shell.Plugins = append(cfg.Shell.Plugins, pl)
			}
		}
	}

	// Merge VS Code extensions
	if len(p.VSCode.Extensions) > 0 {
		extSet := make(map[string]bool)
		for _, e := range cfg.VSCode.Extensions {
			extSet[e] = true
		}
		for _, e := range p.VSCode.Extensions {
			if !extSet[e] {
				cfg.VSCode.Extensions = append(cfg.VSCode.Extensions, e)
			}
		}
	}
}

// ToConfig converts a profile to a standalone config.
func ToConfig(p Profile) *config.Config {
	cfg := &config.Config{
		Tools:  p.Tools,
		Shell:  p.Shell,
		VSCode: p.VSCode,
	}
	return cfg
}

// Describe returns a formatted description of a profile.
func Describe(p Profile) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Profile: %s\n", p.Name))
	sb.WriteString(fmt.Sprintf("  %s\n\n", p.Description))
	sb.WriteString(fmt.Sprintf("  Tools (%d): %s\n", len(p.Tools), strings.Join(p.Tools, ", ")))

	if len(p.Shell.Aliases) > 0 {
		sb.WriteString(fmt.Sprintf("  Aliases (%d):", len(p.Shell.Aliases)))
		for k, v := range p.Shell.Aliases {
			sb.WriteString(fmt.Sprintf(" %s=%s", k, v))
		}
		sb.WriteString("\n")
	}
	if len(p.Shell.Plugins) > 0 {
		sb.WriteString(fmt.Sprintf("  Plugins (%d): %s\n", len(p.Shell.Plugins), strings.Join(p.Shell.Plugins, ", ")))
	}
	if len(p.VSCode.Extensions) > 0 {
		sb.WriteString(fmt.Sprintf("  VS Code extensions (%d): %s\n", len(p.VSCode.Extensions), strings.Join(p.VSCode.Extensions, ", ")))
	}

	return sb.String()
}
