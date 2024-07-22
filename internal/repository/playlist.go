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

func (p *PlaylistRepository) Add(playlistModel *model.Playlist) error {
	playlistModel.UpdatedAt = time.Now()
	playlistModel.CreatedAt = time.Now()

	_, err := p.db.NamedExec(
		`INSERT INTO playlist (playlist_id, playlist_name, user_id, updated_at, created_at)
		VALUES (:playlist_id, :playlist_name, :user_id, :updated_at, :created_at)
		RETURNING id`,
		playlistModel,
	)

	return fmt.Errorf("add playlist to db: %w", err)
}
