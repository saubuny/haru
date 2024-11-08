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
// this will take in a completion number, the individual import functions will handle deciding what that is
func uploadToDB(title string, id int64, startDate time.Time, completion int) {}

// Kitsu also exports in the MAL format
func (cfg dbConfig) importMAL(malXml []byte) error {
	var animeList Myanimelist
	if err := xml.Unmarshal(malXml, &animeList); err != nil {
		return err
	}

	// TODO: convert completion to standard completion type like in old haru, and check if an id already exists in the db first, and if so, overwrite based on date, which should also be added in the db (have an statusUpdatedAt, and a CreatedAt)
	log.Printf("Importing Anime...")
	for _, anime := range animeList.Anime {
		// Kitsu requires us to do this, which is annoying, so we just wont support kitsu :)

		// Add a spinner or something here in the future
		// time.Sleep(time.Duration(time.Millisecond * 500)) // Rate limiting
		// malAnime, err := getAnimeById(anime.SeriesAnimedbID)
		// if err != nil {
		// 	return err
		// }
		// log.Printf("Importing " + malAnime.Data.Title)

		id, _ := strconv.Atoi(anime.SeriesAnimedbID)
		cfg.DB.CreateAnime(cfg.Ctx, database.CreateAnimeParams{
			ID:         int64(id),
			Title:      anime.SeriesTitle,
			Completion: anime.MyStatus,
		})
	}

	return nil
}

func (cfg dbConfig) importHianime(hiXml string) {}
