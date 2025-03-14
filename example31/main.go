package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	logrus_stack "github.com/Gurpartap/logrus-stack"
	"github.com/agclqq/goencryption"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var keyStr = "63dTjxISXlwAso0n"
var ivStr = "a1b2c3d4e5f6g7h8"

var g_mapDBPort map[string]string

type Service struct {
	Name string `yaml:"name"`
}

type AppConfig struct {
	Services    []Service `yaml:"Services"`
	ListenPort  string    `yaml:"ListenPort"`
	HttpsFlag   bool      `yaml:"HttpsFlag"`
	CertFile    string    `yaml:"CertFile"`
	KeyFile     string    `yaml:"KeyFile"`
	EncryptFlag bool      `yaml:"EncryptFlag"`
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

	return true
}

type HttpHandler struct {
	EncryptFlag bool
}

func newHttpHandler(EncryptFlag bool) *HttpHandler {
	return &HttpHandler{
		EncryptFlag: EncryptFlag,
	}
}

func (h *HttpHandler) DisConnectionCallback(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if h.EncryptFlag {
		plainData, err := goencryption.EasyDecrypt("aes/cbc/pkcs7/base64", string(body), keyStr, ivStr)
		if err != nil {
			strErrMsg := fmt.Sprintf("EasyDecrypt Failed:%v", err)
			logrus.Error(strErrMsg)
			http.Error(w, strErrMsg, http.StatusBadRequest)
			return
		}

		logrus.Info(fmt.Sprintf("EncryptFlag:true, ReceiveMessage EnctyptData: %s, PlainData: %s", string(body), plainData))
		body = []byte(plainData)
	} else {
		// 打印接收到的消息
		logrus.Info(fmt.Sprintf("EncryptFlag:false, interface Received Message: %s", string(body)))
	}

	for name, port := range g_mapDBPort {
		targetURL := fmt.Sprintf("http://localhost:%s/disConnection", port)
		_, err := http.Post(targetURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			logrus.Warning(fmt.Sprintf("SendMessage to %s Failed, Url: %s", name, targetURL))
		} else {
			logrus.Info(fmt.Sprintf("SendMessage to %s Success, Url: %s", name, targetURL))
		}
	}

	// 返回目标程序的响应给客户端
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/disConnection" {
		h.DisConnectionCallback(w, r)
		return
	}

	strErrMsg := fmt.Sprintf("Don't support url: %s", r.URL.Path)
	logrus.Warning(strErrMsg)
	http.Error(w, strErrMsg, http.StatusBadRequest)
}

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(false)
	stackLevels := []logrus.Level{logrus.PanicLevel, logrus.FatalLevel}
	logrus.AddHook(logrus_stack.NewHook(stackLevels, stackLevels))
}

func main() {
	config := AppConfig{}
	ParseConfig(&config)

	g_mapDBPort = make(map[string]string)
	for i := 0; i < len(config.Services); i++ {
		name := strings.Split(config.Services[i].Name, ":")[0]
		port := strings.Split(config.Services[i].Name, ":")[1]
		g_mapDBPort[name] = port
	}

	logrus.Info(fmt.Sprintf("HttpServer Config, ListenPort: %s, HttpsFlag: %t, EncryptFlag: %t",
		config.ListenPort, config.HttpsFlag, config.EncryptFlag))

	for k, v := range g_mapDBPort {
		logrus.Info(fmt.Sprintf("Name: %s, Port: %s", k, v))
	}

	if !config.HttpsFlag {
		logrus.Info(fmt.Sprintf("Server started on http://0.0.0.0:%s", config.ListenPort))

		server := &http.Server{
			Addr:    ":" + config.ListenPort,
			Handler: newHttpHandler(config.EncryptFlag),
		}

		err := server.ListenAndServe()
		if err != nil {
			logrus.Warning(fmt.Sprintf("Failed to start HTTP server: %v", err.Error()))
		}
	} else {
		logrus.Info(fmt.Sprintf("Server started on https://0.0.0.0:%s", config.ListenPort))

		server := &http.Server{
			Addr: ":" + config.ListenPort,
			TLSConfig: &tls.Config{
				ClientAuth: tls.NoClientCert,
			},
			Handler: newHttpHandler(config.EncryptFlag),
		}

		err := server.ListenAndServeTLS(config.CertFile, config.KeyFile)
		if err != nil {
			logrus.Warning(fmt.Sprintf("Failed to start HTTPS server: %v", err.Error()))
		}
	}
}
