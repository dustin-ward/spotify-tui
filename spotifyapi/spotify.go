package spotifyapi

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

const (
	PAGE_SIZE = 50
)

func GetLikedTracks() ([]spotify.SavedTrack, error) {
	offset := 0
	found := PAGE_SIZE
	tracks := make([]spotify.SavedTrack, 0, 500)
	for found == PAGE_SIZE && offset < 10*PAGE_SIZE {
		result, err := Client.CurrentUsersTracks(context.Background(),
			spotify.Limit(PAGE_SIZE),
			spotify.Offset(offset),
		)
		if err != nil {
			return []spotify.SavedTrack{}, err
		}

		found = len(result.Tracks)
		tracks = append(tracks, result.Tracks...)
		offset += PAGE_SIZE
	}

	return tracks, nil
}

func PlayPauseTrack(uri *spotify.URI, device spotify.ID) (bool, error) {
	ctx := context.Background()
	playbackStatus, err := Client.PlayerState(ctx)
	if err != nil {
		return false, err
	}

	if playbackStatus.Item != nil && (*playbackStatus.Item).URI == *uri {
		if !playbackStatus.Playing {
			return true, Client.Play(ctx)
		} else {
			return false, Client.Pause(ctx)
		}
	} else if uri == nil {
		return false, Client.Pause(ctx)
	} else {
		return true, Client.PlayOpt(ctx, &spotify.PlayOptions{
			DeviceID: &device,
			URIs:     []spotify.URI{*uri},
		})
	}
}
