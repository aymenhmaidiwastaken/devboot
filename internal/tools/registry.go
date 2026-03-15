package tools

import "github.com/aymenhmaidiwastaken/devboot/internal/platform"

// toolMapping maps friendly tool names to platform-specific package names.
var toolMapping = map[string]map[platform.OS]string{
	"git":       {platform.MacOS: "git", platform.Linux: "git", platform.WSL: "git"},
	"node":      {platform.MacOS: "node", platform.Linux: "nodejs", platform.WSL: "nodejs"},
	"nodejs":    {platform.MacOS: "node", platform.Linux: "nodejs", platform.WSL: "nodejs"},
	"python":    {platform.MacOS: "python@3", platform.Linux: "python3", platform.WSL: "python3"},
	"go":        {platform.MacOS: "go", platform.Linux: "golang", platform.WSL: "golang"},
	"golang":    {platform.MacOS: "go", platform.Linux: "golang", platform.WSL: "golang"},
	"rust":      {platform.MacOS: "rustup", platform.Linux: "rustup", platform.WSL: "rustup"},
	"docker":    {platform.MacOS: "docker", platform.Linux: "docker.io", platform.WSL: "docker.io"},
	"kubectl":   {platform.MacOS: "kubectl", platform.Linux: "kubectl", platform.WSL: "kubectl"},
	"terraform": {platform.MacOS: "terraform", platform.Linux: "terraform", platform.WSL: "terraform"},
	"ripgrep":   {platform.MacOS: "ripgrep", platform.Linux: "ripgrep", platform.WSL: "ripgrep"},
	"fzf":       {platform.MacOS: "fzf", platform.Linux: "fzf", platform.WSL: "fzf"},
	"lazygit":   {platform.MacOS: "lazygit", platform.Linux: "lazygit", platform.WSL: "lazygit"},
	"neovim":    {platform.MacOS: "neovim", platform.Linux: "neovim", platform.WSL: "neovim"},
	"nvim":      {platform.MacOS: "neovim", platform.Linux: "neovim", platform.WSL: "neovim"},
	"tmux":      {platform.MacOS: "tmux", platform.Linux: "tmux", platform.WSL: "tmux"},
	"htop":      {platform.MacOS: "htop", platform.Linux: "htop", platform.WSL: "htop"},
	"jq":        {platform.MacOS: "jq", platform.Linux: "jq", platform.WSL: "jq"},
	"curl":      {platform.MacOS: "curl", platform.Linux: "curl", platform.WSL: "curl"},
	"wget":      {platform.MacOS: "wget", platform.Linux: "wget", platform.WSL: "wget"},
}

// ResolvePackage maps a friendly tool name to the correct package name for the platform.
func ResolvePackage(name string, p platform.Platform) string {
	os := p.OS
	if os == platform.WSL {
		os = platform.WSL
	}

	if mapping, ok := toolMapping[name]; ok {
		if pkg, ok := mapping[os]; ok {
			return pkg
		}
	}
	// Fall through: use the name as-is
	return name
}
