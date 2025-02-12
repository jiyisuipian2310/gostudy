package main

import (
	"fmt"
	"go-plugin-system/plugin"
	"go-plugin-system/plugin/plugins"
)

func main() {
	manager := plugin.NewPluginManager()

	// 预先加载插件
	pluginNames := []string{plugins.SQLInjectionPulginName, plugins.CSSPulginName}
	for _, name := range pluginNames {
		err := manager.LoadPlugin(name)
		if err != nil {
			fmt.Println("加载插件失败:", err)
			return
		}
	}

	// 执行 SQL 注入插件
	result, err := manager.ExecutePlugin(plugins.SQLInjectionPulginName, map[string]interface{}{"url": "http://example.com"})
	if err != nil {
		fmt.Println("执行插件失败:", err)
		return
	}
	fmt.Println("SQL 注入扫描结果:", result)

	// 执行 XSS 漏洞插件
	result, err = manager.ExecutePlugin(plugins.CSSPulginName, map[string]interface{}{"param": "<script>alert(1)</script>"})
	if err != nil {
		fmt.Println("执行插件失败:", err)
		return
	}
	fmt.Println("XSS 扫描结果:", result)
}
