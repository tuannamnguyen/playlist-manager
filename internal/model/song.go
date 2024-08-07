package model

type Song struct {
	ID          int      `json:"song_id" db:"song_id"`
	Name        string   `json:"song_name" db:"song_name"`
	ArtistNames []string `json:"artist_names" db:"artist_names"`
	AlbumName   string   `json:"album_name" db:"album_name"`
	Timestamp
}

type SongIn struct {
	Name        string   `json:"song_name" db:"song_name"`
	ArtistNames []string `json:"artist_names" db:"artist_names"`
	AlbumName   string   `json:"album_name" db:"album_name"`
}
