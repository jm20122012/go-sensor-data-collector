package mqttservice

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type MqttService struct {
	Config MqttConfig
}

func NewMqttService() MqttService {
	config, err := readConfig("mqtt_config.yml")
	if err != nil {
		slog.Error("Error loading mqtt yaml file", "error", err)
	}
	return MqttService{
		Config: *config,
	}
}

func readConfig(filename string) (*MqttConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err)
	}

	var config MqttConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %w", err)
	}

	return &config, nil
}
