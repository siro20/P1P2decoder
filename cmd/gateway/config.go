package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Prometheus    PrometheusConfig    `yaml:"prometheus"`
	HomeAssistant HomeAssistantConfig `yaml:"homeassistant"`
	Html          HtmlConfig          `yaml:"html"`
	Serial        SerialConfig        `yaml:"serial"`
}

// NewConfig returns a new decoded Config struct
func ReadConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	fmt.Printf("Looking for config %s\n", configPath)
	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		fmt.Printf("Failed to decode config: %v\n", err)
		return nil, err
	}

	return config, nil
}
