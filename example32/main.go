package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"

	logrus_stack "github.com/Gurpartap/logrus-stack"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(false)
	stackLevels := []logrus.Level{logrus.PanicLevel, logrus.FatalLevel}
	logrus.AddHook(logrus_stack.NewHook(stackLevels, stackLevels))
}

func main() {
	config := AppConfig{}
	if !ParseConfig(&config) {
		return
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go processPidAndMasterAccount(&config)
	go startServer(&config)
	wg.Wait()
	logrus.Info(fmt.Sprintf("httpServer Quit!"))
}

func startServer(config *AppConfig) {
	if !config.HttpsFlag {
		StartHttpServer(config)
	} else {
		StartHttpsServer(config)
	}
}

func StartHttpServer(config *AppConfig) {
	logrus.Info(fmt.Sprintf("Remote httpServer started on http://0.0.0.0:%s", config.ListenPort))

	server := &http.Server{
		Addr:    ":" + config.ListenPort,
		Handler: newRemoteHttpHandler(),
	}

	err := server.ListenAndServe()
	if err != nil {
		logrus.Warning(fmt.Sprintf("Failed to start HTTP server: %v", err.Error()))
	}
}

func StartHttpsServer(config *AppConfig) {
	logrus.Info(fmt.Sprintf("Remote httpsServer started on https://0.0.0.0:%s", config.ListenPort))

	server := &http.Server{
		Addr: ":" + config.ListenPort,
		TLSConfig: &tls.Config{
			ClientAuth: tls.NoClientCert,
		},
		Handler: newRemoteHttpHandler(),
	}

	err := server.ListenAndServeTLS(config.CertFile, config.KeyFile)
	if err != nil {
		logrus.Warning(fmt.Sprintf("Failed to start HTTPS server: %v", err.Error()))
	}
}

func processPidAndMasterAccount(config *AppConfig) {
	var localAddr string
	if config.Debug {
		localAddr = ":62020"
		logrus.Info(fmt.Sprintf("local httpServer started on http://0.0.0.0:62020"))
	} else {
		localAddr = "127.0.0.1:62020"
		logrus.Info(fmt.Sprintf("local httpServer started on http://127.0.0.1:62020"))
	}

	server := &http.Server{
		Addr:    localAddr,
		Handler: newLocalHttpHandler(),
	}

	err := server.ListenAndServe()
	if err != nil {
		logrus.Warning(fmt.Sprintf("Failed to start HTTP server: %v", err.Error()))
	}
}
