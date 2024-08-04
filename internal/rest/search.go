package rest

import "github.com/labstack/echo/v4"

type SearchHandler struct{}

func NewSearchHandler() *SearchHandler {
	return &SearchHandler{}
}

func (s *SearchHandler) SearchMusicData(c echo.Context) error {
	// TODO: IMPLEMENT THIS
	return nil
}
