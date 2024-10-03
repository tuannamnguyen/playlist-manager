package geniusprovider

import (
	"errors"
	"time"

	"github.com/markbates/goth"
)

type Session struct {
	AuthURL      string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}


func (s Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(goth.NoAuthUrlErrorMessage)
	}

	return s.AuthURL, nil
}