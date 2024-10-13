package model

import "golang.org/x/oauth2"

type ConverterRequestBody struct {
	PlaylistName     string                    `json:"playlist_name"`
	ProviderMetadata ConverterProviderMetadata `json:"provider_metadata,omitempty"`
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
