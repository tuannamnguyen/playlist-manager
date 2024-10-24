package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"

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

func writeCsvRecord(writer *csv.Writer, record []string) error {
	err := writer.Write(record)
	if err != nil {
		return fmt.Errorf("csv write: %s", err)
	}

	return nil
}
