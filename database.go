package main

import (
	"context"
	"database/sql"
	"github.com/saubuny/haru/internal/database"
)

func initDB(schema string) (dbConfig, error) {
	// TODO: Make the database persist
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return dbConfig{}, err
	}

	dbQueries := database.New(db)
	cfg := dbConfig{DB: dbQueries, Ctx: context.Background()}

	if _, err := db.ExecContext(cfg.Ctx, schema); err != nil {
		return dbConfig{}, err
	}

	return cfg, nil
}

// func (cfg)
