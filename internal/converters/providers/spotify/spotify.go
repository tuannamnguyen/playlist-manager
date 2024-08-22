package spotify

import (
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type SpotifyConverter struct {
	// client *spotify.Client
}

func New() *SpotifyConverter {
	// TODO: do this

	return &SpotifyConverter{}
}

func (s *SpotifyConverter) Export(playlistName string, songs []model.SongOutAPI) error {
	// TODO: do this
	return nil
}
