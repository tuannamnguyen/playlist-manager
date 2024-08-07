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
	Name    string
	AlbumID int
}
