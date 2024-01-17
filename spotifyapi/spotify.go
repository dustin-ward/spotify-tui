package spotifyapi

import (
	"context"
	"fmt"
	"os"

	"github.com/zmb3/spotify/v2"
)

var DEVICE_ID spotify.ID

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
			DeviceID: &DEVICE_ID,
			URIs:     []spotify.URI{*uri},
		})
	}

}

func InitDevice() error {
	ctx := context.Background()
	devices, err := spclient.PlayerDevices(ctx)
	if err != nil {
		return err
	}

	target := os.Getenv("SPOTIFY_DEVICE")
	if target == "" {
		return fmt.Errorf("no device set in $SPOTIFY_DEVICE")
	}

	for _, d := range devices {
		if d.Name == target {
			DEVICE_ID = d.ID
			return nil
		}
	}

	return fmt.Errorf("device: '%s' not found", target)
}
