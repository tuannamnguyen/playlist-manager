package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type OAuthHandler struct {
	sessionStore sessions.Store
}

func NewOAuthHandler(store sessions.Store) *OAuthHandler {
	return &OAuthHandler{
		sessionStore: store,
	}
}

func (o *OAuthHandler) LoginHandler(c echo.Context) error {
	provider := c.Param("provider")
	q := c.Request().URL.Query()
	q.Add("provider", provider)
	c.Request().URL.RawQuery = q.Encode()

	if _, err := gothic.CompleteUserAuth(c.Response(), c.Request()); err == nil {
		return c.NoContent(http.StatusOK)
	} else {
		gothic.BeginAuthHandler(c.Response(), c.Request())
		return nil
	}
}

func (o *OAuthHandler) CallbackHandler(c echo.Context) error {
	provider := c.Param("provider")
	q := c.Request().URL.Query()
	q.Add("provider", provider)
	c.Request().URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error complete user auth: %w", err))
	}

	sessionValues := make(map[any]any)
	sessionValues[fmt.Sprintf("%s_access_token", provider)] = user.AccessToken
	sessionValues[fmt.Sprintf("%s_refresh_token", provider)] = user.RefreshToken
	sessionValues[fmt.Sprintf("%s_token_expiry", provider)] = user.ExpiresAt

	store := o.sessionStore

	err = saveSessionValues(c, store, sessionValues)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error saving session values: %w", err))
	}

	return c.JSON(http.StatusOK, user)
}
