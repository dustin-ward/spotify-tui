package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Track struct {
	Name     string
	Artist   string
	Album    string
	Duration string
}

var songs = []list.Item{
	Track{"The End of the World", "Skeeter Davis", "The Essential Skeeter Davis", "2:38"},
	Track{"The Night We Met", "Lord Huron", "Strange Trails", "3:28"},
	Track{"As the World Caves In", "Matt Maltese", "As the World Caves In", "3:39"},
	Track{"Bags", "Clairo", "Immunity", "4:21"},
	Track{"Ruby", "Geskle", "Rose Colored Glasses", "3:22"},
	Track{"Velvet Light", "Jakob", "Velvet Light", "2:22"},
	Track{"affection", "BETWEEN FRIENDS", "we just need some time together", "3:55"},
	Track{"Clueless", "The Marias", "Clueless", "3:47"},
}

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

func (t Track) FilterValue() string { return t.Name }

type trackDelegate struct{}

func (d trackDelegate) Height() int                             { return 1 }
func (d trackDelegate) Spacing() int                            { return 0 }
func (d trackDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d trackDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	t, ok := listItem.(Track)
	if !ok {
		return
	}

	name := t.Name
	if len(name) > 25 {
		name = name[:22] + "..."
	}
	artist := t.Artist
	if len(artist) > 15 {
		artist = artist[:12] + "..."
	}
	album := t.Album
	if len(album) > 25 {
		album = album[:22] + "..."
	}
	str := fmt.Sprintf("%-25s %15s\t%25s %5s", name, artist, album, t.Duration)

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
	choice   *Track
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
			t, ok := m.list.SelectedItem().(Track)
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
		return quitTextStyle.Render(fmt.Sprintf("Now Playing: %s - %s", m.choice.Name, m.choice.Artist))
	}
	if m.quitting {
		return quitTextStyle.Render("Not playing any track.")
	}
	return "\n" + m.list.View()
}

func main() {
	const defaultWidth = 20

	l := list.New(songs, trackDelegate{}, defaultWidth, listHeight)
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
