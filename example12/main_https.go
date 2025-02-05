package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type GlobalConfig struct {
	CertFile string `json:"cert_file" mapstructure:"cert_file"`
	KeyFile  string `json:"key_file" mapstructure:"key_file"`
}

type BodyResponseConfig struct {
	Body     string `json:"body" mapstructure:"body"`
	Response string `json:"response" mapstructure:"response"`
}

type InterfaceConfig struct {
	Name           string               `json:"name" mapstructure:"name"`
	BodyResponse   []BodyResponseConfig `json:"body_response" mapstructure:"body_response"`
	GlobalResponse string               `json:"global_response" mapstructure:"global_response"`
}

type ServiceConfig struct {
	ListenPort int               `json:"listen_port" mapstructure:"listen_port"`
	HttpsFlag  bool              `json:"https_flag" mapstructure:"https_flag"`
	DelayTime  int               `json:"delay_time" mapstructure:"delay_time"`
	Interfaces []InterfaceConfig `json:"interfaces" mapstructure:"interfaces"`
}

type appConfig struct {
	Global   GlobalConfig    `json:"global" mapstructure:"global"`
	Services []ServiceConfig `json:"services" mapstructure:"services"`
}

var (
	appCfg appConfig
)

type delayHandler struct {
	delayTime  int
	Interfaces []InterfaceConfig
}

func (h *delayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.delayTime > 0 {
		time.Sleep(time.Duration(h.delayTime) * time.Second)
	}

	bodydata, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	
	fmt.Printf("bodydata: %s\n", bodydata)

	for _, interfaceCfg := range h.Interfaces {
		if interfaceCfg.Name == r.URL.Path {
			if len(interfaceCfg.BodyResponse) > 0 {
				for _, bodyResponseCfg := range interfaceCfg.BodyResponse {
					if string(bodydata) == bodyResponseCfg.Body {
						fmt.Fprintf(w, "%s\n", bodyResponseCfg.Response)
						return
					}
				}
			}
			fmt.Fprintf(w, "%s\n", interfaceCfg.GlobalResponse)
			return
		}
	}
}

func LoadConfig() error {
	var err error
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.SetConfigName("httpConfig")
	if err = viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			fmt.Printf("No configuration file found\n")
			return err
		} else {
			fmt.Printf("error loading configuration file: %v", err)
			return err
		}
	}

	err = viper.Unmarshal(&appCfg)
	if err != nil {
		fmt.Printf("error parsing configuration file: %v", err)
		return err
	}

	return nil
}

func ShowConfig() {
	fmt.Println("=============================cert==============================")
	fmt.Println("certFile: ", appCfg.Global.CertFile)
	fmt.Println("keyFile: ", appCfg.Global.KeyFile)

	for _, service := range appCfg.Services {
		fmt.Println("\n=============================service=============================")
		fmt.Println("listenPort: ", service.ListenPort)
		fmt.Println("httpsFlag: ", service.HttpsFlag)
		fmt.Println("delayTime: ", service.DelayTime)
		for _, interfaceCfg := range service.Interfaces {
			fmt.Println("\n***************interface***************")
			fmt.Println("    name: ", interfaceCfg.Name)
			fmt.Println("    global_response: ", interfaceCfg.GlobalResponse)
		}
	}
}

func startHTTPServer(service ServiceConfig) {
	var err error

	server := &http.Server{
		Addr: ":" + strconv.Itoa(service.ListenPort),
		Handler: &delayHandler{
			delayTime:  service.DelayTime,
			Interfaces: service.Interfaces,
		},
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func startHTTPSServer(service ServiceConfig, certFile, keyFile string) {
	var err error

	tlsConfig := &tls.Config{
		// 不要求客户端提供证书
		ClientAuth: tls.NoClientCert,
	}

	server := &http.Server{
		Addr:      ":" + strconv.Itoa(service.ListenPort),
		TLSConfig: tlsConfig,
		Handler: &delayHandler{
			delayTime:  service.DelayTime,
			Interfaces: service.Interfaces,
		},
	}

	err = server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to start HTTPS server: %v", err)
	}
}

func main() {
	err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}

	ShowConfig()

	for _, service := range appCfg.Services {
		if service.ListenPort == 0 {
			continue
		}
		if service.HttpsFlag {
			go startHTTPSServer(service, appCfg.Global.CertFile, appCfg.Global.KeyFile)
		} else {
			go startHTTPServer(service)
		}
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
