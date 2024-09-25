package repository

import (
	"context"
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

func (p *PlaylistRepository) Insert(ctx context.Context, playlistModel model.PlaylistIn) error {
	updatedAt := time.Now()
	createdAt := time.Now()

	_, err := p.db.ExecContext(
		ctx,
		`INSERT INTO playlist (playlist_name, user_id, user_name, updated_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING playlist_id`,
		playlistModel.Name, playlistModel.UserID, playlistModel.Username, updatedAt, createdAt,
	)

	if err != nil {
		return &execError{err}
	}

	return nil
}

func (p *PlaylistRepository) SelectAll(ctx context.Context, userID string) ([]model.Playlist, error) {
	var playlists []model.Playlist
	var query string
	var args []interface{}

	if userID != "" {
		query = "SELECT * FROM playlist WHERE user_id = $1"
		args = append(args, userID)
	} else {
		query = "SELECT * FROM playlist"
	}

	err := p.db.SelectContext(ctx, &playlists, query, args...)
	if err != nil {
		return nil, &selectError{err}
	}

	return playlists, nil
}

func (p *PlaylistRepository) SelectWithID(ctx context.Context, id int) (model.Playlist, error) {
	var playlist model.Playlist

	err := p.db.QueryRowxContext(ctx, "SELECT * FROM playlist WHERE playlist_id = $1", id).StructScan(&playlist)
	if err != nil {
		return model.Playlist{}, &structScanError{err}
	}

	return playlist, nil
}

func (p *PlaylistRepository) DeleteByID(ctx context.Context, id int) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM playlist WHERE playlist_id = $1", id)
	if err != nil {
		return &execError{err}
	}

	return nil
}
