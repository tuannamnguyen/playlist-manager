package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

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
	tx, err := ps.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction bulk insert playlist song: %w", err)
	}
	defer func() {
		err = tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("error rolling back transaction bulk insert playlist song: %v\n", err)
		}
	}()

	query := `INSERT INTO playlist_song (playlist_id, song_id)
			VALUES %s ON CONFLICT DO NOTHING`

	valueStrings := make([]string, 0, len(songsID))
	valueArgs := make([]any, 0, len(songsID)*2)
	for _, songID := range songsID {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, playlistID, songID)
	}

	query = sqlx.Rebind(
		sqlx.DOLLAR,
		fmt.Sprintf(query, strings.Join(valueStrings, ",")),
	)

	_, err = tx.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("bulk INSERT playlist song: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commiting transaction playlist song: %w", err)
	}

	return nil
}
