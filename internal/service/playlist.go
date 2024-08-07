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
	InsertAndGetID(ctx context.Context, song model.SongInDB) (int, error)
}

type PlaylistSongRepository interface {
	DeleteWithManyID(ctx context.Context, playlistID int, songsID []int) error
	SelectAllSongsInPlaylist(ctx context.Context, playlistID int) ([]model.SongOutAPI, error)
}

type AlbumRepository interface {
	InsertAndGetID(ctx context.Context, albumName string) (int, error)
}

type ArtistRepository interface {
	BulkInsertAndGetIDs(ctx context.Context, artistNames []string) ([]int, error)
}

type ArtistSongRepository interface {
	Insert(ctx context.Context, songID int, artistIDs []int) error
}

type PlaylistService struct {
	playlistRepo     PlaylistRepository
	songRepo         SongRepository
	playlistSongRepo PlaylistSongRepository
	albumRepo        AlbumRepository
	artistRepo       ArtistRepository
	artistSongRepo ArtistSongRepository
}

func NewPlaylist(playlistRepo PlaylistRepository, songRepo SongRepository, playlistSongRepo PlaylistSongRepository, albumRepo AlbumRepository) *PlaylistService {
	return &PlaylistService{
		playlistRepo:     playlistRepo,
		songRepo:         songRepo,
		playlistSongRepo: playlistSongRepo,
		albumRepo:        albumRepo,
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

func (p *PlaylistService) AddSongsToPlaylist(ctx context.Context, playlistID int, songs []model.SongInAPI) error {
	// TODO: THIS IS JUST A ROUGH DRAFT

	// TODO: need better error handling
	for _, song := range songs {
		albumID, err := p.albumRepo.InsertAndGetID(ctx, song.AlbumName)
		if err != nil {
			return err
		}

		songID, err := p.songRepo.InsertAndGetID(ctx, model.SongInDB{
			Name:    song.Name,
			AlbumID: albumID,
		})
		if err != nil {
			return err
		}

		artistIDs, err := p.artistRepo.BulkInsertAndGetIDs(ctx, song.ArtistNames)
		if err != nil {
			return err
		}

		err = p.artistSongRepo.Insert(ctx, songID, artistIDs)
		if err != nil {
			return err
		}
	}

	return nil
}
