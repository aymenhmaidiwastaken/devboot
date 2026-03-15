package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Action represents a single action devboot performed.
type Action struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`      // "install", "config", "symlink", "alias", "env", "plugin", "extension", "git"
	Section   string    `json:"section"`   // "tools", "shell", "git", "vscode", "dotfiles"
	Target    string    `json:"target"`    // what was acted on (e.g. "node", "alias g=git")
	Detail    string    `json:"detail"`    // extra info (e.g. package name, previous value)
	Status    string    `json:"status"`    // "ok", "skipped", "failed"
	Reversal  string    `json:"reversal"`  // command to undo this action
}

// Store holds the full devboot state.
type Store struct {
	Version     int       `json:"version"`
	LastApply   time.Time `json:"lastApply"`
	Actions     []Action  `json:"actions"`
	ManagedTools []string `json:"managedTools"` // tools installed by devboot
	path        string
}

func statePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "devboot", "state.json")
}

// Load reads the state from disk.
func Load() (*Store, error) {
	p := statePath()
	s := &Store{Version: 1, path: p}

	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	s.path = p
	return s, nil
}

// Save writes the state to disk.
func (s *Store) Save() error {
	os.MkdirAll(filepath.Dir(s.path), 0755)

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// Record adds an action to the state.
func (s *Store) Record(action Action) {
	if action.Timestamp.IsZero() {
		action.Timestamp = time.Now()
	}
	s.Actions = append(s.Actions, action)
}

// RecordTool records a tool installation.
func (s *Store) RecordTool(name, pkg, status string) {
	reversal := ""
	if status == "ok" {
		reversal = "uninstall:" + pkg
		// Track as managed
		found := false
		for _, t := range s.ManagedTools {
			if t == name {
				found = true
				break
			}
		}
		if !found {
			s.ManagedTools = append(s.ManagedTools, name)
		}
	}

	s.Record(Action{
		Type:     "install",
		Section:  "tools",
		Target:   name,
		Detail:   pkg,
		Status:   status,
		Reversal: reversal,
	})
}

// RecordConfig records a configuration change.
func (s *Store) RecordConfig(section, target, detail, reversal string) {
	s.Record(Action{
		Type:     "config",
		Section:  section,
		Target:   target,
		Detail:   detail,
		Status:   "ok",
		Reversal: reversal,
	})
}

// RecordSymlink records a symlink creation.
func (s *Store) RecordSymlink(src, dst string) {
	s.Record(Action{
		Type:     "symlink",
		Section:  "dotfiles",
		Target:   dst,
		Detail:   src,
		Status:   "ok",
		Reversal: "rm:" + dst,
	})
}

// MarkApply updates the last apply time.
func (s *Store) MarkApply() {
	s.LastApply = time.Now()
}

// History returns the last N actions.
func (s *Store) History(n int) []Action {
	if n <= 0 || n > len(s.Actions) {
		return s.Actions
	}
	return s.Actions[len(s.Actions)-n:]
}

// ActionsForSection returns all actions for a given section.
func (s *Store) ActionsForSection(section string) []Action {
	var result []Action
	for _, a := range s.Actions {
		if a.Section == section {
			result = append(result, a)
		}
	}
	return result
}

// RemoveTool removes a tool from managed list.
func (s *Store) RemoveTool(name string) {
	var filtered []string
	for _, t := range s.ManagedTools {
		if t != name {
			filtered = append(filtered, t)
		}
	}
	s.ManagedTools = filtered
}

// IsManagedTool checks if a tool was installed by devboot.
func (s *Store) IsManagedTool(name string) bool {
	for _, t := range s.ManagedTools {
		if t == name {
			return true
		}
	}
	return false
}
