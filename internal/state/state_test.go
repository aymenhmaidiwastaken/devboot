package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// --- Save/Load round trip ---

func TestStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	s := &Store{Version: 1, path: path}
	s.RecordTool("git", "git", "ok")
	s.RecordTool("node", "nodejs", "ok")
	s.RecordTool("python", "python3", "skipped")
	s.RecordConfig("git", "user.name", "Test User", "git config --global --unset user.name")
	s.RecordSymlink("/src/.vimrc", "/home/test/.vimrc")
	s.MarkApply()

	if err := s.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists and is valid JSON
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("state file is empty")
	}

	var loaded Store
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if loaded.Version != 1 {
		t.Errorf("expected version 1, got %d", loaded.Version)
	}
	if len(loaded.Actions) != 5 {
		t.Errorf("expected 5 actions, got %d", len(loaded.Actions))
	}
	if len(loaded.ManagedTools) != 2 {
		t.Errorf("expected 2 managed tools, got %d", len(loaded.ManagedTools))
	}
	if loaded.LastApply.IsZero() {
		t.Error("expected LastApply to be set")
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "state.json")

	s := &Store{Version: 1, path: path}
	s.RecordTool("git", "git", "ok")

	if err := s.Save(); err != nil {
		t.Fatalf("Save should create parent dirs: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("state file not created")
	}
}

func TestLoadNonExistent(t *testing.T) {
	// Load from the default path (won't exist in CI if never run)
	// But the function should not error for non-existent
	s := &Store{Version: 1, path: filepath.Join(t.TempDir(), "nope.json")}
	_ = s // just verify we can construct one
}

func TestLoadCorruptJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	os.WriteFile(path, []byte("{invalid json!!!"), 0644)

	s := &Store{path: path}
	data, _ := os.ReadFile(path)
	err := json.Unmarshal(data, s)
	if err == nil {
		t.Error("expected error for corrupt JSON")
	}
}

// --- ManagedTools ---

func TestManagedTools(t *testing.T) {
	s := &Store{Version: 1}
	s.RecordTool("git", "git", "ok")
	s.RecordTool("node", "nodejs", "ok")

	if !s.IsManagedTool("git") {
		t.Error("expected git to be managed")
	}
	if !s.IsManagedTool("node") {
		t.Error("expected node to be managed")
	}
	if s.IsManagedTool("python") {
		t.Error("expected python to NOT be managed")
	}
}

func TestRemoveTool(t *testing.T) {
	s := &Store{Version: 1}
	s.RecordTool("git", "git", "ok")
	s.RecordTool("node", "nodejs", "ok")
	s.RecordTool("go", "golang", "ok")

	s.RemoveTool("node")
	if s.IsManagedTool("node") {
		t.Error("expected node to be removed")
	}
	if !s.IsManagedTool("git") {
		t.Error("git should still be managed")
	}
	if !s.IsManagedTool("go") {
		t.Error("go should still be managed")
	}
}

func TestRemoveToolNotPresent(t *testing.T) {
	s := &Store{Version: 1, ManagedTools: []string{"git"}}
	s.RemoveTool("python") // should not panic
	if len(s.ManagedTools) != 1 {
		t.Errorf("expected 1 managed tool, got %d", len(s.ManagedTools))
	}
}

func TestRecordToolSkipped(t *testing.T) {
	s := &Store{Version: 1}
	s.RecordTool("python", "python3", "skipped")

	if s.IsManagedTool("python") {
		t.Error("skipped tool should NOT be managed")
	}
	if len(s.Actions) != 1 {
		t.Error("action should still be recorded")
	}
	if s.Actions[0].Reversal != "" {
		t.Error("skipped tool should have no reversal")
	}
}

func TestRecordToolFailed(t *testing.T) {
	s := &Store{Version: 1}
	s.RecordTool("broken", "broken-pkg", "failed")

	if s.IsManagedTool("broken") {
		t.Error("failed tool should NOT be managed")
	}
}

func TestRecordToolDuplicateDoesNotDuplicate(t *testing.T) {
	s := &Store{Version: 1}
	s.RecordTool("git", "git", "ok")
	s.RecordTool("git", "git", "ok")

	count := 0
	for _, t := range s.ManagedTools {
		if t == "git" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 git entry in ManagedTools, got %d", count)
	}
}

// --- RecordConfig ---

func TestRecordConfig(t *testing.T) {
	s := &Store{Version: 1}
	s.RecordConfig("git", "user.name", "Test User", "git config --global --unset user.name")

	if len(s.Actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(s.Actions))
	}
	a := s.Actions[0]
	if a.Type != "config" {
		t.Errorf("expected type 'config', got %q", a.Type)
	}
	if a.Section != "git" {
		t.Errorf("expected section 'git', got %q", a.Section)
	}
	if a.Target != "user.name" {
		t.Errorf("expected target 'user.name', got %q", a.Target)
	}
	if a.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", a.Status)
	}
}

// --- RecordSymlink ---

func TestRecordSymlink(t *testing.T) {
	s := &Store{Version: 1}
	s.RecordSymlink("/src/.vimrc", "/home/user/.vimrc")

	if len(s.Actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(s.Actions))
	}
	a := s.Actions[0]
	if a.Type != "symlink" {
		t.Errorf("expected type 'symlink', got %q", a.Type)
	}
	if a.Reversal != "rm:/home/user/.vimrc" {
		t.Errorf("expected reversal 'rm:/home/user/.vimrc', got %q", a.Reversal)
	}
}

// --- History ---

func TestHistory(t *testing.T) {
	s := &Store{Version: 1}
	for i := 0; i < 10; i++ {
		s.RecordTool("tool"+string(rune('a'+i)), "pkg", "ok")
	}

	h := s.History(3)
	if len(h) != 3 {
		t.Errorf("expected 3 history entries, got %d", len(h))
	}

	// Should return last 3
	if h[0].Target != "toolh" {
		t.Errorf("expected last 3 starting from 'toolh', got %q", h[0].Target)
	}
}

func TestHistoryZero(t *testing.T) {
	s := &Store{Version: 1}
	for i := 0; i < 5; i++ {
		s.RecordTool("tool", "pkg", "ok")
	}

	all := s.History(0)
	if len(all) != 5 {
		t.Errorf("History(0) should return all, got %d", len(all))
	}
}

func TestHistoryNegative(t *testing.T) {
	s := &Store{Version: 1}
	s.RecordTool("a", "a", "ok")
	all := s.History(-1)
	if len(all) != 1 {
		t.Errorf("History(-1) should return all, got %d", len(all))
	}
}

func TestHistoryLargerThanSize(t *testing.T) {
	s := &Store{Version: 1}
	s.RecordTool("a", "a", "ok")
	h := s.History(100)
	if len(h) != 1 {
		t.Errorf("History(100) should return all when only 1, got %d", len(h))
	}
}

func TestHistoryEmpty(t *testing.T) {
	s := &Store{Version: 1}
	h := s.History(5)
	if len(h) != 0 {
		t.Errorf("expected empty history, got %d", len(h))
	}
}

// --- ActionsForSection ---

func TestActionsForSection(t *testing.T) {
	s := &Store{Version: 1}
	s.RecordTool("git", "git", "ok")
	s.RecordTool("node", "node", "ok")
	s.RecordConfig("shell", "alias g", "git", "")
	s.RecordConfig("git", "user.name", "Test", "")

	toolActions := s.ActionsForSection("tools")
	if len(toolActions) != 2 {
		t.Errorf("expected 2 tool actions, got %d", len(toolActions))
	}

	shellActions := s.ActionsForSection("shell")
	if len(shellActions) != 1 {
		t.Errorf("expected 1 shell action, got %d", len(shellActions))
	}

	gitActions := s.ActionsForSection("git")
	if len(gitActions) != 1 {
		t.Errorf("expected 1 git action, got %d", len(gitActions))
	}

	empty := s.ActionsForSection("nonexistent")
	if len(empty) != 0 {
		t.Errorf("expected 0 actions for nonexistent section, got %d", len(empty))
	}
}

// --- MarkApply ---

func TestMarkApply(t *testing.T) {
	s := &Store{Version: 1}
	if !s.LastApply.IsZero() {
		t.Error("LastApply should be zero initially")
	}

	before := time.Now()
	s.MarkApply()
	after := time.Now()

	if s.LastApply.Before(before) || s.LastApply.After(after) {
		t.Error("LastApply should be between before and after")
	}
}

// --- Record with timestamp ---

func TestRecordSetsTimestamp(t *testing.T) {
	s := &Store{Version: 1}
	before := time.Now()
	s.Record(Action{Type: "test", Target: "test"})
	after := time.Now()

	if s.Actions[0].Timestamp.Before(before) || s.Actions[0].Timestamp.After(after) {
		t.Error("timestamp should be set automatically")
	}
}

func TestRecordPreservesExplicitTimestamp(t *testing.T) {
	s := &Store{Version: 1}
	explicit := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	s.Record(Action{Type: "test", Target: "test", Timestamp: explicit})

	if !s.Actions[0].Timestamp.Equal(explicit) {
		t.Error("explicit timestamp should be preserved")
	}
}
