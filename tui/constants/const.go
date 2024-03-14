package constants

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	P          *tea.Program
	WindowSize tea.WindowSizeMsg
	Temp       string

	UrlStyle     = lipgloss.AdaptiveColor{Light: "#6da6fc", Dark: "#4088f5"}
	SpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	UrlRender = lipgloss.NewStyle().Foreground(UrlStyle).Render

	DocStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	DialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
)
