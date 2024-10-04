package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Interface struct {
	Name     string `yaml:"name"`
	Response string `yaml:"response"`
}

type Config struct {
	Interfaces []Interface `yaml:"interfaces"`
}

func ParseConfig(config *Config) bool {
	dataBytes, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("读取config.yaml文件失败: ", err)
		return false
	}

	//fmt.Println("yaml 文件的内容: \n", string(dataBytes))
	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
		return false
	}

	return true
}

type httpMsgStruct struct {
	content string
	config  *Config
}

func (handler *httpMsgStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	text, _ := io.ReadAll(r.Body)
	fmt.Printf("bodydata: %s\n", text)
	defer r.Body.Close()

	for i := 0; i < len(handler.config.Interfaces); i++ {
		if r.URL.Path == handler.config.Interfaces[i].Name {
			fmt.Fprintf(w, handler.config.Interfaces[i].Response)
			return
		}
	}

	http.Error(w, handler.content, http.StatusNotFound)
}

func main() {
	config := Config{}
	ParseConfig(&config)

	for i := 0; i < len(config.Interfaces); i++ {
		fmt.Printf("name: %s, response: %s\n", config.Interfaces[i].Name, config.Interfaces[i].Response)
	}

	var s httpMsgStruct
	s.content = "404 Not Found"
	s.config = &config
	http.ListenAndServe("localhost:8000", &s)
}
