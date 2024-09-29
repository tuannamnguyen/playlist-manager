package repository

import (
	"context"

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
			INSERT INTO song (song_name, album_id, image_url, duration)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT DO NOTHING
			RETURNING song_id
		)
			SELECT song_id FROM ins
			UNION ALL
			SELECT song_id FROM song
			WHERE (song_name, album_id) IN (($5, $6))
			LIMIT 1`,
		song.Name,
		song.AlbumID,
		song.ImageURL,
		song.Duration,
		song.Name,
		song.AlbumID,
	)

	var lastInsertID int
	err := row.Scan(&lastInsertID)
	if err != nil {
		return 0, &rowScanError{err}
	}

	return lastInsertID, nil
}
