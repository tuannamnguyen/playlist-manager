package service

import "github.com/tuannamnguyen/playlist-manager/internal/model"

type PlaylistRepository interface {
	Insert(playlistModel model.Playlist) error
	SelectAll() ([]model.Playlist, error)
	SelectWithID(id string) (model.Playlist, error)
	DeleteByID(id string) error
}

type Playlist struct {
	playlistRepo PlaylistRepository
}

func NewPlaylist(playlistRepo PlaylistRepository) *Playlist {
	return &Playlist{
		playlistRepo: playlistRepo,
	}
}

func (p *Playlist) Add(playlistModel model.Playlist) error {
	return p.playlistRepo.Insert(playlistModel)
}

func (p *Playlist) GetAll() ([]model.Playlist, error) {
	return p.playlistRepo.SelectAll()
}

func (p *Playlist) GetByID(id string) (model.Playlist, error) {
	return p.playlistRepo.SelectWithID(id)
}

func (p *Playlist) DeleteByID(id string) error {
	return p.playlistRepo.DeleteByID(id)
}
