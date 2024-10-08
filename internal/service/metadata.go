package service

import (
	"os"

	"github.com/broxgit/genius"
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
	geniusClient := genius.NewClient(nil, os.Getenv("GENIUS_CLIENT_ACCESS_TOKEN"))

	searchRes, err := geniusClient.Search(artistName)
	if err != nil {
		return "", err
	}

	artist := searchRes.Response.Hits[0].Result.PrimaryArtist

	return artist.URL, nil
}
