package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func setupFileLogger(logFile string) (*os.File, error) {
	// 创建或打开日志文件（追加模式）
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// 设置日志输出到文件
	log.SetOutput(file)

	return file, nil
}

// 全局变量
var (
	logFile string
)

func initFlags() {
	flag.StringVar(&logFile, "log", "", "日志文件路径，如果指定则输出到文件")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()
}

func main() {
	initFlags()
	if logFile != "" {
		setupFileLogger(logFile)
	}

	config, err := loadConfig("watchDogConfig.json") // 读取 JSON 文件
	if err != nil {
		log.Println("read watchDogConfig.json failed: %v", err)
		return
	}

	manager := NewProgramManager() // 创建程序管理器

	var enableApps = getEnabledApps(config)
	for _, app := range enableApps {
		manager.AddProgram(app.Name, app.Path, strings.Fields(app.Param), true)
	}

	manager.setupSignalHandler() // 设置信号处理

	if err := manager.StartAll(); err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}

	go manager.RestartProcessor()

	// 开始监控
	go manager.Monitor()

	// 等待退出信号
	<-manager.stopCh
	log.Println("程序管理器正常退出")
}
