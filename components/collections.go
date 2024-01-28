package components

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

var (
	collectionsStyle = lipgloss.NewStyle().Width(65).Height(30).Padding(2)
)

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

func NewCollectionsModel(spplaylists []spotify.SimplePlaylist) CollectionsModel {
	playlists := make([]list.Item, 0, 100)
	for _, p := range spplaylists {
		playlists = append(playlists, SimplePlaylist{p})
	}

	l := list.New(playlists, list.NewDefaultDelegate(), 58, 24)
	l.Title = "Playlists and Mixes"
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
			t, _ := m.list.SelectedItem().(spotify.SimplePlaylist)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m CollectionsModel) View() string {
	return collectionsStyle.Render(m.list.View())
}
