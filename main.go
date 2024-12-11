package main

import (
	"context"
	_ "embed"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/saubuny/haru/internal/database"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"

	"github.com/kevm/bubbleo/navstack"
)

type dbConfig struct {
	DB  *database.Queries
	Ctx context.Context
}

//go:embed sql/schema/schema.sql
var migrations string

func main() {
	cfg, err := initDB(migrations, "anime.db")
	if err != nil {
		log.Fatalf("Error initalizing DB: %v", err)
	}

	var importFile string
	var importPlatform string

	// Run TUI by default
	app := &cli.App{
		Name:  "Haru",
		Usage: "Track anime",
		Action: func(ctx *cli.Context) error {
			m := initialModel(cfg)
			m.Shell.Navstack.Push(navstack.NavigationItem{Model: m, Title: "Haru"})
			p := tea.NewProgram(m.Shell, tea.WithAltScreen())
			tea.SetWindowTitle("Haru")
			if _, err := p.Run(); err != nil {
				log.Fatalf("Error: %v", err)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "import",
				Aliases: []string{"i"},
				Usage:   "import from another tracking platform",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:        "source",
						Usage:       "file to import",
						Destination: &importFile,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "platform",
						Usage:       "platform to import from (must be one of Hianime or MAL)",
						Destination: &importPlatform,
						Required:    true,
						Action: func(ctx *cli.Context, s string) error {
							validPlatforms := []string{"hianime", "mal"}
							if !slices.Contains(validPlatforms, strings.ToLower(importPlatform)) {
								return cli.Exit("Invalid platform", 1)
							}
							return nil
						},
					},
				},
				Action: func(ctx *cli.Context) error {
					file, err := os.ReadFile(importFile)
					if err != nil {
						return err
					}

					if strings.ToLower(importPlatform) == "mal" {
						err = cfg.importMAL(file)
					} else if strings.ToLower(importPlatform) == "hianime" {
						err = cfg.importHianime(file)
					}
					if err != nil {
						return err
					}

					anime, err := cfg.DB.GetAllAnime(cfg.Ctx) // TMP
					if err != nil {
						return err
					}

					log.Printf("%#v\n", anime)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
