package components

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin-ward/spotify-tui/colours"
	"github.com/zmb3/spotify/v2"
)

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

type playlistDelegate struct{}

func (d playlistDelegate) Height() int                             { return 2 }
func (d playlistDelegate) Spacing() int                            { return 1 }
func (d playlistDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d playlistDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	p, ok := listItem.(SimplePlaylist)
	if !ok {
		return
	}

	defaultStyles := list.NewDefaultItemStyles()
	titleStyle := defaultStyles.NormalTitle
	descStyle := defaultStyles.NormalDesc

	if p.URI == *CurrentContext {
		titleStyle = defaultStyles.SelectedTitle.Foreground(colours.GREEN).UnsetBorderLeft().PaddingLeft(2)
		descStyle = defaultStyles.SelectedDesc.Foreground(colours.GREEN).UnsetBorderLeft().PaddingLeft(2)
	}
	if index == m.Index() {
		titleStyle = defaultStyles.SelectedTitle.Foreground(colours.PURPLE).BorderLeft(true).BorderForeground(colours.PURPLE).PaddingLeft(1)
		descStyle = defaultStyles.SelectedDesc.Foreground(colours.PURPLE).BorderLeft(true).BorderForeground(colours.PURPLE).PaddingLeft(1)
	}

	desc := p.SimplePlaylist.Description
	if desc == "" {
		desc = p.Owner.DisplayName
	}

	fmt.Fprint(w, fmt.Sprintf("%s\n%s",
		titleStyle.Render(p.Name),
		descStyle.Render(desc),
	))
}
