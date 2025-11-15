package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var tasks []string
var lastDeleted string

var textTheme = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFC8DD"))

var quit = lipgloss.NewStyle().Foreground(lipgloss.Color("#464646"))

type model struct {
	textInput textinput.Model
	cursor    int              // which to-do list item our cursor is pointing at
	selected  map[int]struct{} // which to-do items are selected
}

func main() {

	createFiles()
	tasks = read()

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {

	textInput := textinput.New()
	textInput.Placeholder = ":ribbit!"
	textInput.CharLimit = 255
	textInput.Width = 50

	return model{

		//add text input for commands
		textInput: textInput,

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) View() string {
	// The header
	s := "\nWhat's on the menu today?\n\n"

	// Iterate over our choices
	for i, choice := range tasks {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
			choice = lipgloss.NewStyle().Strikethrough(true).Render(choice)
		}

		// Render the row
		s += fmt.Sprintf("%s "+textTheme.Render("[ ]")+" %s\n", cursor, choice)
	}

	s += "\n" + m.textInput.View()

	// The footer
	s += quit.Render("\n\nPress q to quit.\n")

	// Send the UI for rendering
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		//access the command input when the : key is pressed
		case ":":
			m.textInput.Focus()

		// These keys should exit the program.
		case "ctrl+c", "q":
			if !m.textInput.Focused() {
				return m, tea.Quit
			}
		// The "up" and "k" keys move the cursor up
		case "up":
			if !m.textInput.Focused() {
				if m.cursor > 0 {
					m.cursor--
				}
			}

		// The "down" key move the cursor down
		case "down":
			if !m.textInput.Focused() {
				if m.cursor < len(tasks)-1 {
					m.cursor++
				}
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter":
			if !m.textInput.Focused() {
				selectChoice(m)
			} else {
				command(m.textInput.Value())
				m.textInput.Reset()
				m.textInput.Blur()
			}
		case " ":
			if !m.textInput.Focused() {
				selectChoice(m)
			}
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

// read user commands
func command(input string) {
	if input == "" {
		return
	}
	input = input[1:]

	//split the command into arguments
	arguments := strings.Split(input, " ")

	switch arguments[0] {
	case "add":
		//add a task except for the command argument
		addTask(strings.Join(arguments[1:], " "))
	case "undo":
		//add the last deleted task back
		if lastDeleted != "" {
			addTask(lastDeleted)
			lastDeleted = ""
		}

	}
}

// add a task to the list
func addTask(s string) {
	tasks = append(tasks, s)
	write(tasks)
}

// delete a selected choice
func selectChoice(m model) {
	if len(tasks) >= 1 {
		lastDeleted = tasks[m.cursor]
		tasks = append(tasks[:m.cursor], tasks[m.cursor+1:]...)
		write(tasks)
	}
}

func (m model) Init() tea.Cmd {
	//blink da text
	return textinput.Blink
}
