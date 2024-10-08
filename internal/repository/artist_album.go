package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ArtistAlbumRepository struct {
	db *sqlx.DB
}

func NewArtistAlbumRepository(db *sqlx.DB) *ArtistAlbumRepository {
	return &ArtistAlbumRepository{db: db}
}

func (a *ArtistAlbumRepository) Insert(ctx context.Context, artistID int, albumID int) error {
	_, err := a.db.QueryxContext(
		ctx,
		`INSERT INTO artist_album (artist_id, album_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		artistID, albumID,
	)
	if err != nil {
		return &queryError{err}
	}

	return nil
}
