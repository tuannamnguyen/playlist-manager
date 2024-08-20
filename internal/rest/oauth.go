package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type OAuthHandler struct {
}

func NewOAuthHandler() *OAuthHandler {
	return &OAuthHandler{}
}

func (o *OAuthHandler) LoginHandler(c echo.Context) error {
	provider := c.Param("provider")
	q := c.Request().URL.Query()
	q.Add("provider", provider)
	c.Request().URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Response(), c.Request())
	return nil
}
