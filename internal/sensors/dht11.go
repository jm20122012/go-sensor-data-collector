package sensors

type DHT11Message struct {
	Temp_F         float64 `json:"temp_f"`
	Temp_C         float64 `json:"temp_c"`
	Humidity       float64 `json:"humidity"`
	SensorLocation string  `json:"device_id"`
}

func (d *DHT11Message) ParseData(data DHT11Message) {
	d.Temp_F = data.Temp_F
	d.Temp_C = data.Temp_C
	d.Humidity = data.Humidity
	d.SensorLocation = data.SensorLocation
}
