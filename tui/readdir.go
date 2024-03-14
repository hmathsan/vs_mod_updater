package tui

import (
	"fmt"
	"log"
	"strings"
	"vs-mod-updater/model"
	"vs-mod-updater/readdir"
	"vs-mod-updater/tui/constants"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var ()

type responseMsg struct{}

type ReadDirModel struct {
	fileNames []string
	tempPath  string
	sub       chan model.ModInfo
	spinner   spinner.Model
	modsInfo  []model.ModInfo
	completed bool
}

func InitReadDir() tea.Model {

	s := spinner.New()
	s.Style = constants.SpinnerStyle
	s.Spinner = spinner.Jump

	return ReadDirModel{
		fileNames: readdir.FindInstalledMods(),
		sub:       make(chan model.ModInfo),
		spinner:   s,
		modsInfo:  make([]model.ModInfo, 0),
	}
}

func (m ReadDirModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		listenForActivity(m, m.sub),
		waitForActivities(m.sub),
	)
}

func (m ReadDirModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg.(tea.WindowSizeMsg)
		return m, cmd
	case tea.KeyMsg:
		m.completed = true
		return m, tea.Quit
	case model.ModInfo:
		log.Println("Appending mod information", msg.(model.ModInfo))
		m.modsInfo = append(m.modsInfo, msg.(model.ModInfo))
		if len(m.modsInfo) >= len(m.fileNames) {
			m := InitUpdateMod(m)
			return m.Update(nil)
		} else {
			return m, waitForActivities(m.sub)
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case responseMsg:
		return m, waitForActivities(m.sub)
	default:
		return m, cmd
	}
}

func (m ReadDirModel) View() string {
	doc := strings.Builder{}

	{
		title := lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render("Vintage Story Mod Update"),
			lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderTop(true).
				Align(lipgloss.Center).
				Render(constants.UrlRender("https://github.com/hmathsan/vs_mod_updater")),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, title)

		titleUi := lipgloss.Place(
			constants.WindowSize.Width, 9,
			lipgloss.Center, lipgloss.Center,
			row,
		)
		doc.WriteString(titleUi)
	}

	{
		title := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Looking for installed mods in the Vintage Story directory")
		spinner := fmt.Sprintf("%s Mods found so far: %d", m.spinner.View(), len(m.modsInfo))

		ui := lipgloss.JoinVertical(lipgloss.Center, title, spinner)

		dialog := lipgloss.Place(
			constants.WindowSize.Width, 9,
			lipgloss.Center,
			lipgloss.Center,
			constants.DialogBoxStyle.Render(ui),
		)

		doc.WriteString(dialog)
	}

	constants.DocStyle = constants.DocStyle.MaxWidth(constants.WindowSize.Width)

	return fmt.Sprint(constants.DocStyle.Render(doc.String()))

	// s := fmt.Sprintf("\n %s Looking for installed mods. Mods found so far: %d\n\n Press any key to exit\n", m.spinner.View(), len(m.modsInfo))
	// if m.completed {
	// 	s += "\n"
	// }
	// return s
}

func listenForActivity(m ReadDirModel, sub chan model.ModInfo) tea.Cmd {
	return func() tea.Msg {
		go readdir.FetchModFiles(m.fileNames, &m.tempPath, sub)
		return responseMsg{}
	}
}

func waitForActivities(sub chan model.ModInfo) tea.Cmd {
	return func() tea.Msg {
		return model.ModInfo(<-sub)
	}
}
