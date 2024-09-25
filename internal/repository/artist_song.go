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
		return &beginTransactionError{err}
	}
	defer func() {
		err = tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("error rolling back transaction bulk insert artist songs: %v\n", err)
		}
	}()

	query := `INSERT INTO artist_song (song_id, artist_id, artist_insertion_order)
		VALUES %s
		ON CONFLICT DO NOTHING`

	valueStrings := make([]string, 0, len(artistIDs))
	valueArgs := make([]any, 0, len(artistIDs)*3)

	for index, artistID := range artistIDs {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, songID, artistID, index)
	}

	query = sqlx.Rebind(
		sqlx.DOLLAR,
		fmt.Sprintf(query, strings.Join(valueStrings, ",")),
	)

	_, err = tx.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return &execError{err}
	}

	err = tx.Commit()
	if err != nil {
		return &transactionCommitError{err}
	}

	return nil
}
