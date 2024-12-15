package types

import (
	"encoding/xml"
)

// A central completion type to convert all formats into
const (
	Watching    = "Watching"
	PlanToWatch = "Plan To Watch"
	Completed   = "Completed"
	OnHold      = "On Hold"
	Dropped     = "Dropped"
)

type ErrorMsg string

type AnimeData struct {
	MalID  int    `json:"mal_id"`
	URL    string `json:"url"`
	Images struct {
		Jpg struct {
			ImageURL      string `json:"image_url"`
			SmallImageURL string `json:"small_image_url"`
			LargeImageURL string `json:"large_image_url"`
		} `json:"jpg"`
		Webp struct {
			ImageURL      string `json:"image_url"`
			SmallImageURL string `json:"small_image_url"`
			LargeImageURL string `json:"large_image_url"`
		} `json:"webp"`
	} `json:"images"`
	Trailer struct {
		YoutubeID string `json:"youtube_id"`
		URL       string `json:"url"`
		EmbedURL  string `json:"embed_url"`
	} `json:"trailer"`
	Approved bool `json:"approved"`
	Titles   []struct {
		Type  string `json:"type"`
		Title string `json:"title"`
	} `json:"titles"`
	Title         string   `json:"title"`
	TitleEnglish  string   `json:"title_english"`
	TitleJapanese string   `json:"title_japanese"`
	TitleSynonyms []string `json:"title_synonyms"`
	Type          string   `json:"type"`
	Source        string   `json:"source"`
	Episodes      int      `json:"episodes"`
	Status        string   `json:"status"`
	Airing        bool     `json:"airing"`
	Aired         struct {
		From string `json:"from"`
		To   string `json:"to"`
		Prop struct {
			From struct {
				Day   int `json:"day"`
				Month int `json:"month"`
				Year  int `json:"year"`
			} `json:"from"`
			To struct {
				Day   int `json:"day"`
				Month int `json:"month"`
				Year  int `json:"year"`
			} `json:"to"`
			String string `json:"string"`
		} `json:"prop"`
	} `json:"aired"`
	Duration   string  `json:"duration"`
	Rating     string  `json:"rating"`
	Score      float64 `json:"score"`
	ScoredBy   int     `json:"scored_by"`
	Rank       int     `json:"rank"`
	Popularity int     `json:"popularity"`
	Members    int     `json:"members"`
	Favorites  int     `json:"favorites"`
	Synopsis   string  `json:"synopsis"`
	Background string  `json:"background"`
	Season     string  `json:"season"`
	Year       int     `json:"year"`
	Broadcast  struct {
		Day      string `json:"day"`
		Time     string `json:"time"`
		Timezone string `json:"timezone"`
		String   string `json:"string"`
	} `json:"broadcast"`
	Producers []struct {
		MalID int    `json:"mal_id"`
		Type  string `json:"type"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"producers"`
	Licensors []struct {
		MalID int    `json:"mal_id"`
		Type  string `json:"type"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"licensors"`
	Studios []struct {
		MalID int    `json:"mal_id"`
		Type  string `json:"type"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"studios"`
	Genres []struct {
		MalID int    `json:"mal_id"`
		Type  string `json:"type"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"genres"`
	ExplicitGenres []struct {
		MalID int    `json:"mal_id"`
		Type  string `json:"type"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"explicit_genres"`
	Themes []struct {
		MalID int    `json:"mal_id"`
		Type  string `json:"type"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"themes"`
	Demographics []struct {
		MalID int    `json:"mal_id"`
		Type  string `json:"type"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"demographics"`
}

type AnimeListResponse struct {
	Data       []AnimeData `json:"data"`
	Pagination struct {
		LastVisiblePage int  `json:"last_visible_page"`
		HasNextPage     bool `json:"has_next_page"`
		Items           struct {
			Count   int `json:"count"`
			Total   int `json:"total"`
			PerPage int `json:"per_page"`
		} `json:"items"`
	} `json:"pagination"`
}

type Myanimelist struct {
	XMLName xml.Name `xml:"myanimelist"`
	Text    string   `xml:",chardata"`
	Myinfo  struct {
		Text           string `xml:",chardata"`
		UserExportType string `xml:"user_export_type"`
	} `xml:"myinfo"`
	Anime []struct {
		Text              string `xml:",chardata"`
		SeriesTitle       string `xml:"series_title"`
		SeriesAnimedbID   string `xml:"series_animedb_id"`
		MyWatchedEpisodes string `xml:"my_watched_episodes"`
		MyReadVolumes     string `xml:"my_read_volumes"`
		MyStartDate       string `xml:"my_start_date"`
		MyStatus          string `xml:"my_status"`
		MyTimesWatched    string `xml:"my_times_watched"`
		UpdateOnImport    string `xml:"update_on_import"`
		MyFinishDate      string `xml:"my_finish_date"`
		MyScore           string `xml:"my_score"`
	} `xml:"anime"`
}

type HiAnimeList struct {
	XMLName xml.Name `xml:"list"`
	Text    string   `xml:",chardata"`
	Folder  []struct {
		Text string `xml:",chardata"`
		Name string `xml:"name"`
		Data struct {
			Text string `xml:",chardata"`
			Item []struct {
				Text string `xml:",chardata"`
				Name string `xml:"name"`
				Link string `xml:"link"`
			} `xml:"item"`
		} `xml:"data"`
	} `xml:"folder"`
}
