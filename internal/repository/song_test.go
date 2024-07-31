package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

func TestSelectWithManyID(t *testing.T) {
	db, cleanup := setupTestDB(t, "script_test_get_all_song.sql")
	defer cleanup()

	tests := []struct {
		name    string
		songIDs []int
		want    []model.Song
		wantErr bool
	}{
		{
			name:    "success get all songs",
			songIDs: []int{1},
			want: []model.Song{
				{
					ID:       1,
					Name:     "devil in a new dress",
					ArtistID: "kanye west",
					AlbumID:  "mbdtf",
					Timestamp: model.Timestamp{
						UpdatedAt: time.Date(2024, 7, 27, 10, 12, 0, 0, time.UTC),
						CreatedAt: time.Date(2024, 7, 27, 10, 12, 0, 0, time.UTC),
					},
				},
			},
			wantErr: false,
		},
		// Add more test cases here if needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			songRepository := NewSongRepository(db)
			got, err := songRepository.SelectWithManyID(context.Background(), tt.songIDs)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
