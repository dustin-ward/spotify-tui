package components

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin-ward/spotify-tui/spotifyapi"
	"github.com/zmb3/spotify/v2"
)

var (
	playerStyle = lipgloss.NewStyle().Align(lipgloss.Center).Width(65).Height(10).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#5900bf")).PaddingTop(2).PaddingBottom(2)
)

var (
	DEVICE_ID        spotify.ID
	CurrentlyPlaying *spotify.FullTrack
	IsPlaying        bool
)

type PlayerModel struct {
}

func NewPlayerModel() PlayerModel {
	return PlayerModel{}
}

type PlayerStatusMsg string

func PlayPauseTrack(uri *spotify.URI) tea.Cmd {
	return func() tea.Msg {
		var err error
		IsPlaying, err = spotifyapi.PlayPauseTrack(uri, DEVICE_ID)

		if err != nil {
			return PlayerStatusMsg("err " + err.Error())
		}

		if uri != nil {
			if IsPlaying {
				// spotify:track:6rqhFgbbKwnb9MLmUQDhG6
				id := strings.Split(string(*uri), ":")[2]
				track, err := spotifyapi.Client.GetTrack(context.Background(), spotify.ID(id))
				if err != nil {
					return PlayerStatusMsg("Unable to get track: " + err.Error())
				}
				CurrentlyPlaying = track
			} else {
				CurrentlyPlaying = nil
			}
		}

		return PlayerStatusMsg("ok")
	}
}

func (m PlayerModel) Init() tea.Cmd {
	return nil
}

func (m PlayerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m PlayerModel) View() string {
	view := ""
	if CurrentlyPlaying != nil {
		view = CurrentlyPlaying.Name + "\n"
		artists := CurrentlyPlaying.Artists[0].Name
		for i, a := range CurrentlyPlaying.Artists {
			if i == 0 {
				continue
			}
			artists = fmt.Sprintf("%s, %s", artists, a.Name)
		}
		view += artists + "\n"
		view += CurrentlyPlaying.Album.Name + "\n\n"
		view += "◀◀   ▷   ▶▶"
	} else {
		view = "Nothing playing at the moment..."
	}
	return playerStyle.Render(view)
}
