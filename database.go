package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"strconv"

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

// Kitsu also exports in the MAL format
func (cfg dbConfig) importMAL(malXml []byte) error {
	var animeList Myanimelist
	if err := xml.Unmarshal(malXml, &animeList); err != nil {
		return err
	}

	// TODO: convert completion to standard completion type like in old haru, and check if an id already exists in the db first, and if so, overwrite based on date, which should also be added in the db (have an statusUpdatedAt, and a CreatedAt)
	for _, anime := range animeList.Anime {
		// The title isn't included, so we have to fetch the title manually for each anime (very slow, maybe add a progress bar somehow :3)

		// ADD THAT HERE !!!

		id, _ := strconv.Atoi(anime.SeriesAnimedbID)
		cfg.DB.CreateAnime(cfg.Ctx, database.CreateAnimeParams{
			ID:         int64(id),
			Title:      title,
			Completion: anime.MyStatus,
		})
	}

	return nil
}
func (cfg dbConfig) importHianime(hiXml string) {}
