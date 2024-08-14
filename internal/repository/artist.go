package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
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

	inQuery, inQueryArgs, err := sqlx.In(
		`SELECT artist_id FROM artist WHERE artist_name IN (?)`,
		artistNames,
	)
	if err != nil {
		return nil, fmt.Errorf("prepare in query to select artist ID: %w", err)
	}

	query := `WITH ins AS (
			INSERT INTO artist (artist_name)
			VALUES %s
			ON CONFLICT DO NOTHING
			RETURNING artist_id
		)
			SELECT artist_id FROM ins
			UNION ALL
			%s
			LIMIT %s`

	valueStrings := make([]string, 0, len(artistNames))
	valueArgs := make([]any, 0, len(artistNames))

	for _, name := range artistNames {
		valueStrings = append(valueStrings, "(?)")
		valueArgs = append(valueArgs, name)
	}

	query = sqlx.Rebind(
		sqlx.DOLLAR,
		fmt.Sprintf(query, strings.Join(valueStrings, ","), inQuery, strconv.Itoa(len(artistNames))),
	)
	args := append(valueArgs, inQueryArgs...)

	var insertedIDs []int
	err = tx.SelectContext(ctx, &insertedIDs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("bulk INSERT artists: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("commiting transaction bulk insert artists: %w", err)
	}

	return insertedIDs, nil
}
