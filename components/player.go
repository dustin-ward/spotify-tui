package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin-ward/spotify-tui/spotifyapi"
	"github.com/zmb3/spotify/v2"
)

var (
	playerStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("46"))
)

type PlayerModel struct {
	currentlyPlaying *spotify.FullTrack
	isPlaying        bool
}

func NewPlayerModel() PlayerModel {
	return PlayerModel{}
}

type PlayerStatusMsg string

func PlayPauseTrack(uri *spotify.URI) tea.Cmd {
	return func() tea.Msg {
		err := spotifyapi.PlayPauseTrack(uri)

		if err != nil {
			return PlayerStatusMsg("err " + err.Error())
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
