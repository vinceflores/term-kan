package main

import (
	"log"

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

func main() {
	models = []tea.Model{New(), NewForm(todo)}
	m := models[model]
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatal("Boo")
	}
}
