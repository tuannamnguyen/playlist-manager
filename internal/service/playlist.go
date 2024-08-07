package service

import (
	"context"

	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type PlaylistRepository interface {
	Insert(ctx context.Context, playlistModel model.PlaylistIn) error
	SelectAll(ctx context.Context) ([]model.Playlist, error)
	SelectWithID(ctx context.Context, id int) (model.Playlist, error)
	DeleteByID(ctx context.Context, id int) error
}

type SongRepository interface {
	Insert(ctx context.Context, song model.Song) (int, error)
	BulkInsert(ctx context.Context, songs []model.SongIn) ([]int, error)
	SelectWithManyID(ctx context.Context, ID []int) ([]model.Song, error)
	GetIDsFromSongsDetail(ctx context.Context, songs []model.SongIn) ([]int, error)
}

type PlaylistSongRepository interface {
	Insert(ctx context.Context, playlistID int, songID int) error
	BulkInsert(ctx context.Context, playlistID int, songsID []int) error
	SelectAll(ctx context.Context, playlistID int) ([]model.PlaylistSong, error)
	DeleteWithManyID(ctx context.Context, playlistID int, songsID []int) error
	SelectAllSongsInPlaylist(ctx context.Context, playlistID int) ([]model.Song, error)
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

func (p *PlaylistService) Add(ctx context.Context, playlistModel model.PlaylistIn) error {
	return p.playlistRepo.Insert(ctx, playlistModel)
}

func (p *PlaylistService) GetAll(ctx context.Context) ([]model.Playlist, error) {
	return p.playlistRepo.SelectAll(ctx)
}

func (p *PlaylistService) GetByID(ctx context.Context, id int) (model.Playlist, error) {
	return p.playlistRepo.SelectWithID(ctx, id)
}

func (p *PlaylistService) DeleteByID(ctx context.Context, id int) error {
	return p.playlistRepo.DeleteByID(ctx, id)
}

func (p *PlaylistService) AddSongsToPlaylist(ctx context.Context, playlistID int, songs []model.SongIn) error {
	// TODO: REDO THIS

	return nil
}

func (p *PlaylistService) GetAllSongsFromPlaylist(ctx context.Context, playlistID int) ([]model.Song, error) {
	return p.playlistSongRepo.SelectAllSongsInPlaylist(ctx, playlistID)
}

func (p *PlaylistService) DeleteSongsFromPlaylist(ctx context.Context, playlistID int, songsID []int) error {
	return p.playlistSongRepo.DeleteWithManyID(ctx, playlistID, songsID)
}
