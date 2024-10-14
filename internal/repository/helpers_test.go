package repository

import (
	"database/sql"
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
		name string
		args args
		want []model.SongOutAPI
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
					ArtistNames: []string{"pusha t", "kanye west"},
					Timestamp:   fakeTimestamp,
				},
			},
		},
		{
			name: "success",
			args: args{
				rows: []model.SongOutDB{
					{
						ID:         2,
						Name:       "Devil In A New Dress",
						AlbumName:  "My Beautiful Dark Twisted Fantasy",
						ArtistName: "Kanye West",
						ImageURL:   "",
						Duration:   0,
						ISRC:       sql.NullString{},
						Timestamp:  fakeTimestamp,
					},
					{
						ID:         2,
						Name:       "Devil In A New Dress",
						AlbumName:  "My Beautiful Dark Twisted Fantasy",
						ArtistName: "Rick Ross",
						ImageURL:   "",
						Duration:   0,
						ISRC:       sql.NullString{},
						Timestamp:  fakeTimestamp,
					},
					{
						ID:         1,
						Name:       "Location Unknown ◐",
						AlbumName:  "Love Me / Love Me Not",
						ArtistName: "Honne",
						ImageURL:   "",
						Duration:   0,
						ISRC:       sql.NullString{},
						Timestamp:  fakeTimestamp,
					},
					{
						ID:         1,
						Name:       "Location Unknown ◐",
						AlbumName:  "Love Me / Love Me Not",
						ArtistName: "Georgia",
						ImageURL:   "",
						Duration:   0,
						ISRC:       sql.NullString{},
						Timestamp:  fakeTimestamp,
					},
				},
			},
			want: []model.SongOutAPI{
				{
					ID:          2,
					Name:        "Devil In A New Dress",
					AlbumName:   "My Beautiful Dark Twisted Fantasy",
					ArtistNames: []string{"Kanye West", "Rick Ross"},
					ImageURL:    "",
					Duration:    0,
					ISRC:        "",
					Timestamp:   fakeTimestamp,
				},
				{
					ID:          1,
					Name:        "Location Unknown ◐",
					AlbumName:   "Love Me / Love Me Not",
					ArtistNames: []string{"Honne", "Georgia"},
					ImageURL:    "",
					Duration:    0,
					ISRC:        "",
					Timestamp:   fakeTimestamp,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePlaylistSongData(tt.args.rows)
			assert.Equal(t, tt.want, got)
		})
	}
}
