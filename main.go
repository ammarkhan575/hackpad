package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	msg string
}

// Init initializes the model and returns an initial command.
// In this case, we don't have any initial commands to run, so we return nil.
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	var style = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
	Padding(0, 2)

	welcome := style.Render("Welcome to Hackpad!")
	view := ""
	help := "Ctrl+N: New Note | Ctrl+L: List | Esc: back/save | Ctrl+S: save | Ctrl+Q: quit"

	return fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help)
}

func initialModel() model {
	return model{
		msg: "Hello, Welcome to Hackpad",
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
