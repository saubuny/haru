package animeinfo

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/saubuny/haru/navstack"
	"github.com/saubuny/haru/types"
)

var titleStyle = func() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
}()

func (m Model) headerView(name string) string {
	title := titleStyle.Render(name)
	line := strings.Repeat("─", max(0, int(float64(m.width)*0.8)-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

type Model struct {
	width    int
	height   int
	title    string
	viewport viewport.Model
	help     help.Model
	showHelp bool
}

func New() Model {
	help := help.New()
	help.ShowAll = true
	return Model{
		help:     help,
		viewport: viewport.New(0, 0),
		showHelp: true,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case types.AnimeDataMessage:
		m.title = msg.Data.Title
		content := msg.Data.Synopsis
		m.viewport.SetContent(content)
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = int(float64(m.width) * 0.8)
		m.viewport.Height = m.height - 6
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, AnimeInfoKeyMap.Esc):
			pop := navstack.Cmd(navstack.PopNavigation{})
			return m, pop
		case key.Matches(msg, AnimeInfoKeyMap.Help):
			m.showHelp = !m.showHelp
			return m, nil
		case key.Matches(msg, AnimeInfoKeyMap.ModifyEntry):
			m.viewport.SetContent("fuck")
			// Push Modify Menu to stack
			return m, nil
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	render := ""
	render += m.headerView(m.title) + "\n"
	render += m.viewport.View() + "\n"
	if m.showHelp {
		render += m.help.View(AnimeInfoKeyMap)
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, 0, render)
}
