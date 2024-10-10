package applemusicconverter

import (
	"context"
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
	// TODO: implement this later
	return nil
}
