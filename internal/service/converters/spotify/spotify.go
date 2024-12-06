package spotifyconverter

import (
	"context"
	"fmt"
	"log"
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
			spotifyauth.ScopeStreaming,
		),
	)

	client := spotify.New(auth.Client(ctx, token), spotify.WithRetry(true))

	return &SpotifyConverter{client: client}
}

func (s *SpotifyConverter) Export(ctx context.Context, playlistName string, songs []model.SongOutAPI) error {
	playlistID, err := s.createPlaylist(ctx, playlistName)
	if err != nil {
		return fmt.Errorf("create playlist: %w", err)
	}

	var tracksID []spotify.ID
	for _, song := range songs {
		searchQuery := s.formatSearchQuery(song)
		log.Printf("spotify search query: %s", searchQuery)

		result, err := s.client.Search(
			ctx,
			searchQuery,
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
	artistQuery := make([]string, len(song.ArtistNames))
	for i, artist := range song.ArtistNames {
		artistQuery[i] = fmt.Sprintf("artist:%s", artist)
	}
	artists := strings.Join(artistQuery, " ")

	trackQuery := fmt.Sprintf("track:%s", song.Name)
	albumQuery := fmt.Sprintf("album:%s", song.AlbumName)
	isrcQuery := fmt.Sprintf("isrc:%s", song.ISRC)

	var queryParts []string
	if song.ISRC != "" {
		queryParts = []string{isrcQuery}

	} else {
		queryParts = []string{trackQuery, artists, albumQuery}
	}

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
