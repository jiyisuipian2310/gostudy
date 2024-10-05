package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenAddr  string `json:"listenaddr"`
	IdleTimeout int    `json:"idletimeout"`
}

func LoadConfig(path string) (*Config, error) {
	var c Config
	defaultConfig(&c)

	dataBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dataBytes, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func defaultConfig(config *Config) {
	config.ListenAddr = "0.0.0.0:9000"
}
