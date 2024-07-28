package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type MockPlaylistRepository struct {
	mock.Mock
}

func (m *MockPlaylistRepository) Insert(playlistModel model.Playlist) error {
	args := m.Called(playlistModel)
	return args.Error(0)
}

func (m *MockPlaylistRepository) SelectAll() ([]model.Playlist, error) {
	args := m.Called()
	return args.Get(0).([]model.Playlist), args.Error(1)
}

func (m *MockPlaylistRepository) SelectWithID(id string) (model.Playlist, error) {
	args := m.Called(id)
	return args.Get(0).(model.Playlist), args.Error(1)
}

func (m *MockPlaylistRepository) DeleteByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockSongRepository struct {
	mock.Mock
}

func (m *MockSongRepository) Insert(song model.Song) error {
	args := m.Called(song)
	return args.Error(0)
}

type MockPlaylistSongRepository struct {
	mock.Mock
}

func (m *MockPlaylistSongRepository) Insert(playlistID string, songID string) error {
	args := m.Called(playlistID, songID)
	return args.Error(0)
}

func (m *MockPlaylistSongRepository) SelectAll(playlistID string) ([]model.Song, error) {
	args := m.Called(playlistID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.Song), args.Error(1)
}
