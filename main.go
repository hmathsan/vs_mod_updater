package main

import (
	"fmt"
	"os"
	"path/filepath"
	"vs-mod-updater/model"
	"vs-mod-updater/utils"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Stage int64

const (
	ReadDirectory Stage = iota
)

type responseMsg struct{}

type readDirModel struct {
	fileNames []string
	tempPath  string
	sub       chan model.ModInfo
	responses int
	spinner   spinner.Model
	completed bool
}

func main() {
	tempPath := filepath.Join(".", "Temp")
	err := os.MkdirAll(tempPath, os.ModePerm)
	utils.Check(err)

	m := readDirModel{
		fileNames: findInstalledMods(),
		sub:       make(chan model.ModInfo),
		spinner:   spinner.New(),
	}

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}

func (m readDirModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		listenForActivity(m, m.sub),
		waitForActivities(m.sub),
	)
}

func (m readDirModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.completed = true
		return m, tea.Quit
	case model.ModInfo:
		m.responses++
		if m.responses >= len(m.fileNames) {
			close(m.sub)
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

func (m readDirModel) View() string {
	s := fmt.Sprintf("\n %s Looking for installed mods. Mods found so far: %d\n\n Press any key to exit\n", m.spinner.View(), m.responses)
	if m.completed {
		s += "\n"
	}
	return s
}

func listenForActivity(m readDirModel, sub chan model.ModInfo) tea.Cmd {
	return func() tea.Msg {
		go fetchModFiles(m.fileNames, &m.tempPath, sub)
		return responseMsg{}
	}
}

func waitForActivities(sub chan model.ModInfo) tea.Cmd {
	return func() tea.Msg {
		return model.ModInfo(<-sub)
	}
}
