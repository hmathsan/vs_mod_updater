package tui

import (
	"vs-mod-updater/model"
	"vs-mod-updater/tui/constants"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type listItem struct {
	title, description string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.description }
func (i listItem) FilterValue() string { return i.title }

type UpdateModsModel struct {
	upToDateMods     []model.ModInfo
	outdatedMods     []model.ModInfo
	upToDateModsList list.Model
	outdatesModsList list.Model
}

func InitUpdateMod(m ReadDirModel) *UpdateModsModel {
	upToDateMods := make([]model.ModInfo, 0)
	outdatedMods := make([]model.ModInfo, 0)

	for _, modInfo := range m.modsInfo {
		if modInfo.Outdated {
			outdatedMods = append(outdatedMods, modInfo)
		} else {
			upToDateMods = append(upToDateMods, modInfo)
		}
	}

	uListItems := make([]list.Item, len(upToDateMods))
	oListItems := make([]list.Item, len(outdatedMods))

	for i, item := range upToDateMods {
		uListItems[i] = listItem{title: item.Title(), description: item.Description()}
	}

	for i, item := range outdatedMods {
		oListItems[i] = listItem{title: item.Title(), description: item.Description()}
	}

	uList := list.New(uListItems, list.NewDefaultDelegate(), constants.WindowSize.Width/2, constants.WindowSize.Height/2)
	oList := list.New(oListItems, list.NewDefaultDelegate(), constants.WindowSize.Width/2, constants.WindowSize.Height/2)

	return &UpdateModsModel{upToDateMods, outdatedMods, uList, oList}
}

func (m UpdateModsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m UpdateModsModel) View() string {
	return "test"
}

func (m UpdateModsModel) Init() tea.Cmd {
	return nil
}
