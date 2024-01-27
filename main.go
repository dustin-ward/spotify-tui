package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dustin-ward/spotify-tui/components"
	"github.com/dustin-ward/spotify-tui/spotifyapi"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

var (
	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

type mainModel struct {
	list        tea.Model
	player      tea.Model
	collections tea.Model
	focus       int
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func newMainModel(sptracks []spotify.SavedTrack, spplaylists []spotify.SimplePlaylist) mainModel {
	return mainModel{
		list:        components.NewListModel(sptracks, "Liked Songs"),
		player:      components.NewPlayerModel(),
		collections: components.NewCollectionsModel(spplaylists),
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case components.PlayerStatusMsg:
		str := string(msg)
		if str != "ok" {
			log.Println(str)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m mainModel) View() string {
	return docStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.list.View(),
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.player.View(),
				m.collections.View(),
			),
		),
	)
}

func init() {
	err := spotifyapi.Login()
	if err != nil {
		log.Fatal(err)
	}

	target := os.Getenv("SPOTIFY_DEVICE")
	if target == "" {
		log.Fatal("Error: No device set in $SPOTIFY_DEVICE")
	}

	ctx := context.Background()
	devices, err := spotifyapi.Client.PlayerDevices(ctx)
	if err != nil {
		log.Fatal("Error:", err)
	}

	components.DEVICE_ID = ""
	for _, d := range devices {
		if d.Name == target {
			components.DEVICE_ID = d.ID
		}
	}

	if components.DEVICE_ID == "" {
		log.Fatal("Error: Device not found. Make sure spotify is runnning")
	}

	playerState, err := spotifyapi.Client.PlayerState(ctx)
	if err != nil {
		log.Fatal("Error:", err)
	}

	components.CurrentlyPlaying = playerState.Item
	components.IsPlaying = playerState.Playing

	currentUser, err := spotifyapi.Client.CurrentUser(ctx)
	if err != nil {
		log.Fatal("Error:", err)
	}
	uri := spotify.URI(fmt.Sprintf("%s:collection", currentUser.URI))
	components.CurrentContext = &uri
}

func main() {
	log.Println("Getting tracks...")
	sptracks, err := spotifyapi.GetLikedTracks()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Getting collections...")
	spplaylists, err := spotifyapi.GetPlaylists()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tea.NewProgram(newMainModel(sptracks, spplaylists), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
