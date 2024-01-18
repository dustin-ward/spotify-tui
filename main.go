package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dustin-ward/spotify-tui/components"
	"github.com/dustin-ward/spotify-tui/spotifyapi"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

type mainModel struct {
	list         tea.Model
	player       tea.Model
	choice       *components.SavedTrack
	currentTrack *spotify.FullTrack
}

func (m mainModel) Init() tea.Cmd {
	err := spotifyapi.InitDevice()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func newMainModel(sptracks []spotify.SavedTrack) mainModel {
	return mainModel{
		list:   components.NewListModel(sptracks, "Liked Songs"),
		player: components.NewPlayerModel(),
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	m.list.SetWidth(msg.Width)
	// 	return m, nil

	// case tea.KeyMsg:
	// 	switch keypress := msg.String(); keypress {
	// 	case "q", "ctrl+c":
	// 		return m, tea.Quit

	// 	case "enter":
	// 		t, ok := m.list.SelectedItem().(components.SavedTrack)
	// 		if ok {
	// 			m.choice = &t
	// 		}
	// 		return m, components.PlayPauseTrack(&t.URI)
	// 	}

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
	return lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), m.player.View())
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	err := spotifyapi.Login()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Getting tracks...")
	sptracks, err := spotifyapi.GetLikedTracks()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tea.NewProgram(newMainModel(sptracks), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
