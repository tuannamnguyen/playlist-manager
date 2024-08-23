package spotifyconverter

import (
	"context"
	"fmt"

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
	_, err := s.createPlaylist(ctx, playlistName)
	if err != nil {
		return fmt.Errorf("create playlist: %w", err)
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
