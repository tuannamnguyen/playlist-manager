package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type SearchService interface {
	SongSearch(track string, artist string, album string) (model.SongInAPI, error)
}

type SearchHandler struct {
	service SearchService
}

func NewSearchHandler(service SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

func (s *SearchHandler) SearchMusicData(c echo.Context) error {
	type searchBody struct {
		Track  string `json:"track"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
	}

	var reqBody searchBody
	err := c.Bind(&reqBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error binding search body")
	}

	song, err := s.service.SongSearch(reqBody.Track, reqBody.Artist, reqBody.Album)
	if err != nil {
		c.Logger().Error()
		return echo.NewHTTPError(http.StatusInternalServerError, "error searching for song")
	}

	return c.JSON(http.StatusOK, song)
}
