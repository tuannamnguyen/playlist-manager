package applemusicconverter

import (
	"context"
	"fmt"
	"os"

	applemusic "github.com/minchao/go-apple-music"
	"github.com/tuannamnguyen/playlist-manager/internal/model"
)

type AppleMusicConverter struct {
	client *applemusic.Client
}

func New(ctx context.Context, musicUserToken string) *AppleMusicConverter {
	tp := applemusic.Transport{
		Token:          os.Getenv("APPLE_MUSIC_ACCESS_TOKEN"),
		MusicUserToken: musicUserToken,
	}
	client := applemusic.NewClient(tp.Client())

	return &AppleMusicConverter{client: client}
}

func (a *AppleMusicConverter) Export(ctx context.Context, playlistName string, songs []model.SongOutAPI) error {
	libraryPlaylistTracks, err := a.parseSongToPlaylistTracks(ctx, songs)
	if err != nil {
		return fmt.Errorf("parse song models to playlist tracks: %w", err)
	}

	_, _, err = a.client.Me.CreateLibraryPlaylist(
		ctx,
		applemusic.CreateLibraryPlaylist{
			Attributes: applemusic.CreateLibraryPlaylistAttributes{
				Name:        playlistName,
				Description: "",
			},
			Relationships: &applemusic.CreateLibraryPlaylistRelationships{
				Tracks: applemusic.CreateLibraryPlaylistTrackData{
					Data: libraryPlaylistTracks,
				},
			},
		},
		&applemusic.Options{},
	)

	if err != nil {
		return fmt.Errorf("create apple music playlist: %w", err)
	}

	return nil
}

func (a *AppleMusicConverter) parseSongToPlaylistTracks(ctx context.Context, songs []model.SongOutAPI) ([]applemusic.CreateLibraryPlaylistTrack, error) {
	libraryPlaylistTracks := make([]applemusic.CreateLibraryPlaylistTrack, len(songs))

	for i, song := range songs {
		amTrack, err := a.searchAndMatch(ctx, song)
		if err != nil {
			return nil, fmt.Errorf("search and match apple music track: %w", err)
		}

		libraryPlaylistTracks[i] = amTrack
	}

	return libraryPlaylistTracks, nil
}

func (a *AppleMusicConverter) searchAndMatch(ctx context.Context, song model.SongOutAPI) (applemusic.CreateLibraryPlaylistTrack, error) {
	amTrack, _, err := a.client.Catalog.GetSongsByIsrcs(
		ctx,
		"vn",
		[]string{song.ISRC},
		&applemusic.Options{},
	)
	if err != nil {
		return applemusic.CreateLibraryPlaylistTrack{}, err
	}

	return applemusic.CreateLibraryPlaylistTrack{
		Id:   amTrack.Data[0].Id,
		Type: amTrack.Data[0].Type,
	}, nil
}
