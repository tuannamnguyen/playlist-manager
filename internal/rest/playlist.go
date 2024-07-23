package rest

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type PlaylistService interface {
	Add(playlistModel *model.Playlist) error
	GetAll() ([]model.Playlist, error)
}

type PlaylistHandler struct {
	Service PlaylistService
}

func NewPlaylistHandler(svc PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{
		Service: svc,
	}
}

func (p *PlaylistHandler) Add(c echo.Context) error {
	var playlist model.Playlist
	err := c.Bind(&playlist)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error add: %v", err))
	}

	err = p.Service.Add(&playlist)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error add: %v", err))
	}

	return c.JSON(http.StatusCreated, playlist)
}

func (p *PlaylistHandler) GetAll(c echo.Context) error {
	playlists, err := p.Service.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error get all: %v", err))
	}

	return c.JSON(http.StatusOK, playlists)
}
