package rest

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
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
	provider, err := getProvider(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	addQueryParams(c, provider)

	gothic.BeginAuthHandler(c.Response(), c.Request())

	log.Printf("logged in to %s successfully", provider)

	return nil
}

func (o *OAuthHandler) CallbackHandler(c echo.Context) error {
	provider, err := getProvider(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	addQueryParams(c, provider)

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

func (o *OAuthHandler) CheckAuthHandler(c echo.Context) error {
	provider, err := getProvider(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	addQueryParams(c, provider)

	sessionValues, err := getOauthSessionValues(c.Request(), o.sessionStore)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting session values: %v", err))
	}

	userInfo, ok := sessionValues[fmt.Sprintf("%s_user_info", provider)]
	if ok {
		user := userInfo.(goth.User)
		log.Printf("user logged in. email: %s, name: %s, provider: %s", user.Email, user.Name, user.Provider)

		return c.String(http.StatusOK, fmt.Sprintf("user has authenticated with %s", provider))
	}

	return c.String(http.StatusUnauthorized, fmt.Sprintf("not authenticated with %s", provider))
}

func (o *OAuthHandler) LogoutHandler(c echo.Context) error {
	provider, err := getProvider(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	addQueryParams(c, provider)

	err = gothic.Logout(c.Response(), c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error logging out from %v: %v", provider, err))
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (o *OAuthHandler) GetAccessTokenHandler(c echo.Context) error {
	provider, err := getProvider(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	addQueryParams(c, provider)

	sessionValues, err := getOauthSessionValues(c.Request(), o.sessionStore)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting session values: %v", err))
	}

	user := sessionValues[fmt.Sprintf("%s_user_info", provider)].(goth.User)

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": user.AccessToken,
	})
}
