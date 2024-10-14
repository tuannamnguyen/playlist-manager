package service

import (
	"context"
	"errors"

	"github.com/tuannamnguyen/playlist-manager/internal/model"
	applemusicconverter "github.com/tuannamnguyen/playlist-manager/internal/service/converters/applemusic"
	spotifyconverter "github.com/tuannamnguyen/playlist-manager/internal/service/converters/spotify"
)

func getConverter(ctx context.Context, provider string, providerMetadata model.ConverterServiceProviderMetadata) (Converter, error) {
	switch provider {
	case "spotify":
		return spotifyconverter.New(ctx, providerMetadata.Spotify.Token), nil
	case "applemusic":
		return applemusicconverter.New(ctx, providerMetadata.AppleMusic.MusicUserToken), nil
	}

	return nil, errors.New("no converter available")
}
