package model

import "time"

type Playlist struct {
	ID        string    `json:"playlist_id"`
	Name      string    `json:"playlist_name"`
	UserID    string    `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
