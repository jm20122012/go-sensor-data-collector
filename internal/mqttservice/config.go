package mqttservice

type MqttConfig struct {
	BaseTopic string `yaml:"base_topic"`
	Server    string `yaml:"server"`
	Port      int    `yaml:"port"`
}
