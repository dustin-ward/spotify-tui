package spotifyapi

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

const (
	PAGE_SIZE  = 50
	INIT_PAGES = 3
)

func GetLikedTracks() ([]spotify.SavedTrack, error) {
	offset := 0
	found := PAGE_SIZE
	tracks := make([]spotify.SavedTrack, 0, INIT_PAGES*PAGE_SIZE)
	for found == PAGE_SIZE && offset < INIT_PAGES*PAGE_SIZE {
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

func PlayPauseTrack(uri *spotify.URI, device spotify.ID, playbackContext *spotify.URI) (bool, error) {
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
			DeviceID:        &device,
			PlaybackContext: playbackContext,
			PlaybackOffset: &spotify.PlaybackOffset{
				URI: *uri,
			},
			// URIs:            []spotify.URI{*uri},
		})
	}
}
