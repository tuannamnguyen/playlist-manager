package repository

import (
	"encoding/base64"
	"fmt"
	"os"
	"slices"
	"sort"
	"time"

	"cloud.google.com/go/storage"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

func parsePlaylistSongData(rows []model.SongOutDB) []model.SongOutAPI {
	songMap := make(map[int]*model.SongOutAPI)

	for _, row := range rows {
		if song, exists := songMap[row.ID]; exists {
			// Song already exists, just add the artist if it's not already there
			if !slices.Contains(song.ArtistNames, row.ArtistName) {
				song.ArtistNames = append(song.ArtistNames, row.ArtistName)
			}
		} else {
			var ISRC string
			if row.ISRC.Valid {
				ISRC = row.ISRC.String
			} else {
				ISRC = ""
			}

			// Create a new song entry
			songMap[row.ID] = &model.SongOutAPI{
				ID:          row.ID,
				Name:        row.Name,
				ArtistNames: []string{row.ArtistName},
				AlbumName:   row.AlbumName,
				ImageURL:    row.ImageURL,
				Duration:    row.Duration,
				Timestamp:   row.Timestamp,
				ISRC:        ISRC,
			}
		}
	}

	// Convert map to slice
	result := make([]model.SongOutAPI, 0, len(songMap))
	for _, song := range songMap {
		result = append(result, *song)
	}

	// Sort the result slice by ID before returning
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result
}

func (p *PlaylistRepository) mapPlaylistDBToAPI(playlistsOutDB []model.PlaylistOutDB) ([]model.Playlist, error) {
	var playlists []model.Playlist
	for _, playlistOutDB := range playlistsOutDB {
		playlistAPIResponse, err := p.mapSinglePlaylistDBToApiResponse(playlistOutDB)
		if err != nil {
			return nil, err
		}

		playlists = append(playlists, playlistAPIResponse)
	}
	return playlists, nil
}

func (p *PlaylistRepository) mapSinglePlaylistDBToApiResponse(playlistOutDB model.PlaylistOutDB) (model.Playlist, error) {
	var playlistDescription string
	if playlistOutDB.PlaylistDescription.Valid {
		playlistDescription = playlistOutDB.PlaylistDescription.String
	} else {
		playlistDescription = ""
	}

	imageURL, err := p.generateSignedURLFromObjectName(playlistOutDB.ImageName)
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

func (p *PlaylistRepository) generateSignedURLFromObjectName(objectName string) (string, error) {
	bucketName := os.Getenv("GCS_BUCKET_NAME")

	privateKeyB64Encoded := os.Getenv("GCP_PRIVATE_ACCOUNT_PRIVATE_KEY_B64_ENCODED")
	privateKey, err := base64.RawStdEncoding.DecodeString(privateKeyB64Encoded)
	if err != nil {
		return "", fmt.Errorf("decoding base 64 private key: %w", err)
	}

	opts := &storage.SignedURLOptions{
		GoogleAccessID: os.Getenv("GCP_SERVICE_ACCOUNT"),
		PrivateKey:     privateKey,
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		Expires:        time.Now().Add(15 * time.Minute),
	}

	url, err := p.gcsClient.Bucket(bucketName).SignedURL(objectName, opts)
	if err != nil {
		return "", &gcsGetSignedURLError{err}
	}

	return url, nil
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
