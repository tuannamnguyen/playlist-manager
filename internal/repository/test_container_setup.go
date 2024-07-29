package repository

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
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