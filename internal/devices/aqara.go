package devices

type AqaraTempSensorRawMessage struct {
	Battery          float64
	Humidity         float64
	LinkQuality      float64
	PowerOutageCount int
	Pressure         float64
	Temperature      float64
	Voltage          float64
}
