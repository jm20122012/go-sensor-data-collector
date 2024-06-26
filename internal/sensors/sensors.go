package sensors

type ProcessedData interface {
	DHT11Message | AqaraMessage | AvtechResponse | AmbientStationResponse
}

type DataParser[T ProcessedData] interface {
	ParseData(data interface{})
}
