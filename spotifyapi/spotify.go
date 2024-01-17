package spotifyapi

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

func GetLikedTracks() ([]spotify.SavedTrack, error) {
	offset := 0
	found := 50
	tracks := make([]spotify.SavedTrack, 0, 500)
	for found == 50 && offset < 10 {
		result, err := spclient.CurrentUsersTracks(context.Background(),
			spotify.Limit(50),
			spotify.Offset(offset),
		)
		if err != nil {
			return []spotify.SavedTrack{}, err
		}

		found = len(result.Tracks)
		tracks = append(tracks, result.Tracks...)
		offset += 1
	}

	return tracks, nil
}

func PlayPauseTrack(uri *spotify.URI) error {
	ctx := context.Background()
	playbackStatus, err := spclient.PlayerState(ctx)
	if err != nil {
		return err
	}

	if playbackStatus.Item != nil && (*playbackStatus.Item).URI == *uri {
		if !playbackStatus.Playing {
			return spclient.Play(ctx)
		} else {
			return spclient.Pause(ctx)
		}
	} else if uri == nil {
		return spclient.Pause(ctx)
	} else {
		return spclient.PlayOpt(ctx, &spotify.PlayOptions{
			URIs: []spotify.URI{*uri},
		})
	}

}
