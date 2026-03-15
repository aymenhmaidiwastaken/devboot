package ui

import (
	"testing"
)

// The UI functions write to stdout. We mainly verify they don't panic
// and that the constants are properly defined.

func TestIconConstants(t *testing.T) {
	icons := map[string]string{
		"IconSuccess": IconSuccess,
		"IconSkip":    IconSkip,
		"IconInstall": IconInstall,
		"IconFail":    IconFail,
		"IconInfo":    IconInfo,
		"IconArrow":   IconArrow,
	}
	for name, icon := range icons {
		if icon == "" {
			t.Errorf("%s is empty", name)
		}
	}
}

func TestColorConstants(t *testing.T) {
	colors := []string{colorReset, colorGreen, colorYellow, colorRed, colorCyan, colorDim, colorBold}
	for i, c := range colors {
		if c == "" {
			t.Errorf("color constant at index %d is empty", i)
		}
		// All ANSI codes should start with ESC
		if c[0] != '\033' {
			t.Errorf("color constant at index %d doesn't start with ESC: %q", i, c)
		}
	}
}

func TestSuccessDoesNotPanic(t *testing.T) {
	Success("test message")
	Success("")
	Success("message with special chars: <>&\"'")
}

func TestSkipDoesNotPanic(t *testing.T) {
	Skip("test skip")
	Skip("")
}

func TestInstallingDoesNotPanic(t *testing.T) {
	Installing("installing something...")
	Installing("")
}

func TestFailDoesNotPanic(t *testing.T) {
	Fail("something failed")
	Fail("")
}

func TestInfoDoesNotPanic(t *testing.T) {
	Info("info message")
	Info("")
}

func TestSectionDoesNotPanic(t *testing.T) {
	Section("Test Section")
	Section("")
}

func TestWarnDoesNotPanic(t *testing.T) {
	Warn("warning message")
	Warn("")
}

func TestConcurrentOutputDoesNotPanic(t *testing.T) {
	// Verify the mutex doesn't deadlock under concurrent access
	done := make(chan bool, 7)

	go func() { Success("concurrent 1"); done <- true }()
	go func() { Fail("concurrent 2"); done <- true }()
	go func() { Skip("concurrent 3"); done <- true }()
	go func() { Info("concurrent 4"); done <- true }()
	go func() { Installing("concurrent 5"); done <- true }()
	go func() { Section("concurrent 6"); done <- true }()
	go func() { Warn("concurrent 7"); done <- true }()

	for i := 0; i < 7; i++ {
		<-done
	}
}
