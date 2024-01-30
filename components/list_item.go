package components

import (
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
	"github.com/zmb3/spotify/v2"
)

type FullTrack struct {
	spotify.FullTrack
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
		artist = string([]rune(artist)[:8]) + "..."
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
