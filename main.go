package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Application states
type appState int

const (
	stateMenu appState = iota
	stateRecording
	stateCompleted
	stateError
)

// Messages
type tickMsg time.Time
type progressMsg struct {
	current int
	total   int
	percent float64
}
type logMsg string
type completedMsg struct {
	success bool
	message string
}
type processExitMsg struct {
	err error
}

// Model represents the application state
type model struct {
	state      appState
	interval   textinput.Model
	duration   textinput.Model
	output     textinput.Model
	focusIndex int
	inputs     []textinput.Model
	spinner    spinner.Model
	cmd        *exec.Cmd

	// Recording state
	startTime     time.Time
	progress      progressMsg
	logs          []string
	recordingDone bool
	finalMessage  string
	err           error
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D7D7D"))

	cursorStyle = focusedStyle.Copy()

	noStyle = lipgloss.NewStyle()

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	progressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))

	logStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

func initialModel() model {
	m := model{
		state:   stateMenu,
		inputs:  make([]textinput.Model, 3),
		spinner: spinner.New(),
		logs:    make([]string, 0),
	}

	// Setup interval input
	m.interval = textinput.New()
	m.interval.Placeholder = "5"
	m.interval.Focus()
	m.interval.CharLimit = 10
	m.interval.Width = 30
	m.interval.Prompt = "│ "
	m.interval.TextStyle = focusedStyle
	m.interval.PromptStyle = focusedStyle

	// Setup duration input
	m.duration = textinput.New()
	m.duration.Placeholder = "60"
	m.duration.CharLimit = 10
	m.duration.Width = 30
	m.duration.Prompt = "│ "

	// Setup output input
	m.output = textinput.New()
	m.output.Placeholder = "timelapse.mp4"
	m.output.CharLimit = 256
	m.output.Width = 30
	m.output.Prompt = "│ "

	m.inputs[0] = m.interval
	m.inputs[1] = m.duration
	m.inputs[2] = m.output

	m.spinner.Spinner = spinner.Dot
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateMenu:
			return m.updateMenu(msg)
		case stateRecording:
			return m.updateRecording(msg)
		case stateCompleted, stateError:
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		}

	case progressMsg:
		m.progress = msg
		return m, nil

	case logMsg:
		m.logs = append(m.logs, string(msg))
		// Keep only last 5 logs
		if len(m.logs) > 5 {
			m.logs = m.logs[1:]
		}
		return m, nil

	case completedMsg:
		if msg.success {
			m.state = stateCompleted
		} else {
			m.state = stateError
		}
		m.recordingDone = true
		m.finalMessage = msg.message
		return m, nil

	case processExitMsg:
		if msg.err != nil && !m.recordingDone {
			m.state = stateError
			m.err = msg.err
			m.finalMessage = fmt.Sprintf("Process error: %v", msg.err)
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tickMsg:
		return m, tick()
	}

	return m, nil
}

func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit

	case "tab", "shift+tab", "enter", "up", "down":
		s := msg.String()

		if s == "enter" {
			// Start recording
			return m.startRecording()
		}

		// Cycle through inputs
		if s == "up" || s == "shift+tab" {
			m.focusIndex--
		} else {
			m.focusIndex++
		}

		if m.focusIndex > len(m.inputs) {
			m.focusIndex = 0
		} else if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs)
		}

		cmds := make([]tea.Cmd, len(m.inputs))
		for i := 0; i <= len(m.inputs)-1; i++ {
			if i == m.focusIndex {
				cmds[i] = m.inputs[i].Focus()
				m.inputs[i].PromptStyle = focusedStyle
				m.inputs[i].TextStyle = focusedStyle
			} else {
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}
		}

		return m, tea.Batch(cmds...)
	}

	// Handle character input for focused field
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m model) updateRecording(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		// Stop recording
		if m.cmd != nil && m.cmd.Process != nil {
			m.cmd.Process.Signal(os.Interrupt)
		}
		return m, nil
	}
	return m, nil
}

func (m *model) updateInputs(msg tea.KeyMsg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) startRecording() (tea.Model, tea.Cmd) {
	// Validate inputs
	interval := m.inputs[0].Value()
	if interval == "" {
		interval = "5"
	}
	intervalFloat, err := strconv.ParseFloat(interval, 64)
	if err != nil || intervalFloat <= 0 {
		m.state = stateError
		m.finalMessage = "Invalid interval value"
		return m, nil
	}

	duration := m.inputs[1].Value()
	if duration == "" {
		duration = "60"
	}
	durationFloat, err := strconv.ParseFloat(duration, 64)
	if err != nil || durationFloat <= 0 {
		m.state = stateError
		m.finalMessage = "Invalid duration value (use seconds)"
		return m, nil
	}

	output := m.inputs[2].Value()
	if output == "" {
		output = "timelapse.mp4"
	}

	// Change state to recording
	m.state = stateRecording
	m.startTime = time.Now()
	m.recordingDone = false

	// Start the Python subprocess
	return m, tea.Batch(
		m.spinner.Tick,
		tick(),
		m.runTimelapse(intervalFloat, durationFloat, output),
	)
}

func (m model) runTimelapse(interval, duration float64, output string) tea.Cmd {
	return func() tea.Msg {
		// Build command
		cmd := exec.Command(
			"python3",
			"timelapse.py",
			"-i", fmt.Sprintf("%.2f", interval),
			"-d", fmt.Sprintf("%.2f", duration),
			"-o", output,
		)

		// Get stdout and stderr pipes
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return completedMsg{false, fmt.Sprintf("Failed to create stdout pipe: %v", err)}
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return completedMsg{false, fmt.Sprintf("Failed to create stderr pipe: %v", err)}
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			return completedMsg{false, fmt.Sprintf("Failed to start process: %v", err)}
		}

		// Read output in goroutines
		go m.readOutput(stdout)
		go m.readOutput(stderr)

		// Wait for completion
		err = cmd.Wait()

		if err != nil {
			return completedMsg{false, fmt.Sprintf("Recording failed: %v", err)}
		}

		return completedMsg{true, fmt.Sprintf("Timelapse saved to: %s", output)}
	}
}

func (m model) readOutput(pipe io.ReadCloser) tea.Cmd {
	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {
		line := scanner.Text()

		// Parse progress messages
		if strings.Contains(line, "[PROGRESS]") {
			// Format: [PROGRESS] 5/120 (4.2%)
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				counts := strings.Split(parts[1], "/")
				if len(counts) == 2 {
					current, _ := strconv.Atoi(counts[0])
					total, _ := strconv.Atoi(counts[1])
					percent := 0.0
					if len(parts) >= 3 {
						percentStr := strings.Trim(parts[2], "(%))")
						percent, _ = strconv.ParseFloat(percentStr, 64)
					}

					// Send progress message
					go func() {
						tea.Printf("%v", progressMsg{current, total, percent})
					}()
				}
			}
		}

		// Send log message for all lines
		go func(l string) {
			tea.Printf("%v", logMsg(l))
		}(line)
	}

	return nil
}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) View() string {
	switch m.state {
	case stateMenu:
		return m.viewMenu()
	case stateRecording:
		return m.viewRecording()
	case stateCompleted:
		return m.viewCompleted()
	case stateError:
		return m.viewError()
	}
	return ""
}

func (m model) viewMenu() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("⏱  Chronapse Timelapse Recorder"))
	b.WriteString("\n\n")

	// Interval input
	label := "Interval (seconds):"
	if m.focusIndex == 0 {
		label = focusedStyle.Render("▸ " + label)
	} else {
		label = blurredStyle.Render("  " + label)
	}
	b.WriteString(label + "\n")
	b.WriteString(m.inputs[0].View() + "\n\n")

	// Duration input
	label = "Duration (seconds):"
	if m.focusIndex == 1 {
		label = focusedStyle.Render("▸ " + label)
	} else {
		label = blurredStyle.Render("  " + label)
	}
	b.WriteString(label + "\n")
	b.WriteString(m.inputs[1].View() + "\n\n")

	// Output input
	label = "Output file:"
	if m.focusIndex == 2 {
		label = focusedStyle.Render("▸ " + label)
	} else {
		label = blurredStyle.Render("  " + label)
	}
	b.WriteString(label + "\n")
	b.WriteString(m.inputs[2].View() + "\n\n")

	// Start button
	button := "[ Start Recording ]"
	if m.focusIndex == 3 {
		button = focusedStyle.Render("▸ " + button)
	} else {
		button = blurredStyle.Render("  " + button)
	}
	b.WriteString(button + "\n")

	b.WriteString(helpStyle.Render("\nTab: next • Enter: start • Ctrl+C: quit"))

	return "\n" + b.String() + "\n"
}

func (m model) viewRecording() string {
	var b strings.Builder

	elapsed := time.Since(m.startTime)

	b.WriteString(titleStyle.Render("⏱  Recording in Progress"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("%s Recording...\n\n", m.spinner.View()))

	// Progress bar
	if m.progress.total > 0 {
		barWidth := 40
		filled := int(float64(barWidth) * m.progress.percent / 100.0)
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

		b.WriteString(progressStyle.Render(fmt.Sprintf("Progress: [%s] %.1f%%\n", bar, m.progress.percent)))
		b.WriteString(progressStyle.Render(fmt.Sprintf("Frames:   %d / %d\n", m.progress.current, m.progress.total)))
	}

	b.WriteString(fmt.Sprintf("\nElapsed:  %s\n", elapsed.Round(time.Second)))

	// Show recent logs
	if len(m.logs) > 0 {
		b.WriteString("\n" + logStyle.Render("Recent activity:") + "\n")
		for _, log := range m.logs {
			if len(log) > 80 {
				log = log[:77] + "..."
			}
			b.WriteString(logStyle.Render("  "+log) + "\n")
		}
	}

	b.WriteString(helpStyle.Render("\nPress 'q' to stop recording"))

	return "\n" + b.String() + "\n"
}

func (m model) viewCompleted() string {
	var b strings.Builder

	elapsed := time.Since(m.startTime)

	b.WriteString(titleStyle.Render("✓ Recording Complete"))
	b.WriteString("\n\n")

	b.WriteString(successStyle.Render(m.finalMessage))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("Total time: %s\n", elapsed.Round(time.Second)))

	if m.progress.total > 0 {
		b.WriteString(fmt.Sprintf("Frames captured: %d\n", m.progress.current))
	}

	b.WriteString(helpStyle.Render("\nPress 'q' to quit"))

	return "\n" + b.String() + "\n"
}

func (m model) viewError() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("✗ Error"))
	b.WriteString("\n\n")

	b.WriteString(errorStyle.Render(m.finalMessage))
	b.WriteString("\n")

	b.WriteString(helpStyle.Render("\nPress 'q' to quit"))

	return "\n" + b.String() + "\n"
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
