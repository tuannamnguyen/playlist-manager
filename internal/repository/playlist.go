package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type PlaylistRepository struct {
	db              *sqlx.DB
	s3Client        *s3.Client
	s3PresignClient *s3.PresignClient
}

func NewPlaylistRepository(db *sqlx.DB, s3Client *s3.Client, s3PresignedClient *s3.PresignClient) *PlaylistRepository {
	return &PlaylistRepository{db, s3Client, s3PresignedClient}
}

func (p *PlaylistRepository) Insert(ctx context.Context, playlistModel model.PlaylistInDB) error {
	updatedAt := time.Now()
	createdAt := time.Now()

	_, err := p.db.ExecContext(
		ctx,
		`INSERT INTO playlist (playlist_name, user_id, user_name, playlist_description, updated_at, created_at, image_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING playlist_id`,
		playlistModel.Name,
		playlistModel.UserID,
		playlistModel.Username,
		playlistModel.PlaylistDescription,
		updatedAt,
		createdAt,
		playlistModel.ImageName,
	)

	if err != nil {
		return &execError{err}
	}

	return nil
}

func (p *PlaylistRepository) SelectAll(ctx context.Context, userID string) ([]model.Playlist, error) {
	var playlistsOutDB []model.PlaylistOutDB
	var query string
	var args []any

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

	return p.mapPlaylistDBToAPI(ctx, playlistsOutDB)
}

func (p *PlaylistRepository) SelectWithID(ctx context.Context, id int) (model.Playlist, error) {
	var playlist model.PlaylistOutDB

	err := p.db.QueryRowxContext(ctx, "SELECT * FROM playlist WHERE playlist_id = $1", id).StructScan(&playlist)
	if err != nil {
		return model.Playlist{}, &structScanError{err}
	}

	return p.mapSinglePlaylistDBToApiResponse(ctx, playlist)
}

func (p *PlaylistRepository) DeleteByID(ctx context.Context, id int) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM playlist WHERE playlist_id = $1", id)
	if err != nil {
		return &execError{err}
	}

	return nil
}

func (p *PlaylistRepository) AddPlaylistPicture(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	// TODO: update this to use S3

	bucketName := os.Getenv("S3_BUCKET_NAME")

	timestamp := time.Now().Format(time.RFC3339)
	uuid := uuid.New().String()
	objectName := fmt.Sprintf("playlist_cover/%s_%s_%s", timestamp, uuid, header.Filename)

	_, err := p.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
		Body:   file,
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			log.Printf("Error while uploading object to %s. The object is too large.\n"+
				"To upload objects larger than 5GB, use the S3 console (160GB max)\n"+
				"or the multipart upload API (5TB max).", bucketName)
		} else {
			log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
				header.Filename, bucketName, objectName, err)
		}
		return "", err
	}

	err = s3.NewObjectExistsWaiter(p.s3Client).Wait(
		ctx, &s3.HeadObjectInput{Bucket: aws.String(bucketName), Key: aws.String(objectName)}, time.Minute)
	if err != nil {
		log.Printf("Failed attempt to wait for object %s to exist.\n", objectName)
	}

	return objectName, nil
}
