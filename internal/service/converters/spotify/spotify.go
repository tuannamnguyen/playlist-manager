package spotifyconverter

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"

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
		searchQuery := s.formatSearchQuery(song)
		slog.Info(searchQuery)
		slog.Info("url formatted search: ", "query", url.QueryEscape(searchQuery))

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

	chunkedTracksID := chunkBy(tracksID, 100)
	for _, IDs := range chunkedTracksID {
		_, err = s.client.AddTracksToPlaylist(ctx, spotify.ID(playlistID), IDs...)
		if err != nil {
			return fmt.Errorf("add track to playlist: %w", err)
		}
	}

	return nil
}

func (s *SpotifyConverter) formatSearchQuery(song model.SongOutAPI) string {
	// Helper function to wrap terms with spaces in quotes
	wrapInQuotes := func(s string) string {
		if strings.Contains(s, " ") {
			return fmt.Sprintf(`"%s"`, s)
		}
		return s
	}

	// Format artist query
	artistQuery := make([]string, len(song.ArtistNames))
	for i, artist := range song.ArtistNames {
		artistQuery[i] = fmt.Sprintf("artist:%s", wrapInQuotes(artist))
	}
	artists := strings.Join(artistQuery, " ")

	// Format track, album, and year (if available)
	trackQuery := fmt.Sprintf("track:%s", wrapInQuotes(song.Name))
	albumQuery := fmt.Sprintf("album:%s", wrapInQuotes(song.AlbumName))

	// Combine all parts of the query
	queryParts := []string{trackQuery, artists, albumQuery}

	return strings.Join(queryParts, " ")
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
