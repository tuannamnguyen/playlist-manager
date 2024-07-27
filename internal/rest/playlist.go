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
	DeleteByID(id string) error

	AddSongsToPlaylist(playlistID string, songs []model.Song) error
	GetAllSongsFromPlaylist(playlist_id string) ([]model.Song, error)
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error get playlist by ID: %v", err))
	}

	return c.JSON(http.StatusOK, playlist)
}

func (p *PlaylistHandler) DeleteByID(c echo.Context) error {
	id := c.Param("id")

	err := p.Service.DeleteByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error delete playlist by ID: %v", err))
	}

	return c.JSON(http.StatusOK, map[string]string{
		"playlist_id": id,
	})
}

func (p *PlaylistHandler) AddSongsToPlaylist(c echo.Context) error {
	playlistID := c.Param("playlist_id")

	var songs []model.Song
	err := c.Bind(&songs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error binding request body to songs: %v", err))
	}

	err = p.Service.AddSongsToPlaylist(playlistID, songs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error adding songs to playlist: %v", err))
	}

	return c.JSON(http.StatusOK, songs)
}

func (p *PlaylistHandler) GetAllSongsFromPlaylist(c echo.Context) error {
	playlistID := c.Param("playlist_id")

	songs, err := p.Service.GetAllSongsFromPlaylist(playlistID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error get all songs from playlist: %v", err))
	}

	return c.JSON(http.StatusOK, songs)
}
