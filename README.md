<p align="center">
  <h1 align="center">DevBoot</h1>
  <p align="center">
    <strong>Fresh machine to productive in one command.</strong>
  </p>
  <p align="center">
    <a href="#install">Install</a> •
    <a href="#quick-start">Quick Start</a> •
    <a href="#config-reference">Config</a> •
    <a href="#profiles">Profiles</a> •
    <a href="#commands">Commands</a> •
    <a href="#vs-alternatives">Comparison</a>
  </p>
</p>

<p align="center">
  <img src="demo/demo.gif" alt="DevBoot Demo" width="880">
</p>

<p align="center">
  <a href="https://github.com/aymenhmaidiwastaken/devboot/actions"><img src="https://github.com/aymenhmaidiwastaken/devboot/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://github.com/aymenhmaidiwastaken/devboot/releases"><img src="https://img.shields.io/github/v/release/aymenhmaidiwastaken/devboot" alt="Release"></a>
  <a href="https://goreportcard.com/report/github.com/aymenhmaidiwastaken/devboot"><img src="https://goreportcard.com/badge/github.com/aymenhmaidiwastaken/devboot" alt="Go Report Card"></a>
  <a href="LICENSE"><img src="https://img.shields.io/github/license/aymenhmaidiwastaken/devboot" alt="License"></a>
</p>

---

DevBoot sets up a complete development environment from a single YAML config file. Install tools, configure your shell, set up git, install IDE extensions, and sync dotfiles — all in one command. Cross-platform: macOS, Linux, and WSL.

## The Problem

Setting up a new dev machine takes hours. You install Homebrew, git, Node, Python, Docker, configure your shell, set up SSH keys, install VS Code extensions, copy dotfiles...

Existing tools solve part of the problem: chezmoi does dotfiles, devbox does Nix packages, mise handles runtime versions. **No single tool handles the full journey from fresh OS to "ready to code".**

## The Solution

```yaml
# devboot.yaml — your entire dev setup in one file
tools:
  - git
  - node@22
  - go
  - docker
  - ripgrep
  - fzf
  - lazygit

shell:
  type: zsh
  plugins:
    - zsh-autosuggestions
    - zsh-syntax-highlighting
  aliases:
    g: git
    k: kubectl
    dc: docker compose
    dev: npm run dev
  env:
    EDITOR: nvim

git:
  user.name: "Your Name"
  user.email: "you@example.com"
  init.defaultBranch: main
  pull.rebase: true
  aliases:
    co: checkout
    st: status
    lg: "log --oneline --graph --decorate"

vscode:
  extensions:
    - golang.go
    - esbenp.prettier-vscode
    - dbaeumer.vscode-eslint
  settings:
    editor.formatOnSave: true
    editor.tabSize: 2
```

```bash
devboot apply  # one command, everything installed and configured
```

## Install

**curl (Linux/macOS/WSL):**
```bash
curl -fsSL https://raw.githubusercontent.com/aymenhmaidiwastaken/devboot/main/install.sh | sh
```

**Go:**
```bash
go install github.com/aymenhmaidiwastaken/devboot@latest
```

**From source:**
```bash
git clone https://github.com/aymenhmaidiwastaken/devboot.git
cd devboot
make install
```

## Quick Start

```bash
# 1. Create your config (interactive wizard)
devboot init

# 2. Preview what would change
devboot diff

# 3. Apply everything
devboot apply

# 4. Verify your environment
devboot doctor
```

Or start from a profile:

```bash
# Apply a curated profile in one shot
devboot profile apply fullstack
```

## Profiles

DevBoot ships with 7 curated profiles — pre-built tool sets for common workflows:

| Profile | Tools | Description |
|---------|-------|-------------|
| `frontend` | node, bun, git, ripgrep, fzf | Frontend web development |
| `backend` | go, docker, kubectl, jq, curl, git | Backend & API development |
| `fullstack` | node, go, python, docker, kubectl, git, gh | Full-stack web development |
| `devops` | docker, kubectl, helm, terraform, jq, git, gh | DevOps & SRE toolkit |
| `data` | python, git, jq, curl | Data science & ML |
| `rust` | rust, git, ripgrep | Rust development |
| `terminal` | ripgrep, fzf, bat, eza, fd, zoxide, lazygit, neovim, tmux, starship | Power user terminal |

```bash
devboot profile list              # see all profiles
devboot profile show terminal     # inspect a profile
devboot profile apply terminal    # apply it
devboot profile export rust       # export as devboot.yaml
devboot profile search docker     # search by keyword
```

## Commands

| Command | Description |
|---------|-------------|
| `devboot apply [config.yaml]` | Apply full configuration |
| `devboot apply --only tools` | Apply only a specific section |
| `devboot init` | Interactive config wizard |
| `devboot init --plain` | Generate starter template |
| `devboot status` | Dashboard of installed vs configured |
| `devboot diff` | Preview what would change |
| `devboot doctor` | Diagnose environment issues with fix suggestions |
| `devboot export` | Reverse-engineer current machine to YAML |
| `devboot add [tool...]` | Install tools + update config |
| `devboot uninstall [tool...]` | Remove devboot-managed tools |
| `devboot history` | Audit trail of all actions |
| `devboot update` | Update all installed packages |
| `devboot profile <sub>` | Manage curated profiles |
| `devboot dotfiles push` | Push dotfile changes back to repo |
| `devboot version` | Print version info |

## Config Reference

### `tools`

List of tools to install. Supports version pinning with `@`:

```yaml
tools:
  - git
  - node@22
  - python@3.12
  - go
  - rust
  - docker
  - kubectl
  - terraform
  - ripgrep
  - fzf
  - lazygit
  - jq
  - bat
  - eza
  - fd
  - zoxide
  - starship
```

DevBoot automatically resolves the correct package name per platform (e.g., `node` becomes `nodejs` on apt, `node` on Homebrew). Dependencies are resolved automatically — requesting `lazygit` ensures `git` is installed first.

### `shell`

```yaml
shell:
  type: zsh              # zsh, bash, or fish
  plugins:
    - zsh-autosuggestions
    - zsh-syntax-highlighting
    - fzf-tab
  aliases:
    g: git
    k: kubectl
    ll: eza -la
  env:
    EDITOR: nvim
    GOPATH: ~/go
```

Plugins are installed via git clone (no oh-my-zsh dependency). Aliases and env vars are added idempotently — safe to run multiple times.

### `git`

```yaml
git:
  user.name: "Your Name"
  user.email: "you@example.com"
  init.defaultBranch: main
  pull.rebase: true
  sshKey: true            # generate SSH key if missing
  aliases:
    co: checkout
    br: branch
    st: status
    lg: "log --oneline --graph --decorate"
```

### `vscode`

```yaml
vscode:
  extensions:
    - golang.go
    - esbenp.prettier-vscode
    - ms-azuretools.vscode-docker
  settings:
    editor.fontSize: 14
    editor.tabSize: 2
    editor.formatOnSave: true
```

### `neovim`

```yaml
neovim:
  configRepo: https://github.com/aymenhmaidiwastaken/nvim-config.git
```

### `dotfiles`

```yaml
dotfiles:
  repo: https://github.com/aymenhmaidiwastaken/dotfiles.git
  mappings:
    .vimrc: ~/.vimrc
    .tmux.conf: ~/.tmux.conf
    starship.toml: ~/.config/starship.toml
  templates:
    .gitconfig: ~/.gitconfig   # supports {{.Name}}, {{.Email}}, {{.Home}}
```

## Team Sharing

```bash
# Export your setup
devboot export > devboot.yaml

# Share with your team (commit to repo)
git add devboot.yaml

# New team member applies it
devboot apply

# Or use a remote config
devboot apply https://raw.githubusercontent.com/team/configs/main/devboot.yaml
```

## Platform Support

| Platform | Package Manager | Status |
|----------|----------------|--------|
| macOS | Homebrew (auto-installs) | Supported |
| Ubuntu/Debian | apt | Supported |
| Arch Linux | pacman | Supported |
| WSL | apt | Supported |
| Fedora | dnf | Planned |

## How It Works

1. **Dependency resolution** — tools are installed in the correct order (e.g., git before lazygit, docker before kubectl)
2. **Idempotent** — every operation checks current state before acting. Safe to run repeatedly.
3. **State tracking** — DevBoot records every action. Use `devboot history` to audit and `devboot uninstall` to rollback.
4. **Post-install hooks** — tools like Rust automatically run `rustup default stable`, Docker adds your user to the docker group.
5. **Conflict detection** — warns if you have multiple Node version managers (nvm + fnm + volta).

## `devboot doctor`

The doctor command runs comprehensive diagnostics:

- PATH analysis (duplicates, missing directories)
- Version manager conflict detection
- Git authentication checks (SSH key, ssh-agent, GitHub connectivity)
- Shell config validation
- Actionable fix suggestions for every issue

```bash
$ devboot doctor

  → Git
  ✓ user.name: Jane Developer
  ✓ user.email: jane@example.com
  ⚠ SSH key not added to agent
    fix: ssh-add ~/.ssh/id_ed25519
  ✓ GitHub SSH authentication works

  → PATH analysis
  ✓ ~/.local/bin in PATH
  ⚠ 3 duplicate PATH entries
```

## vs Alternatives

| Feature | DevBoot | chezmoi | devbox | mise | Ansible |
|---------|---------|---------|--------|------|---------|
| Tool installation | Yes | No | Yes (Nix) | Yes (runtimes) | Yes |
| Shell config | Yes | Partial | No | No | Yes |
| Git setup + SSH | Yes | No | No | No | Yes |
| IDE extensions | Yes | No | No | No | Partial |
| Dotfiles sync | Yes | Yes | No | No | Yes |
| Interactive wizard | Yes | No | No | No | No |
| Profiles | Yes | No | No | No | Roles |
| Single binary | Yes | Yes | Yes | Yes | No (Python) |
| No dependencies | Yes | Yes | Nix | Yes | Python + SSH |
| Learning curve | Low | Medium | Medium | Low | High |

DevBoot is not a replacement for any of these tools — it's the glue that ties everything together. Use mise for version management, chezmoi for advanced dotfiles, and DevBoot to orchestrate the full setup.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

```bash
# Development
git clone https://github.com/aymenhmaidiwastaken/devboot.git
cd devboot
go build ./...
go test ./...

# Run locally
go run . doctor
go run . init --plain
```

## License

[MIT](LICENSE)
