package spotifyconverter

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/tuannamnguyen/playlist-manager/internal/model"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type SpotifyConverter struct {
	client *spotify.Client
}

func New(ctx context.Context, token *oauth2.Token) *SpotifyConverter {
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(os.Getenv("SPOTIFY_REDIRECT_URL")),
		spotifyauth.WithScopes(
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistReadPrivate,
		),
	)

	client := spotify.New(auth.Client(ctx, token), spotify.WithRetry(true))

	return &SpotifyConverter{client: client}
}

func (s *SpotifyConverter) Export(ctx context.Context, playlistName string, songs []model.SongOutAPI) error {
	// ? what if the playlist already exists
	playlistID, err := s.createPlaylist(ctx, playlistName)
	if err != nil {
		return fmt.Errorf("create playlist: %w", err)
	}

	// ? see if we can use goroutines here
	var tracksID []spotify.ID
	for _, song := range songs {
		// TODO: need to match songs better
		searchQuery := fmt.Sprintf("track:%s artist:%s album:%s", song.Name, song.ArtistNames, song.AlbumName)
		result, err := s.client.Search(
			ctx,
			url.QueryEscape(searchQuery),
			spotify.SearchTypeTrack,
			spotify.Limit(1),
		)
		if err != nil {
			return fmt.Errorf("search for song in spotify: %w", err)
		}

		tracksID = append(tracksID, result.Tracks.Tracks[0].ID)
	}

	// TODO: need to break into 100 songs chunk
	_, err = s.client.AddTracksToPlaylist(ctx, spotify.ID(playlistID), tracksID...)
	if err != nil {
		return fmt.Errorf("add track to playlist: %w", err)
	}

	return nil
}

func (s *SpotifyConverter) createPlaylist(ctx context.Context, playlistName string) (string, error) {
	currentUser, err := s.client.CurrentUser(ctx)
	if err != nil {
		return "", fmt.Errorf("get current user: %w", err)
	}

	playlistDetail, err := s.client.CreatePlaylistForUser(ctx, currentUser.User.ID, playlistName, "", false, false)
	if err != nil {
		return "", err
	}

	return string(playlistDetail.ID), nil
}
