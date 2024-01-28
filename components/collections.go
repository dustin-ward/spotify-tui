package components

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin-ward/spotify-tui/colours"
	"github.com/dustin-ward/spotify-tui/spotifyapi"
	"github.com/zmb3/spotify/v2"
)

var (
	collectionsStyle = lipgloss.NewStyle().Width(65).Height(31).Padding(1)
)

type UpdatePlaylistMsg string

type CollectionsModel struct {
	list list.Model
}

func NewCollectionsModel(userURI spotify.URI) CollectionsModel {
	spplaylists, err := spotifyapi.GetPlaylists()
	if err != nil {
		log.Fatal(err)
	}

	playlists := make([]list.Item, 0, 100)
	playlists = append(playlists, SimplePlaylist{
		spotify.SimplePlaylist{
			Name:        "Liked Songs",
			Description: "Your saved tracks",
			URI:         userURI + ":collection",
		},
	})
	for _, p := range spplaylists {
		playlists = append(playlists, SimplePlaylist{p})
	}

	d := playlistDelegate{}
	l := list.New(playlists, d, 58, 27)
	l.Title = "Playlists and Mixes"
	l.Styles.Title = titleStyle
	l.Styles.FilterCursor = l.Styles.FilterCursor.Foreground(colours.GREEN)
	l.SetShowHelp(false)

	return CollectionsModel{list: l}
}

func (m CollectionsModel) Init() tea.Cmd {
	return nil
}

func (m CollectionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			p, _ := m.list.SelectedItem().(SimplePlaylist)
			return m, func() tea.Msg { return UpdatePlaylistMsg(p.URI) }
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m CollectionsModel) View() string {
	return collectionsStyle.Render(m.list.View())
}
