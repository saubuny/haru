package main

import (
	"context"
	// "database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	// _ "github.com/mattn/go-sqlite3"
	"github.com/saubuny/haru/internal/database"

	_ "embed"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Exit   key.Binding
	Select key.Binding
}

// ShortHelp implements the KeyMap interface.
func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Up, km.Down}
}

// FullHelp implements the KeyMap interface.
func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.Up, km.Down},
		{km.Exit, km.Select},
	}
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	Exit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "exit"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
}

type appConfig struct {
	DB  *database.Queries
	Ctx context.Context
}

type model struct {
	table table.Model
	help  help.Model
}

type topAnimeResponse struct {
	Data []struct {
		MalID  int    `json:"mal_id"`
		URL    string `json:"url"`
		Images struct {
			Jpg struct {
				ImageURL      string `json:"image_url"`
				SmallImageURL string `json:"small_image_url"`
				LargeImageURL string `json:"large_image_url"`
			} `json:"jpg"`
			Webp struct {
				ImageURL      string `json:"image_url"`
				SmallImageURL string `json:"small_image_url"`
				LargeImageURL string `json:"large_image_url"`
			} `json:"webp"`
		} `json:"images"`
		Trailer struct {
			YoutubeID string `json:"youtube_id"`
			URL       string `json:"url"`
			EmbedURL  string `json:"embed_url"`
		} `json:"trailer"`
		Approved bool `json:"approved"`
		Titles   []struct {
			Type  string `json:"type"`
			Title string `json:"title"`
		} `json:"titles"`
		Title         string   `json:"title"`
		TitleEnglish  string   `json:"title_english"`
		TitleJapanese string   `json:"title_japanese"`
		TitleSynonyms []string `json:"title_synonyms"`
		Type          string   `json:"type"`
		Source        string   `json:"source"`
		Episodes      int      `json:"episodes"`
		Status        string   `json:"status"`
		Airing        bool     `json:"airing"`
		Aired         struct {
			From string `json:"from"`
			To   string `json:"to"`
			Prop struct {
				From struct {
					Day   int `json:"day"`
					Month int `json:"month"`
					Year  int `json:"year"`
				} `json:"from"`
				To struct {
					Day   int `json:"day"`
					Month int `json:"month"`
					Year  int `json:"year"`
				} `json:"to"`
				String string `json:"string"`
			} `json:"prop"`
		} `json:"aired"`
		Duration   string  `json:"duration"`
		Rating     string  `json:"rating"`
		Score      float64 `json:"score"`
		ScoredBy   int     `json:"scored_by"`
		Rank       int     `json:"rank"`
		Popularity int     `json:"popularity"`
		Members    int     `json:"members"`
		Favorites  int     `json:"favorites"`
		Synopsis   string  `json:"synopsis"`
		Background string  `json:"background"`
		Season     string  `json:"season"`
		Year       int     `json:"year"`
		Broadcast  struct {
			Day      string `json:"day"`
			Time     string `json:"time"`
			Timezone string `json:"timezone"`
			String   string `json:"string"`
		} `json:"broadcast"`
		Producers []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"producers"`
		Licensors []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"licensors"`
		Studios []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"studios"`
		Genres []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"genres"`
		ExplicitGenres []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"explicit_genres"`
		Themes []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"themes"`
		Demographics []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"demographics"`
	} `json:"data"`
	Pagination struct {
		LastVisiblePage int  `json:"last_visible_page"`
		HasNextPage     bool `json:"has_next_page"`
		Items           struct {
			Count   int `json:"count"`
			Total   int `json:"total"`
			PerPage int `json:"per_page"`
		} `json:"items"`
	} `json:"pagination"`
}

func initialModel() model {
	// show a table of top anime initially, the model will change to whatever type of content the user is searching for
	c := &http.Client{Timeout: 4 * time.Second}
	res, err := c.Get("https://api.jikan.moe/v4/top/anime")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	var topAnime topAnimeResponse
	err = json.NewDecoder(res.Body).Decode(&topAnime)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	columns := []table.Column{
		{Title: "Name", Width: 40},
		{Title: "Rating", Width: 40},
	}

	rows := make([]table.Row, 0)
	for _, anime := range topAnime.Data {
		rows = append(rows, table.Row{anime.Title, anime.Rating})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
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

	return model{
		table: t,
		help:  help,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Exit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Select):
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n " + m.help.View(DefaultKeyMap) + "\n"
}

//go:embed sql/schema/schema.sql
var migrations string

func main() {
	// TODO: Make the database persist
	// db, err := sql.Open("sqlite3", ":memory:")
	// if err != nil {
	// 	log.Fatalf("Error initializing database: %v", err)
	// }

	// dbQueries := database.New(db)
	// cfg := appConfig{DB: dbQueries, Ctx: context.Background()}

	// Should only run once when database is created i think, maybe have a "haru init" command or something to set up the database once it's set up to persist
	// if _, err := db.ExecContext(cfg.Ctx, migrations); err != nil {
	// 	log.Fatalf("Error running migrations: %v", err)
	// }

	// user, err := cfg.DB.CreateUser(cfg.Ctx, "saubuny")
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }
	//
	// userFromSelect, err := cfg.DB.GetUser(cfg.Ctx, 1)
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
