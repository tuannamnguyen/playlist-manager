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

type MockSongRepository struct {
	mock.Mock
}

func (m *MockSongRepository) Insert(ctx context.Context, song model.Song) (int, error) {
	args := m.Called(ctx, song)
	return args.Int(0), args.Error(1)
}

func (m *MockSongRepository) SelectWithManyID(ctx context.Context, ID []int) ([]model.Song, error) {
	args := m.Called(ctx, ID)

	songs, ok := (args.Get(0)).([]model.Song)
	if ok {
		return songs, args.Error(1)
	}

	return nil, args.Error(1)
}

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
