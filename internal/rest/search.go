package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type SearchService interface {
	SongSearch(track string, artist string) (model.Song, error)
}

type SearchHandler struct {
	service SearchService
}

func NewSearchHandler() *SearchHandler {
	return &SearchHandler{}
}

func (s *SearchHandler) SearchMusicData(c echo.Context) error {
	type searchBody struct {
		Track  string `json:"track"`
		Artist string `json:"artist"`
	}

	var reqBody searchBody
	err := c.Bind(&reqBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error binding search body")
	}

	song, err := s.service.SongSearch(reqBody.Track, reqBody.Artist)
	if err != nil {
		c.Logger().Error()
		return echo.NewHTTPError(http.StatusInternalServerError, "error searching for song")
	}

	return c.JSON(http.StatusOK, song)
}
