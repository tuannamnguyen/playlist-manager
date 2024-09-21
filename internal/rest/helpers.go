package rest

import (
	"net/http"

	"github.com/gorilla/sessions"
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
