package devices

type BME280 struct {
	MicrocontrollerType string  `json:"microcontrollerType"`
	TempF               float32 `json:"tempF"`
	TempC               float32 `json:"tempC"`
	Humidity            float32 `json:"humidity"`
	AbsolutePressure    float32 `json:"pressure"`
}
