package db

import (
	"reflect"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/saubuny/haru/internal/database"
)

// This project only really needs to test the importing logic for the database
func TestImportMal1(t *testing.T) {
	// Create test database in memory
	migrations := `CREATE TABLE IF NOT EXISTS anime (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    startDate TEXT NOT NULL,
    updatedDate TEXT NOT NULL,
    completion TEXT NOT NULL
);`

	cfg, err := InitDB(migrations, ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	// Short list of Mal XML
	xml1 := `<?xml version="1.0" encoding="UTF-8" ?>
        <myanimelist>
            <anime>
                <series_animedb_id>853</series_animedb_id>
                <series_title><![CDATA[Ouran Koukou Host Club]]></series_title>
                <my_start_date>2022-01-07</my_start_date>
                <my_status>Dropped</my_status>
            </anime>
            <anime>
					<series_animedb_id>66</series_animedb_id>
					<series_title><![CDATA[Azumanga Daiou The Animation]]></series_title>
					<my_start_date>0000-00-00</my_start_date>
					<my_status>Plan to Watch</my_status>
            </anime>
            <anime>
					<series_animedb_id>21</series_animedb_id>
					<series_title><![CDATA[One Piece]]></series_title>
					<my_start_date>2021-07-06</my_start_date>
					<my_status>Dropped</my_status>
            </anime>
        </myanimelist>
    `

	// Second list of Mal XML with some conflicting entries
	xml2 := `<?xml version="1.0" encoding="UTF-8" ?>
        <myanimelist>
            <anime>
					<series_animedb_id>66</series_animedb_id>
					<series_title><![CDATA[Azumanga Daiou The Animation]]></series_title>
					<my_start_date>0000-00-00</my_start_date>
					<my_status>Plan to Watch</my_status>
            </anime>
            <anime>
					<series_animedb_id>21</series_animedb_id>
					<series_title><![CDATA[One Piece]]></series_title>
					<my_start_date>2024-11-13</my_start_date>
					<my_status>Watching</my_status>
            </anime>
            <anime>
					<series_animedb_id>30276</series_animedb_id>
					<series_title><![CDATA[One Punch Man]]></series_title>
					<my_start_date>2020-02-05</my_start_date>
					<my_status>Completed</my_status>
            </anime>
        </myanimelist> `
	// Import to both to DB
	err = cfg.ImportMAL([]byte(xml1))
	if err != nil {
		t.Fatal(err)
	}

	err = cfg.ImportMAL([]byte(xml2))
	if err != nil {
		t.Fatal(err)
	}

	// Expected result !!
	expected := []database.Anime{
		{
			ID:          21,
			Title:       "One Piece",
			Startdate:   "2024-11-13",
			Updateddate: time.Now().Format("2006-01-02"),
			Completion:  "Watching",
		},
		{
			ID:          66,
			Title:       "Azumanga Daiou The Animation",
			Startdate:   "0000-00-00",
			Updateddate: time.Now().Format("2006-01-02"),
			Completion:  "Plan To Watch",
		},
		{
			ID:          853,
			Title:       "Ouran Koukou Host Club",
			Startdate:   "2022-01-07",
			Updateddate: time.Now().Format("2006-01-02"),
			Completion:  "Dropped",
		},
		{
			ID:          30276,
			Title:       "One Punch Man",
			Startdate:   "2020-02-05",
			Updateddate: time.Now().Format("2006-01-02"),
			Completion:  "Completed",
		},
	}

	// Export from DB and verify output is correct (new import should overwrite old, with no duplicate entries)
	dbState, err := cfg.DB.GetAllAnime(cfg.Ctx)
	if !reflect.DeepEqual(dbState, expected) {
		t.Fatalf("dbState differs from expected input:\n%#v\n%#v\n", dbState, expected)
	}
}
