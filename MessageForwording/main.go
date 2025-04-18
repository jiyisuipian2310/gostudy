package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// 配置结构体
type Config struct {
	ListenPort int
	TargetAddr string
}

// 加载配置
func loadConfig() Config {
	// 默认配置
	config := Config{
		ListenPort: 8080,
		TargetAddr: "192.168.100.5:10000",
	}

	return config
}

// 处理连接
func handleConnection(conn net.Conn, TargetAddr string) {
	defer conn.Close()

	// 连接到目标端口
	targetConn, err := net.Dial("tcp", TargetAddr)
	if err != nil {
		log.Printf("无法连接到目标地址 %s: %v", TargetAddr, err)
		return
	}
	defer targetConn.Close()

	// 启动goroutine从客户端读取并转发到目标
	go func() {
		_, err := io.Copy(targetConn, conn)
		if err != nil {
			log.Printf("转发到目标时出错: %v", err)
		}
	}()

	// 从目标读取并返回给客户端
	_, err = io.Copy(conn, targetConn)
	if err != nil {
		log.Printf("从目标读取时出错: %v", err)
	}
}

// 启动转发服务
func startForwarding(listenPort int, targetAddr string) {
	listenAddr := fmt.Sprintf("0.0.0.0:%d", listenPort)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("无法监听端口 %d: %v", listenPort, err)
	}
	defer listener.Close()

	log.Printf("服务启动，监听端口 %d, 转发到本地地址 %s", listenPort, targetAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("接受连接时出错: %v", err)
			continue
		}

		remoteAddr := conn.RemoteAddr()
		fmt.Printf("Accepted connection from: %s, Forward to address: %s\n", remoteAddr.String(), targetAddr)

		go handleConnection(conn, targetAddr)
	}
}

// 自定义帮助信息
func Help() {
	flag.PrintDefaults()
}

func main() {
	config := Config{}

	flag.IntVar(&config.ListenPort, "listenport", 0, "本地监听端口(9000)")
	flag.StringVar(&config.TargetAddr, "targetaddr", "", "目标地址(127.0.0.1:8080)")
	flag.Usage = Help

	// 解析命令行参数
	flag.Parse()

	if config.ListenPort == 0 || config.TargetAddr == "" {
		flag.Usage()
		os.Exit(1)
	}

	// 启动服务
	startForwarding(config.ListenPort, config.TargetAddr)
}
