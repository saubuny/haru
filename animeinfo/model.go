package animeinfo

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevm/bubbleo/navstack"
	"github.com/kevm/bubbleo/utils"
)

type Model struct {
	Viewport viewport.Model
	ShowHelp bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, AnimeInfoKeyMap.Esc):
			// Push back to table (which means it has to be in its own package i think :0
			pop := utils.Cmdize(navstack.PopNavigation{})
			return m, pop
		case key.Matches(msg, AnimeInfoKeyMap.Help):
			m.ShowHelp = !m.ShowHelp
			return m, nil
		case key.Matches(msg, AnimeInfoKeyMap.Exit):
			return m, tea.Quit
		case key.Matches(msg, AnimeInfoKeyMap.ModifyEntry):

			// Push Modify Menu to stack
			return m, nil
		}
	}

	// m.Viewport, cmd = m.Viewport.Update(msg)
	return m, nil
}

func (m Model) View() string {
	return m.Viewport.View() + "\n"
}
