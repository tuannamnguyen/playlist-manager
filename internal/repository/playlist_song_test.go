package repository

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

func setupTestDB(t *testing.T, initScriptPath string) (*sqlx.DB, func()) {
	var (
		dbUser     = "postgres"
		dbPassword = "password"
	)
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:latest",
		postgres.WithInitScripts(filepath.Join(".", "testdata", initScriptPath)),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err, "failed to start container")

	connectionString, err := postgresContainer.ConnectionString(ctx, "dbname=playlist_manager")
	require.NoError(t, err, "failed to get connection string")

	db, err := sqlx.Connect("pgx", connectionString)
	require.NoError(t, err, "Unable to connect to database")

	cleanup := func() {
		db.Close()
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	}

	return db, cleanup
}

func TestPlaylistSongRepositoryInsert(t *testing.T) {
	db, cleanup := setupTestDB(t, "script_test_insert_playlist_song.sql")
	defer cleanup()

	tests := []struct {
		name       string
		playlistID string
		songID     string
		wantErr    bool
	}{
		{
			name:       "successful insert",
			playlistID: "asdasdasdsaasd",
			songID:     "asiuasubfasuifaufb",
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
		playlistID string
		want       []model.PlaylistSong
		wantErr    bool
	}{
		{
			name:       "get all success",
			playlistID: "asdasdasdsaasd",
			want: []model.PlaylistSong{
				{
					PlaylistID: "asdasdasdsaasd",
					SongID:     "asiuasubfasuifaufb",
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
		playlistID string
		songsID    []string
		wantErr    bool
	}{
		{
			name:       "delete successfully",
			playlistID: "asdasdasdsaasd",
			songsID:    []string{"asiuasubfasuifaufb"},
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
