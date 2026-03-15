package dotfiles

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderFile(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "template.conf")
	os.WriteFile(src, []byte(`# Config for {{.Name}}
email = {{.Email}}
home = {{.Home}}
`), 0644)

	dst := filepath.Join(dir, "output.conf")
	vars := TemplateVars{
		Name:  "Test User",
		Email: "test@example.com",
		Home:  "/home/test",
	}

	if err := RenderFile(src, dst, vars); err != nil {
		t.Fatalf("RenderFile failed: %v", err)
	}

	data, _ := os.ReadFile(dst)
	content := string(data)

	if !strings.Contains(content, "Test User") {
		t.Error("expected rendered content to contain Name")
	}
	if !strings.Contains(content, "test@example.com") {
		t.Error("expected rendered content to contain Email")
	}
	if !strings.Contains(content, "/home/test") {
		t.Error("expected rendered content to contain Home")
	}
}

func TestRenderFileEmptyVars(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "template.conf")
	os.WriteFile(src, []byte(`name = {{.Name}}, email = {{.Email}}`), 0644)

	dst := filepath.Join(dir, "output.conf")
	vars := TemplateVars{Name: "Test"}

	// Empty Email should render as empty string
	if err := RenderFile(src, dst, vars); err != nil {
		t.Fatalf("RenderFile failed: %v", err)
	}

	data, _ := os.ReadFile(dst)
	if !strings.Contains(string(data), "name = Test") {
		t.Error("expected Name to be rendered")
	}
}

func TestDefaultTemplateVars(t *testing.T) {
	vars := DefaultTemplateVars("Test", "test@test.com", "Linux")
	if vars.Name != "Test" {
		t.Errorf("expected Name=Test, got %s", vars.Name)
	}
	if vars.Home == "" {
		t.Error("expected Home to be set")
	}
}
