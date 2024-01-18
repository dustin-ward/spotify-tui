package components

import (
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

const DEFAULT_WIDTH = 20
const LIST_HEIGHT = 40

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("46"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	listStyle         = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("46"))
)

type ListModel struct {
	list list.Model
}

func NewListModel(sptracks []spotify.SavedTrack, title string) ListModel {
	tracksItems := make([]list.Item, 0, 1000)
	for i := 0; i < len(sptracks); i++ {
		tracksItems = append(tracksItems, SavedTrack{sptracks[i]})
	}

	l := list.New(tracksItems, trackDelegate{}, DEFAULT_WIDTH, LIST_HEIGHT)
	l.Title = title
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return ListModel{list: l}
}

type SavedTrack struct {
	spotify.SavedTrack
}

func (t SavedTrack) FilterValue() string {
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
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			t, _ := m.list.SelectedItem().(SavedTrack)
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
