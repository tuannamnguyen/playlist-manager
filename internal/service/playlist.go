package service

import (
	"bytes"
	"context"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type PlaylistRepository interface {
	Insert(ctx context.Context, playlistModel model.PlaylistInDB) error
	SelectAll(ctx context.Context, userID string) ([]model.Playlist, error)
	SelectWithID(ctx context.Context, id int) (model.Playlist, error)
	DeleteByID(ctx context.Context, id int) error
	AddPlaylistPicture(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
}

type SongRepository interface {
	InsertAndGetID(ctx context.Context, song model.SongInDB) (int, error)
}

type PlaylistSongRepository interface {
	BulkInsert(ctx context.Context, playlistID int, songsID []int) error
	GetAll(ctx context.Context, playlistID int, sortBy string, sortOrder string) ([]model.SongOutAPI, error)
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

func (p *PlaylistService) Add(
	ctx context.Context,
	playlistModel model.PlaylistIn,
	imageFile multipart.File,
	imageHeader *multipart.FileHeader,
) error {
	imageName, err := p.playlistRepo.AddPlaylistPicture(ctx, imageFile, imageHeader)
	if err != nil {
		return err
	}

	playlistInDBModel := model.PlaylistInDB{
		Name:                playlistModel.Name,
		PlaylistDescription: playlistModel.PlaylistDescription,
		UserID:              playlistModel.UserID,
		Username:            playlistModel.Username,
		ImageName:           imageName,
	}

	return p.playlistRepo.Insert(ctx, playlistInDBModel)
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
			Name:     song.Name,
			AlbumID:  albumID,
			Duration: song.Duration,
			ImageURL: song.ImageURL,
			ISRC:     song.ISRC,
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

func (p *PlaylistService) GetAllSongsFromPlaylist(
	ctx context.Context,
	playlistID int,
	sortBy string,
	sortOrder string,
) ([]model.SongOutAPI, error) {
	return p.playlistSongRepo.GetAll(ctx, playlistID, sortBy, sortOrder)
}

func (p *PlaylistService) DeleteSongsFromPlaylist(ctx context.Context, playlistID int, songsID []int) error {
	return p.playlistSongRepo.BulkDelete(ctx, playlistID, songsID)
}

func (p *PlaylistService) Convert(
	ctx context.Context,
	provider string,
	providerMetadata model.ConverterServiceProviderMetadata,
	playlistName string,
	songs []model.SongOutAPI,
) error {
	converter, err := getConverter(ctx, provider, providerMetadata)
	if err != nil {
		return err
	}

	return converter.Export(ctx, playlistName, songs)
}

func (p *PlaylistService) ConvertSongsToCsv(songs []model.SongOutAPI) (bytes.Buffer, error) {
	var buffer bytes.Buffer

	header := []string{"Name", "Artists", "Album", "Song Cover URL", "Duration", "ISRC"}

	err := writeCsvRecord(&buffer, header)
	if err != nil {
		return bytes.Buffer{}, err
	}

	for _, song := range songs {
		record := []string{
			song.Name,
			strings.Join(song.ArtistNames, "|"),
			song.AlbumName,
			song.ImageURL,
			strconv.Itoa(song.Duration),
			song.ISRC,
		}

		err := writeCsvRecord(&buffer, record)
		if err != nil {
			return bytes.Buffer{}, err
		}
	}

	return buffer, nil
}

func (p *PlaylistService) ConvertCsvToSongs(file multipart.File) ([]model.SongInAPI, error) {
	records, err := parseCsv(file)
	if err != nil {
		return nil, err
	}

	songs := make([]model.SongInAPI, len(records)-1)

	for i, record := range records {
		// file header so skip
		if i == 0 {
			continue
		}

		duration, err := strconv.Atoi(record[4])
		if err != nil {
			return nil, err
		}

		song := model.SongInAPI{
			Name:        record[0],
			ArtistNames: strings.Split(record[1], "|"),
			AlbumName:   record[2],
			ImageURL:    record[3],
			Duration:    duration,
			ISRC:        record[5],
		}

		songs[i-1] = song
	}

	return songs, nil
}
