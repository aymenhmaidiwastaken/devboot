package profile

import (
	"testing"

	"github.com/aymenhmaidiwastaken/devboot/internal/config"
)

func TestGetProfile(t *testing.T) {
	p, ok := Get("frontend")
	if !ok {
		t.Fatal("expected frontend profile")
	}
	if p.Name != "frontend" {
		t.Errorf("expected name frontend, got %s", p.Name)
	}
	if len(p.Tools) == 0 {
		t.Error("expected tools in frontend profile")
	}
}

func TestGetUnknown(t *testing.T) {
	_, ok := Get("nonexistent")
	if ok {
		t.Error("expected false for unknown profile")
	}
}

func TestList(t *testing.T) {
	names := List()
	if len(names) < 5 {
		t.Errorf("expected at least 5 profiles, got %d", len(names))
	}
}

func TestSearch(t *testing.T) {
	results := Search("docker")
	if len(results) == 0 {
		t.Error("expected search results for 'docker'")
	}

	results = Search("nonexistent_xyz")
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestMergeIntoConfig(t *testing.T) {
	cfg := &config.Config{
		Tools: []string{"git"},
		Shell: config.ShellConfig{
			Aliases: map[string]string{"g": "git"},
		},
	}

	p, _ := Get("frontend")
	MergeIntoConfig(cfg, p)

	// Should have merged tools without duplicates
	gitCount := 0
	for _, t := range cfg.Tools {
		if t == "git" {
			gitCount++
		}
	}
	if gitCount != 1 {
		t.Errorf("expected 1 git entry, got %d", gitCount)
	}

	// Should have merged aliases
	if _, ok := cfg.Shell.Aliases["nr"]; !ok {
		t.Error("expected 'nr' alias from frontend profile")
	}

	// Original alias should be preserved
	if cfg.Shell.Aliases["g"] != "git" {
		t.Error("expected original 'g' alias to be preserved")
	}
}

func TestToConfig(t *testing.T) {
	p, _ := Get("backend")
	cfg := ToConfig(p)

	if len(cfg.Tools) == 0 {
		t.Error("expected tools in config")
	}
}

func TestDescribe(t *testing.T) {
	p, _ := Get("terminal")
	desc := Describe(p)
	if desc == "" {
		t.Error("expected description")
	}
}
