package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

func TestPlaylistSongRepositoryInsert(t *testing.T) {
	db, cleanup := setupTestDB(t, "script_test_insert_playlist_song.sql")
	defer cleanup()

	tests := []struct {
		name       string
		playlistID int
		songID     int
		wantErr    bool
	}{
		{
			name:       "successful insert",
			playlistID: 1,
			songID:     1,
			wantErr:    false,
		},
		// Add more test cases here if needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			playlistSongRepo := NewPlaylistSongRepository(db)

			err := playlistSongRepo.Insert(context.Background(), tt.playlistID, tt.songID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				var insertedPlaylistSong model.PlaylistSong
				err = db.QueryRowx(
					`SELECT playlist_id, song_id
                    FROM playlist_song
                    WHERE playlist_id = $1
                    AND song_id = $2`,
					tt.playlistID,
					tt.songID,
				).StructScan(&insertedPlaylistSong)

				require.NoError(t, err)
				assert.Equal(t, tt.playlistID, insertedPlaylistSong.PlaylistID)
				assert.Equal(t, tt.songID, insertedPlaylistSong.SongID)
			}
		})
	}
}

func TestPlaylistSongRepositorySelectAll(t *testing.T) {
	db, cleanup := setupTestDB(t, "script_test_get_all_song.sql")
	defer cleanup()

	tests := []struct {
		name       string
		playlistID int
		want       []model.PlaylistSong
		wantErr    bool
	}{
		{
			name:       "get all success",
			playlistID: 1,
			want: []model.PlaylistSong{
				{
					PlaylistID: 1,
					SongID:     1,
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
			playlistSongRepo := NewPlaylistSongRepository(db)

			got, err := playlistSongRepo.SelectAll(context.Background(), tt.playlistID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDeleteSongsFromPlaylist(t *testing.T) {
	db, cleanup := setupTestDB(t, "script_test_delete_playlist_song.sql")
	defer cleanup()

	tests := []struct {
		name       string
		playlistID int
		songsID    []int
		wantErr    bool
	}{
		{
			name:       "delete successfully",
			playlistID: 1,
			songsID:    []int{1},
			wantErr:    false,
		},
		// Add more test cases here if needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			playlistSongRepository := NewPlaylistSongRepository(db)

			err := playlistSongRepository.DeleteWithManyID(context.Background(), tt.playlistID, tt.songsID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				songs, err := playlistSongRepository.SelectAll(context.Background(), tt.playlistID)
				assert.NoError(t, err)
				assert.Empty(t, songs)
			}
		})
	}
}

func TestSelectAllSongsInPlaylist(t *testing.T) {
	db, cleanup := setupTestDB(t, "script_test_get_all_song.sql")
	defer cleanup()

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx        context.Context
		playlistID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Song
		wantErr bool
	}{
		{
			name: "get all success",
			fields: fields{
				db: db,
			},
			args: args{
				ctx:        context.Background(),
				playlistID: 1,
			},
			want: []model.Song{
				{
					ID:       1,
					Name:     "devil in a new dress",
					ArtistID: "kanye west",
					AlbumID:  "mbdtf",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PlaylistSongRepository{
				db: tt.fields.db,
			}
			got, err := ps.SelectAllSongsInPlaylist(tt.args.ctx, tt.args.playlistID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PlaylistSongRepository.SelectAllSongsInPlaylist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlaylistSongRepository.SelectAllSongsInPlaylist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlaylistSongBulkInsert(t *testing.T) {
	db, cleanup := setupTestDB(t, "script_test_insert_playlist_song.sql")
	defer cleanup()

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx        context.Context
		playlistID int
		songsID    []int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "insert success",
			fields: fields{
				db: db,
			},
			args: args{
				ctx:        context.Background(),
				playlistID: 1,
				songsID:    []int{1, 2, 3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PlaylistSongRepository{
				db: tt.fields.db,
			}
			if err := ps.BulkInsert(tt.args.ctx, tt.args.playlistID, tt.args.songsID); (err != nil) != tt.wantErr {
				t.Errorf("PlaylistSongRepository.BulkInsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
