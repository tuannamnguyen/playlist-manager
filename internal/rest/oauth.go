package rest

import (
	"fmt"
	"log"
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

	if _, err := gothic.CompleteUserAuth(c.Response(), c.Request()); err != nil {
		return c.NoContent(http.StatusOK)
	}
	log.Println("oauth2 user already logged in")

	gothic.BeginAuthHandler(c.Response(), c.Request())
	return nil
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

	session, err := o.sessionStore.Get(c.Request(), "oauth-session")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting session store: %w", err))
	}

	session.Values[fmt.Sprintf("%s_access_token", provider)] = user.AccessToken
	session.Values[fmt.Sprintf("%s_refresh_token", provider)] = user.RefreshToken
	session.Values[fmt.Sprintf("%s_token_expiry", provider)] = user.ExpiresAt

	err = session.Save(c.Request(), c.Response())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error saving session data: %w", err))
	}

	return c.JSON(http.StatusOK, user)
}
