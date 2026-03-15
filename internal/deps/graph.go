package deps

import (
	"fmt"
	"strings"
)

// ToolInfo holds metadata about a tool including dependencies and post-install hooks.
type ToolInfo struct {
	Name         string
	Dependencies []string // tools that must be installed first
	PostInstall  []string // shell commands to run after install
	BinName      string   // binary name to check (if different from Name)
	Description  string
	Category     string
}

// Registry holds all known tool metadata.
var Registry = map[string]ToolInfo{
	"git": {
		Name: "git", BinName: "git", Category: "vcs",
		Description: "Distributed version control system",
	},
	"node": {
		Name: "node", BinName: "node", Category: "runtime",
		Description:  "JavaScript runtime",
		PostInstall:  []string{"npm install -g npm@latest"},
	},
	"python": {
		Name: "python", BinName: "python3", Category: "runtime",
		Description:  "Python programming language",
		PostInstall:  []string{"python3 -m pip install --upgrade pip"},
	},
	"go": {
		Name: "go", BinName: "go", Category: "runtime",
		Description: "Go programming language",
	},
	"rust": {
		Name: "rust", BinName: "rustc", Category: "runtime",
		Description:  "Rust programming language",
		PostInstall:  []string{"rustup default stable"},
	},
	"docker": {
		Name: "docker", BinName: "docker", Category: "devops",
		Description:  "Container runtime",
		PostInstall:  []string{"sudo usermod -aG docker $USER || true"},
	},
	"kubectl": {
		Name: "kubectl", BinName: "kubectl", Category: "devops",
		Description:  "Kubernetes CLI",
		Dependencies: []string{"docker"},
	},
	"terraform": {
		Name: "terraform", BinName: "terraform", Category: "devops",
		Description: "Infrastructure as code",
	},
	"ripgrep": {
		Name: "ripgrep", BinName: "rg", Category: "cli",
		Description: "Fast recursive grep",
	},
	"fzf": {
		Name: "fzf", BinName: "fzf", Category: "cli",
		Description: "Fuzzy finder",
	},
	"lazygit": {
		Name: "lazygit", BinName: "lazygit", Category: "cli",
		Description:  "Terminal UI for git",
		Dependencies: []string{"git"},
	},
	"neovim": {
		Name: "neovim", BinName: "nvim", Category: "editor",
		Description: "Hyperextensible text editor",
	},
	"tmux": {
		Name: "tmux", BinName: "tmux", Category: "cli",
		Description: "Terminal multiplexer",
	},
	"jq": {
		Name: "jq", BinName: "jq", Category: "cli",
		Description: "Command-line JSON processor",
	},
	"htop": {
		Name: "htop", BinName: "htop", Category: "cli",
		Description: "Interactive process viewer",
	},
	"curl": {
		Name: "curl", BinName: "curl", Category: "cli",
		Description: "URL transfer tool",
	},
	"wget": {
		Name: "wget", BinName: "wget", Category: "cli",
		Description: "Network downloader",
	},
	"starship": {
		Name: "starship", BinName: "starship", Category: "shell",
		Description: "Cross-shell prompt",
	},
	"zoxide": {
		Name: "zoxide", BinName: "zoxide", Category: "cli",
		Description: "Smarter cd command",
	},
	"bat": {
		Name: "bat", BinName: "bat", Category: "cli",
		Description: "Cat with syntax highlighting",
	},
	"eza": {
		Name: "eza", BinName: "eza", Category: "cli",
		Description: "Modern replacement for ls",
	},
	"fd": {
		Name: "fd", BinName: "fd", Category: "cli",
		Description: "Simple fast find alternative",
	},
	"gh": {
		Name: "gh", BinName: "gh", Category: "cli",
		Description:  "GitHub CLI",
		Dependencies: []string{"git"},
	},
	"helm": {
		Name: "helm", BinName: "helm", Category: "devops",
		Description:  "Kubernetes package manager",
		Dependencies: []string{"kubectl"},
	},
	"bun": {
		Name: "bun", BinName: "bun", Category: "runtime",
		Description: "Fast JavaScript runtime & bundler",
	},
	"deno": {
		Name: "deno", BinName: "deno", Category: "runtime",
		Description: "Secure JavaScript/TypeScript runtime",
	},
}

// Resolve returns an ordered list of tools to install, respecting dependencies.
// Returns an error if there are circular dependencies or unknown deps.
func Resolve(tools []string) ([]string, error) {
	// Parse tool specs (strip version)
	names := make([]string, len(tools))
	specMap := make(map[string]string) // name -> original spec
	for i, spec := range tools {
		name := spec
		if idx := strings.Index(spec, "@"); idx > 0 {
			name = spec[:idx]
		}
		names[i] = name
		specMap[name] = spec
	}

	// Build requested set
	requested := make(map[string]bool)
	for _, n := range names {
		requested[n] = true
	}

	// Topological sort with dependency injection
	var order []string
	visited := make(map[string]bool)
	visiting := make(map[string]bool)

	var visit func(name string) error
	visit = func(name string) error {
		if visited[name] {
			return nil
		}
		if visiting[name] {
			return fmt.Errorf("circular dependency detected involving %q", name)
		}
		visiting[name] = true

		if info, ok := Registry[name]; ok {
			for _, dep := range info.Dependencies {
				if err := visit(dep); err != nil {
					return err
				}
			}
		}

		visiting[name] = false
		visited[name] = true

		// Use original spec if this was explicitly requested
		if spec, ok := specMap[name]; ok {
			order = append(order, spec)
		} else {
			// This is an injected dependency
			order = append(order, name)
		}
		return nil
	}

	for _, name := range names {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	return order, nil
}

// GetPostInstall returns post-install commands for a tool.
func GetPostInstall(name string) []string {
	if info, ok := Registry[name]; ok {
		return info.PostInstall
	}
	return nil
}

// GetBinName returns the binary name for a tool.
func GetBinName(name string) string {
	if info, ok := Registry[name]; ok && info.BinName != "" {
		return info.BinName
	}
	return name
}

// GetDependencies returns the dependency list for a tool.
func GetDependencies(name string) []string {
	if info, ok := Registry[name]; ok {
		return info.Dependencies
	}
	return nil
}

// Categories returns all tools grouped by category.
func Categories() map[string][]ToolInfo {
	result := make(map[string][]ToolInfo)
	for _, info := range Registry {
		result[info.Category] = append(result[info.Category], info)
	}
	return result
}

// AllToolNames returns all known tool names sorted.
func AllToolNames() []string {
	var names []string
	for name := range Registry {
		names = append(names, name)
	}
	return names
}

// Conflicts checks for version manager conflicts.
func Conflicts(tools []string) []string {
	var warnings []string

	nodeManagers := 0
	for _, t := range tools {
		name := t
		if idx := strings.Index(t, "@"); idx > 0 {
			name = t[:idx]
		}
		switch name {
		case "nvm", "fnm", "volta", "mise":
			nodeManagers++
		}
	}
	if nodeManagers > 1 {
		warnings = append(warnings, "multiple Node.js version managers detected (nvm, fnm, volta, mise) — this can cause PATH conflicts")
	}

	return warnings
}
