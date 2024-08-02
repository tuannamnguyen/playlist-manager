package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

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

func (m *MockSongRepository) BulkInsert(ctx context.Context, songs []model.Song) ([]int, error) {
	args := m.Called(ctx, songs)

	if args.Get(0) != nil {
		return args.Get(0).([]int), args.Error(1)
	}

	return nil, args.Error(1)
}
