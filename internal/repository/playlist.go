package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type PlaylistRepository struct {
	db          *sqlx.DB
	minioClient *minio.Client
}

func NewPlaylistRepository(db *sqlx.DB, minioClient *minio.Client) *PlaylistRepository {
	return &PlaylistRepository{db, minioClient}
}

func (p *PlaylistRepository) Insert(ctx context.Context, playlistModel model.PlaylistInDB) error {
	updatedAt := time.Now()
	createdAt := time.Now()

	_, err := p.db.ExecContext(
		ctx,
		`INSERT INTO playlist (playlist_name, user_id, user_name, playlist_description, updated_at, created_at, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING playlist_id`,
		playlistModel.Name,
		playlistModel.UserID,
		playlistModel.Username,
		playlistModel.PlaylistDescription,
		updatedAt,
		createdAt,
		playlistModel.ImageURL,
	)

	if err != nil {
		return &execError{err}
	}

	return nil
}

func (p *PlaylistRepository) SelectAll(ctx context.Context, userID string) ([]model.Playlist, error) {
	var playlistsOutDB []model.PlaylistOutDB
	var query string
	var args []interface{}

	if userID != "" {
		query = "SELECT * FROM playlist WHERE user_id = $1"
		args = append(args, userID)
	} else {
		query = "SELECT * FROM playlist"
	}

	err := p.db.SelectContext(ctx, &playlistsOutDB, query, args...)
	if err != nil {
		return nil, &selectError{err}
	}

	playlists := mapPlaylistDBToAPI(playlistsOutDB)

	return playlists, nil
}

func (p *PlaylistRepository) SelectWithID(ctx context.Context, id int) (model.Playlist, error) {
	var playlist model.PlaylistOutDB

	err := p.db.QueryRowxContext(ctx, "SELECT * FROM playlist WHERE playlist_id = $1", id).StructScan(&playlist)
	if err != nil {
		return model.Playlist{}, &structScanError{err}
	}

	return mapSinglePlaylistDBToApiResponse(playlist), nil
}

func (p *PlaylistRepository) DeleteByID(ctx context.Context, id int) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM playlist WHERE playlist_id = $1", id)
	if err != nil {
		return &execError{err}
	}

	return nil
}

func (p *PlaylistRepository) AddPlaylistPicture(ctx context.Context, playlistID string, file multipart.File, header *multipart.FileHeader) error {
	bucketName := "playlist-cover"
	contentType := header.Header.Get("Content-Type")

	timestamp := time.Now().Format(time.RFC3339)
	uuid := uuid.New().String()
	filename := fmt.Sprintf("%s/%s_%s_%s", playlistID, timestamp, uuid, header.Filename)

	_, err := p.minioClient.PutObject(
		ctx,
		bucketName,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		return &putObjectError{err}
	}

	return nil
}
