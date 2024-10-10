package service

import (
	"context"
	"errors"

	applemusicconverter "github.com/tuannamnguyen/playlist-manager/internal/service/converters/applemusic"
	spotifyconverter "github.com/tuannamnguyen/playlist-manager/internal/service/converters/spotify"

	"golang.org/x/oauth2"
)

func getConverter(ctx context.Context, provider string, token *oauth2.Token, appleMusicUserToken string) (Converter, error) {
	if provider == "spotify" {
		return spotifyconverter.New(ctx, token), nil
	} else if provider == "applemusic" {
		return applemusicconverter.New(ctx, appleMusicUserToken), nil
	}

	return nil, errors.New("no converter available")
}
