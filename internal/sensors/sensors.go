package sensors

type RawData interface {
	DHT11Message | AqaraMessage | AvtechResponse | AmbientStationResponse
}

type DataParser[T RawData] interface {
	ParseData(data T)
}
