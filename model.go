package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*MAIN MODEL*/

type Model struct {
	focused  status
	lists    []list.Model
	loaded   bool
	quitting bool
}

func New() *Model {
	return &Model{}
}

func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m Model) MoveToNext() tea.Msg {
	if len(m.lists[m.focused].Items()) > 0 {
		selectedItem := m.lists[m.focused].SelectedItem()
		selectedTask := selectedItem.(Task)
		m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
		selectedTask.Next()
		m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	}
	return nil
}

func (m Model) RemoveItem() tea.Msg {
	tasks := m.lists[m.focused].Items()
	if len(tasks) > 0 {
		selectedItem := m.lists[m.focused].SelectedItem()
		selectedTask := selectedItem.(Task)
		m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	}
	return nil
}

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
			return m, m.MoveToNext
		case "backspace":
			return m, m.RemoveItem
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
