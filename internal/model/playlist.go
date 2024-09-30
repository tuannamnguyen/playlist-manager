package model

import "database/sql"

type Playlist struct {
	ID                  int    `json:"playlist_id"`
	Name                string `json:"playlist_name"`
	PlaylistDescription string `json:"playlist_description"`
	UserID              string `json:"user_id"`
	Username            string `json:"user_name"`
	ImageURL            string `json:"image_url"`
	Timestamp
}

type PlaylistOutDB struct {
	ID                  int            `db:"playlist_id"`
	Name                string         `db:"playlist_name"`
	PlaylistDescription sql.NullString `db:"playlist_description"`
	UserID              string         `db:"user_id"`
	Username            string         `db:"user_name"`
	ImageURL            string         `db:"image_url"`
	Timestamp
}

type PlaylistInDB struct {
	Name                string `db:"playlist_name"`
	PlaylistDescription string `db:"playlist_description"`
	UserID              string `db:"user_id"`
	Username            string `db:"user_name"`
	ImageURL            string `db:"image_url"`
}

type PlaylistIn struct {
	Name                string `json:"playlist_name" validate:"required"`
	PlaylistDescription string `json:"playlist_description"`
	UserID              string `json:"user_id" validate:"required"`
	Username            string `json:"user_name" validate:"required"`
}
