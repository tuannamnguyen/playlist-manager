package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

// Root structure for the JSON data
type SearchResponse struct {
	Tracks []Track `json:"tracks"`
}

// Track structure for individual track information
type Track struct {
	Source string `json:"source"`
	Status string `json:"status"`
	Data   Data   `json:"data"`
	Type   string `json:"type"`
}

// Data structure for the detailed track information
type Data struct {
	ExternalID  string   `json:"externalId"`
	PreviewURL  string   `json:"previewUrl"`
	Name        string   `json:"name"`
	ArtistNames []string `json:"artistNames"`
	AlbumName   string   `json:"albumName"`
	ImageURL    string   `json:"imageUrl"`
	ISRC        string   `json:"isrc"`
	Duration    int      `json:"duration"`
	URL         string   `json:"url"`
}

type SearchRequest struct {
	Track   string   `json:"track"`
	Artist  string   `json:"artist"`
	Type    string   `json:"type"`
	Sources []string `json:"sources"`
}

type SearchRepository struct {
	httpClient *http.Client
}

func NewSearchRepository(httpClient *http.Client) *SearchRepository {
	return &SearchRepository{httpClient: httpClient}
}

func (s *SearchRepository) Song(track string, artist string) (model.SongIn, error) {
	searchReqBody := SearchRequest{
		Track:   track,
		Artist:  artist,
		Type:    "track",
		Sources: []string{"spotify"},
	}
	searchReqBodyEncoded, err := json.Marshal(searchReqBody)
	if err != nil {
		return model.SongIn{}, fmt.Errorf("marshalling search request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/public/search", os.Getenv("MUSIC_API_ENDPOINT")), bytes.NewBuffer(searchReqBodyEncoded))
	if err != nil {
		return model.SongIn{}, fmt.Errorf("making search request body: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", os.Getenv("MUSIC_API_CLIENT_ID")))
	req.Header.Set("Content-Type", echo.MIMEApplicationJSON)

	res, err := s.httpClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		return model.SongIn{}, fmt.Errorf("fetching info from music api: %w", err)
	}
	defer res.Body.Close()

	var searchRes SearchResponse
	err = json.NewDecoder(res.Body).Decode(&searchRes)
	if err != nil {
		return model.SongIn{}, fmt.Errorf("decoding music api response: %w", err)
	}

	return model.SongIn{
		Name:        searchRes.Tracks[0].Data.Name,
		ArtistNames: searchRes.Tracks[0].Data.ArtistNames,
		AlbumName:   searchRes.Tracks[0].Data.AlbumName,
	}, nil
}
