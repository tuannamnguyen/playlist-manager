package rest

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
	spotifyconverter "github.com/tuannamnguyen/playlist-manager/internal/service/converters/spotify"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func saveSessionValues(c echo.Context, store sessions.Store, sessionValues map[any]any) error {
	session, err := store.Get(c.Request(), "oauth-session")
	if err != nil {
		return err
	}

	session.Values = sessionValues

	err = session.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}
	return nil
}

func getSessionValues(c echo.Context, store sessions.Store) (map[any]any, error) {
	session, err := store.Get(c.Request(), "oauth-session")
	if err != nil {
		return nil, err
	}

	return session.Values, nil
}

func spotifyConvertHandler(c echo.Context, token *oauth2.Token, songs []model.SongOutAPI) error {
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(os.Getenv("SPOTIFY_REDIRECT_URL")),
		spotifyauth.WithScopes(
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistReadPrivate,
		),
	)

	client := spotify.New(auth.Client(c.Request().Context(), token), spotify.WithRetry(true))

	err := spotifyconverter.New(client).Export(c.Request().Context(), "test playlist", songs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error converting playlist to spotify: %v", err)
	}

	return c.NoContent(http.StatusOK)
}
