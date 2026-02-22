package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	newFileInput           textinput.Model
	createFileInputVisible bool
	file                   *os.File
	noteTextArea           textarea.Model
	fileList               list.Model
	showingList            bool
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var (
	vaultDirectory string
	cursorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	docStyle       = lipgloss.NewStyle().Margin(1, 2)
)

// Init initializes the model and returns an initial command.
// In this case, we don't have any initial commands to run, so we return nil.
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.noteTextArea.SetWidth(msg.Width - h)
		m.noteTextArea.SetHeight(msg.Height - v - 5)
		m.fileList.SetWidth(msg.Width - h)
		m.fileList.SetHeight(msg.Height - v - 5)
		return m, nil

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+q":
			return m, tea.Quit

		case "esc":
			if m.showingList {
				m.showingList = false
				return m, nil
			}

			if m.file != nil {
				m.file = nil
				m.noteTextArea.SetValue("")
				return m, nil
			}

			if m.showingList {
				if m.fileList.FilterState() == list.Filtering {
					break
				}
			}

		case "ctrl+l":
			noteList := listFiles()
			m.fileList.SetItems(noteList)
			m.showingList = true
			return m, nil

		case "ctrl+n":
			m.createFileInputVisible = true
			return m, nil

		case "ctrl+s":
			if m.file == nil {
				break
			}

			if err := m.file.Truncate(0); err != nil {
				fmt.Println("Cannot save file :(")
				return m, nil
			}

			if _, err := m.file.Seek(0, 0); err != nil {
				fmt.Println("Cannot save file :(")
				return m, nil
			}
			if _, err := m.file.WriteString(m.noteTextArea.Value()); err != nil {
				fmt.Println("Cannot save file :(")
				return m, nil
			}
			if err := m.file.Close(); err != nil {
				fmt.Println("Cannot close the file :(")
			}
			m.file = nil
			m.noteTextArea.SetValue("")

			return m, nil

		case "enter":
			// if the file is already created, we don't need to create it again
			// we can just break out of the switch statement and let the user start writing in the text area
			// allow to writing in new line when the file is already created
			if m.file != nil {
				break
			}

			if m.showingList {
				// we are typecasting the selected item to our custom item type so that we can access the title and description of the file
				selectedItem, ok := m.fileList.SelectedItem().(item)
				if ok {
					filePath := fmt.Sprintf("%s/%s", vaultDirectory, selectedItem.title)
					content, err := os.ReadFile(filePath)
					if err != nil {
						log.Fatalf("%v", err)
						return m, nil
					}
					m.noteTextArea.SetValue(string(content))
					f, err := os.OpenFile(filePath, os.O_RDWR, 0644)
					if err != nil {
						log.Fatalf("%v", err)
						return m, nil
					}
					m.file = f
					m.showingList = false
				}
				return m, nil
			}
			// create a file with the name in the input
			m.createFileInputVisible = false
			fileName := m.newFileInput.Value()
			if fileName != "" {
				filePath := fmt.Sprintf("%s/%s.md", vaultDirectory, fileName)
				if _, err := os.Stat(filePath); err == nil {
					return m, nil
				}
				f, err := os.Create(filePath)
				if err != nil {
					log.Fatalf("%v", err)
				}
				m.file = f
				m.createFileInputVisible = false
				m.newFileInput.SetValue("")

			}
			return m, nil
		}
	}
	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}

	if m.file != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
	}

	if m.showingList {
		m.fileList, cmd = m.fileList.Update(msg)
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
	help := "Ctrl+N: New Note | Ctrl+L: List | Esc: back | Ctrl+S: save | Ctrl+Q: quit"
	view := ""
	if m.createFileInputVisible {
		view = m.newFileInput.View()
	}

	if m.file != nil {
		view = m.noteTextArea.View()
	}

	if m.showingList {
		view = m.fileList.View()
	}
	return fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help)
}

func initializeModel() model {
	err := os.MkdirAll(vaultDirectory, 0750)
	if err != nil {
		log.Fatal(err)
	}

	// initialize new file input
	ti := textinput.New()
	ti.Placeholder = "Enter file name"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Cursor.Style = cursorStyle
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)

	// initialize text are
	ta := textarea.New()
	ta.Placeholder = "Start writing your note..."
	ta.Focus()
	ta.CharLimit = 10000
	ta.Cursor.Style = cursorStyle

	// file list
	noteList := listFiles()
	finalList := list.New(noteList, list.NewDefaultDelegate(), 0, 0)
	finalList.Title = "Your Notes"

	return model{
		newFileInput:           ti,
		createFileInputVisible: false,
		noteTextArea:           ta,
		fileList:               finalList,
	}
}

// init is a special function in Go that is automatically called when the package is initialized.
// It's often used to set up initial state or perform any necessary setup before the main function runs.
// In this case, we don't have any specific initialization logic, so the function is empty.
func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	vaultDirectory = homeDir + "/.hackpad/vault"
}

func main() {
	p := tea.NewProgram(initializeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func listFiles() []list.Item {
	items := make([]list.Item, 0)
	entries, err := os.ReadDir(vaultDirectory)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range entries {
		if !file.IsDir() {
			if info, err := file.Info(); err == nil {
				items = append(items, item{
					title: file.Name(),
					desc:  fmt.Sprintf("Last edited: %s", info.ModTime().Format("Jan 2, 2006 at 3:04pm")),
				})
			}
		}
	}
	return items
}
