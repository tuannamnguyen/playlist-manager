package repository

import (
	"context"
	"fmt"
	"time"

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

func (s *SongRepository) Insert(ctx context.Context, song model.Song) error {
	song.UpdatedAt = time.Now()
	song.CreatedAt = time.Now()

	_, err := s.db.NamedExecContext(
		ctx,
		`INSERT INTO song (song_id, song_name, artist_id, album_id, updated_at, created_at)
		VALUES (:song_id, :song_name, :artist_id, :album_id, :updated_at, :created_at)`,
		&song,
	)
	if err != nil {
		return fmt.Errorf("INSERT song into db: %w", err)
	}

	return nil
}
