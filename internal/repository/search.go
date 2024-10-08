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
	Album   string   `json:"album"`
	Sources []string `json:"sources"`
}

type SearchRepository struct {
	httpClient *http.Client
}

func NewSearchRepository(httpClient *http.Client) *SearchRepository {
	return &SearchRepository{httpClient: httpClient}
}

func (s *SearchRepository) Song(track string, artist string, album string) ([]model.SongInAPI, error) {
	searchReqBody := SearchRequest{
		Track:  track,
		Artist: artist,
		Album:  album,
		Type:   "track",
		Sources: []string{
			"spotify",
			"appleMusic",
			"tidal",
			"amazonMusic",
		},
	}
	searchReqBodyEncoded, err := json.Marshal(searchReqBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling search request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/public/search", os.Getenv("MUSIC_API_ENDPOINT")), bytes.NewBuffer(searchReqBodyEncoded))
	if err != nil {
		return nil, &requestMarshalError{err}
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", os.Getenv("MUSIC_API_CLIENT_ID")))
	req.Header.Set("Content-Type", echo.MIMEApplicationJSON)

	res, err := s.httpClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		return nil, fmt.Errorf("fetching info from music api: %w", err)
	}
	defer res.Body.Close()

	var searchRes SearchResponse
	err = json.NewDecoder(res.Body).Decode(&searchRes)
	if err != nil {
		return nil, &responseDecodeError{err}
	}

	return transformSearchAPIResponse(searchRes), nil
}
