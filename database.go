package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/saubuny/haru/internal/database"
)

func initDB(schema string) dbConfig {
	// TODO: Make the database persist
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	dbQueries := database.New(db)
	cfg := dbConfig{DB: dbQueries, Ctx: context.Background()}

	if _, err := db.ExecContext(cfg.Ctx, schema); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	return cfg
}

// func (cfg)
