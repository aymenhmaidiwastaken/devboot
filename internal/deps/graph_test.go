package deps

import (
	"sort"
	"testing"
)

// --- Resolve ---

func TestResolveSimple(t *testing.T) {
	order, err := Resolve([]string{"git", "fzf", "jq"})
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if len(order) != 3 {
		t.Errorf("expected 3 tools, got %d: %v", len(order), order)
	}
}

func TestResolveDependencies(t *testing.T) {
	order, err := Resolve([]string{"lazygit", "kubectl"})
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	indexOf := func(name string) int {
		for i, n := range order {
			if n == name {
				return i
			}
		}
		return -1
	}

	if indexOf("git") >= indexOf("lazygit") {
		t.Errorf("git (%d) should come before lazygit (%d)", indexOf("git"), indexOf("lazygit"))
	}
	if indexOf("docker") >= indexOf("kubectl") {
		t.Errorf("docker (%d) should come before kubectl (%d)", indexOf("docker"), indexOf("kubectl"))
	}
}

func TestResolveChainedDependencies(t *testing.T) {
	// helm -> kubectl -> docker
	order, err := Resolve([]string{"helm"})
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	indexOf := func(name string) int {
		for i, n := range order {
			if n == name {
				return i
			}
		}
		return -1
	}

	dockerIdx := indexOf("docker")
	kubectlIdx := indexOf("kubectl")
	helmIdx := indexOf("helm")

	if dockerIdx == -1 || kubectlIdx == -1 || helmIdx == -1 {
		t.Fatalf("expected docker, kubectl, helm in order, got %v", order)
	}

	if dockerIdx >= kubectlIdx {
		t.Error("docker should come before kubectl")
	}
	if kubectlIdx >= helmIdx {
		t.Error("kubectl should come before helm")
	}
}

func TestResolveWithVersion(t *testing.T) {
	order, err := Resolve([]string{"node@22", "git"})
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	found := false
	for _, o := range order {
		if o == "node@22" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected node@22 in order, got %v", order)
	}
}

func TestResolveNoDependencies(t *testing.T) {
	order, err := Resolve([]string{"git", "jq", "curl"})
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if len(order) != 3 {
		t.Errorf("expected 3, got %d", len(order))
	}
}

func TestResolveEmpty(t *testing.T) {
	order, err := Resolve([]string{})
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if len(order) != 0 {
		t.Errorf("expected 0, got %d", len(order))
	}
}

func TestResolveSingleTool(t *testing.T) {
	order, err := Resolve([]string{"git"})
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if len(order) != 1 || order[0] != "git" {
		t.Errorf("expected [git], got %v", order)
	}
}

func TestResolveUnknownTool(t *testing.T) {
	// Unknown tools should still resolve (no deps in registry)
	order, err := Resolve([]string{"unknown-tool-xyz"})
	if err != nil {
		t.Fatalf("Resolve failed for unknown tool: %v", err)
	}
	if len(order) != 1 || order[0] != "unknown-tool-xyz" {
		t.Errorf("expected [unknown-tool-xyz], got %v", order)
	}
}

func TestResolveDuplicateTools(t *testing.T) {
	order, err := Resolve([]string{"git", "git", "git"})
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	// Topological sort with visited set should deduplicate
	gitCount := 0
	for _, o := range order {
		if o == "git" {
			gitCount++
		}
	}
	if gitCount != 1 {
		t.Errorf("expected 1 git entry, got %d in %v", gitCount, order)
	}
}

func TestResolveDependencyAlreadyRequested(t *testing.T) {
	// lazygit depends on git; if both are requested, git should appear once
	order, err := Resolve([]string{"git", "lazygit"})
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	gitCount := 0
	for _, o := range order {
		if o == "git" {
			gitCount++
		}
	}
	if gitCount != 1 {
		t.Errorf("expected 1 git entry, got %d", gitCount)
	}
}

// --- GetBinName ---

func TestGetBinName(t *testing.T) {
	tests := map[string]string{
		"ripgrep":  "rg",
		"neovim":   "nvim",
		"python":   "python3",
		"git":      "git",
		"rust":     "rustc",
		"docker":   "docker",
		"bat":      "bat",
		"fd":       "fd",
		"unknown":  "unknown",
		"lazygit":  "lazygit",
		"starship": "starship",
	}
	for tool, expected := range tests {
		got := GetBinName(tool)
		if got != expected {
			t.Errorf("GetBinName(%q) = %q, want %q", tool, got, expected)
		}
	}
}

// --- GetPostInstall ---

func TestGetPostInstall(t *testing.T) {
	hooks := GetPostInstall("rust")
	if len(hooks) == 0 {
		t.Error("expected post-install hooks for rust")
	}

	hooks = GetPostInstall("docker")
	if len(hooks) == 0 {
		t.Error("expected post-install hooks for docker")
	}

	hooks = GetPostInstall("git")
	if len(hooks) != 0 {
		t.Error("expected no post-install hooks for git")
	}

	hooks = GetPostInstall("unknown-tool")
	if len(hooks) != 0 {
		t.Error("expected no hooks for unknown tool")
	}
}

// --- GetDependencies ---

func TestGetDependencies(t *testing.T) {
	deps := GetDependencies("lazygit")
	if len(deps) != 1 || deps[0] != "git" {
		t.Errorf("lazygit deps: expected [git], got %v", deps)
	}

	deps = GetDependencies("kubectl")
	if len(deps) != 1 || deps[0] != "docker" {
		t.Errorf("kubectl deps: expected [docker], got %v", deps)
	}

	deps = GetDependencies("helm")
	if len(deps) != 1 || deps[0] != "kubectl" {
		t.Errorf("helm deps: expected [kubectl], got %v", deps)
	}

	deps = GetDependencies("git")
	if len(deps) != 0 {
		t.Errorf("git deps: expected empty, got %v", deps)
	}

	deps = GetDependencies("unknown")
	if deps != nil {
		t.Errorf("unknown deps: expected nil, got %v", deps)
	}
}

// --- Conflicts ---

func TestConflicts(t *testing.T) {
	warnings := Conflicts([]string{"nvm", "fnm"})
	if len(warnings) == 0 {
		t.Error("expected conflict warning for nvm + fnm")
	}

	warnings = Conflicts([]string{"nvm", "fnm", "volta"})
	if len(warnings) == 0 {
		t.Error("expected conflict warning for 3 node managers")
	}

	warnings = Conflicts([]string{"git", "node"})
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}

	warnings = Conflicts([]string{})
	if len(warnings) != 0 {
		t.Errorf("expected no warnings for empty, got %v", warnings)
	}

	// Single manager is fine
	warnings = Conflicts([]string{"nvm", "git"})
	if len(warnings) != 0 {
		t.Errorf("single manager should be fine, got %v", warnings)
	}
}

// --- Categories ---

func TestCategories(t *testing.T) {
	cats := Categories()
	if len(cats) == 0 {
		t.Fatal("expected categories")
	}

	expectedCats := []string{"cli", "runtime", "devops", "editor", "vcs"}
	for _, cat := range expectedCats {
		if _, ok := cats[cat]; !ok {
			t.Errorf("expected category %q", cat)
		}
	}
}

func TestCategoriesContainAllTools(t *testing.T) {
	cats := Categories()
	total := 0
	for _, tools := range cats {
		total += len(tools)
	}
	if total != len(Registry) {
		t.Errorf("categories contain %d tools, registry has %d", total, len(Registry))
	}
}

// --- AllToolNames ---

func TestAllToolNames(t *testing.T) {
	names := AllToolNames()
	if len(names) != len(Registry) {
		t.Errorf("expected %d names, got %d", len(Registry), len(names))
	}
}

func TestAllToolNamesContainsKnownTools(t *testing.T) {
	names := AllToolNames()
	nameSet := make(map[string]bool)
	for _, n := range names {
		nameSet[n] = true
	}

	for _, tool := range []string{"git", "node", "python", "docker", "kubectl"} {
		if !nameSet[tool] {
			t.Errorf("expected %q in AllToolNames", tool)
		}
	}
}

// --- Registry consistency ---

func TestRegistryAllHaveBinName(t *testing.T) {
	for name, info := range Registry {
		if info.BinName == "" {
			t.Errorf("Registry[%q] has empty BinName", name)
		}
	}
}

func TestRegistryAllHaveCategory(t *testing.T) {
	for name, info := range Registry {
		if info.Category == "" {
			t.Errorf("Registry[%q] has empty Category", name)
		}
	}
}

func TestRegistryAllHaveDescription(t *testing.T) {
	for name, info := range Registry {
		if info.Description == "" {
			t.Errorf("Registry[%q] has empty Description", name)
		}
	}
}

func TestRegistryDependenciesExist(t *testing.T) {
	for name, info := range Registry {
		for _, dep := range info.Dependencies {
			if _, ok := Registry[dep]; !ok {
				t.Errorf("Registry[%q] depends on %q which is not in the registry", name, dep)
			}
		}
	}
}

func TestRegistryNoCyclesInDeps(t *testing.T) {
	// Attempt to resolve all tools at once - should not panic or error
	var all []string
	for name := range Registry {
		all = append(all, name)
	}
	sort.Strings(all)

	_, err := Resolve(all)
	if err != nil {
		t.Errorf("Resolve(all tools) failed: %v — possible cycle in dependencies", err)
	}
}
