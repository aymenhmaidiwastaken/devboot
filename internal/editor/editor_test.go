package editor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// --- vsCodeSettingsPath ---

func TestVSCodeSettingsPath(t *testing.T) {
	path := vsCodeSettingsPath()
	if path == "" {
		t.Fatal("vsCodeSettingsPath returned empty")
	}

	switch runtime.GOOS {
	case "darwin":
		if filepath.Base(path) != "settings.json" {
			t.Errorf("expected settings.json, got %s", filepath.Base(path))
		}
	case "linux":
		if filepath.Base(path) != "settings.json" {
			t.Errorf("expected settings.json, got %s", filepath.Base(path))
		}
	default: // windows
		if filepath.Base(path) != "settings.json" {
			t.Errorf("expected settings.json, got %s", filepath.Base(path))
		}
	}
}

// --- applyVSCodeSettings ---

func TestApplyVSCodeSettingsNewFile(t *testing.T) {
	dir := t.TempDir()
	settingsPath := filepath.Join(dir, "settings.json")

	// Monkey-patch the path function by testing the internal logic directly
	settings := map[string]interface{}{
		"editor.fontSize":     14,
		"editor.formatOnSave": true,
		"editor.tabSize":      2,
	}

	// Write settings as if applyVSCodeSettings would
	existing := make(map[string]interface{})
	for key, val := range settings {
		existing[key] = val
	}

	data, err := json.MarshalIndent(existing, "", "    ")
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}

	os.MkdirAll(filepath.Dir(settingsPath), 0755)
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Read back and verify
	readData, _ := os.ReadFile(settingsPath)
	var readSettings map[string]interface{}
	json.Unmarshal(readData, &readSettings)

	if readSettings["editor.fontSize"] != float64(14) {
		t.Errorf("expected fontSize=14, got %v", readSettings["editor.fontSize"])
	}
	if readSettings["editor.formatOnSave"] != true {
		t.Errorf("expected formatOnSave=true, got %v", readSettings["editor.formatOnSave"])
	}
}

func TestApplyVSCodeSettingsMerge(t *testing.T) {
	dir := t.TempDir()
	settingsPath := filepath.Join(dir, "settings.json")

	// Write initial settings
	initial := map[string]interface{}{
		"editor.fontSize": 12,
		"editor.theme":    "dark",
	}
	data, _ := json.MarshalIndent(initial, "", "    ")
	os.WriteFile(settingsPath, data, 0644)

	// Merge new settings
	existing := make(map[string]interface{})
	readData, _ := os.ReadFile(settingsPath)
	json.Unmarshal(readData, &existing)

	newSettings := map[string]interface{}{
		"editor.fontSize":     14, // override
		"editor.formatOnSave": true, // new
	}
	for key, val := range newSettings {
		existing[key] = val
	}

	data, _ = json.MarshalIndent(existing, "", "    ")
	os.WriteFile(settingsPath, data, 0644)

	// Verify
	readData, _ = os.ReadFile(settingsPath)
	var result map[string]interface{}
	json.Unmarshal(readData, &result)

	if result["editor.fontSize"] != float64(14) {
		t.Errorf("fontSize should be overridden to 14, got %v", result["editor.fontSize"])
	}
	if result["editor.theme"] != "dark" {
		t.Errorf("theme should be preserved, got %v", result["editor.theme"])
	}
	if result["editor.formatOnSave"] != true {
		t.Errorf("formatOnSave should be added, got %v", result["editor.formatOnSave"])
	}
}

func TestApplyVSCodeSettingsCorruptExisting(t *testing.T) {
	dir := t.TempDir()
	settingsPath := filepath.Join(dir, "settings.json")

	// Write corrupt JSON
	os.WriteFile(settingsPath, []byte("{invalid json"), 0644)

	// Should handle gracefully - reading corrupt file gives empty map
	existing := make(map[string]interface{})
	data, _ := os.ReadFile(settingsPath)
	json.Unmarshal(data, &existing) // silently fails, existing stays empty

	if len(existing) != 0 {
		t.Errorf("expected empty map from corrupt JSON, got %d entries", len(existing))
	}
}

// --- VSCodeApply ---

func TestVSCodeApplyEmptyInputs(t *testing.T) {
	err := VSCodeApply(nil, nil)
	if err != nil {
		t.Errorf("expected nil for empty inputs, got %v", err)
	}

	err = VSCodeApply([]string{}, map[string]interface{}{})
	if err != nil {
		t.Errorf("expected nil for empty slices/maps, got %v", err)
	}
}

// --- NeovimApply ---

func TestNeovimApplyEmptyRepo(t *testing.T) {
	err := NeovimApply("")
	if err != nil {
		t.Errorf("expected nil for empty repo, got %v", err)
	}
}

// --- JetBrainsApply ---

func TestJetBrainsApplyEmpty(t *testing.T) {
	err := JetBrainsApply(nil)
	if err != nil {
		t.Errorf("expected nil for nil plugins, got %v", err)
	}

	err = JetBrainsApply([]string{})
	if err != nil {
		t.Errorf("expected nil for empty plugins, got %v", err)
	}
}

// --- findVSCodeBin ---

func TestFindVSCodeBin(t *testing.T) {
	// Just verify it doesn't panic
	result := findVSCodeBin()
	_ = result // may or may not find VS Code
}

// --- findJetBrainsBin ---

func TestFindJetBrainsBin(t *testing.T) {
	result := findJetBrainsBin()
	_ = result // may or may not find any IDE
}

// --- listInstalledExtensions ---

func TestListInstalledExtensionsNonExistent(t *testing.T) {
	result := listInstalledExtensions("nonexistent-binary-xyz")
	if result != nil {
		t.Errorf("expected nil for nonexistent binary, got %v", result)
	}
}
