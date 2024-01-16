package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dustin-ward/spotify-tui/spotifyapi"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

type SavedTrack struct {
	spotify.SavedTrack
}

const listHeight = 20

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

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
	if len(name) > 25 {
		name = name[:22] + "..."
	}
	artist := t.Artists[0].Name
	for i := 1; i < len(t.Artists); i++ {
		artist = fmt.Sprintf("%s, %s", artist, t.Artists[i].Name)
	}
	if len(artist) > 15 {
		artist = artist[:12] + "..."
	}
	album := t.Album.Name
	if len(album) > 25 {
		album = album[:22] + "..."
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

type model struct {
	list     list.Model
	choice   *SavedTrack
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != nil {
		return quitTextStyle.Render(fmt.Sprintf("Now Playing: %s - %s", m.choice.Name, m.choice.Artists[0].Name))
	}
	if m.quitting {
		return quitTextStyle.Render("Not playing any track.")
	}
	return "\n" + m.list.View()
}

func main() {
	err := spotifyapi.Login()
	if err != nil {
		log.Fatal(err)
	}

	sptracks, err := spotifyapi.GetLikedTracks()
	if err != nil {
		log.Fatal(err)
	}

	tracksItems := make([]list.Item, 0, 1000)
	for i := 0; i < len(sptracks); i++ {
		tracksItems = append(tracksItems, SavedTrack{sptracks[i]})
	}

	const defaultWidth = 20

	l := list.New(tracksItems, trackDelegate{}, defaultWidth, listHeight)
	l.Title = "Liked Songs"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
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
