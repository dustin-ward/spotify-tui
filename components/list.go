package components

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin-ward/spotify-tui/spotifyapi"
	"github.com/mattn/go-runewidth"
	"github.com/zmb3/spotify/v2"
)

const LIST_WIDTH = 92
const LIST_HEIGHT = 42

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#5900bf"))
	playingItemStyle  = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#00bf06"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	listStyle         = lipgloss.NewStyle().MarginRight(1).PaddingRight(2)
)

type ListModel struct {
	list list.Model
}

type FullTrack struct {
	spotify.FullTrack
}

func NewListModel(uri string) ListModel {
	CurrentContext = (*spotify.URI)(&uri)

	sptracks, title, err := spotifyapi.GetTracks(uri)
	if err != nil {
		log.Fatal(err)
	}

	tracksItems := make([]list.Item, 0, 1000)
	for i := 0; i < len(sptracks); i++ {
		tracksItems = append(tracksItems, FullTrack{sptracks[i]})
	}

	l := list.New(tracksItems, trackDelegate{}, LIST_WIDTH, LIST_HEIGHT)
	l.Title = title
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return ListModel{list: l}
}

func (t FullTrack) FilterValue() string {
	ret := t.Name
	for _, a := range t.Artists {
		ret = fmt.Sprintf("%s %s", ret, a.Name)
	}
	ret += " " + t.Album.Name
	return ret
}

type trackDelegate struct{}

func (d trackDelegate) Height() int                             { return 1 }
func (d trackDelegate) Spacing() int                            { return 0 }
func (d trackDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d trackDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	t, ok := listItem.(FullTrack)
	if !ok {
		return
	}

	name := t.Name
	nameWidth := runewidth.StringWidth(name)
	nameRunes := utf8.RuneCountInString(name)
	if nameWidth > 25 {
		name = string([]rune(name)[:20]) + "..."
	}
	artist := t.Artists[0].Name
	for i := 1; i < len(t.Artists); i++ {
		artist = fmt.Sprintf("%s, %s", artist, t.Artists[i].Name)
	}
	artistWidth := runewidth.StringWidth(artist)
	artistRunes := utf8.RuneCountInString(artist)
	if artistWidth > 15 {
		artist = string([]rune(artist)[:10]) + "..."
	}
	album := t.Album.Name
	albumWidth := runewidth.StringWidth(album)
	albumRunes := utf8.RuneCountInString(album)
	if albumWidth > 25 {
		album = string([]rune(album)[:20]) + "..."
	}
	fmtString := fmt.Sprintf("| %%-%ds | %%%ds | %%%ds | %%%ds |",
		25-(nameWidth-nameRunes),
		15-(artistWidth-artistRunes),
		25-(albumWidth-albumRunes),
		5,
	)
	str := fmt.Sprintf(fmtString, name, artist, album, fmtDuration(t.TimeDuration()))

	isCurrent := false
	fn := itemStyle.Render
	if IsPlaying && CurrentlyPlaying != nil && t.ID == CurrentlyPlaying.ID {
		isCurrent = true
		fn = func(s ...string) string {
			return playingItemStyle.Render("▶ " + strings.Join(s, " "))
		}
	}
	if index == m.Index() {
		fn = func(s ...string) string {
			arrow := "> "
			if isCurrent {
				arrow = "▶ "
			}
			return selectedItemStyle.Render(arrow + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func fmtDuration(t time.Duration) string {
	s := int(t.Seconds())
	m := s / 60
	s %= 60
	return fmt.Sprintf("%d:%02d", m, s)
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	return listStyle.Render(m.list.View())
}
