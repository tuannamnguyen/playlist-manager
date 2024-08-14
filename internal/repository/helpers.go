package repository

import (
	"slices"

	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

func parsePlaylistSongData(rows []model.SongOutDB) ([]model.SongOutAPI, error) {
	songMap := make(map[int]*model.SongOutAPI)

	for _, row := range rows {
		if song, exists := songMap[row.ID]; exists {
			// Song already exists, just add the artist if it's not already there
			if !slices.Contains(song.ArtistNames, row.ArtistName) {
				song.ArtistNames = append(song.ArtistNames, row.ArtistName)
			}
		} else {
			// Create a new song entry
			songMap[row.ID] = &model.SongOutAPI{
				ID:          row.ID,
				Name:        row.Name,
				ArtistNames: []string{row.ArtistName},
				AlbumName:   row.AlbumName,
				Timestamp:   row.Timestamp,
			}
		}
	}

	// Convert map to slice
	result := make([]model.SongOutAPI, 0, len(songMap))
	for _, song := range songMap {
		result = append(result, *song)
	}

	return result, nil
}
