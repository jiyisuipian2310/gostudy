package plugin

import (
	"fmt"
	"reflect"
)

// 插件注册表，存储插件类型
var pluginRegistry = make(map[string]reflect.Type)

// 插件管理器
type PluginManager struct {
	plugins map[string]Plugin
}

// 创建新的插件管理器
func NewPluginManager() *PluginManager {
	return &PluginManager{plugins: make(map[string]Plugin)}
}

// 注册插件
func RegisterPlugin(plugin Plugin) {
	pluginType := reflect.TypeOf(plugin).Elem()
	fmt.Printf("pluginType.String()=%s\n", pluginType.String())
	pluginRegistry[plugin.Name()] = pluginType
}

// 加载插件
func (pm *PluginManager) LoadPlugin(name string) error {
	pluginType, exists := pluginRegistry[name]
	if !exists {
		return fmt.Errorf("插件 %s 未注册", name)
	}
	pluginInstance := reflect.New(pluginType).Interface().(Plugin)
	pm.plugins[name] = pluginInstance
	return nil
}

// 执行插件
func (pm *PluginManager) ExecutePlugin(name string, args map[string]interface{}) (map[string]interface{}, error) {
	plugin, exists := pm.plugins[name]
	if !exists {
		return nil, fmt.Errorf("插件 %s 未加载", name)
	}
	return plugin.Execute(args)
}

func init() {
	fmt.Printf("manager.go init ...\n")
}
