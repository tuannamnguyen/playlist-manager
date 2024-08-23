package spotifyconverter

import (
	"context"
	"fmt"
	"net/url"

	"github.com/tuannamnguyen/playlist-manager/internal/model"
	"github.com/zmb3/spotify/v2"
)

type SpotifyConverter struct {
	client *spotify.Client
}

func New(client *spotify.Client) *SpotifyConverter {
	return &SpotifyConverter{client: client}
}

func (s *SpotifyConverter) Export(ctx context.Context, playlistName string, songs []model.SongOutAPI) error {
	// TODO: not done

	// ? what if the playlist already exists
	playlistID, err := s.createPlaylist(ctx, playlistName)
	if err != nil {
		return fmt.Errorf("create playlist: %w", err)
	}

	// ? see if we can use goroutines here
	var searchResults []*spotify.SearchResult
	for _, song := range songs {
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

		searchResults = append(searchResults, result)
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
