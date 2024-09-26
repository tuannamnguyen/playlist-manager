package service

import "github.com/tuannamnguyen/playlist-manager/internal/model"

type SearchRepository interface {
	Song(track string, artist string, album string) (model.SongInAPI, error)
}

func NewSearch(searchRepository SearchRepository) *SearchService {
	return &SearchService{sr: searchRepository}
}

type SearchService struct {
	sr SearchRepository
}

func (s *SearchService) SongSearch(track string, artist string, album string) (model.SongInAPI, error) {
	return s.sr.Song(track, artist, album)
}
