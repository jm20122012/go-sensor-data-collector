package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type OuterConfig struct {
	Config Config
}

type Config struct {
	DebugLevel string `yaml:"debug_level"`
	MqttBroker string `yaml:"mqtt_broker"`
	MqttPort   int    `yaml:"mqtt_port"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err)
	}

	var cfg OuterConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %w", err)
	}

	return &cfg.Config, nil
}
