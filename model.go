package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/saubuny/haru/navstack"
)

// TODO: Put this in own package i think

// I guess we need a surface level model :0
type Model struct {
	Navstack *navstack.Model
}

func New(model tea.Model) Model {
	nav := navstack.New(model)
	return Model{
		Navstack: &nav,
	}
}

func (m Model) Init() tea.Cmd {
	return m.Navstack.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := m.Navstack.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.Navstack.View()
}
