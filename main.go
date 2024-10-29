package main

import (
	"context"
	_ "embed"
	"log"
	"os"

	"github.com/saubuny/haru/internal/database"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

type dbConfig struct {
	DB  *database.Queries
	Ctx context.Context
}

//go:embed sql/schema/schema.sql
var migrations string

// === TODO ===
// 1. CLI interaction
// 2. Display DB information in TUI table (have different tabs?)

func main() {
	_, err := initDB(migrations)
	if err != nil {
		log.Fatalf("Error initalizing DB: %v", err)
	}

	// Run TUI by default
	app := &cli.App{
		Name:  "Haru",
		Usage: "Track anime",
		Action: func(ctx *cli.Context) error {
			p := tea.NewProgram(initialModel(), tea.WithAltScreen())
			tea.SetWindowTitle("Haru")
			if _, err := p.Run(); err != nil {
				log.Fatalf("Error: %v", err)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	// cfg.importMAL()

	// user, err := cfg.DB.CreateUser(cfg.Ctx, "saubuny")
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }
	//
	// userFromSelect, err := cfg.DB.GetUser(cfg.Ctx, 1)
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }
}
