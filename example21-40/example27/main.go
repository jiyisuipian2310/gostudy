package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/astaxie/beego"
)

func GetConfigString(key, defValue string) string {
	val := beego.AppConfig.String(key)
	if val == "" {
		return defValue
	}

	return val
}

func GetConfigInt(key string, defValue int) int {
	val, err := beego.AppConfig.Int(key)
	if err != nil {
		return defValue
	}

	return val
}

func GetConfigBool(key string, defValue bool) bool {
	val, err := beego.AppConfig.Bool(key)
	if err != nil {
		return defValue
	}

	return val
}

func main() {
	if false {
		main1()
	} else {
		main2()
	}
}

func main1() {
	if err := beego.LoadAppConfig("ini", filepath.Join("./", "conf", "config1.conf")); err != nil {
		log.Fatalln("load config file error", err.Error())
	}

	name := GetConfigString("common::name", "zhangsan")
	fmt.Printf("name: %s\n", name)

	age := GetConfigInt("common::age", 90)
	fmt.Printf("age: %d\n", age)

	debug := GetConfigBool("common::debug", false)
	fmt.Printf("debug: %v\n", debug)
}

func main2() {
	if err := beego.LoadAppConfig("ini", filepath.Join("./", "conf", "config2.conf")); err != nil {
		log.Fatalln("load config file error", err.Error())
	}

	name := GetConfigString("name", "zhangsan")
	fmt.Printf("name: %s\n", name)

	age := GetConfigInt("age", 100)
	fmt.Printf("age: %d\n", age)

	debug := GetConfigBool("debug", false)
	fmt.Printf("debug: %v\n", debug)
}
