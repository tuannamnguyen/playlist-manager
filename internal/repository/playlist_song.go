package repository

import (
	"context"
	"fmt"
	"log"
	"time"

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

func (ps *PlaylistSongRepository) Insert(ctx context.Context, playlistID string, songID string) error {
	playlistSong := model.PlaylistSong{
		PlaylistID: playlistID,
		SongID:     songID,
		Timestamp: model.Timestamp{
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
	}

	_, err := ps.db.NamedExecContext(
		ctx,
		`INSERT INTO playlist_song (playlist_id, song_id, updated_at ,created_at)
		VALUES (:playlist_id, :song_id, :updated_at, :created_at)`,
		&playlistSong,
	)
	if err != nil {
		return fmt.Errorf("INSERT playlist song into db: %w", err)
	}

	return nil
}

func (ps *PlaylistSongRepository) SelectAll(ctx context.Context, playlistID string) ([]model.PlaylistSong, error) {
	var playlistSongs []model.PlaylistSong

	rows, err := ps.db.QueryxContext(
		ctx,
		"SELECT * FROM playlist_song WHERE playlist_id = $1",
		playlistID,
	)
	if err != nil {
		return nil, fmt.Errorf("SELECT all songs from playlist_song table: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var song model.PlaylistSong
		if err := rows.StructScan(&song); err != nil {
			return nil, fmt.Errorf("scan song to struct: %w", err)
		}

		playlistSongs = append(playlistSongs, song)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("playlist song query iteration: %w", err)
	}

	return playlistSongs, nil

}

func (ps *PlaylistSongRepository) DeleteWithManyID(ctx context.Context, playlistID string, songsID []string) error {
	query, args, err := sqlx.In("DELETE FROM playlist_song WHERE playlist_id = (?) AND song_id IN (?)", playlistID, songsID)
	if err != nil {
		return fmt.Errorf("prepare delete songs in playlist query: %w", err)
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	log.Println(query, args)

	_, err = ps.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("DELETE songs from playlist_song table: %w", err)
	}

	return nil

}
