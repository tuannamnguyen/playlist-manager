package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type PlaylistRepository struct {
	db *sqlx.DB
}

func NewPlaylistRepository(db *sqlx.DB) *PlaylistRepository {
	return &PlaylistRepository{db}
}

func (p *PlaylistRepository) Insert(playlistModel *model.Playlist) error {
	playlistModel.UpdatedAt = time.Now()
	playlistModel.CreatedAt = time.Now()

	_, err := p.db.NamedExec(
		`INSERT INTO playlist (playlist_id, playlist_name, user_id, updated_at, created_at)
		VALUES (:playlist_id, :playlist_name, :user_id, :updated_at, :created_at)
		RETURNING playlist_id`,
		playlistModel,
	)

	if err != nil {
		return fmt.Errorf("INSERT playlist into db: %w", err)
	}

	return nil
}

func (p *PlaylistRepository) SelectAll() ([]model.Playlist, error) {
	var playlists []model.Playlist

	rows, err := p.db.Queryx("SELECT * FROM playlist")
	if err != nil {
		return nil, fmt.Errorf("SELECT playlist from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var playlist model.Playlist
		if err := rows.StructScan(&playlist); err != nil {
			return nil, fmt.Errorf("scan playlist to struct: %w", err)

		}
		playlists = append(playlists, playlist)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("playlist query iteration: %w", err)
	}

	return playlists, nil
}
