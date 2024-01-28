package components

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin-ward/spotify-tui/spotifyapi"
	"github.com/zmb3/spotify/v2"
)

var (
	collectionsStyle = lipgloss.NewStyle().Width(65).Height(30).Padding(2)
)

type UpdatePlaylistMsg string

type CollectionsModel struct {
	list list.Model
}

type SimplePlaylist struct {
	spotify.SimplePlaylist
}

func (p SimplePlaylist) Title() string { return p.Name }
func (p SimplePlaylist) Description() string {
	if p.SimplePlaylist.Description == "" {
		return p.Owner.DisplayName
	}
	return p.SimplePlaylist.Description
}
func (p SimplePlaylist) FilterValue() string {
	return p.Name + " " + p.SimplePlaylist.Description + " " + p.Owner.DisplayName
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

	l := list.New(playlists, list.NewDefaultDelegate(), 58, 27)
	l.Title = "Playlists and Mixes"
	l.Styles.Title = titleStyle
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
