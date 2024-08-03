package rest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type PlaylistService interface {
	Add(ctx context.Context, playlistModel model.PlaylistIn) error
	GetAll(ctx context.Context) ([]model.Playlist, error)
	GetByID(ctx context.Context, id int) (model.Playlist, error)
	DeleteByID(ctx context.Context, id int) error

	AddSongsToPlaylist(ctx context.Context, playlistID int, songs []model.SongIn) error
	GetAllSongsFromPlaylist(ctx context.Context, playlistID int) ([]model.Song, error)
	DeleteSongsFromPlaylist(ctx context.Context, playlistID int, songsID []int) error
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
	var playlist model.PlaylistIn
	err := c.Bind(&playlist)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error binding playlist: %v", err))
	}

	err = p.Service.Add(c.Request().Context(), playlist)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error add playlist: %v", err))
	}

	return c.JSON(http.StatusCreated, playlist)
}

func (p *PlaylistHandler) GetAll(c echo.Context) error {
	playlists, err := p.Service.GetAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error get all playlists: %v", err))
	}

	return c.JSON(http.StatusOK, playlists)
}

func (p *PlaylistHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error converting ID to int: %v", err))
	}

	playlist, err := p.Service.GetByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error get playlist by ID: %v", err))
	}

	return c.JSON(http.StatusOK, playlist)
}

func (p *PlaylistHandler) DeleteByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error converting ID to int: %v", err))
	}

	err = p.Service.DeleteByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error delete playlist by ID: %v", err))
	}

	return c.JSON(http.StatusOK, map[string]int{
		"playlist_id": id,
	})
}

func (p *PlaylistHandler) AddSongsToPlaylist(c echo.Context) error {
	playlistID, err := strconv.Atoi(c.Param("playlist_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error converting ID to int: %v", err))
	}

	var songs []model.SongIn
	err = c.Bind(&songs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error binding request body to songs: %v", err))
	}

	err = p.Service.AddSongsToPlaylist(c.Request().Context(), playlistID, songs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error adding songs to playlist: %v", err))
	}

	return c.JSON(http.StatusOK, songs)
}

func (p *PlaylistHandler) GetAllSongsFromPlaylist(c echo.Context) error {
	playlistID, err := strconv.Atoi(c.Param("playlist_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error converting ID to int: %v", err))
	}

	songs, err := p.Service.GetAllSongsFromPlaylist(c.Request().Context(), playlistID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error get all songs from playlist: %v", err))
	}

	return c.JSON(http.StatusOK, songs)
}

func (p *PlaylistHandler) DeleteSongsFromPlaylist(c echo.Context) error {
	playlistID, err := strconv.Atoi(c.Param("playlist_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error converting ID to int: %v", err))
	}

	var reqBody map[string][]int
	err = c.Bind(&reqBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error binding list of songs ID: %v", err))
	}

	songsID := reqBody["songs_id"]

	err = p.Service.DeleteSongsFromPlaylist(c.Request().Context(), playlistID, songsID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error delete songs from playlist: %v", err))
	}

	return c.JSON(http.StatusOK, reqBody)
}
