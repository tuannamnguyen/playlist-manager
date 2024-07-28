package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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

		mockSongRepo.On("Insert", songs[0]).Return(nil)
		mockSongRepo.On("Insert", songs[1]).Return(nil)
		mockPlaylistSongRepo.On("Insert", playlistID, "song1").Return(nil)
		mockPlaylistSongRepo.On("Insert", playlistID, "song2").Return(nil)

		err := playlistService.AddSongsToPlaylist(playlistID, songs)
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
		mockPlaylistSongRepo.On("SelectAll", "abcd").Return([]model.Song{
			{ID: "test_id", Name: "test_name", ArtistID: "test_artist_id", AlbumID: "test_album_id"},
		}, nil)

		playlistID := "abcd"

		songs, err := playlistService.GetAllSongsFromPlaylist(playlistID)
		assert.NoError(t, err)
		assert.Len(t, songs, 1)

		mockPlaylistSongRepo.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {
		mockPlaylistSongRepo.On("SelectAll", "defg").Return(nil, errors.New("test error"))

		playlistID := "defg"

		songs, err := playlistService.GetAllSongsFromPlaylist(playlistID)
		assert.EqualError(t, err, "test error")
		assert.Len(t, songs, 0)

		mockPlaylistSongRepo.AssertExpectations(t)
	})
}
