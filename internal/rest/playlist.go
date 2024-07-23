package rest

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type PlaylistService interface {
	Add(playlistModel model.Playlist) error
	GetAll() ([]model.Playlist, error)
	GetByID(id string) (model.Playlist, error)
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error add playlist: %v", err))
	}

	err = p.Service.Add(playlist)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error add playlist: %v", err))
	}

	return c.JSON(http.StatusCreated, playlist)
}

func (p *PlaylistHandler) GetAll(c echo.Context) error {
	playlists, err := p.Service.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error get all playlists: %v", err))
	}

	return c.JSON(http.StatusOK, playlists)
}

func (p *PlaylistHandler) GetByID(c echo.Context) error {
	id := c.Param("id")

	playlist, err := p.Service.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error get playlist by ID: %w", err))
	}

	return c.JSON(http.StatusOK, playlist)
}
