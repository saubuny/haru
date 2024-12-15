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
	"github.com/kevm/bubbleo/menu"
	"github.com/kevm/bubbleo/navstack"
	"github.com/kevm/bubbleo/shell"
	// "github.com/kevm/bubbleo/utils"
	"github.com/saubuny/haru/animeinfo"
	"github.com/saubuny/haru/color"
	"github.com/saubuny/haru/internal/database"
)

const (
	Completion = "Completion"
	StartDate  = "Start Date"
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
	Width        int
	Height       int
	AnimeTitle   string // Needed for
	DefaultWidth int
	DBConfig     dbConfig
	PreviousRows AnimeListResponse

	// Different bubbles !!
	Shell       shell.Model
	ModifyMenu  menu.Model
	SearchInput textinput.Model
	Table       table.Model
	Help        help.Model
	// Viewport    viewport.Model

	AnimeInfo animeinfo.Model

	// These are for toggling visbility and controls of different bubbles. Might not be the best way to do things...
	// TODO: I bet you could use a navstack for all of these. maybe try that out :D
	ShowModifyMenu   bool
	FocusSearchInput bool
	ShowDBInfo       bool
	ShowHelp         bool
	ShowAnimeInfo    bool
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

	// Perhaps change color based on completion?
	choices := []menu.Choice{
		{
			Title:       Completion,
			Description: "Change completion status of anime",
			Model:       color.Model{RGB: "#00FF00", Sample: Completion},
		},
		{
			Title:       StartDate,
			Description: "Change start date of anime",
			Model:       color.Model{RGB: "#00FF00", Sample: StartDate},
		},
	}

	return Model{
		Table:       tb,
		Help:        help,
		SearchInput: ti,
		ShowHelp:    true,
		DBConfig:    cfg,
		ModifyMenu:  menu.New("Modify Entry", choices, nil),
		Shell:       shell.New(),
		AnimeInfo: animeinfo.Model{
			Viewport: viewport.New(0, 0),
			ShowHelp: true,
		},
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

// func getAnimeByIdCmd(id string) tea.Cmd {
// 	return func() tea.Msg {
// 		anime, err := getAnimeById(id)
// 		if err != nil {
// 			return ErrorMessage(err.Error())
// 		}
// 		return AnimeDataMessage(anime)
// 	}
// }

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

func (m Model) modifyStartDate(date string) tea.Cmd {
	// Make sure date is valid
	log.Printf(date)
	return func() tea.Msg {
		return nil
	}
}

func (m Model) modifyCompletion(completion string) tea.Cmd {
	// Depending on input, either remove or add to DB
	log.Printf(completion)
	return func() tea.Msg {
		return nil
	}
}

type AnimeListMessage AnimeListResponse
type AnimeDBListMessage []database.Anime
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
		return m, nil
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
		return m, nil
	// case AnimeDataMessage:
	// 	// Make own bubble
	// 	m.ShowAnimeInfo = true
	// 	m.AnimeTitle = msg.Data.Title
	//
	// 	content := msg.Data.Synopsis
	// 	m.Viewport = viewport.New(m.DefaultWidth, m.Table.Height()+4)
	// 	m.Viewport.SetContent(content)
	// 	return m, nil
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.DefaultWidth = int(float64(m.Width) * 0.8)

		// Height of help + search bar
		// Yucky hardcoded values :(
		m.Table.SetHeight(m.Height - 11)
		m.SearchInput.Width = m.DefaultWidth / 3
		// m.Viewport.Width = m.DefaultWidth

		// This thing takes in the whole message for some reason, making it difficult to change the height
		// m.ModifyMenu.SetSize(tea.WindowSizeMsg{
		// 	Height: msg.Height - 11,
		// 	Width:  msg.Width,
		// })
		return m, nil
	// case color.ColorSelected:
	// 	return m, tea.Quit
	// case selector.SelectedMsg:
	// 	switch msg {
	// 	case Completion:
	// 		// TODO: Change to using KevM/bubbleo
	// 		// Opens up another selector, which will send out its own selected message, which i can check here
	// 		// return m, modifyCompletion()
	// 	case StartDate:
	// 		// Make typing active, look at the search bar for how i did it earlier (although that'd be months ago atp) :3
	// 		// blah
	// 	case Watching, PlanToWatch, Completed, OnHold, Dropped:
	// 		return m, m.modifyCompletion(string(msg))
	// 	}
	case tea.KeyMsg:
		// This causes a tiny bit of code duplication, but the separation is worth it
		// if m.ShowModifyMenu {
		// 	// Blah
		// } else if m.ShowAnimeInfo {
		// 	m.Table.Blur()
		// 	switch {
		// 	case key.Matches(msg, AnimeInfoKeyMap.Esc):
		// 		// if m.ShowModifyMenu {
		// 		// 	m.ShowModifyMenu = false
		// 		// 	return m, nil
		// 		// }
		// 		m.Table.Focus()
		// 		m.ShowAnimeInfo = false
		// 		return m, nil
		// 	case key.Matches(msg, AnimeInfoKeyMap.Help):
		// 		m.ShowHelp = !m.ShowHelp
		// 		return m, nil
		// 	case key.Matches(msg, AnimeInfoKeyMap.Exit):
		// 		return m, tea.Quit
		// 	case key.Matches(msg, AnimeInfoKeyMap.Select):
		// 		m.ShowModifyMenu = true
		// 		return m, nil
		// 	}
		// } else {
		switch {
		case key.Matches(msg, DefaultKeyMap.Tab):
			if m.FocusSearchInput {
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
			m.FocusSearchInput = !m.FocusSearchInput
			if m.FocusSearchInput {
				m.Table.Blur()
				m.SearchInput.Focus()
			} else {
				m.Table.Focus()
				m.SearchInput.Blur()
			}
		case key.Matches(msg, DefaultKeyMap.Exit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Select):
			if m.FocusSearchInput {
				// Search for anime with new cmd
				val := m.SearchInput.Value()
				m.SearchInput.Reset()
				m.Table.Focus()
				m.SearchInput.Blur()
				m.FocusSearchInput = !m.FocusSearchInput
				if m.ShowDBInfo {
					return m, m.DBConfig.searchDBByNameCmd(val)
				}
				return m, searchAnimeByNameCmd(val)
			}
			// Replace with navstack ?
			// return m, getAnimeByIdCmd(m.Table.SelectedRow()[0])
			return m, m.Shell.Navstack.Push(navstack.NavigationItem{
				Title: "Anime Info",
				Model: m.AnimeInfo,
			})
		}
		// }
	}
	m.Table, cmd = m.Table.Update(msg)
	m.SearchInput, cmd = m.SearchInput.Update(msg)
	updatedanimeinfo, cmd := m.AnimeInfo.Update(msg)
	m.AnimeInfo = updatedanimeinfo.(animeinfo.Model)
	// updatedmenu, cmd := m.ModifyMenu.Update(msg)
	// m.ModifyMenu = updatedmenu.(menu.Model)
	return m, cmd
}

func (m Model) View() string {
	if m.Width == 0 {
		return ""
	}

	render := ""
	// if m.ShowAnimeInfo {
	// 	render += m.headerView(m.AnimeTitle) + "\n"
	// 	if m.ShowModifyMenu {
	// 		render += m.ModifyMenu.View() + "\n"
	// 	} else {
	// 		render += m.Viewport.View() + "\n"
	// 	}
	// } else {
	// }
	render += baseStyle.Render(m.SearchInput.View()) + "\n" + baseStyle.Render(m.Table.View()) + "\n"
	if m.ShowHelp {
		// if m.ShowModifyMenu {
		// 	render += m.Help.View(ModifyInfoKeyMap) + "\n"
		// } else if m.ShowAnimeInfo {
		// 	render += m.Help.View(AnimeInfoKeyMap) + "\n"
		// } else {
		render += m.Help.View(DefaultKeyMap) + "\n"
		// }
	}
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, 0, render)
}
