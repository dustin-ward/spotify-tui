package spotifyapi

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

func GetLikedTracks() ([]spotify.SavedTrack, error) {
	tracks, err := spclient.CurrentUsersTracks(context.Background())
	if err != nil {
		return []spotify.SavedTrack{}, err
	}

	return tracks.Tracks, nil
}
