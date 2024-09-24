package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type AlbumRepository struct {
	db *sqlx.DB
}

func NewAlbumRepository(db *sqlx.DB) *AlbumRepository {
	return &AlbumRepository{db: db}
}

func (a *AlbumRepository) InsertAndGetID(ctx context.Context, albumName string) (int, error) {
	row := a.db.QueryRowxContext(
		ctx,
		`WITH ins AS (
			INSERT INTO album (album_name)
			VALUES ($1)
			ON CONFLICT DO NOTHING
			RETURNING album_id
		)
		 SELECT album_id FROM ins
		 UNION ALL
		 SELECT album_id FROM album
		 WHERE album_name = $2
		 LIMIT 1`,
		albumName, albumName,
	)

	var lastInsertID int
	err := row.Scan(&lastInsertID)
	if err != nil {
		return 0, &rowScanError{err}
	}

	return lastInsertID, nil
}
