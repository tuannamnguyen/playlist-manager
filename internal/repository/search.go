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
	Sources []string `json:"string"`
}

type SearchRepository struct {
	httpClient *http.Client
}

func NewSearchRepository(httpClient *http.Client) *SearchRepository {
	return &SearchRepository{httpClient: httpClient}
}

func (s *SearchRepository) Song(track string, artist string) (model.Song, error) {
	searchReqBody := SearchRequest{
		Track:   track,
		Artist:  artist,
		Type:    "track",
		Sources: []string{"spotify"},
	}
	searchReqBodyEncoded, err := json.Marshal(searchReqBody)
	if err != nil {
		return model.Song{}, fmt.Errorf("marshalling search request body: %w", err)
	}

	res, err := s.httpClient.Post(os.Getenv("MUSIC_API_ENDPOINT"), echo.MIMEApplicationJSON, bytes.NewBuffer(searchReqBodyEncoded))
	if err != nil {
		return model.Song{}, fmt.Errorf("fetching info from music api: %w", err)
	}
	defer res.Body.Close()

	var searchRes SearchResponse
	err = json.NewDecoder(res.Body).Decode(&searchRes)
	if err != nil {
		return model.Song{}, fmt.Errorf("decoding music api response: %w", err)
	}

	return model.Song{
		Name:     searchRes.Tracks[0].Data.Name,
		ArtistID: searchRes.Tracks[0].Data.ArtistNames[0], // TODO: UPDATE SONGS MODEL IN DATABASE LATER TO ALLOW MULTIPLE ARTISTS PER SONG
		AlbumID:  searchRes.Tracks[0].Data.AlbumName,
	}, nil
}
