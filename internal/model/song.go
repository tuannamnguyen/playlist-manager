package model

type Song struct {
	ID       string `json:"song_id" db:"song_id"`
	Name     string `json:"song_name" db:"song_name"`
	ArtistID string `json:"artist_id" db:"artist_id"`
	AlbumID  string `json:"album_id" db:"album_id"`
	timestamp
}
