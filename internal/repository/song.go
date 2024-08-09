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
		`WITH ins AS (
			INSERT INTO song (song_name, album_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
			RETURNING song_id
		)
			SELECT song_id FROM ins
			UNION ALL
			SELECT song_id FROM song
			WHERE (song_name, album_id) IN (($3, $4))
			LIMIT 1`,
		song.Name, song.AlbumID, song.Name, song.AlbumID,
	)

	var lastInsertID int
	err := row.Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("scanning last inserted song ID: %w", err)
	}

	return lastInsertID, nil
}
