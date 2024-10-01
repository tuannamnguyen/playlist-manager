package model

import "database/sql"

type SongInAPI struct {
	Name        string   `json:"song_name"`
	ArtistNames []string `json:"artist_names"`
	AlbumName   string   `json:"album_name"`
	Duration    int      `json:"duration"`
	ImageURL    string   `json:"image_url"`
	ISRC        string   `json:"isrc"`
}

type SongInDB struct {
	Name     string `db:"song_name"`
	AlbumID  int    `db:"album_id"`
	ImageURL string `db:"image_url"`
	Duration int    `db:"duration"`
	ISRC     string `db:"isrc"`
}

type SongOutAPI struct {
	ID          int      `json:"song_id"`
	Name        string   `json:"song_name"`
	ArtistNames []string `json:"artist_names"`
	AlbumName   string   `json:"album_name"`
	ImageURL    string   `json:"image_url"`
	Duration    int      `json:"duration"`
	ISRC        string   `json:"isrc"`
	Timestamp
}

type SongOutDB struct {
	ID         int            `db:"song_id"`
	Name       string         `db:"song_name"`
	AlbumName  string         `db:"album_name"`
	ArtistName string         `db:"artist_name"`
	ImageURL   string         `db:"image_url"`
	Duration   int            `db:"duration"`
	ISRC       sql.NullString `db:"isrc"`
	Timestamp
}
