package readdir

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"vs-mod-updater/model"
	"vs-mod-updater/utils"
)

const apiBaseUrl = "https://mods.vintagestory.at/api/mod/"
const modsDirStr = "C:/Users/mathe/Documentos/mods-folder-test/Mods"

func FindInstalledMods() []string {

	modsDir, err := os.ReadDir(modsDirStr)
	utils.Check(err)

	fileNames := make([]string, 0)

	log.Println("Looking for installed mods.")
	for _, entry := range modsDir {
		fileName := entry.Name()
		fileExtension := fileName[len(fileName)-3:]

		if fileExtension != "zip" {
			continue
		} else {
			fileNames = append(fileNames, fileName)
		}
	}

	log.Printf("Found a total of %d mods installed\n", len(fileNames))

	return fileNames
}

func FetchModFiles(fileNames []string, tempPath *string, sub chan model.ModInfo) {
	defer close(sub)
	var infos = []chan model.ModInfo{}

	for i, fileName := range fileNames {
		infos = append(infos, make(chan model.ModInfo))
		go unzipModFilesParallel(fileName, tempPath, infos[i])
	}

	for i := range infos {
		for info1 := range infos[i] {
			sub <- info1
		}
	}
}

func unzipModFilesParallel(fileName string, tempPath *string, mInfoChan chan model.ModInfo) model.ModInfo {
	defer close(mInfoChan)

	archive, err := zip.OpenReader(modsDirStr + "/" + fileName)
	utils.Check(err)
	defer archive.Close()

	var m model.ModInfo

	for _, f := range archive.File {
		if f.Name != "modinfo.json" {
			continue
		}

		modName := fileName[:len(fileName)-4]
		modTempPath := *tempPath + "/" + modName

		err := os.MkdirAll(modTempPath, os.ModePerm)
		utils.Check(err)

		destinationFile, err := os.OpenFile(modTempPath+"/"+f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		utils.Check(err)
		defer destinationFile.Close()

		zippedFile, err := f.Open()
		utils.Check(err)
		defer zippedFile.Close()

		log.Println("unzipping modinfo.json from", fileName)

		_, err = io.Copy(destinationFile, zippedFile)
		utils.Check(err)

		jsonFile, err := os.ReadFile(modTempPath + "/" + f.Name)
		utils.Check(err)

		var modInfo model.ModInfoJson
		json.Unmarshal([]byte(jsonFile), &modInfo)

		apiModInfoResp := GetLatestVersionFromApi(modInfo.ModId)
		mInfoChan <- model.NewModInfo(modInfo.Name, modInfo.Version, modInfo.ModId, apiModInfoResp)

		break
	}

	return m
}

func GetLatestVersionFromApi(modId string) *model.ApiModInfo {
	resp, err := http.Get(apiBaseUrl + modId)
	utils.Check(err)

	body, err := io.ReadAll(resp.Body)
	utils.Check(err)

	var apiModInfoResp model.ApiModInfoResp
	json.Unmarshal(body, &apiModInfoResp)

	return &apiModInfoResp.Mod
}
