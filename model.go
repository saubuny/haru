package main

import (
	"net/http"
	"time"

	"encoding/json"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

type Model struct {
	Table   table.Model
	Help    help.Model
	Width   int
	Height  int
	Loading bool
}

func initialModel() Model {
	t := table.New(
		table.WithColumns([]table.Column{}),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	help := help.New()
	help.ShowAll = true

	return Model{
		Table: t,
		Help:  help,
	}
}

func getTopAnime() tea.Msg {
	c := &http.Client{Timeout: 4 * time.Second}
	res, err := c.Get("https://api.jikan.moe/v4/top/anime")
	if err != nil {
		// return an error here as an interface?
		log.Fatalf("Error: %v", err)
	}

	var topAnime TopAnimeResponse
	err = json.NewDecoder(res.Body).Decode(&topAnime)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return jsonMessage(topAnime)
}

type jsonMessage TopAnimeResponse // TODO: find a generic json format that fits all jikan responses

// Start off with listing top anime
func (m Model) Init() tea.Cmd {
	return getTopAnime
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case jsonMessage:
		columns := []table.Column{
			{Title: "Name", Width: 40},
			{Title: "Rating", Width: 40},
		}

		rows := make([]table.Row, 0)
		for _, anime := range msg.Data {
			rows = append(rows, table.Row{anime.Title, anime.Rating})
		}

		m.Table.SetColumns(columns)
		m.Table.SetRows(rows)
	case tea.WindowSizeMsg:
		// TODO: set a minimum width/height
		m.Width = msg.Width
		m.Height = msg.Height
		m.Table.SetHeight(m.Height - 8)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Exit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Select):
			return m, tea.Batch(
				// This doesn't do anything right now for some reason
				tea.Printf("Let's go to %s!", m.Table.SelectedRow()[1]),
			)
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Width == 0 {
		return ""
	}

	// TODO: show spinner when loading http reqs

	render := baseStyle.Render(m.Table.View()) + "\n " + m.Help.View(DefaultKeyMap) + "\n"
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, 0, render)
}
