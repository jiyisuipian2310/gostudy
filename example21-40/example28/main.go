package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

func handleConn(clientconn net.Conn) {
	defer func() {
		clientconn.Close()
	}()

	fmt.Printf("Receive client %s connect\n", clientconn.RemoteAddr().String())

	reader := bufio.NewReader(clientconn)
	writer := bufio.NewWriter(clientconn)

	//make([]byte, 10) <==> make([]byte, 10, 10)
	clientdata := make([]byte, 10)
	fmt.Printf("len(clientdata): %d, cap(clientdata): %d\n", len(clientdata), cap(clientdata))

	//设置5秒读超时
	clientconn.SetReadDeadline(time.Now().Add(5 * time.Second))

	//至少读取 10 个字节的数据， 如果5秒后读不到10个字节，超时退出
	_, err := io.ReadFull(reader, clientdata)
	if err != nil {
		fmt.Printf("ReadFull error: %s\n", err.Error())
		return
	}

	receivedData := string(clientdata)
	fmt.Printf("Received data: %s\n", receivedData)

	// 在接收到的数据后面添加 '_hello'
	sendData := receivedData + "_hello"

	// 发送修改后的数据回客户端
	_, err = writer.Write([]byte(sendData))
	writer.Flush()
}

func main() {
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error: ", err)
		return
	}

	defer listener.Close()

	for {
		netConn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}

		conn := netConn.(*net.TCPConn)
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(30)

		go handleConn(netConn)
	}
}
