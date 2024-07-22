package service

import "github.com/tuannamnguyen/playlist-manager/internal/model"

type PlaylistRepository interface {
	Add(playlistModel *model.Playlist) error
}

type Playlist struct {
	playlistRepo PlaylistRepository
}

func NewPlaylist(playlistRepo PlaylistRepository) *Playlist {
	return &Playlist{
		playlistRepo: playlistRepo,
	}
}

func (p *Playlist) Add(playlistModel *model.Playlist) error {
	err := p.playlistRepo.Add(playlistModel)
	return err
}
