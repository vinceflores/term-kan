package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
