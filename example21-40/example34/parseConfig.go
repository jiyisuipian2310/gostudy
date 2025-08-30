package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	ListenPort string `yaml:"ListenPort"`
}

func init() {
}

func ParseConfig(config *AppConfig) bool {
	dataBytes, err := os.ReadFile("httpServer.yaml")
	if err != nil {
		logrus.Error(fmt.Sprintf("read config.yaml failed: %s", err.Error()))
		return false
	}

	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		logrus.Error(fmt.Sprintf("parse yaml failed: %s", err.Error()))
		return false
	}

	logrus.Info(fmt.Sprintf("HttpServer Config, ListenPort: %s", config.ListenPort))

	return true
}
