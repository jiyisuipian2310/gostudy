package tcpserver

import (
	"bufio"
	"example8/config"
	"fmt"
	"io"
	"net"
)

type clientHandler struct {
	conn   net.Conn
	cfg    *config.Config
	writer *bufio.Writer
	reader *bufio.Reader
	srcIP  string
	err    error
}

func newClientHandler(connection net.Conn, cfg *config.Config) *clientHandler {
	client := &clientHandler{
		conn:   connection,
		cfg:    cfg,
		writer: bufio.NewWriter(connection),
		reader: bufio.NewReader(connection),
		srcIP:  connection.RemoteAddr().String(),
	}

	return client
}

func (client *clientHandler) DisplayClientInfo() {
	fmt.Printf("Receive connection from client Addr: %s\n", client.srcIP)
}

func (client *clientHandler) ProcessMessage() {
	defer func() {
		if client.err == io.EOF {
			fmt.Printf("客户端 %s 关闭连接\n", client.srcIP)
		}
		client.conn.Close()
	}()

	for {
		var buf [4096]byte
		n, err := client.reader.Read(buf[:]) // 读取数据
		client.err = err
		if err != nil {
			if err != io.EOF {
				fmt.Println("read from client failed, err: ", err)
			}
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("收到Client端发来的数据:", recvStr)
		client.conn.Write([]byte(recvStr)) // 发送数据
	}
}
