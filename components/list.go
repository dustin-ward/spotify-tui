package components

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin-ward/spotify-tui/colours"
	"github.com/dustin-ward/spotify-tui/spotifyapi"
	"github.com/zmb3/spotify/v2"
)

const LIST_WIDTH = 92
const LIST_HEIGHT = 42

var (
	titleStyle        = list.DefaultStyles().Title.Background(colours.PURPLE).Foreground(colours.GREEN).MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(colours.PURPLE)
	playingItemStyle  = lipgloss.NewStyle().PaddingLeft(2).Foreground(colours.GREEN)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	listStyle         = lipgloss.NewStyle().MarginRight(1).PaddingRight(2)
)

var (
	CurrentPage int
	LastPage    bool
)

type ListModel struct {
	list list.Model
}

func NewListModel(uri string) ListModel {
	CurrentContext = (*spotify.URI)(&uri)
	CurrentPage = 1

	sptracks, title, err := spotifyapi.GetTracksFromPlaylist(uri)
	if err != nil {
		log.Fatal(err)
	}
	LastPage = len(sptracks) < spotifyapi.PAGE_SIZE

	tracksItems := make([]list.Item, 0, 1000)
	for i := 0; i < len(sptracks); i++ {
		tracksItems = append(tracksItems, FullTrack{sptracks[i]})
	}

	l := list.New(tracksItems, trackDelegate{}, LIST_WIDTH, LIST_HEIGHT)
	l.Title = fmt.Sprintf("Playlist: %s", title)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return ListModel{list: l}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case UpdatePlaylistMsg:
		uri := string(msg)
		return NewListModel(uri), nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			t, _ := m.list.SelectedItem().(FullTrack)
			return m, PlayPauseTrack(&t.URI)
		}
	}

	if !LastPage && m.list.Paginator.OnLastPage() {
		orig_items := m.list.Items()
		new_items, err := spotifyapi.GetTrackPage(string(*CurrentContext), CurrentPage)
		if err != nil {
			log.Println("unable to fetch tracks")
		} else {
			if len(new_items) < spotifyapi.PAGE_SIZE {
				LastPage = true
			}

			for _, t := range new_items {
				orig_items = append(orig_items, FullTrack{t})
			}

			CurrentPage++
			cmd = m.list.SetItems(orig_items)
			cmds = append(cmds, cmd)
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m ListModel) View() string {
	return listStyle.Render(m.list.View())
}
