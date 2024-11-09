package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 4

const (
	todo status = iota
	inProgress
	done
)

/*STYLING*/
var (
	ColumnStyle = lipgloss.NewStyle().
			Padding(1, 2)
	FocusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

/*CUSTOM ITEM*/
type Task struct {
	status      status
	title       string
	description string
}

// implement the List.Item interface

func (t Task) FilterValue() string {
	return t.title
}

// Getters
func (t Task) Title() string       { return t.title }
func (t Task) Description() string { return t.description }

/*MAIN MODEL*/

type Model struct {
	focused  status
	lists    []list.Model
	err      error
	loaded   bool
	quitting bool
}

func New() *Model {
	return &Model{}
}

// TODO go to next list
func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	}else {
		m.focused++
	}
}
// TODO go to previous list
func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	}else {
		m.focused--
	}
}

// TODO call this in tea.WindowSizeMsg
func (m *Model) InitLists(width, height int) {
	defaultLists := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height/2)
	defaultLists.SetShowHelp(false)
	m.lists = []list.Model{defaultLists, defaultLists, defaultLists}
	// init To Do
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sweets", description: "chocolate"},
		Task{status: todo, title: "read book", description: "50 laws of power"},
	})
	// inProgress
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: inProgress, title: "stay cool", description: "as a cucumber"},
		Task{status: inProgress, title: "stay cool", description: "as a cucumber"},
	})
	// done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: done, title: "Make columns", description: "its go"},
		Task{status: done, title: "Make columns", description: "its go"},
	})

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		if !m.loaded {
			ColumnStyle.Width(msg.Width/divisor)
			FocusedStyle.Width(msg.Width/divisor)
			ColumnStyle.Height(msg.Height-divisor)
			FocusedStyle.Height(msg.Height-divisor)
			m.InitLists(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Prev()
		case "right", "l":
			m.Next()
		}

	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)

	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		todoView := m.lists[todo].View()
		inProgressview := m.lists[inProgress].View()
		doneView := m.lists[done].View()
		switch m.focused {
		case inProgress: 
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			ColumnStyle.Render(todoView),
			FocusedStyle.Render(inProgressview),
			ColumnStyle.Render(doneView),
		)	
		case done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				ColumnStyle.Render(todoView),
				ColumnStyle.Render(inProgressview),
				FocusedStyle.Render(doneView),
			)	
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				FocusedStyle.Render(todoView),
				ColumnStyle.Render(inProgressview),
				ColumnStyle.Render(doneView),
			)
		}

	} else {
		return "Loading..."
	}
}

func main() {
	m := New()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatal("Boo")
	}

}
