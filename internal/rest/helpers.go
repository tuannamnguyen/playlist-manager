package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
	"golang.org/x/oauth2"
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

	log.Println("oauth session saved")

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
	var providerParam model.ProviderParam
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

func getProviderMetadata(
	provider string,
	sessionValues map[any]any,
	reqBody model.ConverterRequestData,
) (providerMetadata model.ConverterServiceProviderMetadata) {
	switch provider {
	case "spotify":
		user := (sessionValues[fmt.Sprintf("%s_user_info", provider)]).(goth.User)
		token := &oauth2.Token{
			AccessToken:  user.AccessToken,
			RefreshToken: user.RefreshToken,
			Expiry:       user.ExpiresAt,
		}
		providerMetadata = model.ConverterServiceProviderMetadata{
			Spotify: model.SpotifyMetadata{
				Token: token,
			},
		}

	case "applemusic":
		providerMetadata = model.ConverterServiceProviderMetadata{
			AppleMusic: model.AppleMusicMetadata{
				MusicUserToken: reqBody.ProviderMetadata.AppleMusic.MusicUserToken,
			},
		}
	}

	return providerMetadata
}
