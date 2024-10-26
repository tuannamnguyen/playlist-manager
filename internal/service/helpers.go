package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

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

func writeCsvRecord(w io.Writer, record []string) error {
	writer := csv.NewWriter(w)

	err := writer.Write(record)
	if err != nil {
		return fmt.Errorf("csv write: %s", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("csv flush: %s", err)
	}
	return nil
}

func parseCsv(r io.Reader) ([][]string, error) {
	reader := csv.NewReader(r)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csv read all records: %s", err)
	}

	return records, nil
}
