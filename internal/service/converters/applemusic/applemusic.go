package applemusicconverter

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
	var id, songType string

	if song.ISRC != "" {
		amTrack, _, err := a.client.Catalog.GetSongsByIsrcs(ctx, "vn", []string{song.ISRC}, &applemusic.Options{})
		if err != nil {
			return applemusic.CreateLibraryPlaylistTrack{}, fmt.Errorf("get songs by ISRC: %w", err)
		}
		if len(amTrack.Data) == 0 {
			return applemusic.CreateLibraryPlaylistTrack{}, fmt.Errorf("no track found for ISRC: %s", song.ISRC)
		}

		id = amTrack.Data[0].Id
		songType = amTrack.Data[0].Type
	} else {
		artistSearch := strings.Join(song.ArtistNames, " ")
		searchTerm := fmt.Sprintf("%s %s %s", song.Name, artistSearch, song.AlbumName)

		log.Printf("apple music search term: %s", searchTerm)

		searchResult, _, err := a.client.Catalog.Search(ctx, "vn", &applemusic.SearchOptions{
			Term:  searchTerm,
			Types: "songs",
		})
		if err != nil {
			return applemusic.CreateLibraryPlaylistTrack{}, fmt.Errorf("search songs: %w", err)
		}

		songs := searchResult.Results.Songs
		if songs == nil {
			return applemusic.CreateLibraryPlaylistTrack{}, fmt.Errorf("no track found for search term: %s", searchTerm)
		}

		id = searchResult.Results.Songs.Data[0].Id
		songType = searchResult.Results.Songs.Data[0].Type
	}

	return applemusic.CreateLibraryPlaylistTrack{
		Id:   id,
		Type: songType,
	}, nil
}
