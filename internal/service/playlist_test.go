package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
	"github.com/tuannamnguyen/playlist-manager/internal/service/mocks"
)

func TestAddSongsToPlaylist(t *testing.T) {
	mockPlaylistRepo := new(mocks.MockPlaylistRepository)
	mockSongRepo := new(mocks.MockSongRepository)
	mockPlaylistSongRepo := new(mocks.MockPlaylistSongRepository)

	playlistService := NewPlaylist(mockPlaylistRepo, mockSongRepo, mockPlaylistSongRepo)

	tests := []struct {
		name       string
		playlistID int
		songs      []model.Song
		wantErr    bool
	}{
		{
			name:       "Success",
			playlistID: 1,
			songs: []model.Song{
				{ID: 1},
				{ID: 2},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, song := range tt.songs {
				mockSongRepo.On("Insert", mock.Anything, song).Return(song.ID, nil)
				mockPlaylistSongRepo.On("Insert", mock.Anything, tt.playlistID, song.ID).Return(nil)
			}

			err := playlistService.AddSongsToPlaylist(context.Background(), tt.playlistID, tt.songs)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockSongRepo.AssertExpectations(t)
			mockPlaylistSongRepo.AssertExpectations(t)
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
				{ID: 1, Name: "test_name", ArtistID: "test_artist_id", AlbumID: "test_album_id"},
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
