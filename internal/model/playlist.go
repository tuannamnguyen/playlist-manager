package model

type Playlist struct {
	ID       int    `json:"playlist_id" db:"playlist_id"`
	Name     string `json:"playlist_name" db:"playlist_name"`
	UserID   string `json:"user_id" db:"user_id"`
	Username string `json:"user_name" db:"user_name"`
	Timestamp
}

type PlaylistIn struct {
	Name     string `json:"playlist_name" db:"playlist_name" validate:"required"`
	UserID   string `json:"user_id" db:"user_id" validate:"required"`
	Username string `json:"user_name" db:"user_name" validate:"required"`
}
