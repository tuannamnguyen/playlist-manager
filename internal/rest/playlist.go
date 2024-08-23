package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
	spotifyconverter "github.com/tuannamnguyen/playlist-manager/internal/service/converters/providers/spotify"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
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
}

type PlaylistHandler struct {
	Service      PlaylistService
	sessionStore sessions.Store
}

func NewPlaylistHandler(svc PlaylistService, store sessions.Store) *PlaylistHandler {
	return &PlaylistHandler{
		Service:      svc,
		sessionStore: store,
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

	var songs []model.SongInAPI
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

	err = p.Service.DeleteSongsFromPlaylist(c.Request().Context(), playlistID, songsID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error delete songs from playlist: %v", err))
	}

	return c.JSON(http.StatusOK, reqBody)
}

func (p *PlaylistHandler) SpotifyConvertHandler(c echo.Context) error {
	playlistID, err := strconv.Atoi(c.Param("playlist_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error converting ID to int: %v", err))
	}

	songs, err := p.Service.GetAllSongsFromPlaylist(c.Request().Context(), playlistID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting all songs from playlist: %v", err))
	}

	session, err := p.sessionStore.Get(c.Request(), "oauth-session")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting session store: %v", err))
	}

	accessToken := (session.Values["spotify_access_token"]).(string)
	refreshToken := (session.Values["spotify_refresh_token"]).(string)
	expiry := (session.Values["spotify_token_expiry"]).(time.Time)

	token := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       expiry,
	}

	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(os.Getenv("SPOTIFY_REDIRECT_URL")),
		spotifyauth.WithScopes(
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistReadPrivate,
		),
	)

	client := spotify.New(auth.Client(c.Request().Context(), token))

	err = spotifyconverter.New(client).Export(c.Request().Context(), "test playlist", songs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error converting playlist to spotify: %v", err)
	}

	return nil
}
