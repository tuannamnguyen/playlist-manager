package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ArtistSongRepository struct {
	db *sqlx.DB
}

func NewArtistSongRepository(db *sqlx.DB) *ArtistSongRepository {
	return &ArtistSongRepository{db: db}
}

func (as *ArtistSongRepository) Insert(ctx context.Context, songID int, artistIDs []int) error {
	tx, err := as.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction insert arist song: %w", err)
	}
	defer func() {
		err = tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("error rolling back transaction bulk insert artist songs: %v\n", err)
		}
	}()

	query := `INSERT INTO artist_song (song_id, artist_id)
		VALUES %s
		ON CONFLICT DO NOTHING`

	valueStrings := make([]string, 0, len(artistIDs))
	valueArgs := make([]any, 0, len(artistIDs)*2)

	for _, artistID := range artistIDs {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, songID, artistID)
	}

	query = sqlx.Rebind(
		sqlx.DOLLAR,
		fmt.Sprintf(query, strings.Join(valueStrings, ",")),
	)

	_, err = as.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("bulk INSERT artist song: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commiting transaction artist song: %w", err)
	}

	return nil
}
