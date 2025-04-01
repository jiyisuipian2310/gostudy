package main

import (
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
	wg.Add(1)
	go StartHttpServer(&config)
	wg.Wait()
	logrus.Info(fmt.Sprintf("httpServer Quit!"))
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
