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

	t.Run("Success", func(t *testing.T) {
		songs := []model.Song{
			{ID: "song1"},
			{ID: "song2"},
		}
		playlistID := "playlist1"

		mockSongRepo.On("Insert", mock.Anything, songs[0]).Return(nil)
		mockSongRepo.On("Insert", mock.Anything, songs[1]).Return(nil)
		mockPlaylistSongRepo.On("Insert", mock.Anything, playlistID, "song1").Return(nil)
		mockPlaylistSongRepo.On("Insert", mock.Anything, playlistID, "song2").Return(nil)

		err := playlistService.AddSongsToPlaylist(context.Background(), playlistID, songs)
		assert.NoError(t, err)

		mockSongRepo.AssertExpectations(t)
		mockPlaylistSongRepo.AssertExpectations(t)
	})
}

func TestGetAllSongsFromPlaylist(t *testing.T) {
	mockPlaylistRepo := new(mocks.MockPlaylistRepository)
	mockSongRepo := new(mocks.MockSongRepository)
	mockPlaylistSongRepo := new(mocks.MockPlaylistSongRepository)

	playlistService := NewPlaylist(mockPlaylistRepo, mockSongRepo, mockPlaylistSongRepo)

	t.Run("success", func(t *testing.T) {
		mockSong := []model.Song{
			{ID: "test_id", Name: "test_name", ArtistID: "test_artist_id", AlbumID: "test_album_id"},
		}

		mockPlaylistSongRepo.On("SelectAll", mock.Anything, "abcd").Return([]model.PlaylistSong{
			{PlaylistID: "abcd", SongID: "test_id"},
		}, nil).Once()
		mockSongRepo.On("SelectWithManyID", mock.Anything, mock.AnythingOfType("[]string")).Return(mockSong, nil).Once()

		playlistID := "abcd"
		songs, err := playlistService.GetAllSongsFromPlaylist(context.Background(), playlistID)
		assert.NoError(t, err)
		assert.Equal(t, mockSong, songs)

		mockSongRepo.AssertExpectations(t)
		mockPlaylistSongRepo.AssertExpectations(t)
	})
}

func TestDeleteSongsFromPlaylist(t *testing.T) {
	mockPlaylistRepo := new(mocks.MockPlaylistRepository)
	mockSongRepo := new(mocks.MockSongRepository)
	mockPlaylistSongRepo := new(mocks.MockPlaylistSongRepository)

	playlistService := NewPlaylist(mockPlaylistRepo, mockSongRepo, mockPlaylistSongRepo)

	t.Run("delete successfully", func(t *testing.T) {
		playlistID := "abcd"
		songsID := []string{"abc", "def"}

		mockPlaylistSongRepo.On("DeleteWithManyID", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]string")).Return(nil).Once()
		err := playlistService.DeleteSongsFromPlaylist(context.Background(), playlistID, songsID)
		assert.NoError(t, err)

		mockPlaylistSongRepo.AssertExpectations(t)
	})
}
