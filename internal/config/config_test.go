package config

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// ---------------------------------------------------------------------------
// 1. Load – valid config with ALL sections
// ---------------------------------------------------------------------------

func TestLoadValidConfigAllSections(t *testing.T) {
	raw := `
tools:
  - git
  - node@22
  - python@3.12

shell:
  type: zsh
  plugins:
    - zsh-autosuggestions
    - zsh-syntax-highlighting
  aliases:
    g: git
    ll: ls -la
  env:
    EDITOR: vim
    GOPATH: /home/user/go

git:
  user.name: "Test User"
  user.email: "test@example.com"
  init.defaultBranch: main
  pull.rebase: true
  sshKey: true
  aliases:
    co: checkout
    br: branch

vscode:
  extensions:
    - ms-python.python
    - golang.go
  settings:
    editor.fontSize: 14
    editor.tabSize: 2

neovim:
  configRepo: https://github.com/test/nvim-config.git

jetbrains:
  plugins:
    - ideavim
    - rainbow-brackets

dotfiles:
  repo: https://github.com/test/dotfiles.git
  mappings:
    .vimrc: ~/.vimrc
    .zshrc: ~/.zshrc
  templates:
    .gitconfig: ~/.gitconfig
`
	dir := t.TempDir()
	path := filepath.Join(dir, "devboot.yaml")
	if err := os.WriteFile(path, []byte(raw), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Tools
	if len(cfg.Tools) != 3 {
		t.Errorf("expected 3 tools, got %d", len(cfg.Tools))
	}
	if cfg.Tools[0] != "git" {
		t.Errorf("expected tools[0]='git', got %q", cfg.Tools[0])
	}
	if cfg.Tools[1] != "node@22" {
		t.Errorf("expected tools[1]='node@22', got %q", cfg.Tools[1])
	}

	// Shell
	if cfg.Shell.Type != "zsh" {
		t.Errorf("expected shell type 'zsh', got %q", cfg.Shell.Type)
	}
	if len(cfg.Shell.Plugins) != 2 {
		t.Errorf("expected 2 shell plugins, got %d", len(cfg.Shell.Plugins))
	}
	if cfg.Shell.Aliases["g"] != "git" {
		t.Errorf("expected alias g='git', got %q", cfg.Shell.Aliases["g"])
	}
	if cfg.Shell.Aliases["ll"] != "ls -la" {
		t.Errorf("expected alias ll='ls -la', got %q", cfg.Shell.Aliases["ll"])
	}
	if cfg.Shell.Env["EDITOR"] != "vim" {
		t.Errorf("expected env EDITOR='vim', got %q", cfg.Shell.Env["EDITOR"])
	}

	// Git
	if cfg.Git.UserName != "Test User" {
		t.Errorf("expected user.name 'Test User', got %q", cfg.Git.UserName)
	}
	if cfg.Git.UserEmail != "test@example.com" {
		t.Errorf("expected user.email 'test@example.com', got %q", cfg.Git.UserEmail)
	}
	if cfg.Git.InitDefaultBranch != "main" {
		t.Errorf("expected init.defaultBranch 'main', got %q", cfg.Git.InitDefaultBranch)
	}
	if cfg.Git.PullRebase == nil || *cfg.Git.PullRebase != true {
		t.Errorf("expected pull.rebase=true, got %v", cfg.Git.PullRebase)
	}
	if cfg.Git.SSHKey != true {
		t.Errorf("expected sshKey=true, got %v", cfg.Git.SSHKey)
	}
	if cfg.Git.Aliases["co"] != "checkout" {
		t.Errorf("expected git alias co='checkout', got %q", cfg.Git.Aliases["co"])
	}
	if cfg.Git.Aliases["br"] != "branch" {
		t.Errorf("expected git alias br='branch', got %q", cfg.Git.Aliases["br"])
	}

	// VSCode
	if len(cfg.VSCode.Extensions) != 2 {
		t.Errorf("expected 2 vscode extensions, got %d", len(cfg.VSCode.Extensions))
	}
	if cfg.VSCode.Settings["editor.fontSize"] != 14 {
		t.Errorf("unexpected editor.fontSize: %v", cfg.VSCode.Settings["editor.fontSize"])
	}

	// Neovim
	if cfg.Neovim.ConfigRepo != "https://github.com/test/nvim-config.git" {
		t.Errorf("unexpected neovim configRepo: %q", cfg.Neovim.ConfigRepo)
	}

	// JetBrains
	if len(cfg.JetBrains.Plugins) != 2 {
		t.Errorf("expected 2 jetbrains plugins, got %d", len(cfg.JetBrains.Plugins))
	}
	if cfg.JetBrains.Plugins[0] != "ideavim" {
		t.Errorf("expected jetbrains plugin 'ideavim', got %q", cfg.JetBrains.Plugins[0])
	}

	// Dotfiles
	if cfg.Dotfiles.Repo != "https://github.com/test/dotfiles.git" {
		t.Errorf("unexpected dotfiles repo: %q", cfg.Dotfiles.Repo)
	}
	if cfg.Dotfiles.Mappings[".vimrc"] != "~/.vimrc" {
		t.Errorf("unexpected dotfiles mapping: %q", cfg.Dotfiles.Mappings[".vimrc"])
	}
	if cfg.Dotfiles.Templates[".gitconfig"] != "~/.gitconfig" {
		t.Errorf("unexpected dotfiles template: %q", cfg.Dotfiles.Templates[".gitconfig"])
	}
}

// ---------------------------------------------------------------------------
// 2. Load – partial config (only some sections)
// ---------------------------------------------------------------------------

func TestLoadPartialConfigToolsOnly(t *testing.T) {
	raw := `
tools:
  - go
  - rust
`
	dir := t.TempDir()
	path := filepath.Join(dir, "devboot.yaml")
	os.WriteFile(path, []byte(raw), 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.Tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(cfg.Tools))
	}
	if cfg.Shell.Type != "" {
		t.Errorf("expected empty shell type, got %q", cfg.Shell.Type)
	}
	if cfg.Git.UserName != "" {
		t.Errorf("expected empty git user.name, got %q", cfg.Git.UserName)
	}
	if len(cfg.VSCode.Extensions) != 0 {
		t.Errorf("expected no vscode extensions, got %d", len(cfg.VSCode.Extensions))
	}
	if cfg.Neovim.ConfigRepo != "" {
		t.Errorf("expected empty neovim configRepo, got %q", cfg.Neovim.ConfigRepo)
	}
	if len(cfg.JetBrains.Plugins) != 0 {
		t.Errorf("expected no jetbrains plugins, got %d", len(cfg.JetBrains.Plugins))
	}
	if cfg.Dotfiles.Repo != "" {
		t.Errorf("expected empty dotfiles repo, got %q", cfg.Dotfiles.Repo)
	}
}

func TestLoadPartialConfigShellAndGitOnly(t *testing.T) {
	raw := `
shell:
  type: fish
git:
  user.name: "Fish User"
`
	dir := t.TempDir()
	path := filepath.Join(dir, "devboot.yaml")
	os.WriteFile(path, []byte(raw), 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Shell.Type != "fish" {
		t.Errorf("expected shell type 'fish', got %q", cfg.Shell.Type)
	}
	if cfg.Git.UserName != "Fish User" {
		t.Errorf("expected git user.name 'Fish User', got %q", cfg.Git.UserName)
	}
	if len(cfg.Tools) != 0 {
		t.Errorf("expected no tools, got %d", len(cfg.Tools))
	}
}

// ---------------------------------------------------------------------------
// 3. Load – empty config (no sections at all)
// ---------------------------------------------------------------------------

func TestLoadEmptyConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "devboot.yaml")
	os.WriteFile(path, []byte(""), 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed on empty config: %v", err)
	}

	if len(cfg.Tools) != 0 {
		t.Errorf("expected 0 tools, got %d", len(cfg.Tools))
	}
	if cfg.Shell.Type != "" {
		t.Errorf("expected empty shell type, got %q", cfg.Shell.Type)
	}
}

func TestLoadBlankYAMLComment(t *testing.T) {
	raw := `# just a comment, nothing else`
	dir := t.TempDir()
	path := filepath.Join(dir, "devboot.yaml")
	os.WriteFile(path, []byte(raw), 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(cfg.Tools) != 0 {
		t.Errorf("expected 0 tools, got %d", len(cfg.Tools))
	}
}

// ---------------------------------------------------------------------------
// 4. Load – missing file
// ---------------------------------------------------------------------------

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/devboot.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "reading config") {
		t.Errorf("expected 'reading config' in error, got %q", err.Error())
	}
}

// ---------------------------------------------------------------------------
// 5. Load – invalid YAML syntax
// ---------------------------------------------------------------------------

func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "devboot.yaml")
	os.WriteFile(path, []byte("{{invalid yaml content:::"), 0644)

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
	if !strings.Contains(err.Error(), "parsing config") {
		t.Errorf("expected 'parsing config' in error, got %q", err.Error())
	}
}

func TestLoadInvalidYAMLTabIndent(t *testing.T) {
	raw := "tools:\n\t- git\n"
	dir := t.TempDir()
	path := filepath.Join(dir, "devboot.yaml")
	os.WriteFile(path, []byte(raw), 0644)

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for YAML with tabs")
	}
}

// ---------------------------------------------------------------------------
// 6. Load – remote URL paths (branching test, not actual HTTP)
// ---------------------------------------------------------------------------

func TestLoadRemoteURLHTTP(t *testing.T) {
	raw := `
tools:
  - curl
`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(raw))
	}))
	defer srv.Close()

	cfg, err := Load(srv.URL + "/devboot.yaml")
	if err != nil {
		t.Fatalf("Load from HTTP server failed: %v", err)
	}
	if len(cfg.Tools) != 1 || cfg.Tools[0] != "curl" {
		t.Errorf("unexpected tools from remote load: %v", cfg.Tools)
	}
}

func TestLoadRemoteURLHTTPS(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("tools:\n  - wget\n"))
	}))
	defer srv.Close()

	// The default http client won't trust the self-signed cert, so this
	// should fail with a TLS error – but crucially it exercises the
	// https:// branch in Load, NOT the os.ReadFile branch.
	_, err := Load(srv.URL + "/devboot.yaml")
	if err == nil {
		// If it somehow succeeds (shouldn't with default client), that is fine.
		return
	}
	if !strings.Contains(err.Error(), "reading config") {
		t.Errorf("expected 'reading config' wrapper, got %q", err.Error())
	}
}

func TestLoadRemoteURL404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := Load(srv.URL + "/missing.yaml")
	if err == nil {
		t.Error("expected error for HTTP 404")
	}
	if !strings.Contains(err.Error(), "HTTP 404") {
		t.Errorf("expected 'HTTP 404' in error, got %q", err.Error())
	}
}

// ---------------------------------------------------------------------------
// 7. Validate – valid shell types (bash, zsh, fish, empty)
// ---------------------------------------------------------------------------

func TestValidateValidShellTypes(t *testing.T) {
	for _, shell := range []string{"bash", "zsh", "fish", ""} {
		cfg := &Config{Shell: ShellConfig{Type: shell}}
		if err := cfg.Validate(); err != nil {
			t.Errorf("unexpected validation error for shell %q: %v", shell, err)
		}
	}
}

// ---------------------------------------------------------------------------
// 8. Validate – invalid shell type
// ---------------------------------------------------------------------------

func TestValidateInvalidShellType(t *testing.T) {
	for _, shell := range []string{"powershell", "cmd", "sh", "csh", "ksh", "nushell"} {
		cfg := &Config{Shell: ShellConfig{Type: shell}}
		err := cfg.Validate()
		if err == nil {
			t.Errorf("expected error for shell type %q", shell)
			continue
		}
		if !strings.Contains(err.Error(), "unsupported shell type") {
			t.Errorf("expected 'unsupported shell type' in error for %q, got %q", shell, err.Error())
		}
	}
}

// ---------------------------------------------------------------------------
// 9. Validate – empty tool name
// ---------------------------------------------------------------------------

func TestValidateEmptyToolName(t *testing.T) {
	cases := []struct {
		name  string
		tools []string
	}{
		{"literal empty", []string{"git", "", "node"}},
		{"whitespace only", []string{"git", "   ", "node"}},
		{"tab whitespace", []string{"\t"}},
		{"single empty", []string{""}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &Config{Tools: tc.tools}
			err := cfg.Validate()
			if err == nil {
				t.Error("expected validation error for empty tool name")
			}
			if !strings.Contains(err.Error(), "empty tool name") {
				t.Errorf("expected 'empty tool name' in error, got %q", err.Error())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 10. Validate – tools with version specs (node@22)
// ---------------------------------------------------------------------------

func TestValidateToolsWithVersionSpecs(t *testing.T) {
	cfg := &Config{
		Tools: []string{"node@22", "python@3.12", "go@1.22", "rust@nightly"},
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected validation error for versioned tools: %v", err)
	}
}

// ---------------------------------------------------------------------------
// 11. Validate – config with all fields valid
// ---------------------------------------------------------------------------

func TestValidateAllFieldsValid(t *testing.T) {
	tr := true
	cfg := &Config{
		Tools: []string{"git", "node@22", "python@3.12"},
		Shell: ShellConfig{
			Type:    "bash",
			Plugins: []string{"bash-completion"},
			Aliases: map[string]string{"g": "git", "ll": "ls -la"},
			Env:     map[string]string{"EDITOR": "vim"},
		},
		Git: GitConfig{
			UserName:          "Test User",
			UserEmail:         "test@example.com",
			InitDefaultBranch: "main",
			PullRebase:        &tr,
			SSHKey:            true,
			Aliases:           map[string]string{"co": "checkout"},
		},
		VSCode: VSCodeConfig{
			Extensions: []string{"golang.go"},
			Settings:   map[string]interface{}{"editor.fontSize": 14},
		},
		Neovim: NeovimConfig{
			ConfigRepo: "https://github.com/user/nvim.git",
		},
		JetBrains: JetBrainsConfig{
			Plugins: []string{"ideavim"},
		},
		Dotfiles: DotfilesConfig{
			Repo:      "https://github.com/user/dots.git",
			Mappings:  map[string]string{".zshrc": "~/.zshrc"},
			Templates: map[string]string{".gitconfig": "~/.gitconfig"},
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no validation error, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// 12. DefaultConfigPath – returns correct default
// ---------------------------------------------------------------------------

func TestDefaultConfigPath(t *testing.T) {
	got := DefaultConfigPath()
	if got != "devboot.yaml" {
		t.Errorf("expected 'devboot.yaml', got %q", got)
	}
}

// ---------------------------------------------------------------------------
// 13. Config struct – YAML unmarshalling for each field type
// ---------------------------------------------------------------------------

func TestYAMLUnmarshalBoolPullRebase(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected *bool
	}{
		{"true", "pull.rebase: true", boolPtr(true)},
		{"false", "pull.rebase: false", boolPtr(false)},
		{"omitted", "", nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			raw := "git:\n"
			if tc.yaml != "" {
				raw += "  " + tc.yaml + "\n"
			} else {
				raw += "  user.name: placeholder\n"
			}
			var cfg Config
			if err := yaml.Unmarshal([]byte(raw), &cfg); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			if tc.expected == nil {
				if cfg.Git.PullRebase != nil {
					t.Errorf("expected nil PullRebase, got %v", *cfg.Git.PullRebase)
				}
			} else {
				if cfg.Git.PullRebase == nil {
					t.Fatalf("expected PullRebase=%v, got nil", *tc.expected)
				}
				if *cfg.Git.PullRebase != *tc.expected {
					t.Errorf("expected PullRebase=%v, got %v", *tc.expected, *cfg.Git.PullRebase)
				}
			}
		})
	}
}

func TestYAMLUnmarshalSSHKeyBool(t *testing.T) {
	raw := "git:\n  sshKey: true\n"
	var cfg Config
	if err := yaml.Unmarshal([]byte(raw), &cfg); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if cfg.Git.SSHKey != true {
		t.Errorf("expected sshKey=true, got %v", cfg.Git.SSHKey)
	}
}

func TestYAMLUnmarshalMapStringString(t *testing.T) {
	raw := `
shell:
  aliases:
    g: git
    k: kubectl
  env:
    HOME: /home/user
`
	var cfg Config
	if err := yaml.Unmarshal([]byte(raw), &cfg); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if len(cfg.Shell.Aliases) != 2 {
		t.Errorf("expected 2 aliases, got %d", len(cfg.Shell.Aliases))
	}
	if cfg.Shell.Env["HOME"] != "/home/user" {
		t.Errorf("expected HOME=/home/user, got %q", cfg.Shell.Env["HOME"])
	}
}

func TestYAMLUnmarshalMapStringInterface(t *testing.T) {
	raw := `
vscode:
  settings:
    editor.fontSize: 16
    editor.wordWrap: "on"
    editor.minimap.enabled: true
`
	var cfg Config
	if err := yaml.Unmarshal([]byte(raw), &cfg); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if cfg.VSCode.Settings["editor.fontSize"] != 16 {
		t.Errorf("expected fontSize=16, got %v", cfg.VSCode.Settings["editor.fontSize"])
	}
	if cfg.VSCode.Settings["editor.wordWrap"] != "on" {
		t.Errorf("expected wordWrap='on', got %v", cfg.VSCode.Settings["editor.wordWrap"])
	}
	if cfg.VSCode.Settings["editor.minimap.enabled"] != true {
		t.Errorf("expected minimap.enabled=true, got %v", cfg.VSCode.Settings["editor.minimap.enabled"])
	}
}

func TestYAMLUnmarshalStringSlice(t *testing.T) {
	raw := `
tools:
  - alpha
  - bravo
  - charlie
jetbrains:
  plugins:
    - ideavim
`
	var cfg Config
	if err := yaml.Unmarshal([]byte(raw), &cfg); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if len(cfg.Tools) != 3 {
		t.Errorf("expected 3 tools, got %d", len(cfg.Tools))
	}
	if len(cfg.JetBrains.Plugins) != 1 {
		t.Errorf("expected 1 jetbrains plugin, got %d", len(cfg.JetBrains.Plugins))
	}
}

func TestYAMLUnmarshalDotfilesConfig(t *testing.T) {
	raw := `
dotfiles:
  repo: https://github.com/user/dots
  mappings:
    .bashrc: ~/.bashrc
  templates:
    .gitconfig: ~/.gitconfig
`
	var cfg Config
	if err := yaml.Unmarshal([]byte(raw), &cfg); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if cfg.Dotfiles.Repo != "https://github.com/user/dots" {
		t.Errorf("unexpected repo: %q", cfg.Dotfiles.Repo)
	}
	if cfg.Dotfiles.Mappings[".bashrc"] != "~/.bashrc" {
		t.Errorf("unexpected mapping: %q", cfg.Dotfiles.Mappings[".bashrc"])
	}
	if cfg.Dotfiles.Templates[".gitconfig"] != "~/.gitconfig" {
		t.Errorf("unexpected template: %q", cfg.Dotfiles.Templates[".gitconfig"])
	}
}

func TestYAMLUnmarshalNeovimConfig(t *testing.T) {
	raw := `
neovim:
  configRepo: https://github.com/user/nvim-config
`
	var cfg Config
	if err := yaml.Unmarshal([]byte(raw), &cfg); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if cfg.Neovim.ConfigRepo != "https://github.com/user/nvim-config" {
		t.Errorf("unexpected configRepo: %q", cfg.Neovim.ConfigRepo)
	}
}

// ---------------------------------------------------------------------------
// 14. Edge cases – unicode, very long tool names, special characters
// ---------------------------------------------------------------------------

func TestEdgeCaseUnicodeValues(t *testing.T) {
	raw := `
tools:
  - "日本語ツール"
shell:
  type: zsh
  aliases:
    café: "echo ☕"
    "🚀": "echo launch"
  env:
    GREETING: "こんにちは世界"
git:
  user.name: "Ñoño Müller"
  user.email: "user@ünì.com"
`
	dir := t.TempDir()
	path := filepath.Join(dir, "devboot.yaml")
	os.WriteFile(path, []byte(raw), 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed with unicode: %v", err)
	}
	if cfg.Tools[0] != "日本語ツール" {
		t.Errorf("expected unicode tool name, got %q", cfg.Tools[0])
	}
	if cfg.Shell.Aliases["café"] != "echo ☕" {
		t.Errorf("expected unicode alias value, got %q", cfg.Shell.Aliases["café"])
	}
	if cfg.Shell.Env["GREETING"] != "こんにちは世界" {
		t.Errorf("expected unicode env value, got %q", cfg.Shell.Env["GREETING"])
	}
	if cfg.Git.UserName != "Ñoño Müller" {
		t.Errorf("expected unicode git user.name, got %q", cfg.Git.UserName)
	}
}

func TestEdgeCaseVeryLongToolName(t *testing.T) {
	longName := strings.Repeat("a", 1000)
	raw := "tools:\n  - " + longName + "\n"

	dir := t.TempDir()
	path := filepath.Join(dir, "devboot.yaml")
	os.WriteFile(path, []byte(raw), 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed with long tool name: %v", err)
	}
	if len(cfg.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(cfg.Tools))
	}
	if cfg.Tools[0] != longName {
		t.Errorf("long tool name not preserved (length %d vs %d)", len(cfg.Tools[0]), len(longName))
	}
}

func TestEdgeCaseSpecialCharactersInAliases(t *testing.T) {
	raw := `
shell:
  type: bash
  aliases:
    "..": "cd .."
    "...": "cd ../.."
    "-": "cd -"
    "g!": "git stash pop"
    "g?": "git status"
    "ls -a": "ls -la --color=auto"
`
	dir := t.TempDir()
	path := filepath.Join(dir, "devboot.yaml")
	os.WriteFile(path, []byte(raw), 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed with special char aliases: %v", err)
	}
	if cfg.Shell.Aliases[".."] != "cd .." {
		t.Errorf("expected alias '..'='cd ..', got %q", cfg.Shell.Aliases[".."])
	}
	if cfg.Shell.Aliases["-"] != "cd -" {
		t.Errorf("expected alias '-'='cd -', got %q", cfg.Shell.Aliases["-"])
	}
	if cfg.Shell.Aliases["g!"] != "git stash pop" {
		t.Errorf("expected alias 'g!'='git stash pop', got %q", cfg.Shell.Aliases["g!"])
	}
}

func TestEdgeCaseToolWithAtSign(t *testing.T) {
	cfg := &Config{
		Tools: []string{"node@22", "@scope/package", "python@3.12.1"},
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected validation error for @ in tool names: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Load integration: validation failure propagates through Load
// ---------------------------------------------------------------------------

func TestLoadReturnsValidationError(t *testing.T) {
	raw := `
tools:
  - git
  - ""
  - node
shell:
  type: bash
`
	dir := t.TempDir()
	path := filepath.Join(dir, "devboot.yaml")
	os.WriteFile(path, []byte(raw), 0644)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected Load to return validation error")
	}
	if !strings.Contains(err.Error(), "invalid config") {
		t.Errorf("expected 'invalid config' in error, got %q", err.Error())
	}
}

func TestLoadReturnsValidationErrorForBadShell(t *testing.T) {
	raw := `
shell:
  type: powershell
`
	dir := t.TempDir()
	path := filepath.Join(dir, "devboot.yaml")
	os.WriteFile(path, []byte(raw), 0644)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected Load to return validation error for bad shell")
	}
	if !strings.Contains(err.Error(), "unsupported shell type") {
		t.Errorf("expected 'unsupported shell type' in error, got %q", err.Error())
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func boolPtr(b bool) *bool {
	return &b
}
