package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Config 主配置结构体
type Config struct {
	Apps []AppConfig `json:"Apps"`
}

// AppConfig 应用配置结构体
type AppConfig struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Param  string `json:"param"`
	Enable bool   `json:"enable"`
}

// loadConfig 加载和解析配置文件
func loadConfig(filename string) (*Config, error) {
	// 读取文件内容
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 解析 JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v", err)
	}

	printConfig(&config)

	return &config, nil
}

// printConfig 打印配置信息
func printConfig(config *Config) {
	var info = fmt.Sprintf("总共配置了 %d 个应用:\n", len(config.Apps))
	for i, app := range config.Apps {
		info = info + fmt.Sprintf("应用%d:\n\t名称: %s\n\t路径: %s\n\t参数: %s\n\t启用: %v\n",
			i+1, app.Name, app.Path, app.Param, app.Enable)
	}
	log.Println(info)
}

// getEnabledApps 获取所有启用的应用
func getEnabledApps(config *Config) []AppConfig {
	var enabledApps []AppConfig
	for _, app := range config.Apps {
		if app.Enable {
			enabledApps = append(enabledApps, app)
		}
	}
	return enabledApps
}

// findAppByName 根据名称查找应用
func findAppByName(config *Config, name string) *AppConfig {
	for _, app := range config.Apps {
		if app.Name == name {
			return &app
		}
	}
	return nil
}
