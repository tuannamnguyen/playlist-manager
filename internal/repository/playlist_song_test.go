package repository

import (
	"context"
	"log"
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

func TestPlaylistSongRepositoryInsert(t *testing.T) {
	var (
		dbUser     = "postgres"
		dbPassword = "password"
	)
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:latest",
		postgres.WithInitScripts(filepath.Join(".", "testdata", "script_test_insert_playlist_song.sql")),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connectionString, err := postgresContainer.ConnectionString(ctx, "dbname=playlist_manager")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	// setup DB

	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	t.Run("test insert song in playlist", func(t *testing.T) {
		playlistSongRepo := NewPlaylistSongRepository(db)
		playlistID := "asdasdasdsaasd"
		songID := "asiuasubfasuifaufb"

		err = playlistSongRepo.Insert(context.Background(), playlistID, songID)
		if assert.NoError(t, err) {
			var insertedPlaylistSong model.PlaylistSong
			err = db.QueryRowx(
				`SELECT playlist_id, song_id
			FROM playlist_song
			WHERE playlist_id = $1
			AND song_id = $2`,
				playlistID,
				songID,
			).StructScan(&insertedPlaylistSong)
			if err != nil {
				log.Fatalf("test: error querying playlist song: %v", err)
			}

			assert.Equal(t, playlistID, insertedPlaylistSong.PlaylistID)
			assert.Equal(t, songID, insertedPlaylistSong.SongID)
		} else {
			t.Errorf("expected no error but got: %s", err)
		}

	})
}

func TestPlaylistSongRepositorySelectAll(t *testing.T) {
	var (
		dbUser     = "postgres"
		dbPassword = "password"
	)
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:latest",
		postgres.WithInitScripts(filepath.Join(".", "testdata", "script_test_get_all_song.sql")),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connectionString, err := postgresContainer.ConnectionString(ctx, "dbname=playlist_manager")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	// setup DB

	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	t.Run("get all success", func(t *testing.T) {
		playlistSongRepo := NewPlaylistSongRepository(db)
		playlistID := "asdasdasdsaasd"

		parsedUpdatedAt, err := time.Parse(time.DateTime, "2024-07-27 10:12:00")
		require.NoError(t, err)

		parsedCreatedAt, err := time.Parse(time.DateTime, "2024-07-27 10:12:00")
		require.NoError(t, err)

		expectedSongs := []model.PlaylistSong{
			{PlaylistID: "asdasdasdsaasd", SongID: "asiuasubfasuifaufb", Timestamp: model.Timestamp{
				UpdatedAt: parsedUpdatedAt,
				CreatedAt: parsedCreatedAt,
			}},
		}

		playlistSongs, err := playlistSongRepo.SelectAll(context.Background(), playlistID)
		require.NoError(t, err)
		assert.Equal(t, expectedSongs, playlistSongs)
	})
}

func TestDeleteSongsFromPlaylist(t *testing.T) {
	var (
		dbUser     = "postgres"
		dbPassword = "password"
	)
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:latest",
		postgres.WithInitScripts(filepath.Join(".", "testdata", "script_test_delete_playlist_song.sql")),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connectionString, err := postgresContainer.ConnectionString(ctx, "dbname=playlist_manager")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	// setup DB

	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	t.Run("delete successfully", func(t *testing.T) {
		playlistSongRepository := NewPlaylistSongRepository(db)
		playlistID := "asdasdasdsaasd"
		songsID := []string{"asiuasubfasuifaufb"}

		err := playlistSongRepository.DeleteWithManyID(context.Background(), playlistID, songsID)
		require.NoError(t, err)

		songs, err := playlistSongRepository.SelectAll(context.Background(), playlistID)
		assert.NoError(t, err)

		assert.Nil(t, songs)
	})

}
