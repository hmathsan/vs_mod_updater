package model

import (
	"fmt"
	"log"
)

type (
	ModInfoJson struct {
		Name    string
		Version string
		ModId   string
	}

	ModInfo struct {
		Name             string
		InstalledVersion string
		ModId            string
		ApiModInfo       *ApiModInfo
		Outdated         bool
	}

	ApiModInfoResp struct {
		StatusCode string     `json:"statuscode"`
		Mod        ApiModInfo `json:"mod"`
	}

	ApiModInfo struct {
		Name     string
		ModId    string `json:"modid"`
		Releases []ApiModReleases
	}

	ApiModReleases struct {
		ReleaseId  int    `json:"releaseid"`
		MainFile   string `json:"mainfile"`
		ModVersion string `json:"modversion"`
		Tags       []string
	}
)

func (i ModInfo) Title() string { return i.Name }

func (i ModInfo) Description() string {
	var latestVersion string

	if len(i.ApiModInfo.Releases) <= 0 {
		latestVersion = i.InstalledVersion
	} else {
		latestVersion = i.ApiModInfo.Releases[0].ModVersion
	}

	return fmt.Sprintf("Installed version: v%s | Latest version: v%s", i.InstalledVersion, latestVersion)
}

func NewModInfo(name string, installedVersion string, modId string, apiModInfo *ApiModInfo) ModInfo {
	outdated := false

	log.Println("Name:", name, "installedVersion:", installedVersion, "modId:", modId, "apiInfo:", *apiModInfo)
	if len(apiModInfo.Releases) > 0 && installedVersion != apiModInfo.Releases[0].ModVersion {
		outdated = true
	}

	return ModInfo{name, installedVersion, modId, apiModInfo, outdated}
}
