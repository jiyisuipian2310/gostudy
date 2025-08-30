package main

import (
	"fmt"
	"go-plugin-system/plugin"
	_ "go-plugin-system/plugin/plugins"
)

func main() {
	manager := plugin.NewPluginManager()

	// 执行 SQL 注入插件
	result, err := manager.ExecutePlugin("plugins.SQLInjectionScanner", map[string]interface{}{"url": "http://example.com"})
	if err != nil {
		fmt.Println("执行插件失败:", err)
		return
	}
	fmt.Println("SQL 注入扫描结果:", result)

	// 执行 XSS 漏洞插件
	result, err = manager.ExecutePlugin("plugins.XSSScanner", map[string]interface{}{"param": "<script>alert(1)</script>"})
	if err != nil {
		fmt.Println("执行插件失败:", err)
		return
	}
	fmt.Println("XSS 扫描结果:", result)
}
