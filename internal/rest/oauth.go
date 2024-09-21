package rest

import (
	"fmt"
	"net/http"
	"os"

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

	sessionValues, err := getOauthSessionValues(c.Request(), o.sessionStore)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting session values: %v", err))
	}

	_, ok := sessionValues[fmt.Sprintf("%s_user_info", provider)]
	if ok {
		return c.String(http.StatusOK, "user has already logged in")
	}

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

	sessionValues := make(map[any]any)
	sessionValues[fmt.Sprintf("%s_user_info", provider)] = user

	store := o.sessionStore

	err = saveOauthSessionValues(c.Request(), c.Response(), store, sessionValues)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error saving session values: %w", err))
	}

	return c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL"))
}

func (o *OAuthHandler) LogoutHandler(c echo.Context) error {
	provider := c.Param("provider")
	q := c.Request().URL.Query()
	q.Add("provider", provider)
	c.Request().URL.RawQuery = q.Encode()

	err := gothic.Logout(c.Response(), c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error logging out from %v: %v", provider, err))
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
