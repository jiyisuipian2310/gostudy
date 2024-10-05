package main

import (
	"example7/config"
	"fmt"
)

/*
这个示例很有代表性：
	执行 go run .\main.go， 提示错误：./main.go:16:9: undefined: config.loadConfig

已经确认包路径导入没有问题，loadConfig 函数也已经定义，但还是提示 undefined。
解决方法：将 loadConfig 函数的首字母改为大写
*/

func main() {
	fmt.Println("main Hello, world!")

	config.LoadConfig()
}
