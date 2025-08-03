package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DeviceID string         `yaml:"device_id"`
	NATS     NATSConfig     `yaml:"nats"`
	Sensors  []SensorConfig `yaml:"sensors"`
}

type NATSConfig struct {
	URL string `yaml:"url"`
}

type SensorConfig struct {
	ID        string        `yaml:"id"`
	Type      string        `yaml:"type"`
	Frequency time.Duration `yaml:"frequency"`
	Min       float64       `yaml:"min"`
	Max       float64       `yaml:"max"`
	Unit      string        `yaml:"unit"`
	Enabled   bool          `yaml:"enabled"`
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}