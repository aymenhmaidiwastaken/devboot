package git

import (
	"testing"

	"github.com/aymenhmaidiwastaken/devboot/internal/config"
)

func TestConfigureEmptyConfig(t *testing.T) {
	// Empty config should be a no-op
	cfg := config.GitConfig{}
	err := Configure(cfg)
	if err != nil {
		t.Errorf("empty config should not error: %v", err)
	}
}

func TestGetGitConfigReturnsEmptyForNonexistent(t *testing.T) {
	val := getGitConfig("devboot.nonexistent.key.xyz12345")
	if val != "" {
		t.Errorf("expected empty for nonexistent key, got %q", val)
	}
}

func TestGetGitConfigUserName(t *testing.T) {
	// This reads the actual global git config
	val := getGitConfig("user.name")
	// We can't assert a specific value, just that it doesn't panic
	_ = val
}

func TestGetGitConfigUserEmail(t *testing.T) {
	val := getGitConfig("user.email")
	_ = val
}

func TestConfigureWithAliasesOnly(t *testing.T) {
	cfg := config.GitConfig{
		Aliases: map[string]string{
			"devboot-test-alias": "status",
		},
	}

	// This will actually set the git alias in the global config
	// We test and then clean up
	err := Configure(cfg)
	if err != nil {
		t.Errorf("Configure with aliases failed: %v", err)
	}

	// Verify it was set
	val := getGitConfig("alias.devboot-test-alias")
	if val != "status" {
		t.Errorf("expected alias 'status', got %q", val)
	}

	// Clean up
	setGitConfig("alias.devboot-test-alias", "")
}

func TestSetGitConfig(t *testing.T) {
	key := "devboot.test.value"
	value := "test123"

	err := setGitConfig(key, value)
	if err != nil {
		t.Fatalf("setGitConfig failed: %v", err)
	}

	got := getGitConfig(key)
	if got != value {
		t.Errorf("expected %q, got %q", value, got)
	}

	// Clean up by unsetting
	setGitConfig(key, "")
}

func TestConfigureIdempotent(t *testing.T) {
	cfg := config.GitConfig{
		Aliases: map[string]string{
			"devboot-idem-test": "log --oneline",
		},
	}

	// Run twice — second should be a no-op
	Configure(cfg)
	err := Configure(cfg)
	if err != nil {
		t.Errorf("second Configure should not error: %v", err)
	}

	// Clean up
	setGitConfig("alias.devboot-idem-test", "")
}

func TestConfigureWithPullRebase(t *testing.T) {
	tr := true
	cfg := config.GitConfig{
		PullRebase: &tr,
	}

	err := Configure(cfg)
	if err != nil {
		t.Errorf("Configure with pull.rebase failed: %v", err)
	}
}
