package sensors

type AvtechResponse struct {
	Sensor []SensorData `json:"sensor"`
}

type SensorData struct {
	Label string `json:"label"`
	TempF string `json:"tempf"`
	TempC string `json:"tempc"`
	HighF string `json:"highf"`
	HighC string `json:"highc"`
	LowF  string `json:"lowf"`
	LowC  string `json:"lowc"`
}
