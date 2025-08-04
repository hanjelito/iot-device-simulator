// Package config provides structures for handling application configuration.
package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level configuration structure, loaded from a YAML file.
// It includes device settings, NATS connection info, and a list of sensor configurations.
type Config struct {
	DeviceID string         `yaml:"device_id"`
	NATS     NATSConfig     `yaml:"nats"`
	Sensors  []SensorConfig `yaml:"sensors"`
}

// NATSConfig holds the configuration for connecting to the NATS server.
type NATSConfig struct {
	URL string `yaml:"url"`
}

// SensorConfig defines the configuration for a single simulated sensor.
// This includes its identity, behavior, and operational parameters.
type SensorConfig struct {
	ID        string        `yaml:"id"`
	Type      string        `yaml:"type"`
	Frequency time.Duration `yaml:"frequency"`
	Min       float64       `yaml:"min"`
	Max       float64       `yaml:"max"`
	Unit      string        `yaml:"unit"`
	Enabled   bool          `yaml:"enabled"`
}

// Load reads a YAML configuration file from the given path and decodes it into a Config struct.
// It returns the populated Config struct or an error if the file cannot be read or parsed.
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