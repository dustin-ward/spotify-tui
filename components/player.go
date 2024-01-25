package components

import (
	"context"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin-ward/spotify-tui/spotifyapi"
	"github.com/zmb3/spotify/v2"
)

var (
	playerStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("46"))
)

var (
	DEVICE_ID        spotify.ID
	CurrentlyPlaying *spotify.FullTrack
)

type PlayerModel struct {
	isPlaying bool
}

func NewPlayerModel() PlayerModel {
	return PlayerModel{}
}

type PlayerStatusMsg string

func PlayPauseTrack(uri *spotify.URI) tea.Cmd {
	return func() tea.Msg {
		playing, err := spotifyapi.PlayPauseTrack(uri, DEVICE_ID)

		if err != nil {
			return PlayerStatusMsg("err " + err.Error())
		}

		if uri != nil {
			if playing {
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
	return "TEST"
}
