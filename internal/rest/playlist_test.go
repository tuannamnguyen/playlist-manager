package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
	"github.com/tuannamnguyen/playlist-manager/internal/rest/mocks"
)

func TestAddSongToPlaylist(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// create fake request body
		mockSongs := []model.Song{
			{
				ID:       "id",
				Name:     "name",
				ArtistID: "artist_id",
				AlbumID:  "album_id",
			},
		}
		requestBody, err := json.Marshal(mockSongs)
		assert.NoError(t, err)

		// create fake request and response recorder
		req, err := http.NewRequest(
			http.MethodPost,
			"/api/playlists/abcd/songs",
			bytes.NewReader(requestBody),
		)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// setup echo route handler
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/api/playlists/:playlist_id/songs")
		c.SetParamNames("playlist_id")
		c.SetParamValues("abcd")

		mockPlaylistService := new(mocks.MockPlaylistService)
		mockPlaylistService.On("AddSongsToPlaylist", "abcd", mockSongs).Return(nil)

		handler := NewPlaylistHandler(mockPlaylistService)
		err = handler.AddSongsToPlaylist(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockPlaylistService.AssertExpectations(t)
	})
}

func TestGetAllSongsFromPlaylist(t *testing.T) {
	mockPlaylistService := new(mocks.MockPlaylistService)

	t.Run("success get all songs from playlist", func(t *testing.T) {
		mockPlaylistService.On("GetAllSongsFromPlaylist", "abcd").Return([]model.Song{}, nil)

		req, err := http.NewRequest(
			http.MethodGet,
			"/api/playlists/abcd/songs",
			nil,
		)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/api/playlists/:playlist_id/songs")
		c.SetParamNames("playlist_id")
		c.SetParamValues("abcd")

		handler := NewPlaylistHandler(mockPlaylistService)
		err = handler.GetAllSongsFromPlaylist(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockPlaylistService.AssertExpectations(t)
	})

	t.Run("error get all songs", func(t *testing.T) {
		mockPlaylistService.On("GetAllSongsFromPlaylist", "defg").Return(nil, errors.New("test error"))

		req, err := http.NewRequest(
			http.MethodGet,
			"/api/playlists/defg/songs",
			nil,
		)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/api/playlists/:playlist_id/songs")
		c.SetParamNames("playlist_id")
		c.SetParamValues("defg")

		handler := NewPlaylistHandler(mockPlaylistService)
		err = handler.GetAllSongsFromPlaylist(c)
		require.EqualError(t, err, echo.NewHTTPError(http.StatusInternalServerError, "error get all songs from playlist: test error").Error())

		mockPlaylistService.AssertExpectations(t)
	})

}
