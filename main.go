package main

import (
	"os"
	"path/filepath"
	"vs-mod-updater/tui"
	"vs-mod-updater/tui/constants"
	"vs-mod-updater/utils"
)

func main() {
	constants.Temp = filepath.Join(".", "Temp")
	err := os.MkdirAll(constants.Temp, os.ModePerm)
	utils.Check(err)

	tui.StartTea()
}
