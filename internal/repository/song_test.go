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

func TestSelectWithManyID(t *testing.T) {
	ctx := context.Background()
	var (
		dbUser     = "postgres"
		dbPassword = "password"
	)

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

	t.Run("success get all songs", func(t *testing.T) {
		songRepository := NewSongRepository(db)
		songsDetail, err := songRepository.SelectWithManyID(context.Background(), []string{"asiuasubfasuifaufb"})
		require.NoError(t, err)

		parsedUpdatedAt, err := time.Parse(time.DateTime, "2024-07-27 10:12:00")
		require.NoError(t, err)

		parsedCreatedAt, err := time.Parse(time.DateTime, "2024-07-27 10:12:00")
		require.NoError(t, err)

		expectedSongs := []model.Song{
			{
				ID:       "asiuasubfasuifaufb",
				Name:     "devil in a new dress",
				ArtistID: "kanye west",
				AlbumID:  "mbdtf",
				Timestamp: model.Timestamp{
					UpdatedAt: parsedCreatedAt,
					CreatedAt: parsedUpdatedAt,
				},
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedSongs, songsDetail)
	})
}
