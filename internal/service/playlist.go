package service

import (
	"context"
	"log"

	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type PlaylistRepository interface {
	Insert(ctx context.Context, playlistModel model.Playlist) error
	SelectAll(ctx context.Context) ([]model.Playlist, error)
	SelectWithID(ctx context.Context, id string) (model.Playlist, error)
	DeleteByID(ctx context.Context, id string) error
}

type SongRepository interface {
	Insert(ctx context.Context, song model.Song) error
	SelectWithManyID(ctx context.Context, ID []string) ([]model.Song, error)
}

type PlaylistSongRepository interface {
	Insert(ctx context.Context, playlistID string, songID string) error
	SelectAll(ctx context.Context, playlistID string) ([]model.PlaylistSong, error)
	DeleteWithManyID(ctx context.Context, playlistID string, songsID []string) error
}

type PlaylistService struct {
	playlistRepo     PlaylistRepository
	songRepo         SongRepository
	playlistSongRepo PlaylistSongRepository
}

func NewPlaylist(playlistRepo PlaylistRepository, songRepo SongRepository, playlistSongRepo PlaylistSongRepository) *PlaylistService {
	return &PlaylistService{
		playlistRepo:     playlistRepo,
		songRepo:         songRepo,
		playlistSongRepo: playlistSongRepo,
	}
}

func (p *PlaylistService) Add(ctx context.Context, playlistModel model.Playlist) error {
	return p.playlistRepo.Insert(ctx, playlistModel)
}

func (p *PlaylistService) GetAll(ctx context.Context) ([]model.Playlist, error) {
	return p.playlistRepo.SelectAll(ctx)
}

func (p *PlaylistService) GetByID(ctx context.Context, id string) (model.Playlist, error) {
	return p.playlistRepo.SelectWithID(ctx, id)
}

func (p *PlaylistService) DeleteByID(ctx context.Context, id string) error {
	return p.playlistRepo.DeleteByID(ctx, id)
}

func (p *PlaylistService) AddSongsToPlaylist(ctx context.Context, playlistID string, songs []model.Song) error {
	for _, song := range songs {
		err := p.songRepo.Insert(ctx, song)
		if err != nil {
			return err
		}

		log.Println("inserted song in song table")

		err = p.playlistSongRepo.Insert(ctx, playlistID, song.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PlaylistService) GetAllSongsFromPlaylist(ctx context.Context, playlistID string) ([]model.Song, error) {
	playlistSongs, err := p.playlistSongRepo.SelectAll(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	var ID []string
	for _, song := range playlistSongs {
		ID = append(ID, song.SongID)
	}

	songsDetail, err := p.songRepo.SelectWithManyID(ctx, ID)

	return songsDetail, err
}

func (p *PlaylistService) DeleteSongsFromPlaylist(ctx context.Context, playlistID string, songsID []string) error {
	return p.playlistSongRepo.DeleteWithManyID(ctx, playlistID, songsID)
}
