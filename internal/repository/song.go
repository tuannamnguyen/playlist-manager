package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
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
	// start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction insert songs: %w", err)
	}
	defer func() {
		err = tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("error rolling back transaction bulk insert songs: %v\n", err)
		}
	}()

	// prepare query
	createdAt := time.Now()
	updatedAt := time.Now()

	query := `
		INSERT INTO song (song_name, artist_id, album_id, updated_at, created_at)
		VALUES %s
		RETURNING song_id
	`
	valueStrings := make([]string, 0, len(songs))
	valueArgs := make([]any, 0, len(songs)*5)
	for _, song := range songs {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, song.Name, song.ArtistID, song.AlbumID, createdAt, updatedAt)
	}
	query = sqlx.Rebind(
		sqlx.DOLLAR,
		fmt.Sprintf(query, strings.Join(valueStrings, ",")),
	)

	// Execute the query
	rows, err := tx.QueryContext(ctx, query, valueArgs...)
	if err != nil {
		return nil, fmt.Errorf("bulk INSERT songs: %w", err)
	}
	defer rows.Close()

	// get inserted IDs
	var insertedIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scanning inserted songs id: %w", err)
		}
		insertedIDs = append(insertedIDs, id)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("commiting transaction: %w", err)
	}

	return insertedIDs, nil

}

func (s *SongRepository) GetIDsFromSongsDetail(ctx context.Context, songs []model.Song) ([]int, error) {
	query := `SELECT song_id
		FROM song
		WHERE (song_name, artist_id, album_id) IN (%s)`
	valueStrings := make([]string, 0, len(songs))
	valueArgs := make([]any, 0, len(songs)*3)
	for _, song := range songs {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, song.Name, song.ArtistID, song.AlbumID)
	}

	query = sqlx.Rebind(
		sqlx.DOLLAR,
		fmt.Sprintf(query, strings.Join(valueStrings, ",")),
	)

	rows, err := s.db.QueryxContext(ctx, query, valueArgs...)
	if err != nil {
		return nil, fmt.Errorf("SELECT song id from song detail: %w", err)
	}
	defer rows.Close()

	var songIDs []int
	for rows.Next() {
		var songID int
		if err := rows.Scan(&songID); err != nil {
			return nil, fmt.Errorf("scanning song ID: %w", err)
		}

		songIDs = append(songIDs, songID)
	}

	return songIDs, nil
}
