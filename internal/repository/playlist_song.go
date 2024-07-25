package repository

import (
	"fmt"
	"time"

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

func (ps *PlaylistSongRepository) Insert(playlistID string, songID string) error {
	playlistSong := model.PlaylistSong{
		PlaylistID: playlistID,
		SongID:     songID,
		Timestamp: model.Timestamp{
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
	}

	_, err := ps.db.NamedExec(
		`INSERT INTO playlist_song (playlist_id, song_id, updated_at ,created_at)
		VALUES (:playlist_id, :song_id, :updated_at, :created_at)`,
		&playlistSong,
	)
	if err != nil {
		return fmt.Errorf("INSERT playlist song into db: %w", err)
	}

	return nil
}
