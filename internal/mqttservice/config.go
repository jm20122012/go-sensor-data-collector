package mqttservice

type MqttConfig struct {
	BaseTopic   string                    `yaml:"base_topic"`
	Server      string                    `yaml:"server"`
	Port        int                       `yaml:"port"`
	DeviceTypes map[string]DeviceTypeData `yaml:"device_types"`
}

type DeviceTypeData struct {
	Topic   string   `yaml:"topic"`
	Devices []string `yaml:"devices"`
}
