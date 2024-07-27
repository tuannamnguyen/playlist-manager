package model

type Playlist struct {
	ID     string `json:"playlist_id" db:"playlist_id"`
	Name   string `json:"playlist_name" db:"playlist_name"`
	UserID string `json:"user_id" db:"user_id"`
	Timestamp
}
