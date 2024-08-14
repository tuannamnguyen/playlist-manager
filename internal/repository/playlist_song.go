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
		return fmt.Errorf("begin transaction bulk insert playlist song: %w", err)
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
		return fmt.Errorf("bulk INSERT playlist song: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commiting transaction playlist song: %w", err)
	}

	return nil
}

func (ps *PlaylistSongRepository) GetAll(ctx context.Context, playlistID int) ([]model.SongOutAPI, error) {
	query := `SELECT pls.song_id, s.song_name, al.album_name, ar.artist_name, pls.created_at, pls.updated_at
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
				WHERE pl.playlist_id = $1`

	var rows []model.SongOutDB
	err := ps.db.SelectContext(ctx, &rows, query, playlistID)
	if err != nil {
		return nil, fmt.Errorf("SELECT all songs in playlist: %w", err)
	}

	return parsePlaylistSongData(rows)
}
