package service

import (
	"net/http"
	"os"

	"github.com/tuannamnguyen/playlist-manager/internal/model"
	spotifyconverter "github.com/tuannamnguyen/playlist-manager/internal/service/converters/spotify"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func spotifyConvert(r *http.Request, token *oauth2.Token, songs []model.SongOutAPI) error {
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(os.Getenv("SPOTIFY_REDIRECT_URL")),
		spotifyauth.WithScopes(
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistReadPrivate,
		),
	)

	client := spotify.New(auth.Client(r.Context(), token), spotify.WithRetry(true))

	err := spotifyconverter.New(client).Export(r.Context(), "test playlist", songs)
	return err
}
