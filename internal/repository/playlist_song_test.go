package repository

import (
	"context"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

var (
	dbUser     = "postgres"
	dbPassword = "password"
)

func TestPlaylistSongRepositoryInsert(t *testing.T) {
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
	log.Println(connectionString)

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

		err = playlistSongRepo.Insert(playlistID, songID)
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
