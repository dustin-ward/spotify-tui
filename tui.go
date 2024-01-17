package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dustin-ward/spotify-tui/spotifyapi"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("46"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type mainModel struct {
	list         list.Model
	choice       *SavedTrack
	currentTrack *spotify.FullTrack
	quitting     bool
}

func (m mainModel) Init() tea.Cmd {
	err := spotifyapi.InitDevice()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func newMainModel(sptracks []spotify.SavedTrack) mainModel {
	return mainModel{list: newListModel(sptracks)}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			t, ok := m.list.SelectedItem().(SavedTrack)
			if ok {
				m.choice = &t
			}
			return m, playPauseTrack(&t.URI)
		}

	case playerStatusMsg:
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
	return "\n" + m.list.View()
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

func fmtDuration(t time.Duration) string {
	s := int(t.Seconds())
	m := s / 60
	s %= 60
	return fmt.Sprintf("%d:%02d", m, s)
}

type playerStatusMsg string

func playPauseTrack(uri *spotify.URI) tea.Cmd {
	return func() tea.Msg {
		err := spotifyapi.PlayPauseTrack(uri)

		if err != nil {
			return playerStatusMsg("err " + err.Error())
		}

		return playerStatusMsg("ok")
	}
}
