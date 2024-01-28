package spotifyapi

import (
	"context"
	"strings"

	"github.com/zmb3/spotify/v2"
)

const (
	PAGE_SIZE  = 50
	INIT_PAGES = 1
)

func GetTracks(uri string) ([]spotify.FullTrack, string, error) {
	ctx := context.Background()

	var title string
	offset := 0
	found := PAGE_SIZE
	tracks := make([]spotify.FullTrack, 0, INIT_PAGES*PAGE_SIZE)
	for found == PAGE_SIZE && offset < INIT_PAGES*PAGE_SIZE {
		found = 0

		if strings.HasSuffix(uri, "collection") {
			result, err := Client.CurrentUsersTracks(
				ctx,
				spotify.Limit(PAGE_SIZE),
				spotify.Offset(offset),
			)
			if err != nil {
				return []spotify.FullTrack{}, "", err
			}

			title = "Liked Songs"
			for _, t := range result.Tracks {
				tracks = append(tracks, t.FullTrack)
				found++
			}

		} else {
			id := strings.Split(uri, ":")[2]
			result, err := Client.GetPlaylist(
				ctx,
				spotify.ID(id),
				spotify.Limit(PAGE_SIZE),
				spotify.Offset(offset),
			)
			if err != nil {
				return []spotify.FullTrack{}, "", err
			}

			title = result.Name
			for _, t := range result.Tracks.Tracks {
				tracks = append(tracks, t.Track)
				found++
			}
		}

		offset += PAGE_SIZE
	}

	return tracks, title, nil
}

func GetPlaylists() ([]spotify.SimplePlaylist, error) {
	offset := 0
	found := PAGE_SIZE
	playlists := make([]spotify.SimplePlaylist, 0, INIT_PAGES*PAGE_SIZE)
	for found == PAGE_SIZE && offset < INIT_PAGES*PAGE_SIZE {
		result, err := Client.CurrentUsersPlaylists(context.Background(),
			spotify.Limit(PAGE_SIZE),
			spotify.Offset(offset),
		)
		if err != nil {
			return []spotify.SimplePlaylist{}, err
		}

		found = len(result.Playlists)
		playlists = append(playlists, result.Playlists...)
		offset += PAGE_SIZE
	}

	return playlists, nil
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
