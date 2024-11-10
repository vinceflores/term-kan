package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
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

/* Model Management */
var models []tea.Model

const (
	model status = iota
	form
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

func NewTask(status status, title, description string) Task {
	return Task{
		status:      status,
		title:       title,
		description: description,
	}
}

// implement the List.Item interface

func (t Task) FilterValue() string {
	return t.title
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
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
	} else {
		m.focused++
	}
}

// TODO go to previous list
func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m Model) MoveToNext() tea.Cmd {
	if len(m.lists[m.focused].Items()) > 0 {
		selectedItem := m.lists[m.focused].SelectedItem()
		selectedTask := selectedItem.(Task)
		m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
		selectedTask.Next()
		m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	}
	return nil
}

// TODO call this in tea.WindowSizeMsg
func (m *Model) InitLists(width, height int) {
	defaultLists := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height/2)
	defaultLists.SetShowHelp(false)
	m.lists = []list.Model{defaultLists, defaultLists, defaultLists}
	// init To Do
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "Add Task here", description: "Anything At all"},
	})
	// inProgress
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{})
	// done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{})

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		if !m.loaded {
			ColumnStyle.Width(msg.Width / divisor)
			FocusedStyle.Width(msg.Width / divisor)
			ColumnStyle.Height(msg.Height - divisor)
			FocusedStyle.Height(msg.Height - divisor)
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
		case "enter":
			return m, m.MoveToNext()
		case "n":
			models[model] = m
			models[form] = NewForm(m.focused)
			return models[form].Update(nil)
		}
	case Task:
		task := msg
		return m, m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)
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

/* Form Model */
type Form struct {
	focused     status
	title       textinput.Model
	description textarea.Model
}

func NewForm(focused status) *Form {
	form := &Form{}
	labelStyle := lipgloss.NewStyle().PaddingRight(2)

	input := textinput.New()
	input.Prompt = "Title:"
	input.Placeholder = "Add a Title"
	input.PromptStyle = labelStyle

	description := textarea.New()
	// description.Prompt = "Description: "
	description.Placeholder = "Add a description"

	form.title = input
	form.description = description

	form.focused = focused
	form.title.Focus()
	return form
}

func (m Form) CreateTask() tea.Msg {
	// TODO create a new Task
	task := NewTask(m.focused, m.title.Value(), m.description.Value())
	return task
}

func (m Form) Init() tea.Cmd { return nil }
func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.title.Focused() {
				m.title.Blur()
				m.description.Focus()
				return m, textarea.Blink
			} else {
				models[form] = m
				return models[model], m.CreateTask
			}
		}
	}
	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
		return m, cmd
	} else {
		m.description, cmd = m.description.Update(msg)
		return m, cmd
	}
}

func (m Form) View() string {

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"Create title",
		m.title.View(),
		m.description.View(),
	)
}

func main() {
	models = []tea.Model{New(), NewForm(todo)}
	m := models[model]
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatal("Boo")
	}

}
