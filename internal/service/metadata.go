package service

import (
	"fmt"
	"net/url"

	lyrics "github.com/rhnvrm/lyric-api-go"
)

type MetadataService struct {
	lyric lyrics.Lyric
}

func NewMetadataService() *MetadataService {
	lyric := lyrics.New()

	return &MetadataService{lyric: lyric}
}

func (m *MetadataService) FetchLyrics(artists string, song string) (string, error) {
	return m.lyric.Search(artists, song)
}

func (m *MetadataService) ArtistInformation(artistName string) (string, error) {
	escapedArtistName := url.PathEscape(artistName)
	url := fmt.Sprintf("https://genius.com/artists/%s", escapedArtistName)

	return url, nil
}
