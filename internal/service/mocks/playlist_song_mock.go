package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type MockPlaylistSongRepository struct {
	mock.Mock
}

func (m *MockPlaylistSongRepository) Insert(ctx context.Context, playlistID int, songID int) error {
	args := m.Called(ctx, playlistID, songID)
	return args.Error(0)
}

func (m *MockPlaylistSongRepository) SelectAll(ctx context.Context, playlistID int) ([]model.PlaylistSong, error) {
	args := m.Called(ctx, playlistID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.PlaylistSong), args.Error(1)
}

func (m *MockPlaylistSongRepository) DeleteWithManyID(ctx context.Context, playlistID int, songsID []int) error {
	args := m.Called(ctx, playlistID, songsID)

	return args.Error(0)
}

func (m *MockPlaylistSongRepository) SelectAllSongsInPlaylist(ctx context.Context, playlistID int) ([]model.Song, error) {
	args := m.Called(ctx, playlistID)

	songs, ok := (args.Get(0)).([]model.Song)

	if !ok {
		return nil, args.Error(1)
	}

	return songs, args.Error(1)
}

func (m *MockPlaylistSongRepository) BulkInsert(ctx context.Context, playlistID int, songsID []int) error {
	args := m.Called(ctx, playlistID, songsID)

	return args.Error(0)
}
