package main

import (
	"example8/tcpserver"
	"fmt"
)

func main() {
	tcpServer, err := tcpserver.NewTcpServer()
	if err != nil {
		fmt.Printf("NewTcpServer error:%v\n", err)
		return
	}

	tcpServer.DisplayCfgInfo()

	tcpServer.StartServer()
}
