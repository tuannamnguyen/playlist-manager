package repository

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	dbName     = "playlist_manager"
	dbUser     = "postgres"
	dbPassword = "password"
)

func TestPlaylistSongRepositoryInsert(t *testing.T) {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:latest",
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

	containerIP, err := postgresContainer.ContainerIP(ctx)
	if err != nil {
		log.Fatalf("failed to get container IP: %s", err)
	}

	containerPort, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("failed to get container port: %s", err)
	}

	// setup DB
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		containerIP,
		containerPort,
		dbUser,
		dbPassword,
		dbName,
	)
	db, err := sqlx.Connect("pgx", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()
}
