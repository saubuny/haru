package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/saubuny/haru/internal/database"

	_ "embed"

	tea "github.com/charmbracelet/bubbletea"
)

type dbConfig struct {
	DB  *database.Queries
	Ctx context.Context
}

//go:embed sql/schema/schema.sql
var migrations string

func main() {
	initDB(migrations)

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
	tea.SetWindowTitle("Haru")
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
