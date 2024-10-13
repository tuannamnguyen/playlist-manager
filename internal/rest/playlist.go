package rest

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
	"golang.org/x/oauth2"
)

type PlaylistService interface {
	// playlist operations
	Add(ctx context.Context, playlistModel model.PlaylistIn, imageFile multipart.File, imageHeader *multipart.FileHeader) error
	GetAll(ctx context.Context, userID string) ([]model.Playlist, error)
	GetByID(ctx context.Context, id int) (model.Playlist, error)
	DeleteByID(ctx context.Context, id int) error

	// playlist-song operations
	AddSongsToPlaylist(ctx context.Context, playlistID int, songs []model.SongInAPI) error
	GetAllSongsFromPlaylist(ctx context.Context, playlistID int, sortBy string, sortOrder string) ([]model.SongOutAPI, error)
	DeleteSongsFromPlaylist(ctx context.Context, playlistID int, songsID []int) error

	// convert operation
	Convert(ctx context.Context, provider string, token *oauth2.Token, playlistName string, songs []model.SongOutAPI) error
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
	playlist := model.PlaylistIn{
		Name:                c.FormValue("playlist_name"),
		PlaylistDescription: c.FormValue("playlist_description"),
		UserID:              c.FormValue("user_id"),
		Username:            c.FormValue("user_name"),
	}

	if err := c.Validate(playlist); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	header, err := c.FormFile("playlist_cover_image")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	file, err := header.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer file.Close()

	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	fileType := http.DetectContentType(buff)
	log.Println(fileType)

	if fileType != "image/jpeg" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file type for playlist cover")
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = p.service.Add(c.Request().Context(), playlist, file, header)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, playlist)
}

func (p *PlaylistHandler) GetAll(c echo.Context) error {
	userID := c.QueryParam("user_id")

	playlists, err := p.service.GetAll(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, playlists)
}

func (p *PlaylistHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	playlist, err := p.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, playlist)
}

func (p *PlaylistHandler) DeleteByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = p.service.DeleteByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]int{
		"playlist_id": id,
	})
}

func (p *PlaylistHandler) AddSongsToPlaylist(c echo.Context) error {
	playlistID, err := strconv.Atoi(c.Param("playlist_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var songs []model.SongInAPI
	err = c.Bind(&songs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = p.service.AddSongsToPlaylist(c.Request().Context(), playlistID, songs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, songs)
}

func (p *PlaylistHandler) GetAllSongsFromPlaylist(c echo.Context) error {
	type QueryParams struct {
		SortBy    string `query:"sort_by" validate:"omitempty,oneof=s.song_name al.album_name pls.created_at"`
		SortOrder string `query:"sort_order" validate:"required_with=SortBy,omitempty,oneof=ASC DESC"`
	}
	var qParams QueryParams

	err := c.Bind(&qParams)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(qParams); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	playlistID, err := strconv.Atoi(c.Param("playlist_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	songs, err := p.service.GetAllSongsFromPlaylist(c.Request().Context(), playlistID, qParams.SortBy, qParams.SortOrder)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, songs)
}

func (p *PlaylistHandler) DeleteSongsFromPlaylist(c echo.Context) error {
	playlistID, err := strconv.Atoi(c.Param("playlist_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var reqBody map[string][]int
	err = c.Bind(&reqBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	songsID := reqBody["songs_id"]

	err = p.service.DeleteSongsFromPlaylist(c.Request().Context(), playlistID, songsID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, reqBody)
}

func (p *PlaylistHandler) ConvertHandler(c echo.Context) error {
	type reqBody struct {
		PlaylistName     string `json:"playlist_name"`
		ProviderMetadata struct {
			AppleMusic struct {
				MusicUserToken string `json:"musicUserToken"`
			} `json:"applemusic"`
		} `json:"provider_metadata,omitempty"`
	}

	playlistID, err := strconv.Atoi(c.Param("playlist_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var newReqBody reqBody
	err = c.Bind(&newReqBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	songs, err := p.service.GetAllSongsFromPlaylist(c.Request().Context(), playlistID, "", "")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	sessionValues, err := getOauthSessionValues(c.Request(), p.sessionStore)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	provider := c.Param("provider")
	user := (sessionValues[fmt.Sprintf("%s_user_info", provider)]).(goth.User)
	token := &oauth2.Token{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		Expiry:       user.ExpiresAt,
	}

	err = p.service.Convert(c.Request().Context(), provider, token, newReqBody.PlaylistName, songs)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Converted successfully",
	})
}
