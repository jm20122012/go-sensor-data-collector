package devices

type BMP280 struct {
	MicrocontrollerType string  `json:"microcontrollerType"`
	TempF               float32 `json:"tempF"`
	TempC               float32 `json:"tempC"`
	AbsolutePressure    float32 `json:"pressure"`
}
