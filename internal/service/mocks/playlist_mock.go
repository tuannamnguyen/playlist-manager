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

func (m *MockPlaylistRepository) SelectWithID(ctx context.Context, id int) (model.Playlist, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Playlist), args.Error(1)
}

func (m *MockPlaylistRepository) DeleteByID(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
