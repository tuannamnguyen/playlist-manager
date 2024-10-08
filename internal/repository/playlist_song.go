package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type PlaylistSongRepository struct {
	db *sqlx.DB
}

func NewPlaylistSongRepository(db *sqlx.DB) *PlaylistSongRepository {
	return &PlaylistSongRepository{
		db: db,
	}
}

func (ps *PlaylistSongRepository) BulkInsert(ctx context.Context, playlistID int, songsID []int) error {
	tx, err := ps.db.BeginTxx(ctx, nil)
	if err != nil {
		return &beginTransactionError{err}
	}
	defer func() {
		err = tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("error rolling back transaction bulk insert playlist song: %v\n", err)
		}
	}()

	query := `INSERT INTO playlist_song (playlist_id, song_id)
			VALUES %s ON CONFLICT DO NOTHING`

	valueStrings := make([]string, 0, len(songsID))
	valueArgs := make([]any, 0, len(songsID)*2)
	for _, songID := range songsID {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, playlistID, songID)
	}

	query = sqlx.Rebind(
		sqlx.DOLLAR,
		fmt.Sprintf(query, strings.Join(valueStrings, ",")),
	)

	_, err = tx.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return &execError{err}
	}

	err = tx.Commit()
	if err != nil {
		return &transactionCommitError{err}
	}

	return nil
}

func (ps *PlaylistSongRepository) GetAll(ctx context.Context, playlistID int, sortBy string, sortOrder string) ([]model.SongOutAPI, error) {
	var query string

	if sortBy == "" && sortOrder == "" {
		query = `SELECT pls.song_id, s.song_name, s.image_url, s.duration, s.isrc, al.album_name, ar.artist_name, pls.created_at, pls.updated_at
				FROM playlist_song AS pls
				JOIN playlist AS pl
				ON pl.playlist_id = pls.playlist_id
				JOIN song AS s
				ON pls.song_id = s.song_id
				JOIN album AS al
				ON al.album_id = s.album_id
				JOIN artist_song AS ars
				ON s.song_id = ars.song_id
				JOIN artist AS ar
				ON ars.artist_id = ar.artist_id
				WHERE pl.playlist_id = $1
				ORDER BY song_id, ar.artist_id`
	} else {
		query = fmt.Sprintf(`SELECT pls.song_id, s.song_name, s.image_url, s.duration, s.isrc, al.album_name, ar.artist_name, pls.created_at, pls.updated_at
				FROM playlist_song AS pls
				JOIN playlist AS pl
				ON pl.playlist_id = pls.playlist_id
				JOIN song AS s
				ON pls.song_id = s.song_id
				JOIN album AS al
				ON al.album_id = s.album_id
				JOIN artist_song AS ars
				ON s.song_id = ars.song_id
				JOIN artist AS ar
				ON ars.artist_id = ar.artist_id
				WHERE pl.playlist_id = $1
				ORDER BY %s %s`, sortBy, sortOrder)
	}

	var rows []model.SongOutDB
	err := ps.db.SelectContext(ctx, &rows, query, playlistID)
	if err != nil {
		return nil, &selectError{err}
	}

	return parsePlaylistSongData(rows), nil
}

func (ps *PlaylistSongRepository) BulkDelete(ctx context.Context, playlistID int, songsID []int) error {
	query, args, err := sqlx.In("DELETE FROM playlist_song WHERE playlist_id = (?) AND song_id IN (?)", playlistID, songsID)
	if err != nil {
		return fmt.Errorf("prepare delete songs in playlist query: %w", err)
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	_, err = ps.db.ExecContext(ctx, query, args...)
	if err != nil {
		return &execError{err}
	}

	return nil
}
