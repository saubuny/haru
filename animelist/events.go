package animelist

import (
	"github.com/saubuny/haru/internal/database"
	"github.com/saubuny/haru/types"
)

type AnimeListMessage types.AnimeListResponse
type AnimeDBListMessage []database.Anime

type AnimeDataResponse struct {
	Data types.AnimeData `json:"data"`
}
