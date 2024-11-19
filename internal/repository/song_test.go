package repository

import (
	"context"
	"testing"

	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

func TestSongRepositoryInsertAndGetID(t *testing.T) {
	db, cleanup := setupTestDB(t, "test_init_script.sql")
	defer cleanup()

	type args struct {
		ctx  context.Context
		song model.SongInDB
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "insert successfully",
			args: args{
				ctx: context.Background(),
				song: model.SongInDB{
					Name:     "All I Want",
					AlbumID:  1,
					ImageURL: "https://example.com",
					Duration: 1234,
					ISRC:     "ABC123",
				},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SongRepository{
				db: db,
			}
			got, err := s.InsertAndGetID(tt.args.ctx, tt.args.song)
			if (err != nil) != tt.wantErr {
				t.Errorf("SongRepository.InsertAndGetID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SongRepository.InsertAndGetID() = %v, want %v", got, tt.want)
			}
		})
	}
}
