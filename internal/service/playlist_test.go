package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
	"github.com/tuannamnguyen/playlist-manager/internal/service/mocks"
)

func TestAddSongsToPlaylist(t *testing.T) {
	mockPlaylistRepo := new(mocks.MockPlaylistRepository)
	mockSongRepo := new(mocks.MockSongRepository)
	mockPlaylistSongRepo := new(mocks.MockPlaylistSongRepository)

	playlist := NewPlaylist(mockPlaylistRepo, mockSongRepo, mockPlaylistSongRepo)

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

		err := playlist.AddSongsToPlaylist(playlistID, songs)
		assert.NoError(t, err)

		mockSongRepo.AssertExpectations(t)
		mockPlaylistSongRepo.AssertExpectations(t)
	})
}
