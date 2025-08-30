package plugin

import (
	"fmt"
	"reflect"
)

var plugins = make(map[string]Plugin)

// 插件管理器
type PluginManager struct {
}

// 创建新的插件管理器
func NewPluginManager() *PluginManager {
	return &PluginManager{}
}

// 注册插件
func RegisterPlugin(plugin Plugin) {
	pluginType := reflect.TypeOf(plugin).Elem()
	fmt.Printf("pluginType.String()=%s\n", pluginType.String())
	plugins[pluginType.String()] = plugin
}

// 执行插件
func (pm *PluginManager) ExecutePlugin(name string, args map[string]interface{}) (map[string]interface{}, error) {
	plugin, exists := plugins[name]
	if !exists {
		return nil, fmt.Errorf("插件 %s 注册", name)
	}
	return plugin.Execute(args)
}

func init() {
	fmt.Printf("manager.go init ...\n")
}
