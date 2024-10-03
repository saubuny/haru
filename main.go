package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/saubuny/haru/internal/database"

	_ "embed"

	tea "github.com/charmbracelet/bubbletea"
)

type appConfig struct {
	DB  *database.Queries
	Ctx context.Context
}

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	return model{
		choices:  []string{"Buy Carrots", "Buy Celery", "Buy whatever that is"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

//go:embed sql/schema/schema.sql
var migrations string

func main() {
	// TODO: Make the database persist
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	dbQueries := database.New(db)
	cfg := appConfig{DB: dbQueries, Ctx: context.Background()}

	// Should only run once when database is created i think, maybe have a "haru init" command or something to set up the database once it's set up to persist
	if _, err := db.ExecContext(cfg.Ctx, migrations); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	// user, err := cfg.DB.CreateUser(cfg.Ctx, "saubuny")
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }
	//
	// userFromSelect, err := cfg.DB.GetUser(cfg.Ctx, 1)
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
