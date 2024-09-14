package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/saubuny/haru/internal/database"

	_ "embed"
)

type appConfig struct {
	DB  *database.Queries
	Ctx context.Context
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

	user, err := cfg.DB.CreateUser(cfg.Ctx, "saubuny")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("User: %v %v", user.ID, user.Name)

	userFromSelect, err := cfg.DB.GetUser(cfg.Ctx, 1)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("User: %v %v", userFromSelect.ID, userFromSelect.Name)
}
