package tui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			MarginBottom(1)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#06B6D4"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981"))

	failStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444"))

	skipStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))

	spinnerChars = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	progressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))
)

// TaskStatus represents the state of a task.
type TaskStatus int

const (
	TaskPending TaskStatus = iota
	TaskRunning
	TaskDone
	TaskSkipped
	TaskFailed
)

// Task is a single item being processed.
type Task struct {
	Name    string
	Section string
	Status  TaskStatus
	Message string
}

// ApplyModel is the bubbletea model for apply progress.
type ApplyModel struct {
	tasks      []Task
	mu         sync.Mutex
	sections   []string
	done       bool
	startTime  time.Time
	frame      int
	quitting   bool
	err        error
	workFunc   func(m *ApplyModel) error
	finished   chan struct{}
}

type tickMsg time.Time
type workDoneMsg struct{ err error }

// NewApplyModel creates a new apply TUI model.
func NewApplyModel(workFunc func(m *ApplyModel) error) *ApplyModel {
	return &ApplyModel{
		workFunc:  workFunc,
		startTime: time.Now(),
		finished:  make(chan struct{}),
	}
}

// AddTask adds a task to the display.
func (m *ApplyModel) AddTask(section, name string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Add section if new
	found := false
	for _, s := range m.sections {
		if s == section {
			found = true
			break
		}
	}
	if !found {
		m.sections = append(m.sections, section)
	}

	idx := len(m.tasks)
	m.tasks = append(m.tasks, Task{
		Name:    name,
		Section: section,
		Status:  TaskPending,
	})
	return idx
}

// UpdateTask updates a task's status.
func (m *ApplyModel) UpdateTask(idx int, status TaskStatus, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if idx >= 0 && idx < len(m.tasks) {
		m.tasks[idx].Status = status
		if message != "" {
			m.tasks[idx].Message = message
		}
	}
}

// SetRunning marks a task as running.
func (m *ApplyModel) SetRunning(idx int) {
	m.UpdateTask(idx, TaskRunning, "")
}

// SetDone marks a task as done.
func (m *ApplyModel) SetDone(idx int, msg string) {
	m.UpdateTask(idx, TaskDone, msg)
}

// SetSkipped marks a task as skipped.
func (m *ApplyModel) SetSkipped(idx int, msg string) {
	m.UpdateTask(idx, TaskSkipped, msg)
}

// SetFailed marks a task as failed.
func (m *ApplyModel) SetFailed(idx int, msg string) {
	m.UpdateTask(idx, TaskFailed, msg)
}

func (m *ApplyModel) Init() tea.Cmd {
	return tea.Batch(
		m.tick(),
		m.doWork(),
	)
}

func (m *ApplyModel) tick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *ApplyModel) doWork() tea.Cmd {
	return func() tea.Msg {
		err := m.workFunc(m)
		return workDoneMsg{err: err}
	}
}

func (m *ApplyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

	case tickMsg:
		m.frame++
		if m.done {
			return m, tea.Quit
		}
		return m, m.tick()

	case workDoneMsg:
		m.done = true
		m.err = msg.err
		return m, nil // wait for next tick to quit
	}

	return m, nil
}

func (m *ApplyModel) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("  devboot apply"))
	sb.WriteString("\n")

	m.mu.Lock()
	tasks := make([]Task, len(m.tasks))
	copy(tasks, m.tasks)
	sections := make([]string, len(m.sections))
	copy(sections, m.sections)
	m.mu.Unlock()

	for _, section := range sections {
		sb.WriteString(sectionStyle.Render(fmt.Sprintf("\n  ▸ %s", section)))
		sb.WriteString("\n")

		for _, task := range tasks {
			if task.Section != section {
				continue
			}

			var icon, line string
			switch task.Status {
			case TaskPending:
				icon = dimStyle.Render("○")
				line = dimStyle.Render(task.Name)
			case TaskRunning:
				spinner := spinnerChars[m.frame%len(spinnerChars)]
				icon = progressStyle.Render(spinner)
				line = progressStyle.Render(task.Name)
			case TaskDone:
				icon = successStyle.Render("✓")
				msg := task.Name
				if task.Message != "" {
					msg = task.Message
				}
				line = successStyle.Render(msg)
			case TaskSkipped:
				icon = skipStyle.Render("‣")
				msg := task.Name + " (already installed)"
				if task.Message != "" {
					msg = task.Message
				}
				line = skipStyle.Render(msg)
			case TaskFailed:
				icon = failStyle.Render("✗")
				msg := task.Name
				if task.Message != "" {
					msg += ": " + task.Message
				}
				line = failStyle.Render(msg)
			}

			sb.WriteString(fmt.Sprintf("    %s %s\n", icon, line))
		}
	}

	elapsed := time.Since(m.startTime).Round(time.Millisecond)

	if m.done {
		sb.WriteString("\n")
		if m.err != nil {
			sb.WriteString(failStyle.Render(fmt.Sprintf("  ✗ completed with errors (%s)", elapsed)))
		} else {
			sb.WriteString(successStyle.Render(fmt.Sprintf("  ✓ all done! (%s)", elapsed)))
		}
		sb.WriteString("\n\n")
	} else {
		sb.WriteString(dimStyle.Render(fmt.Sprintf("\n  elapsed: %s", elapsed)))
		sb.WriteString("\n")
	}

	return sb.String()
}

// RunApplyTUI runs the apply command with a TUI.
func RunApplyTUI(workFunc func(m *ApplyModel) error) error {
	model := NewApplyModel(workFunc)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}
	if fm, ok := finalModel.(*ApplyModel); ok && fm.err != nil {
		return fm.err
	}
	return nil
}
