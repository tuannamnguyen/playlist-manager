package model

type PlaylistSong struct {
	PlaylistID int `json:"playlist_id" db:"playlist_id"`
	SongID     int `json:"song_id" db:"song_id"`
	Timestamp
}
