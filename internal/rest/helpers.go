package rest

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func saveSessionValues(c echo.Context, store sessions.Store, sessionValues map[any]any) error {
	session, err := store.Get(c.Request(), "oauth-session")
	if err != nil {
		return err
	}

	session.Values = sessionValues

	err = session.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}
	return nil
}
