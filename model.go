package main

import (
	"net/http"
	"time"

	"encoding/json"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

type Model struct {
	TextInput textinput.Model
	Table     table.Model
	Help      help.Model
	Width     int
	Height    int
	Typing    bool
	ShowHelp  bool
}

func initialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Insert Peak Here..."
	ti.Blur()
	ti.CharLimit = 60

	tb := table.New(
		table.WithColumns([]table.Column{}),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	tbStyle := table.DefaultStyles()
	tbStyle.Header = tbStyle.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	tbStyle.Selected = tbStyle.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	tb.SetStyles(tbStyle)

	help := help.New()
	help.ShowAll = true

	return Model{
		Table:     tb,
		Help:      help,
		TextInput: ti,
		ShowHelp:  true,
	}
}

func searchAnimeByName(searchString string) tea.Cmd {
	return func() tea.Msg {
		c := &http.Client{Timeout: 4 * time.Second}
		res, err := c.Get("https://api.jikan.moe/v4/anime?q=" + searchString)
		if err != nil {
			// return an error here as a message?
			log.Fatalf("Error: %v", err)
		}

		var topAnime TopAnimeResponse
		err = json.NewDecoder(res.Body).Decode(&topAnime)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		return jsonMessage(topAnime)
	}
}

func getTopAnime() tea.Msg {
	c := &http.Client{Timeout: 4 * time.Second}
	res, err := c.Get("https://api.jikan.moe/v4/top/anime")
	if err != nil {
		// return an error here as a message?
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

func (m Model) Init() tea.Cmd {
	return tea.Batch(getTopAnime)
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
		m.TextInput.Width = m.Table.Width()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Help):
			m.ShowHelp = !m.ShowHelp
		case key.Matches(msg, DefaultKeyMap.Esc):
			m.Typing = !m.Typing
			if m.Typing {
				m.Table.Blur()
				m.TextInput.Focus()
			} else {
				m.Table.Focus()
				m.TextInput.Blur()
			}
		case key.Matches(msg, DefaultKeyMap.Exit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Select):
			if m.Typing {
				// Search for anime with new cmd
				val := m.TextInput.Value()
				tea.Batch()
				m.TextInput.Reset()
				m.Table.Focus()
				m.TextInput.Blur()
				m.Typing = !m.Typing
				return m, searchAnimeByName(val)
			}
			return m, tea.Batch(
				// This doesn't do anything right now for some reason
				tea.Printf("Let's go to %s!", m.Table.SelectedRow()[1]),
			)
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Width == 0 {
		return ""
	}

	render := m.TextInput.View() + "\n" + baseStyle.Render(m.Table.View()) + "\n"
	if m.ShowHelp {
		render += m.Help.View(DefaultKeyMap) + "\n"
	}
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, 0, render)
}
