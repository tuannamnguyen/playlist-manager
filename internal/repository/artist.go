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

type ArtistRepository struct {
	db *sqlx.DB
}

func NewArtistRepository(db *sqlx.DB) *ArtistRepository {
	return &ArtistRepository{db}
}

func (a *ArtistRepository) BulkInsertAndGetIDs(ctx context.Context, artistNames []string) ([]int, error) {
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction insert bulk artists: %w", err)
	}
	defer func() {
		err = tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("error rolling back transaction bulk insert artists: %v\n", err)
		}
	}()

	query := `INSERT INTO artist (artist_name)
			VALUES %s
			RETURNING artist_id`

	valueStrings := make([]string, 0, len(artistNames))
	valueArgs := make([]any, 0, len(artistNames))

	for _, name := range artistNames {
		valueStrings = append(valueStrings, "(?)")
		valueArgs = append(valueArgs, name)
	}

	query = sqlx.Rebind(
		sqlx.DOLLAR,
		fmt.Sprintf(query, strings.Join(valueStrings, ",")),
	)

	rows, err := tx.QueryxContext(ctx, query, valueArgs...)
	if err != nil {
		return nil, fmt.Errorf("bulk INSERT artists: %w", err)
	}
	defer rows.Close()

	var insertedIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scanning inserted artist IDs: %w", err)
		}

		insertedIDs = append(insertedIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("artist query iteration: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("commiting transaction bulk insert artists: %w", err)
	}

	return insertedIDs, nil
}
