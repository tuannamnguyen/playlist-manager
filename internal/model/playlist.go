package model

import "time"

type Playlist struct {
	ID        string    `json:"playlist_id" db:"playlist_id"`
	Name      string    `json:"playlist_name" db:"playlist_name"`
	UserID    string    `json:"user_id" db:"user_id"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
