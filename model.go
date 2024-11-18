package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

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
	PreviousRows  AnimeListResponse
	ModifyEntry   bool
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

// Extracted into separate function because it needs to be used in the kitsu import function as well
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

func (cfg dbConfig) searchDBByNameCmd(searchString string) tea.Cmd {
	return func() tea.Msg {
		fullAnime, err := cfg.DB.GetAllAnime(cfg.Ctx)
		if err != nil {
			return ErrorMessage(err.Error())
		}

		// I know nothing about search algorithms, so i'm doing this simple thing for now. maybe i can do something better in the future, idk
		newAnime := []database.Anime{}
		for _, anime := range fullAnime {
			if strings.Contains(strings.ToLower(anime.Title), strings.ToLower(searchString)) {
				newAnime = append(newAnime, anime)
			}
		}

		return AnimeDBListMessage(newAnime)
	}
}

func searchAnimeByNameCmd(searchString string) tea.Cmd {
	// Breaks from spaces, have to replace them
	newSearchString := ""
	for _, r := range searchString {
		if unicode.IsSpace(r) {
			newSearchString += "%20"
		} else {
			newSearchString += string(r)
		}
	}

	return func() tea.Msg {
		c := &http.Client{Timeout: 4 * time.Second}
		res, err := c.Get("https://api.jikan.moe/v4/anime?q=" + newSearchString)
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

func (m Model) getTopAnime() tea.Msg {
	// Show cached results if they've already been saved
	if m.PreviousRows.Data != nil {
		return AnimeListMessage(m.PreviousRows)
	}

	c := &http.Client{Timeout: 4 * time.Second}
	res, err := c.Get("https://api.jikan.moe/v4/top/anime")
	if err != nil {
		return ErrorMessage(err.Error())
	}

	var topAnime AnimeListResponse
	err = json.NewDecoder(res.Body).Decode(&topAnime)
	if err != nil {
		return ErrorMessage(err.Error())
	}

	return AnimeListMessage(topAnime)
}

func (cfg dbConfig) modifyEntry() tea.Cmd {
	// Depending on input, either remove or add to DB
	// the menu should be able
	return func() tea.Msg {
		return nil
	}
}

type AnimeListMessage AnimeListResponse
type AnimeDBListMessage []database.Anime
type AnimeDataMessage AnimeDataResponse
type ErrorMessage string

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.getTopAnime)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case ErrorMessage:
		log.Fatalf("Error: %v", msg)
	case AnimeDBListMessage:
		columns := []table.Column{
			{Title: "Id", Width: 10},
			{Title: "Name", Width: 40},
			{Title: "Completion", Width: 20},
			{Title: "Start Date", Width: 10},
		}

		rows := make([]table.Row, 0)
		for _, anime := range msg {
			rows = append(rows, table.Row{strconv.Itoa(int(anime.ID)), anime.Title, anime.Completion, anime.Startdate})
		}

		m.Table.SetColumns(columns)
		m.Table.SetRows(rows)
		m.Table.SetCursor(0)
	case AnimeListMessage:
		// There is a crash if the number of columns differs between both tabs
		columns := []table.Column{
			{Title: "Id", Width: 10},
			{Title: "Name", Width: 40},
			{Title: "Rating", Width: 30},
			{Title: "Score", Width: 10},
		}

		rows := make([]table.Row, 0)
		for _, anime := range msg.Data {
			rows = append(rows, table.Row{strconv.Itoa(anime.MalID), anime.Title, anime.Rating, fmt.Sprintf("%v", anime.Score)})
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
		// I really do think hardcoded values like this are bad but i dont know an alternative
		m.Table.SetHeight(m.Height - 11)
		m.TextInput.Width = m.DefaultWidth / 3
		m.Viewport.Width = m.DefaultWidth
	case tea.KeyMsg:
		// This causes just a TINY bit of code duplication, but the separation is worth it
		if m.ShowAnimeInfo {
			m.Table.Blur()
			switch {
			case key.Matches(msg, AnimeInfoKeyMap.Esc):
				if m.ModifyEntry {
					m.ModifyEntry = false
					return m, nil
				}
				m.Table.Focus()
				m.ShowAnimeInfo = false
				return m, nil
			case key.Matches(msg, AnimeInfoKeyMap.Help):
				m.ShowHelp = !m.ShowHelp
				return m, nil
			case key.Matches(msg, AnimeInfoKeyMap.Exit):
				return m, tea.Quit
			case key.Matches(msg, AnimeInfoKeyMap.Select):
				if m.ModifyEntry {
					return m, nil
				}
				m.ModifyEntry = true
				// We need a popup, for changing all the list info. sort of like a menu. something like that probably exists, but i dont know how to layer things with lipgloss, or if that's even possible. if it isn't, we can just render this instead of the description. if we're editing the info, the esc button can go back to the description instead of the table !!

				// Nope. i have to make the menu all by myself. that's very difficult. uh oh.

				// We can use a list bubble, it will have two layers
				// Clicking on one will open either another list or a date picker for changing your options. pressing enter will save, exit to escape. there will also be a "Remove" option for completion to remove it from your list. if this option is selected, the other options will be ignored.

				// I can use a selector actually! seems much more simple and better for what im working on. the same library has a prompt! i could use that prompt instead of the one im using now perhaps, but first i have to check if my current prompt crashing like it is is an HTTP issue or a library issue. let me do that now. ITS AN HTTP ERROR. IT DOESNT HAPPEN IN THE DB SEARCH. WHAT THE FREAK. let me try something
			}
		}
		switch {
		case key.Matches(msg, DefaultKeyMap.Tab):
			if m.Typing {
				return m, nil
			}

			if m.ShowDBInfo {
				m.ShowDBInfo = false
				return m, m.getTopAnime
			}

			m.ShowDBInfo = true
			return m, m.DBConfig.showDBAnime()
		case key.Matches(msg, DefaultKeyMap.Help):
			m.ShowHelp = !m.ShowHelp
			return m, nil
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
				m.TextInput.Reset()
				m.Table.Focus()
				m.TextInput.Blur()
				m.Typing = !m.Typing
				if m.ShowDBInfo {
					return m, m.DBConfig.searchDBByNameCmd(val)
				}
				return m, searchAnimeByNameCmd(val)
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
		render += m.headerView(m.AnimeTitle) + "\n"
		if m.ModifyEntry {
			render += "\n"
		} else {
			render += m.Viewport.View() + "\n"
		}
	} else {
		render += baseStyle.Render(m.TextInput.View()) + "\n" + baseStyle.Render(m.Table.View()) + "\n"
	}
	if m.ShowHelp {
		if m.ShowAnimeInfo {
			render += m.Help.View(AnimeInfoKeyMap) + "\n"
		} else {
			render += m.Help.View(DefaultKeyMap) + "\n"
		}
	}
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, 0, render)
}
