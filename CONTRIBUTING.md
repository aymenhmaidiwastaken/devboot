# Contributing to DevBoot

Thanks for your interest in contributing! Here's how to get started.

## Development Setup

```bash
git clone https://github.com/aymenhmaidiwastaken/devboot.git
cd devboot
go build ./...
go test ./...
```

**Requirements:** Go 1.22+

## Making Changes

1. Fork the repo and create a branch from `main`
2. Make your changes
3. Add or update tests
4. Run `go test ./...` and `go vet ./...`
5. Open a pull request

## Project Structure

```
cmd/              CLI commands (cobra)
internal/
  config/         YAML config parser + validation
  tools/          Package manager backends (brew/apt/pacman)
  deps/           Tool registry + dependency graph
  shell/          Shell detection, aliases, plugins
  git/            Git config + SSH key generation
  editor/         VS Code, Neovim, JetBrains
  dotfiles/       Repo sync, symlinks, templates
  profile/        Built-in profiles
  state/          Action tracking + rollback
  tui/            Bubbletea interactive UI
  platform/       OS/distro detection
  ui/             Terminal output helpers
```

## Adding a New Tool to the Registry

Edit `internal/deps/graph.go` and add an entry to the `Registry` map:

```go
"mytool": {
    Name: "mytool", BinName: "mytool", Category: "cli",
    Description:  "What it does",
    Dependencies: []string{"git"},           // optional
    PostInstall:  []string{"mytool setup"},  // optional
},
```

Then add the package mapping in `internal/tools/registry.go`:

```go
"mytool": {platform.MacOS: "mytool", platform.Linux: "mytool", platform.WSL: "mytool"},
```

## Adding a New Profile

Edit `internal/profile/profile.go` and add to `BuiltinProfiles`:

```go
"myprofile": {
    Name:        "myprofile",
    Description: "Description here",
    Tools:       []string{"tool1", "tool2"},
    Shell: config.ShellConfig{
        Aliases: map[string]string{"alias": "command"},
    },
},
```

## Code Style

- Run `go vet ./...` before submitting
- Follow standard Go conventions
- Keep functions focused and small
- Add tests for new functionality

## Reporting Issues

- Use GitHub Issues
- Include `devboot doctor` output
- Include your OS and shell
