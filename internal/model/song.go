package model

type SongOutAPI struct {
	ID          int      `json:"song_id"`
	Name        string   `json:"song_name"`
	ArtistNames []string `json:"artist_names"`
	AlbumName   string   `json:"album_name"`
	Timestamp
}

type SongInAPI struct {
	Name        string   `json:"song_name"`
	ArtistNames []string `json:"artist_names"`
	AlbumName   string   `json:"album_name"`
}

type SongInDB struct {
	Name    string `db:"song_name"`
	AlbumID int    `db:"album_id"`
}

type SongOutDB struct {
	ID         int    `db:"song_id"`
	Name       string `db:"song_name"`
	ArtistName string `db:"artist_name"`
	Timestamp
}
