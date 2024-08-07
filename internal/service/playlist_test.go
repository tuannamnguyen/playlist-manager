package service

import (
	"context"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
	"github.com/tuannamnguyen/playlist-manager/internal/service/mocks"
)

func TestAddSongsToPlaylist(t *testing.T) {
	mockPlaylistRepo := new(mocks.MockPlaylistRepository)

	mockSongRepo := new(mocks.MockSongRepository)
	mockSongRepo.On("BulkInsert", mock.Anything, []model.SongIn{
		{
			Name:        "devil in a new dress",
			ArtistNames: []string{"kanye west"},
			AlbumName:   "mbdtf",
		},
		{
			Name:        "runaway",
			ArtistNames: []string{"kanye west"},
			AlbumName:   "mbdtf",
		},
	}).Return([]int{1, 2}, nil)
	mockSongRepo.On("BulkInsert", mock.Anything, []model.SongIn{
		{
			Name:        "devil in a new dress",
			ArtistNames: []string{"kanye west"},
			AlbumName:   "mbdtf",
		},
		{
			Name:        "runaway",
			ArtistNames: []string{"kanye west"},
			AlbumName:   "mbdtf",
		}, {
			Name:        "devil in a new dress",
			ArtistNames: []string{"kanye west"},
			AlbumName:   "mbdtf",
		},
		{
			Name:        "runaway",
			ArtistNames: []string{"kanye west"},
			AlbumName:   "mbdtf",
		},
	}).Return(
		nil, &pgconn.PgError{Code: pgerrcode.UniqueViolation, Message: "duplicated values"},
	)
	mockSongRepo.On("GetIDsFromSongsDetail", mock.Anything, []model.SongIn{{
		Name:        "devil in a new dress",
		ArtistNames: []string{"kanye west"},
		AlbumName:   "mbdtf",
	},
		{
			Name:        "runaway",
			ArtistNames: []string{"kanye west"},
			AlbumName:   "mbdtf",
		}, {
			Name:        "devil in a new dress",
			ArtistNames: []string{"kanye west"},
			AlbumName:   "mbdtf",
		},
		{
			Name:        "runaway",
			ArtistNames: []string{"kanye west"},
			AlbumName:   "mbdtf",
		}}).Return([]int{1, 2}, nil)

	mockPlaylistSongRepo := new(mocks.MockPlaylistSongRepository)
	mockPlaylistSongRepo.On("BulkInsert", mock.Anything, 1, []int{1, 2}).Return(nil)

	playlistService := NewPlaylist(mockPlaylistRepo, mockSongRepo, mockPlaylistSongRepo)

	tests := []struct {
		name       string
		playlistID int
		songs      []model.SongIn
		wantErr    bool
	}{
		{
			name:       "Success",
			playlistID: 1,
			songs: []model.SongIn{
				{
					Name:        "devil in a new dress",
					ArtistNames: []string{"kanye west"},
					AlbumName:   "mbdtf",
				},
				{
					Name:        "runaway",
					ArtistNames: []string{"kanye west"},
					AlbumName:   "mbdtf",
				},
			},
			wantErr: false,
		},
		{
			name:       "Duplicated success",
			playlistID: 1,
			songs: []model.SongIn{
				{
					Name:        "devil in a new dress",
					ArtistNames: []string{"kanye west"},
					AlbumName:   "mbdtf",
				},
				{
					Name:        "runaway",
					ArtistNames: []string{"kanye west"},
					AlbumName:   "mbdtf",
				},
				{
					Name:        "devil in a new dress",
					ArtistNames: []string{"kanye west"},
					AlbumName:   "mbdtf",
				},
				{
					Name:        "runaway",
					ArtistNames: []string{"kanye west"},
					AlbumName:   "mbdtf",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := playlistService.AddSongsToPlaylist(context.Background(), tt.playlistID, tt.songs)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAllSongsFromPlaylist(t *testing.T) {
	mockPlaylistRepo := new(mocks.MockPlaylistRepository)
	mockSongRepo := new(mocks.MockSongRepository)
	mockPlaylistSongRepo := new(mocks.MockPlaylistSongRepository)

	playlistService := NewPlaylist(mockPlaylistRepo, mockSongRepo, mockPlaylistSongRepo)

	tests := []struct {
		name       string
		playlistID int
		wantSongs  []model.Song
		wantErr    bool
	}{
		{
			name:       "Success",
			playlistID: 1,
			wantSongs: []model.Song{
				{ID: 1, Name: "test_name", ArtistNames: []string{"test_artist_id"}, AlbumName: "test_album_id"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPlaylistSongRepo.On("SelectAllSongsInPlaylist", mock.Anything, 1).Return(tt.wantSongs, nil)

			songs, err := playlistService.GetAllSongsFromPlaylist(context.Background(), tt.playlistID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantSongs, songs)
			}
		})
	}
}

func TestDeleteSongsFromPlaylist(t *testing.T) {
	mockPlaylistRepo := new(mocks.MockPlaylistRepository)
	mockSongRepo := new(mocks.MockSongRepository)
	mockPlaylistSongRepo := new(mocks.MockPlaylistSongRepository)

	playlistService := NewPlaylist(mockPlaylistRepo, mockSongRepo, mockPlaylistSongRepo)

	tests := []struct {
		name       string
		playlistID int
		songsID    []int
		wantErr    bool
	}{
		{
			name:       "Delete Successfully",
			playlistID: 1,
			songsID:    []int{1, 2},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPlaylistSongRepo.On("DeleteWithManyID", mock.Anything, tt.playlistID, tt.songsID).Return(nil).Once()

			err := playlistService.DeleteSongsFromPlaylist(context.Background(), tt.playlistID, tt.songsID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockPlaylistSongRepo.AssertExpectations(t)
		})
	}
}
