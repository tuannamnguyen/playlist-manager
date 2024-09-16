package service

import (
	"context"

	"github.com/tuannamnguyen/playlist-manager/internal/model"
	"golang.org/x/oauth2"
)

type PlaylistRepository interface {
	Insert(ctx context.Context, playlistModel model.PlaylistIn) error
	SelectAll(ctx context.Context, userID string) ([]model.Playlist, error)
	SelectWithID(ctx context.Context, id int) (model.Playlist, error)
	DeleteByID(ctx context.Context, id int) error
}

type SongRepository interface {
	InsertAndGetID(ctx context.Context, song model.SongInDB) (int, error)
}

type PlaylistSongRepository interface {
	BulkInsert(ctx context.Context, playlistID int, songsID []int) error
	GetAll(ctx context.Context, playlistID int) ([]model.SongOutAPI, error)
	BulkDelete(ctx context.Context, playlistID int, songsID []int) error
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

type ArtistAlbumRepository interface {
	Insert(ctx context.Context, artistID int, albumID int) error
}

type Converter interface {
	Export(ctx context.Context, playlistName string, songs []model.SongOutAPI) error
}

type PlaylistService struct {
	playlistRepo     PlaylistRepository
	songRepo         SongRepository
	playlistSongRepo PlaylistSongRepository
	albumRepo        AlbumRepository
	artistRepo       ArtistRepository
	artistSongRepo   ArtistSongRepository
	artistAlbumRepo  ArtistAlbumRepository
}

func NewPlaylist(
	playlistRepo PlaylistRepository,
	songRepo SongRepository,
	playlistSongRepo PlaylistSongRepository,
	albumRepo AlbumRepository,
	artistRepo ArtistRepository,
	artistSongRepo ArtistSongRepository,
	artistAlbumRepo ArtistAlbumRepository,
) *PlaylistService {
	return &PlaylistService{
		playlistRepo:     playlistRepo,
		songRepo:         songRepo,
		playlistSongRepo: playlistSongRepo,
		albumRepo:        albumRepo,
		artistRepo:       artistRepo,
		artistSongRepo:   artistSongRepo,
		artistAlbumRepo:  artistAlbumRepo,
	}
}

func (p *PlaylistService) Add(ctx context.Context, playlistModel model.PlaylistIn) error {
	return p.playlistRepo.Insert(ctx, playlistModel)
}

func (p *PlaylistService) GetAll(ctx context.Context, userID string) ([]model.Playlist, error) {
	return p.playlistRepo.SelectAll(ctx, userID)
}

func (p *PlaylistService) GetByID(ctx context.Context, id int) (model.Playlist, error) {
	return p.playlistRepo.SelectWithID(ctx, id)
}

func (p *PlaylistService) DeleteByID(ctx context.Context, id int) error {
	return p.playlistRepo.DeleteByID(ctx, id)
}

func (p *PlaylistService) AddSongsToPlaylist(ctx context.Context, playlistID int, songs []model.SongInAPI) error {
	var songsID []int
	for _, song := range songs {
		albumID, err := p.albumRepo.InsertAndGetID(ctx, song.AlbumName)
		if err != nil {
			return err
		}

		artistIDs, err := p.artistRepo.BulkInsertAndGetIDs(ctx, song.ArtistNames)
		if err != nil {
			return err
		}

		err = p.artistAlbumRepo.Insert(ctx, artistIDs[0], albumID)
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

		err = p.artistSongRepo.Insert(ctx, songID, artistIDs)
		if err != nil {
			return err
		}

		songsID = append(songsID, songID)

	}

	err := p.playlistSongRepo.BulkInsert(ctx, playlistID, songsID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PlaylistService) GetAllSongsFromPlaylist(ctx context.Context, playlistID int) ([]model.SongOutAPI, error) {
	return p.playlistSongRepo.GetAll(ctx, playlistID)
}

func (p *PlaylistService) DeleteSongsFromPlaylist(ctx context.Context, playlistID int, songsID []int) error {
	return p.playlistSongRepo.BulkDelete(ctx, playlistID, songsID)
}

func (p *PlaylistService) Convert(ctx context.Context, provider string, token *oauth2.Token, songs []model.SongOutAPI) error {
	converter, err := getConverter(ctx, provider, token)
	if err != nil {
		return err
	}

	return converter.Export(ctx, "test playlist", songs)
}
