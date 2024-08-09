package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type SongRepository struct {
	db *sqlx.DB
}

func NewSongRepository(db *sqlx.DB) *SongRepository {
	return &SongRepository{
		db: db,
	}
}

func (s *SongRepository) InsertAndGetID(ctx context.Context, song model.SongInDB) (int, error) {
	row := s.db.QueryRowxContext(
		ctx,
		`INSERT INTO song (song_name, album_id)
		VALUES ($1, $2)
		RETURNING song_id`,
		song.Name, song.AlbumID,
	)

	var lastInsertID int
	err := row.Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("scanning last inserted song ID: %w", err)
	}

	return lastInsertID, nil
}
