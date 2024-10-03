package service

import (
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
