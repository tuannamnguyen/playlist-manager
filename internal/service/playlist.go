package service

import (
	"log"

	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type PlaylistRepository interface {
	Insert(playlistModel model.Playlist) error
	SelectAll() ([]model.Playlist, error)
	SelectWithID(id string) (model.Playlist, error)
	DeleteByID(id string) error
}

type SongRepository interface {
	Insert(song model.Song) error
}

type PlaylistSongRepository interface {
	Insert(playlistID string, songID string) error
}

type Playlist struct {
	playlistRepo     PlaylistRepository
	songRepo         SongRepository
	playlistSongRepo PlaylistSongRepository
}

func NewPlaylist(playlistRepo PlaylistRepository, songRepo SongRepository, playlistSongRepo PlaylistSongRepository) *Playlist {
	return &Playlist{
		playlistRepo:     playlistRepo,
		songRepo:         songRepo,
		playlistSongRepo: playlistSongRepo,
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

func (p *Playlist) AddSongsToPlaylist(playlistID string, songs []model.Song) error {
	for _, song := range songs {
		err := p.songRepo.Insert(song)
		if err != nil {
			return err
		}

		log.Println("inserted song in song table")

		err = p.playlistSongRepo.Insert(playlistID, song.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
