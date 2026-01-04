package repository

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

func parsePlaylistSongData(rows []model.SongOutDB) []model.SongOutAPI {
	songMap := make(map[int]int)
	var result []model.SongOutAPI

	for _, row := range rows {
		if idx, exists := songMap[row.ID]; exists {
			result[idx].ArtistNames = append(result[idx].ArtistNames, row.ArtistName)
		} else {
			var ISRC string
			if row.ISRC.Valid {
				ISRC = row.ISRC.String
			} else {
				ISRC = ""
			}

			songMap[row.ID] = len(result)
			result = append(result, model.SongOutAPI{
				ID:          row.ID,
				Name:        row.Name,
				AlbumName:   row.AlbumName,
				ArtistNames: []string{row.ArtistName},
				ImageURL:    row.ImageURL,
				Duration:    row.Duration,
				ISRC:        ISRC,
				Timestamp:   row.Timestamp,
			})
		}
	}

	return result
}

func (p *PlaylistRepository) mapPlaylistDBToAPI(ctx context.Context, playlistsOutDB []model.PlaylistOutDB) ([]model.Playlist, error) {
	var playlists []model.Playlist
	for _, playlistOutDB := range playlistsOutDB {
		playlistAPIResponse, err := p.mapSinglePlaylistDBToApiResponse(ctx, playlistOutDB)
		if err != nil {
			return nil, err
		}

		playlists = append(playlists, playlistAPIResponse)
	}
	return playlists, nil
}

func (p *PlaylistRepository) mapSinglePlaylistDBToApiResponse(ctx context.Context, playlistOutDB model.PlaylistOutDB) (model.Playlist, error) {
	var playlistDescription string
	if playlistOutDB.PlaylistDescription.Valid {
		playlistDescription = playlistOutDB.PlaylistDescription.String
	} else {
		playlistDescription = ""
	}

	imageURL, err := p.generateSignedURLFromObjectName(ctx, playlistOutDB.ImageName)
	if err != nil {
		return model.Playlist{}, err
	}

	playlistAPIResponse := model.Playlist{
		ID:                  playlistOutDB.ID,
		Name:                playlistOutDB.Name,
		PlaylistDescription: playlistDescription,
		UserID:              playlistOutDB.UserID,
		Username:            playlistOutDB.Username,
		Timestamp:           playlistOutDB.Timestamp,
		ImageURL:            imageURL,
	}
	return playlistAPIResponse, nil
}

func (p *PlaylistRepository) generateSignedURLFromObjectName(ctx context.Context, objectName string) (string, error) {
	bucketName := os.Getenv("S3_BUCKET_NAME")

	request, err := p.s3PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}, func(po *s3.PresignOptions) {
		po.Expires = time.Duration(15 * int64(time.Minute))
	})

	if err != nil {
		return "", err
	}

	return request.URL, nil
}

func transformSearchAPIResponse(searchRes SearchResponse) []model.SongInAPI {
	result := make([]model.SongInAPI, len(searchRes.Tracks))

	for i, track := range searchRes.Tracks {
		result[i] = model.SongInAPI{
			Name:        track.Data.Name,
			ArtistNames: track.Data.ArtistNames,
			AlbumName:   track.Data.AlbumName,
			Duration:    track.Data.Duration,
			ImageURL:    track.Data.ImageURL,
			ISRC:        track.Data.ISRC,
		}
	}

	return result
}
