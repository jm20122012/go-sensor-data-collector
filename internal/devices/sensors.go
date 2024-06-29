package devices

type ProcessedData interface {
	DHT11Message | AqaraTempSensorRawMessage | AvtechResponse | AmbientStationResponse
}

type DataParser[T ProcessedData] interface {
	ParseData(data interface{})
}
