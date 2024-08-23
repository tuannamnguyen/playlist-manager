package spotifyconverter

import (
	"github.com/tuannamnguyen/playlist-manager/internal/model"
	"github.com/zmb3/spotify/v2"
)

type SpotifyConverter struct {
	client *spotify.Client
}

func New(client *spotify.Client) *SpotifyConverter {
	return &SpotifyConverter{client: client}
}

func (s *SpotifyConverter) Export(playlistName string, songs []model.SongOutAPI) error {
	// TODO: do this
	return nil
}
