package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	collectionsStyle = lipgloss.NewStyle().Align(lipgloss.Center).Width(65).Height(30).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#5900bf")).PaddingTop(2).PaddingBottom(2)
)

type CollectionsModel struct {
}

func NewCollectionsModel() CollectionsModel {
	return CollectionsModel{}
}

func (m CollectionsModel) Init() tea.Cmd {
	return nil
}

func (m CollectionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m CollectionsModel) View() string {
	return collectionsStyle.Render("COLLECTIONSMODEL")
}
