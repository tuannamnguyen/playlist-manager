package repository

import (
	"slices"
	"sort"

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

	// Sort the result slice by ID before returning
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result
}

func mapPlaylistDBToAPI(playlistsOutDB []model.PlaylistOutDB) []model.Playlist {
	var playlists []model.Playlist
	for _, playlistOutDB := range playlistsOutDB {
		playlistAPIResponse := mapSinglePlaylistDBToApiResponse(playlistOutDB)

		playlists = append(playlists, playlistAPIResponse)
	}
	return playlists
}

func mapSinglePlaylistDBToApiResponse(playlistOutDB model.PlaylistOutDB) model.Playlist {
	var playlistDescription string
	if playlistOutDB.PlaylistDescription.Valid {
		playlistDescription = playlistOutDB.PlaylistDescription.String
	} else {
		playlistDescription = ""
	}

	playlistAPIResponse := model.Playlist{
		ID:                  playlistOutDB.ID,
		Name:                playlistOutDB.Name,
		PlaylistDescription: playlistDescription,
		UserID:              playlistOutDB.UserID,
		Username:            playlistOutDB.Username,
		Timestamp:           playlistOutDB.Timestamp,
	}
	return playlistAPIResponse
}
