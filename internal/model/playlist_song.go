package model

type PlaylistSong struct {
	PlaylistID string `json:"playlist_id" db:"playlist_id"`
	SongID     string `json:"song_id" db:"song_id"`
	timestamp
}
