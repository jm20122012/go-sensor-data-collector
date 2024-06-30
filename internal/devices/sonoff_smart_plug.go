package devices

/*
	{
		"linkquality": 33,
		"state": "ON",
		"update": {
			"installed_version": 4352,
			"latest_version": 4099,
			"state": "idle"
		}
	}
*/
type SonoffSmartPlugRawMsg struct {
	LinkQuality int         `json:"linkquality"`
	State       string      `json:"state"`
	Update      UpdateItems `json:"update"`
}

type UpdateItems struct {
	InstalledVersion int    `json:"installed_version"`
	LatestVersion    int    `json:"latest_version"`
	State            string `json:"state"`
}
