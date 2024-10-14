package model

import "golang.org/x/oauth2"

type ConverterRequestData struct {
	PlaylistName     string                    `json:"playlist_name" validate:"required"`
	ProviderMetadata ConverterProviderMetadata `json:"provider_metadata,omitempty"`
	ProviderParam
}

type ConverterProviderMetadata struct {
	AppleMusic struct {
		MusicUserToken string `json:"musicUserToken"`
	} `json:"applemusic"`
}

type ConverterServiceProviderMetadata struct {
	AppleMusic AppleMusicMetadata
	Spotify    SpotifyMetadata
}

type AppleMusicMetadata struct {
	MusicUserToken string
}

type SpotifyMetadata struct {
	Token *oauth2.Token
}

type ProviderParam struct {
	Provider string `param:"provider" validate:"required,oneof=spotify applemusic"`
}
