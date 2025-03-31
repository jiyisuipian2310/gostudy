package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var g_mapDBPort map[string]string

type Service struct {
	Name string `yaml:"name"`
}

type AppConfig struct {
	Services     []Service `yaml:"Services"`
	ListenPort   string    `yaml:"ListenPort"`
	HttpsFlag    bool      `yaml:"HttpsFlag"`
	CAFile       string    `yaml:"CAFile"`
	CertFile     string    `yaml:"CertFile"`
	KeyFile      string    `yaml:"KeyFile"`
	Debug        bool      `yaml:"Debug"`
	VerifyClient bool      `yaml:"VerifyClient"`
}

func init() {
	g_mapDBPort = make(map[string]string)
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

	for i := 0; i < len(config.Services); i++ {
		name := strings.Split(config.Services[i].Name, ":")[0]
		port := strings.Split(config.Services[i].Name, ":")[1]
		g_mapDBPort[name] = port
	}

	logrus.Info(fmt.Sprintf("HttpServer Config, ListenPort: %s", config.ListenPort))
	logrus.Info(fmt.Sprintf("HttpServer Config, HttpsFlag: %t", config.HttpsFlag))
	logrus.Info(fmt.Sprintf("HttpServer Config, VerifyClient: %t", config.VerifyClient))
	logrus.Info(fmt.Sprintf("HttpServer Config, CAFile: %s", config.CAFile))
	logrus.Info(fmt.Sprintf("HttpServer Config, CertFile: %s", config.CertFile))
	logrus.Info(fmt.Sprintf("HttpServer Config, KeyFile: %s", config.KeyFile))

	for k, v := range g_mapDBPort {
		logrus.Info(fmt.Sprintf("Name: %s, Port: %s", k, v))
	}

	return true
}
