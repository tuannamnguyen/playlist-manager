package rest

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func saveOauthSessionValues(req *http.Request, res http.ResponseWriter, store sessions.Store, sessionValues map[any]any) error {
	session, err := store.Get(req, "oauth-session")
	if err != nil {
		return err
	}

	session.Values = sessionValues

	err = session.Save(req, res)
	if err != nil {
		return err
	}
	return nil
}

func getOauthSessionValues(req *http.Request, store sessions.Store) (map[any]any, error) {
	session, err := store.Get(req, "oauth-session")
	if err != nil {
		return nil, err
	}

	return session.Values, nil
}

func getProvider(c echo.Context) (string, error) {
	var providerParam ProviderParam
	err := c.Bind(&providerParam)
	if err != nil {
		return "", err
	}
	provider := providerParam.Provider
	return provider, nil
}

func addQueryParams(c echo.Context, provider string) {
	q := c.Request().URL.Query()
	q.Add("provider", provider)
	c.Request().URL.RawQuery = q.Encode()
}
