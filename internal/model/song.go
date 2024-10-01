package model

type SongInAPI struct {
	Name        string   `json:"song_name"`
	ArtistNames []string `json:"artist_names"`
	AlbumName   string   `json:"album_name"`
	Duration    int      `json:"duration"`
	ImageURL    string   `json:"image_url"`
}

type SongInDB struct {
	Name     string `db:"song_name"`
	AlbumID  int    `db:"album_id"`
	ImageURL string `db:"image_url"`
	Duration int    `db:"duration"`
}

type SongOutAPI struct {
	ID          int      `json:"song_id"`
	Name        string   `json:"song_name"`
	ArtistNames []string `json:"artist_names"`
	AlbumName   string   `json:"album_name"`
	ImageURL    string   `json:"image_url"`
	Duration    int      `json:"duration"`
	Timestamp
}

type SongOutDB struct {
	ID         int    `db:"song_id"`
	Name       string `db:"song_name"`
	AlbumName  string `db:"album_name"`
	ArtistName string `db:"artist_name"`
	ImageURL   string `db:"image_url"`
	Duration   int    `db:"duration"`
	Timestamp
}
