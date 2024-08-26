package service

import (
	"context"
	"errors"

	spotifyconverter "github.com/tuannamnguyen/playlist-manager/internal/service/converters/spotify"

	"golang.org/x/oauth2"
)

func getConverter(ctx context.Context, provider string, token *oauth2.Token) (Converter, error) {
	if provider == "spotify" {
		return spotifyconverter.New(ctx, token), nil
	}

	return nil, errors.New("no converter available")
}
