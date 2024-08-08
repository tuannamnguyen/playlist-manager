package repository

import (
	"context"
	"fmt"

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
		`INSERT INTO album (album_name)
		VALUES ($1)
		RETURNING album_id`,
		albumName,
	)

	var lastInsertID int
	err := row.Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("scanning last inserted album ID: %w", err)
	}

	return lastInsertID, nil
}
