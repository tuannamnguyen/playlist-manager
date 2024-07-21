package repository

import (
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
	_, err := p.db.NamedExec(
		`INSERT INTO playlist (playlist_id, playlist_name, user_id)
		VALUES (:ID, :Name, :UserID)
		RETURNING id
		`,
		playlistModel,
	)

	return err
}
