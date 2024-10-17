package main

import (
	"net/http"
	"strconv"
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
	TextInput     textinput.Model
	Table         table.Model
	Help          help.Model
	Width         int
	Height        int
	Typing        bool
	ShowHelp      bool
	ShowAnimeInfo bool
	AnimeInfo     string
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

func getAnimeById(id string) tea.Cmd {
	return func() tea.Msg {
		c := &http.Client{Timeout: 4 * time.Second}
		res, err := c.Get("https://api.jikan.moe/v4/anime/" + id)
		if err != nil {
			// return an error here as a message?
			log.Fatalf("Error: %v", err)
		}

		var anime AnimeDataResponse
		err = json.NewDecoder(res.Body).Decode(&anime)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		return AnimeDataMessage(anime)
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

		var anime AnimeListResponse
		err = json.NewDecoder(res.Body).Decode(&anime)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		return AnimeListMessage(anime)
	}
}

func getTopAnime() tea.Msg {
	c := &http.Client{Timeout: 4 * time.Second}
	res, err := c.Get("https://api.jikan.moe/v4/top/anime")
	if err != nil {
		// return an error here as a message?
		log.Fatalf("Error: %v", err)
	}

	var topAnime AnimeListResponse
	err = json.NewDecoder(res.Body).Decode(&topAnime)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return AnimeListMessage(topAnime)
}

type AnimeListMessage AnimeListResponse
type AnimeDataMessage AnimeDataResponse

func (m Model) Init() tea.Cmd {
	return tea.Batch(getTopAnime)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case AnimeListMessage:
		columns := []table.Column{
			{Title: "Id", Width: 10},
			{Title: "Name", Width: 40},
			{Title: "Rating", Width: 40},
		}

		rows := make([]table.Row, 0)
		for _, anime := range msg.Data {
			rows = append(rows, table.Row{strconv.Itoa(anime.MalID), anime.Title, anime.Rating})
		}

		m.Table.SetColumns(columns)
		m.Table.SetRows(rows)
	case AnimeDataMessage:
		m.ShowAnimeInfo = true
		content := msg.Data.Title + "\n" + msg.Data.Synopsis
		m.AnimeInfo = baseStyle.MaxWidth(m.Width-10).Render(content) + "\n"
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
			if m.ShowAnimeInfo {
				m.ShowAnimeInfo = false
				return m, nil
			}
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
			if m.ShowAnimeInfo {
				return m, nil
			}
			return m, getAnimeById(m.Table.SelectedRow()[0])
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

	render := ""
	if m.ShowAnimeInfo {
		render += m.AnimeInfo
	} else {
		render += m.TextInput.View() + "\n" + baseStyle.Render(m.Table.View()) + "\n"
	}
	if m.ShowHelp {
		render += m.Help.View(DefaultKeyMap) + "\n"
	}
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, 0, render)
}
