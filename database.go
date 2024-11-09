package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"log"
	"strconv"
	"time"

	"github.com/saubuny/haru/internal/database"
)

func initDB(schema string) (dbConfig, error) {
	// db, err := sql.Open("sqlite3", ":memory:")
	db, err := sql.Open("sqlite3", "anime.db")
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

// put data inside haru struct, and upload to database
// TODO: deal with duplicates. if an ID already exists, use an update query instead.
func (cfg dbConfig) uploadToDB(title string, id int, startDate string, completion string) {
	cfg.DB.CreateAnime(cfg.Ctx, database.CreateAnimeParams{
		ID:          int64(id),
		Title:       title,
		Startdate:   startDate,
		Updateddate: time.Now().Format("0000-00-00"), // sqlc >:(
		Completion:  completion,
	})
}

// Kitsu also exports in the MAL format
func (cfg dbConfig) importMAL(malXml []byte) error {
	var animeList Myanimelist
	if err := xml.Unmarshal(malXml, &animeList); err != nil {
		return err
	}

	log.Printf("Importing Anime...")
	for _, anime := range animeList.Anime {
		id, _ := strconv.Atoi(anime.SeriesAnimedbID)

		// Different platforms use different naming
		completion := anime.MyStatus
		if anime.MyStatus == "Plan to Watch" {
			completion = PlanToWatch
		} else if anime.MyStatus == "On-Hold" {
			completion = OnHold
		}

		cfg.uploadToDB(anime.SeriesTitle, id, anime.MyStartDate, completion)
	}

	return nil
}

func (cfg dbConfig) importHianime(hiXml []byte) error {
	return nil
}
