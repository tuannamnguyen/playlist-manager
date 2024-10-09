package applemusicconverter

import (
	"context"

	applemusic "github.com/minchao/go-apple-music"
)

type AppleMusicConverter struct {
	client *applemusic.Client
}

func New(ctx context.Context) *AppleMusicConverter {
	// TODO: update this later
	tp := applemusic.Transport{}
	client := applemusic.NewClient(tp.Client())

	return &AppleMusicConverter{client: client}
}
