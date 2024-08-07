package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
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

func TestSongInsert(t *testing.T) {
	db, cleanup := setupTestDB(t, "script_test_insert_song.sql")
	defer cleanup()

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		song model.Song
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        int
		wantErrCode string
	}{
		{
			name: "insert song success",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				song: model.Song{
					Name:     "devil in a new dress",
					ArtistID: "kanye west",
					AlbumID:  "mbdtf",
				},
			},
			want: 1,
		},
		// Run the tests in this exact order to ensure duplicate
		{
			name: "insert song duplicate",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				song: model.Song{
					Name:     "devil in a new dress",
					ArtistID: "kanye west",
					AlbumID:  "mbdtf",
				},
			},
			want:        0,
			wantErrCode: pgerrcode.UniqueViolation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SongRepository{
				db: tt.fields.db,
			}
			got, err := s.Insert(tt.args.ctx, tt.args.song)
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code != tt.wantErrCode {
					t.Errorf("SongRepository.Insert() error = %v, got code: %v, want code: %v", err, pgErr.Code, tt.wantErrCode)
				}
			}

			if got != tt.want {
				t.Errorf("SongRepository.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSongBulkInsert(t *testing.T) {
	db, cleanup := setupTestDB(t, "script_test_insert_song.sql")
	defer cleanup()

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		songs []model.SongIn
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "bulk insert success",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				songs: []model.SongIn{
					{
						Name:        "devil in a new dress",
						ArtistNames: []string{"kanye west"},
						AlbumName:   "mbdtf",
					},
					{
						Name:        "runaway",
						ArtistNames: []string{"kanye west"},
						AlbumName:   "mbdtf",
					},
				},
			},
			want:    []int{1, 2},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SongRepository{
				db: tt.fields.db,
			}
			got, err := s.BulkInsert(tt.args.ctx, tt.args.songs)
			if (err != nil) != tt.wantErr {
				t.Errorf("SongRepository.BulkInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SongRepository.BulkInsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSongGetIDsFromSongsDetail(t *testing.T) {
	db, cleanup := setupTestDB(t, "script_test_get_id_of_song.sql")
	defer cleanup()

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		songs []model.SongIn
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "get IDs successfully",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				songs: []model.SongIn{
					{
						Name:        "devil in a new dress",
						ArtistNames: []string{"kanye west"},
						AlbumName:   "mbdtf",
					},
					{
						Name:        "runaway",
						ArtistNames: []string{"kanye west"},
						AlbumName:   "mbdtf",
					},
				},
			},
			want: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SongRepository{
				db: tt.fields.db,
			}
			got, err := s.GetIDsFromSongsDetail(tt.args.ctx, tt.args.songs)
			if (err != nil) != tt.wantErr {
				t.Errorf("SongRepository.GetIDsFromSongsDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SongRepository.GetIDsFromSongsDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}
