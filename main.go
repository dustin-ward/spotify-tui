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
	docStyle       = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	focusedStyle   = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#00bf06"))
	unfocusedStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#5900bf"))
)

type appFocus int

const (
	focusList appFocus = iota
	focusPlayer
	focusCollections
)

type mainModel struct {
	list        tea.Model
	player      tea.Model
	collections tea.Model
	focus       appFocus
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func newMainModel() mainModel {
	return mainModel{
		list:        components.NewListModel("Liked Songs"),
		player:      components.NewPlayerModel(),
		collections: components.NewCollectionsModel(),
		focus:       focusList,
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case components.PlayerStatusMsg:
		str := string(msg)
		if str != "ok" {
			log.Println(str)
		}

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "ctrl+l":
			switch m.focus {
			case focusList:
				m.focus = focusCollections
			case focusPlayer:
			case focusCollections:
			}
		case "ctrl+h":
			switch m.focus {
			case focusList:
			case focusPlayer:
				m.focus = focusList
			case focusCollections:
				m.focus = focusList
			}
		case "ctrl+k":
			switch m.focus {
			case focusList:
			case focusPlayer:
			case focusCollections:
				m.focus = focusPlayer
			}
		case "ctrl+j":
			switch m.focus {
			case focusList:
			case focusPlayer:
				m.focus = focusCollections
			case focusCollections:
			}
		}
	}

	var cmd tea.Cmd
	switch m.focus {
	case focusList:
		m.list, cmd = m.list.Update(msg)
	case focusPlayer:
		m.player, cmd = m.player.Update(msg)
	case focusCollections:
		m.collections, cmd = m.collections.Update(msg)
	}
	return m, cmd
}

func (m mainModel) View() string {
	listStyle := unfocusedStyle
	playerStyle := unfocusedStyle
	collectionsStyle := unfocusedStyle
	switch m.focus {
	case focusList:
		listStyle = focusedStyle
	case focusPlayer:
		playerStyle = focusedStyle
	case focusCollections:
		collectionsStyle = focusedStyle
	}

	return docStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			listStyle.Render(m.list.View()),
			lipgloss.JoinVertical(
				lipgloss.Left,
				playerStyle.Render(m.player.View()),
				collectionsStyle.Render(m.collections.View()),
			),
		),
	)
}

func init() {
	target := os.Getenv("SPOTIFY_DEVICE")
	if target == "" {
		log.Fatal("Error: No device set in $SPOTIFY_DEVICE")
	}

	err := spotifyapi.Login()
	if err != nil {
		log.Fatal(err)
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
	if _, err := tea.NewProgram(newMainModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
