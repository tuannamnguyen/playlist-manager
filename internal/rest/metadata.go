package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MetadataService interface {
	FetchLyrics(artists string, song string) (string, error)
}

type MetadataHandler struct {
	ms MetadataService
}

func NewMetadataHandler(ms MetadataService) *MetadataHandler {
	return &MetadataHandler{ms: ms}
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
