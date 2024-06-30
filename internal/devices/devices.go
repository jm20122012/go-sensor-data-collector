package devices

type Devices struct {
	DeviceName     string `json:"device_name"`
	DeviceLocation string `json:"device_location"`
	DeviceTypeID   int32  `json:"device_type_id"`
	DeviceID       int32  `json:"device_id"`
	DeviceType     string `json:"device_type"`
	MqttTopic      string `json:"mqtt_topic"`
}
