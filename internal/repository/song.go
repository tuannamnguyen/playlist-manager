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

func (s *SongRepository) Insert(ctx context.Context, song model.Song) (int, error) {
	song.UpdatedAt = time.Now()
	song.CreatedAt = time.Now()

	var lastInsertID int
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO song (song_name, artist_id, album_id, updated_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING song_id`,
		song.Name, song.ArtistID, song.AlbumID, song.UpdatedAt, song.CreatedAt,
	).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("INSERT song into db: %w", err)
	}

	return lastInsertID, nil
}

func (s *SongRepository) SelectWithManyID(ctx context.Context, ID []int) ([]model.Song, error) {
	var songs []model.Song
	query, args, err := sqlx.In("SELECT * FROM song WHERE song_id IN (?);", ID)
	if err != nil {
		return nil, fmt.Errorf("prepare select songs with many ID query: %w", err)
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("SELECT all songs detail from song table: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var songDetail model.Song
		if err := rows.StructScan(&songDetail); err != nil {
			return nil, fmt.Errorf("scan song detail to struct: %w", err)
		}

		songs = append(songs, songDetail)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("songs detail query iteration: %w", err)
	}

	return songs, nil
}

func (s *SongRepository) BulkInsert(ctx context.Context, songs []model.Song) ([]int, error) {
	// TODO: IMPLEMENT THIS
	return nil, nil
}
