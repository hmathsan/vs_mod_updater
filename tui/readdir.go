package tui

import (
	"fmt"
	"log"
	"vs-mod-updater/model"
	"vs-mod-updater/readdir"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

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
	return ReadDirModel{
		fileNames: readdir.FindInstalledMods(),
		sub:       make(chan model.ModInfo),
		spinner:   spinner.New(),
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
	switch msg.(type) {
	case tea.KeyMsg:
		m.completed = true
		return m, tea.Quit
	case model.ModInfo:
		log.Println("Appending mod information", msg.(model.ModInfo))
		m.modsInfo = append(m.modsInfo, msg.(model.ModInfo))
		if len(m.modsInfo) >= len(m.fileNames) {
			return m, tea.Quit
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
		return m, nil
	}
}

func (m ReadDirModel) View() string {
	s := fmt.Sprintf("\n %s Looking for installed mods. Mods found so far: %d\n\n Press any key to exit\n", m.spinner.View(), len(m.modsInfo))
	if m.completed {
		s += "\n"
	}
	return s
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
