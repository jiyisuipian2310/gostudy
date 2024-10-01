package main

/*
1. 此示例展示了如何使用json包来解析json数据，并将其转换为结构体。
2. 此示例展示里如何发送http和https请求，并获取响应数据，以及如何解析json数据。
3. 此示例展示了如何对发送和接收的数据进行加密和解密
*/

import (
	"encoding/json"
	"example1/httpapi"
	"fmt"

	"github.com/agclqq/goencryption"
	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"

	"os"
)

type Student struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

type ResourceInfo struct {
	ResourceIp       string `json:"resourceIp"`
	ResourcePort     string `json:"resourcePort"`
	ResourceAccount  string `json:"resourceAccount"`
	ResourcePassword string `json:"resourcePassword"`
}

type Global struct {
	SendHttpMsg bool
}

type Mysql struct {
	Url  string
	Port int
}

type Redis struct {
	Host string
	Port int
}

type Config struct {
	Global Global `json:"global"`
	Mysql Mysql   `json:"mysql"`
	Redis Redis   `json:"redis"`
}

/*定义加解密的key值和向量值为常量*/
const g_strkeyStr = "63dTjxISXlwAso0n"
const g_strivStr = "a1b2c3d4e5f6g7h8"

func main() {
	dataBytes, err := os.ReadFile("test.yaml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return
	}

	fmt.Println("yaml 文件的内容: \n", string(dataBytes))

	config := Config{}
	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
		return
	}
	fmt.Printf("mysql Url config -> %+v\n", config.Mysql.Url)
	fmt.Printf("mysql Port config -> %+v\n", config.Mysql.Port)
	fmt.Printf("redis Host config -> %+v\n", config.Redis.Host)
	fmt.Printf("redis Port config -> %+v\n", config.Redis.Port)
	
	fmt.Printf("global config -> %+v\n", config.Global.SendHttpMsg)

	var stuobj Student
	stuobj.Name = "yull"
	stuobj.Age = 25
	stuobj.Gender = "male"

	stuobjJson, err := json.Marshal(stuobj)
	if err != nil {
		return
	}

	logrus.Info(fmt.Sprintf("plain data: %s", stuobjJson))

	/*加密数据*/
	encryptstring, err := goencryption.EasyEncrypt("aes/cbc/pkcs7/base64", string(stuobjJson), g_strkeyStr, g_strivStr)
	logrus.Info(fmt.Sprintf("encrypt data: %s", encryptstring))

	/*解密数据*/
	goencryption.EasyDecrypt("aes/cbc/pkcs7/base64", encryptstring, g_strkeyStr, g_strivStr)
	logrus.Info(fmt.Sprintf("decrypt data: %s", string(stuobjJson)))

	if config.Global.SendHttpMsg == false {
		return 
	}

	var httpmethod int = 1
	var responsedata_enc string
	if httpmethod == 1 {
		var httpaddr string = "https://192.168.104.100:12345/go/gettokeninfo"
		responsedata_enc, err = httpapi.SendAndRecvHttpsPostMsg(httpaddr, encryptstring)
	} else {
		var httpaddr string = "http://192.168.104.100:12346/go/gettokeninfo"
		responsedata_enc, err = httpapi.SendAndRecvHttpPostMsg(httpaddr, encryptstring)
	}

	if err != nil {
		logrus.Error(fmt.Sprintf("SendAndRecvHttpsPostMsg error: %s", err.Error()))
		return
	}

	logrus.Info(fmt.Sprintf("Responsedata data: %s", responsedata_enc))

	responsedata_dec, err := goencryption.EasyDecrypt("aes/cbc/pkcs7/base64", responsedata_enc, g_strkeyStr, g_strivStr)
	if err != nil {
		logrus.Error(fmt.Sprintf("EasyDecrypt error: %s", err.Error()))
		return
	}

	//responsedata_dec := RtKOL5b1lf3dkwRvRgwhkIJ/dGH45r/n+HqRVQiutiQy8TgCbbApsx1GU4YDc1WfE7gd8FFLfnsdpL9ffZDOiiqfVhJl1TuzkTFESFbCwA2Swtatn0uEMiv3waXGlroCD39Cv1OEMUb54dFvq0JdlIvlO+S/CN/+JyFRqLOPhSY=
	logrus.Info(fmt.Sprintf("Responsedata data: %s", responsedata_dec))

	var resInfo = new(ResourceInfo)
	json.Unmarshal([]byte(responsedata_dec), &resInfo)

	logrus.Info(fmt.Sprintf("ResourceIp: %s", resInfo.ResourceIp))
	logrus.Info(fmt.Sprintf("ResourcePort: %s", resInfo.ResourcePort))
	logrus.Info(fmt.Sprintf("ResourceUser: %s", resInfo.ResourceAccount))
	logrus.Info(fmt.Sprintf("ResourcePwd: %s", resInfo.ResourcePassword))
}
