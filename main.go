package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	newFileInput textinput.Model
	createFileInputVisible bool
}

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
)

// Init initializes the model and returns an initial command.
// In this case, we don't have any initial commands to run, so we return nil.
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		
		case "ctrl+n":
			m.createFileInputVisible = true
			return m, nil
		}
	}
	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	} 

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

func (m model) View() string {
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 2)

	welcome := style.Render("Welcome to Hackpad!")
	view := ""
	if m.createFileInputVisible {
		view = m.newFileInput.View()
	}
	help := "Ctrl+N: New Note | Ctrl+L: List | Esc: back/save | Ctrl+S: save | Ctrl+Q: quit"

	return fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help)
}

func initialModel() model {
	// initialize new file input
	ti := textinput.New()
	ti.Placeholder = "Enter file name"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Cursor.Style = cursorStyle
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)

	return model{
		newFileInput: ti,
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
