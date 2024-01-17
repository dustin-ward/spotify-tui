package main

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

const DEFAULT_WIDTH = 20
const LIST_HEIGHT = 40

func newListModel(sptracks []spotify.SavedTrack) list.Model {
	tracksItems := make([]list.Item, 0, 1000)
	for i := 0; i < len(sptracks); i++ {
		tracksItems = append(tracksItems, SavedTrack{sptracks[i]})
	}

	l := list.New(tracksItems, trackDelegate{}, DEFAULT_WIDTH, LIST_HEIGHT)
	l.Title = "Liked Songs"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return l
}

type SavedTrack struct {
	spotify.SavedTrack
}

func (t SavedTrack) FilterValue() string { return t.Name }

type trackDelegate struct{}

func (d trackDelegate) Height() int                             { return 1 }
func (d trackDelegate) Spacing() int                            { return 0 }
func (d trackDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d trackDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	t, ok := listItem.(SavedTrack)
	if !ok {
		return
	}

	name := t.Name
	if utf8.RuneCountInString(name) > 25 {
		name = string([]rune(name)[:22]) + "..."
	}
	artist := t.Artists[0].Name
	for i := 1; i < len(t.Artists); i++ {
		artist = fmt.Sprintf("%s, %s", artist, t.Artists[i].Name)
	}
	if utf8.RuneCountInString(artist) > 15 {
		artist = string([]rune(artist)[:12]) + "..."
	}
	album := t.Album.Name
	if utf8.RuneCountInString(album) > 25 {
		album = string([]rune(album)[:22]) + "..."
	}
	str := fmt.Sprintf("| %-25s | %15s | %25s | %5s |", name, artist, album, fmtDuration(t.TimeDuration()))

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
