package ui

import (
	"fmt"
	"sync"
)

// Status icons for terminal output.
const (
	IconSuccess = "✓"
	IconSkip    = "✓"
	IconInstall = "↓"
	IconFail    = "✗"
	IconInfo    = "●"
	IconArrow   = "→"
)

// Colors using ANSI escape codes.
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
	colorDim    = "\033[2m"
	colorBold   = "\033[1m"
)

var mu sync.Mutex

func Success(msg string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Printf("  %s%s%s %s\n", colorGreen, IconSuccess, colorReset, msg)
}

func Skip(msg string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Printf("  %s%s already installed%s %s\n", colorDim, IconSkip, colorReset, msg)
}

func Installing(msg string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Printf("  %s%s%s %s\n", colorCyan, IconInstall, colorReset, msg)
}

func Fail(msg string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Printf("  %s%s%s %s\n", colorRed, IconFail, colorReset, msg)
}

func Info(msg string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Printf("  %s%s%s %s\n", colorCyan, IconInfo, colorReset, msg)
}

func Section(title string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Printf("\n%s%s %s%s\n", colorBold, IconArrow, title, colorReset)
}

func Warn(msg string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Printf("  %s⚠ %s%s\n", colorYellow, msg, colorReset)
}
