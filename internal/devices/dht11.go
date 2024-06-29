package devices

type DHT11Message struct {
	TempF          float64 `json:"temp_f"`
	TempC          float64 `json:"temp_c"`
	Humidity       float64 `json:"humidity"`
	SensorLocation string  `json:"device_id"`
}
