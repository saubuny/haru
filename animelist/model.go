package animelist

import (
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	"encoding/json"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/saubuny/haru/db"
	"github.com/saubuny/haru/types"

	"github.com/saubuny/haru/internal/database"
)

var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

type Model struct {
	width          int
	height         int
	cachedTopAnime types.AnimeListResponse

	// This allows us to turn off everything related to this model when we pop it on or off the stack
	focus bool

	dbConfig    db.DBConfig
	animeTable  table.Model
	searchInput textinput.Model
	help        help.Model
}

func InitialModel(db db.DBConfig) Model {
	ti := textinput.New()
	ti.Placeholder = "Insert Peak Here..."
	ti.Blur()
	ti.CharLimit = 60

	tb := table.New(
		table.WithColumns([]table.Column{}),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(7), // TODO: why is this hardcoded
	)

	// TODO: Mess with these styles
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
		animeTable:  tb,
		help:        help,
		searchInput: ti,
		dbConfig:    db,
	}
}

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
			return types.ErrorMsg(err.Error())
		}

		var anime types.AnimeListResponse
		err = json.NewDecoder(res.Body).Decode(&anime)
		if err != nil {
			return types.ErrorMsg(err.Error())
		}

		return AnimeListMessage(anime)
	}
}

// TODO: Show spinner on HTTP requests :3
func (m Model) getTopAnime() tea.Msg {
	// Show cached results if they've already been saved
	if m.cachedTopAnime.Data != nil {
		return AnimeListMessage(m.cachedTopAnime)
	}

	c := &http.Client{Timeout: 4 * time.Second}
	res, err := c.Get("https://api.jikan.moe/v4/top/anime")
	if err != nil {
		return types.ErrorMsg(err.Error())
	}

	var topAnime types.AnimeListResponse
	err = json.NewDecoder(res.Body).Decode(&topAnime)
	if err != nil {
		return types.ErrorMsg(err.Error())
	}

	return AnimeListMessage(topAnime)
}

func showDBAnime(db db.DBConfig) tea.Cmd {
	return func() tea.Msg {
		anime, err := db.DB.GetAllAnime(db.Ctx)
		if err != nil {
			return types.ErrorMsg(err.Error())
		}

		return AnimeDBListMessage(anime)
	}
}

func searchDBByNameCmd(db db.DBConfig, searchString string) tea.Cmd {
	return func() tea.Msg {
		fullAnime, err := db.DB.GetAllAnime(db.Ctx)
		if err != nil {
			return types.ErrorMsg(err.Error())
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

func (m Model) Init() tea.Cmd {
	return m.getTopAnime
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case types.ErrorMsg:
		log.Fatalf("Error: %v", msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// I have no idea what these hardcoded values are or what past me was thinking
		// Height of help + search bar
		m.animeTable.SetHeight(m.height - 11)
		m.searchInput.Width = int(float64(m.width)*0.8) / 3
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, AnimeListKeyMap.Exit):
			return m, tea.Quit
		case key.Matches(msg, AnimeListKeyMap.Help):
			// Toggle Help
			return m, nil
		case key.Matches(msg, AnimeListKeyMap.Esc):
			// Toggle Search
			return m, nil
		case key.Matches(msg, AnimeListKeyMap.Tab):
			// Toggle DB and MAL view
			return m, nil
		case key.Matches(msg, AnimeListKeyMap.Select):
			// Either search or open anime info (view navstack)
			return m, nil
		}
	}

	m.animeTable, cmd = m.animeTable.Update(msg)
	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}
	render := ""

	render += baseStyle.Render(m.searchInput.View()) + "\n"
	render += baseStyle.Render(m.animeTable.View()) + "\n"

	return lipgloss.Place(m.width, m.height, lipgloss.Center, 0, render)
}