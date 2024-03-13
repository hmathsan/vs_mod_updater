package model

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

func NewModInfo(name string, installedVersion string, modId string, apiModInfo *ApiModInfo) ModInfo {
	return ModInfo{name, installedVersion, modId, apiModInfo}
}
