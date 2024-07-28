package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type MockPlaylistRepository struct {
	mock.Mock
}

func (m *MockPlaylistRepository) Insert(ctx context.Context, playlistModel model.Playlist) error {
	args := m.Called(ctx, playlistModel)
	return args.Error(0)
}

func (m *MockPlaylistRepository) SelectAll(ctx context.Context) ([]model.Playlist, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Playlist), args.Error(1)
}

func (m *MockPlaylistRepository) SelectWithID(ctx context.Context, id string) (model.Playlist, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Playlist), args.Error(1)
}

func (m *MockPlaylistRepository) DeleteByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockSongRepository struct {
	mock.Mock
}

func (m *MockSongRepository) Insert(ctx context.Context, song model.Song) error {
	args := m.Called(ctx, song)
	return args.Error(0)
}

type MockPlaylistSongRepository struct {
	mock.Mock
}

func (m *MockPlaylistSongRepository) Insert(ctx context.Context, playlistID string, songID string) error {
	args := m.Called(ctx, playlistID, songID)
	return args.Error(0)
}

func (m *MockPlaylistSongRepository) SelectAll(ctx context.Context, playlistID string) ([]model.Song, error) {
	args := m.Called(ctx, playlistID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.Song), args.Error(1)
}
