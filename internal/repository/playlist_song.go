package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PlaylistSongRepository struct {
	db *sqlx.DB
}

func NewPlaylistSongRepository(db *sqlx.DB) *PlaylistSongRepository {
	return &PlaylistSongRepository{
		db: db,
	}
}

func (ps *PlaylistSongRepository) BulkInsert(ctx context.Context, playlistID int, songsID []int) error {
	tx, err := ps.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction insert songs into playlist: %w", err)
	}
	defer func() {
		err := tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("error rolling back transaction bulk insert songs in playlist: %v\n", err)
		}
	}()

	createdAt := time.Now()
	updatedAt := time.Now()

	query := `
		INSERT INTO playlist_song (playlist_id, song_id, updated_at, created_at)
		VALUES
		%s
	`
	valueStrings := make([]string, 0, len(songsID))
	valueArgs := make([]any, 0, len(songsID)*4)

	for _, songID := range songsID {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, playlistID, songID, createdAt, updatedAt)
	}
	query = sqlx.Rebind(
		sqlx.DOLLAR,
		fmt.Sprintf(query, strings.Join(valueStrings, ",")),
	)

	_, err = tx.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("bulk INSERT songs in playlist: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction insert songs into playlist: %w", err)
	}

	return nil
}
