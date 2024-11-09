package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/saubuny/haru/internal/database"
)

var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

var titleStyle = func() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
}()

func (m Model) headerView(name string) string {
	title := titleStyle.Render(name)
	line := strings.Repeat("─", max(0, m.DefaultWidth-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

type Model struct {
	TextInput     textinput.Model
	Table         table.Model
	Help          help.Model
	Width         int
	Height        int
	Typing        bool
	ShowHelp      bool
	ShowAnimeInfo bool
	Viewport      viewport.Model
	AnimeTitle    string
	DefaultWidth  int
	DBConfig      dbConfig
	ShowDBInfo    bool
}

func initialModel(cfg dbConfig) Model {
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
		DBConfig:  cfg,
	}
}

// Extracted into separate function because it needs to be used in the import functions as well
func getAnimeById(id string) (AnimeDataResponse, error) {
	c := &http.Client{Timeout: 4 * time.Second}
	res, err := c.Get("https://api.jikan.moe/v4/anime/" + id)
	if err != nil {
		return AnimeDataResponse{}, err
	}

	var anime AnimeDataResponse
	err = json.NewDecoder(res.Body).Decode(&anime)
	if err != nil {
		return AnimeDataResponse{}, nil
	}

	return anime, nil
}

func getAnimeByIdCmd(id string) tea.Cmd {
	return func() tea.Msg {
		anime, err := getAnimeById(id)
		if err != nil {
			return ErrorMessage(err.Error())
		}
		return AnimeDataMessage(anime)
	}
}

// Breaks if there is a space in the search string ??????
func searchAnimeByNameCmd(searchString string) tea.Cmd {
	return func() tea.Msg {
		c := &http.Client{Timeout: 4 * time.Second}
		res, err := c.Get("https://api.jikan.moe/v4/anime?q=" + searchString)
		if err != nil {
			return ErrorMessage(err.Error())
		}

		var anime AnimeListResponse
		err = json.NewDecoder(res.Body).Decode(&anime)
		if err != nil {
			return ErrorMessage(err.Error())
		}

		return AnimeListMessage(anime)
	}
}

func (cfg dbConfig) showDBAnime() tea.Cmd {
	return func() tea.Msg {
		anime, err := cfg.DB.GetAllAnime(cfg.Ctx)
		if err != nil {
			return ErrorMessage(err.Error())
		}

		return AnimeDBListMessage(anime)
	}
}

func getTopAnime() tea.Msg {
	c := &http.Client{Timeout: 4 * time.Second}
	res, err := c.Get("https://api.jikan.moe/v4/top/anime")
	if err != nil {
		// return an error here as a message? log.Fatalf("Error: %v", err)
	}

	var topAnime AnimeListResponse
	err = json.NewDecoder(res.Body).Decode(&topAnime)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return AnimeListMessage(topAnime)
}

type AnimeListMessage AnimeListResponse
type AnimeDBListMessage []database.Anime
type AnimeDataMessage AnimeDataResponse
type ErrorMessage string

func (m Model) Init() tea.Cmd {
	return tea.Batch(getTopAnime)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case ErrorMessage:
		log.Fatalf("Error: %v", msg)
	case AnimeDBListMessage:
		if m.ShowDBInfo {
			m.ShowDBInfo = false
			return m, getTopAnime
		}

		m.ShowDBInfo = true

		columns := []table.Column{
			{Title: "Id", Width: 10},
			{Title: "Name", Width: 40},
			{Title: "Completion", Width: 30},
			{Title: "Start Date", Width: 30},
		}

		rows := make([]table.Row, 0)
		for _, anime := range msg {
			rows = append(rows, table.Row{strconv.Itoa(int(anime.ID)), anime.Title, anime.Completion, anime.Startdate})
		}

		m.Table.SetColumns(columns)
		m.Table.SetRows(rows)
		m.Table.SetCursor(0)
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
		m.Table.SetCursor(0)
	case AnimeDataMessage:
		m.ShowAnimeInfo = true
		m.AnimeTitle = msg.Data.Title

		content := msg.Data.Synopsis
		m.Viewport = viewport.New(m.DefaultWidth, m.Table.Height()+4)
		m.Viewport.SetContent(content)
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.DefaultWidth = int(float64(m.Width) * 0.8)

		// Height of help + search bar
		m.Table.SetHeight(m.Height - 11)
		m.TextInput.Width = m.DefaultWidth / 3
		m.Viewport.Width = m.DefaultWidth
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.ChangeTab):
			if m.ShowAnimeInfo {
				return m, nil
			}
			return m, m.DBConfig.showDBAnime()
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
				return m, searchAnimeByNameCmd(val)
			}
			if m.ShowAnimeInfo {
				return m, nil
			}
			return m, getAnimeByIdCmd(m.Table.SelectedRow()[0])
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	m.TextInput, cmd = m.TextInput.Update(msg)
	m.Viewport, cmd = m.Viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Width == 0 {
		return ""
	}

	render := ""
	if m.ShowAnimeInfo {
		render += m.headerView(m.AnimeTitle) + "\n" + m.Viewport.View() + "\n"
	} else {
		render += baseStyle.Render(m.TextInput.View()) + "\n" + baseStyle.Render(m.Table.View()) + "\n"
	}
	if m.ShowHelp {
		render += m.Help.View(DefaultKeyMap) + "\n"
	}
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, 0, render)
}
