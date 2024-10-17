package rest

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type MetadataService interface {
	FetchLyrics(artists string, song string) (string, error)
	ArtistInformation(artistName string) (string, error)
}

type MetadataHandler struct {
	ms           MetadataService
	sessionStore sessions.Store
}

func NewMetadataHandler(ms MetadataService, store sessions.Store) *MetadataHandler {
	return &MetadataHandler{ms: ms, sessionStore: store}
}

func (m *MetadataHandler) GetLyrics(c echo.Context) error {
	var reqBody struct {
		Song    string `json:"song"`
		Artists string `json:"artists"`
	}

	err := c.Bind(&reqBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	lyrics, err := m.ms.FetchLyrics(reqBody.Artists, reqBody.Song)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"lyrics": lyrics,
	})
}

func (m *MetadataHandler) GetArtistInformation(c echo.Context) error {
	artistName := c.QueryParam("artist_name")

	artistURL, err := m.ms.ArtistInformation(artistName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"redirect_url": artistURL,
	})
}
