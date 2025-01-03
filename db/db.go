package db

import (
	"context"
	"database/sql"
	"encoding/xml"
	"log"
	"strconv"
	"time"

	"github.com/saubuny/haru/internal/database"
	"github.com/saubuny/haru/types"
)

type DBConfig struct {
	DB  *database.Queries
	Ctx context.Context
}

// TODO: Func for udpating data

// TODO: input custom db location (like ~/.haru/anime.db)
// TODO: have separate table for manga !!
func InitDB(schema string, location string) (DBConfig, error) {
	// db, err := sql.Open("sqlite3", ":memory:")
	db, err := sql.Open("sqlite3", location)
	if err != nil {
		return DBConfig{}, err
	}

	dbQueries := database.New(db)
	cfg := DBConfig{DB: dbQueries, Ctx: context.Background()}

	if _, err := db.ExecContext(cfg.Ctx, schema); err != nil {
		return DBConfig{}, err
	}

	return cfg, nil
}

func (cfg DBConfig) UploadToDB(title string, id int, startDate string, completion string) error {
	// Check if ID already exists, create new anime if it does
	_, err := cfg.DB.GetAnime(cfg.Ctx, int64(id))
	if err == sql.ErrNoRows {
		_, err = cfg.DB.CreateAnime(cfg.Ctx, database.CreateAnimeParams{
			ID:          int64(id),
			Title:       title,
			Startdate:   startDate,
			Updateddate: time.Now().Format("2006-01-02"), // sqlc made the naming weird >:(
			Completion:  completion,
		})

		return nil
	}

	if err != nil {
		return err
	}

	// The current anime being imported should overwrite the old one !!
	err = cfg.DB.UpdateAnime(cfg.Ctx, database.UpdateAnimeParams{
		Startdate:   startDate,
		Updateddate: time.Now().Format("2006-01-02"),
		Completion:  completion,
		ID:          int64(id),
	})

	if err != nil {
		return err
	}

	return nil
}

// Kitsu also exports in the MAL format
// Kitsu needs to make an HTTP request for EVERY entry. we can use a COOL BUBBLES PROGRESS BAR FOR THAT :fire:
func (cfg DBConfig) ImportMAL(malXml []byte) error {
	var animeList types.Myanimelist
	if err := xml.Unmarshal(malXml, &animeList); err != nil {
		return err
	}

	log.Printf("Importing Anime...")
	for _, anime := range animeList.Anime {
		id, _ := strconv.Atoi(anime.SeriesAnimedbID)

		// Different platforms use different naming
		completion := anime.MyStatus
		if anime.MyStatus == "Plan to Watch" {
			completion = types.PlanToWatch
		} else if anime.MyStatus == "On-Hold" {
			completion = types.OnHold
		}

		cfg.UploadToDB(anime.SeriesTitle, id, anime.MyStartDate, completion)
	}

	return nil
}

func (cfg DBConfig) ImportHianime(hiXml []byte) error {
	return nil
}
