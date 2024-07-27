package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type MockPlaylistService struct {
	mock.Mock
}

func (m *MockPlaylistService) Add(playlistModel model.Playlist) error {
	args := m.Called(playlistModel)
	return args.Error(0)
}

func (m *MockPlaylistService) GetAll() ([]model.Playlist, error) {
	args := m.Called()
	return args.Get(0).([]model.Playlist), args.Error(1)
}

func (m *MockPlaylistService) GetByID(id string) (model.Playlist, error) {
	args := m.Called(id)
	return args.Get(0).(model.Playlist), args.Error(1)
}

func (m *MockPlaylistService) DeleteByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPlaylistService) AddSongsToPlaylist(playlistID string, songs []model.Song) error {
	args := m.Called(playlistID, songs)
	return args.Error(0)
}

func (m *MockPlaylistService) GetAllSongsFromPlaylist(playlistID string) ([]model.Song, error) {
	args := m.Called(playlistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.Song), args.Error(1)
}
