package service

import "github.com/tuannamnguyen/playlist-manager/internal/model"

type SearchRepository interface {
	Song(track string, artist string) (model.SongIn, error)
}

func NewSearch(searchRepository SearchRepository) *SearchService {
	return &SearchService{sr: searchRepository}
}

type SearchService struct {
	sr SearchRepository
}

func (s *SearchService) SongSearch(track string, artist string) (model.SongIn, error) {
	return s.sr.Song(track, artist)
}
