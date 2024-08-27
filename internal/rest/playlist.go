package rest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
	"golang.org/x/oauth2"
)

type PlaylistService interface {
	// playlist operations
	Add(ctx context.Context, playlistModel model.PlaylistIn) error
	GetAll(ctx context.Context) ([]model.Playlist, error)
	GetByID(ctx context.Context, id int) (model.Playlist, error)
	DeleteByID(ctx context.Context, id int) error

	// playlist-song operations
	AddSongsToPlaylist(ctx context.Context, playlistID int, songs []model.SongInAPI) error
	GetAllSongsFromPlaylist(ctx context.Context, playlistID int) ([]model.SongOutAPI, error)
	DeleteSongsFromPlaylist(ctx context.Context, playlistID int, songsID []int) error

	// convert operation
	Convert(ctx context.Context, provider string, token *oauth2.Token, songs []model.SongOutAPI) error
}

type PlaylistHandler struct {
	service      PlaylistService
	sessionStore sessions.Store
}

func NewPlaylistHandler(svc PlaylistService, store sessions.Store) *PlaylistHandler {
	return &PlaylistHandler{
		service:      svc,
		sessionStore: store,
	}
}

func (p *PlaylistHandler) Add(c echo.Context) error {
	var playlist model.PlaylistIn
	err := c.Bind(&playlist)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error binding playlist: %v", err))
	}

	err = p.service.Add(c.Request().Context(), playlist)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error add playlist: %v", err))
	}

	return c.JSON(http.StatusCreated, playlist)
}

func (p *PlaylistHandler) GetAll(c echo.Context) error {
	playlists, err := p.service.GetAll(c.Request().Context())
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

	playlist, err := p.service.GetByID(c.Request().Context(), id)
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

	err = p.service.DeleteByID(c.Request().Context(), id)
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

	var songs []model.SongInAPI
	err = c.Bind(&songs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error binding request body to songs: %v", err))
	}

	err = p.service.AddSongsToPlaylist(c.Request().Context(), playlistID, songs)
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

	songs, err := p.service.GetAllSongsFromPlaylist(c.Request().Context(), playlistID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting all songs from playlist: %v", err))
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

	err = p.service.DeleteSongsFromPlaylist(c.Request().Context(), playlistID, songsID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error delete songs from playlist: %v", err))
	}

	return c.JSON(http.StatusOK, reqBody)
}

func (p *PlaylistHandler) ConvertHandler(c echo.Context) error {
	playlistID, err := strconv.Atoi(c.Param("playlist_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error converting ID to int: %v", err))
	}

	provider := c.Param("provider")

	songs, err := p.service.GetAllSongsFromPlaylist(c.Request().Context(), playlistID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting all songs from playlist: %v", err))
	}

	sessionValues, err := getSessionValues(c, p.sessionStore)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting session values: %v", err))
	}

	accessToken := (sessionValues[fmt.Sprintf("%s_access_token", provider)]).(string)
	refreshToken := (sessionValues[fmt.Sprintf("%s_refresh_token", provider)]).(string)
	expiry := (sessionValues[fmt.Sprintf("%s_token_expiry", provider)]).(time.Time)

	token := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       expiry,
	}

	return p.service.Convert(c.Request().Context(), provider, token, songs)
}
