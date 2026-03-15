package dotfiles

import (
	"bytes"
	"os"
	"os/user"
	"path/filepath"
	"text/template"

	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
)

// TemplateVars holds variables available in dotfile templates.
type TemplateVars struct {
	Name     string
	Email    string
	Home     string
	Hostname string
	OS       string
	User     string
}

// DefaultTemplateVars returns a TemplateVars populated with current system info.
func DefaultTemplateVars(name, email, osName string) TemplateVars {
	home, _ := os.UserHomeDir()
	hostname, _ := os.Hostname()
	username := ""
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	return TemplateVars{
		Name:     name,
		Email:    email,
		Home:     home,
		Hostname: hostname,
		OS:       osName,
		User:     username,
	}
}

// RenderFile processes a template file and writes the result to dst.
func RenderFile(src, dst string, vars TemplateVars) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	tmpl, err := template.New(filepath.Base(src)).
		Delims("{{", "}}").
		Option("missingkey=zero").
		Parse(string(data))
	if err != nil {
		// Not a template, skip rendering
		return nil
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(dst), 0755)
	if err := os.WriteFile(dst, buf.Bytes(), 0644); err != nil {
		return err
	}

	ui.Success("rendered " + filepath.Base(src))
	return nil
}
