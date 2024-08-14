package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

var fakeTimestamp = model.Timestamp{
	UpdatedAt: time.Date(2024, 7, 27, 10, 12, 0, 0, time.UTC),
	CreatedAt: time.Date(2024, 7, 27, 10, 12, 0, 0, time.UTC),
}

func TestParsePlaylistSongData(t *testing.T) {
	type args struct {
		rows []model.SongOutDB
	}
	tests := []struct {
		name    string
		args    args
		want    []model.SongOutAPI
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				rows: []model.SongOutDB{
					{
						ID:         1,
						Name:       "Song 1",
						AlbumName:  "Album 1",
						ArtistName: "Artist 1",
						Timestamp:  fakeTimestamp,
					},
					{
						ID:         4,
						Name:       "devil in a new dress",
						AlbumName:  "mbdtf",
						ArtistName: "kanye west",
						Timestamp:  fakeTimestamp,
					},
					{
						ID:         4,
						Name:       "devil in a new dress",
						AlbumName:  "mbdtf",
						ArtistName: "rick ross",
						Timestamp:  fakeTimestamp,
					},
					{
						ID:         5,
						Name:       "runaway",
						AlbumName:  "mbdtf",
						ArtistName: "pusha t",
						Timestamp:  fakeTimestamp,
					},
					{
						ID:         5,
						Name:       "runaway",
						AlbumName:  "mbdtf",
						ArtistName: "kanye west",
						Timestamp:  fakeTimestamp,
					},
				},
			},
			want: []model.SongOutAPI{
				{
					ID:          1,
					Name:        "Song 1",
					AlbumName:   "Album 1",
					ArtistNames: []string{"Artist 1"},
					Timestamp:   fakeTimestamp,
				},
				{
					ID:          4,
					Name:        "devil in a new dress",
					AlbumName:   "mbdtf",
					ArtistNames: []string{"kanye west", "rick ross"},
					Timestamp:   fakeTimestamp,
				},
				{
					ID:          5,
					Name:        "runaway",
					AlbumName:   "mbdtf",
					ArtistNames: []string{"pusha t", "kanye west"}, // TODO: we'll need ordering here in the future
					Timestamp:   fakeTimestamp,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePlaylistSongData(tt.args.rows)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
