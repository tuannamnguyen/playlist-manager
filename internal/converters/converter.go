package converters

import "github.com/tuannamnguyen/playlist-manager/internal/model"

type Converter interface {
	Export(playlistName string, songs []model.SongOutAPI) error
}
