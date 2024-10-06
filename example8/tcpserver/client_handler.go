package tcpserver

import (
	"bufio"
	"example8/config"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type clientHandler struct {
	conn       net.Conn
	cfg        *config.Config
	writer     *bufio.Writer
	reader     *bufio.Reader
	middleware middleware
	srcIP      string
	err        error
	command    string
	param      string
	mutex      *sync.Mutex
}

func newClientHandler(connection net.Conn, cfg *config.Config, m middleware) *clientHandler {
	client := &clientHandler{
		conn:       connection,
		cfg:        cfg,
		middleware: m,
		writer:     bufio.NewWriter(connection),
		reader:     bufio.NewReader(connection),
		srcIP:      connection.RemoteAddr().String(),
		mutex:      &sync.Mutex{},
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
		} else {
			fmt.Println("read from client failed, err: ", client.err)
		}
		client.conn.Close()
	}()

	for {
		var buf [4096]byte
		n, err := client.reader.Read(buf[:]) // 读取数据
		client.err = err
		if err != nil {
			break
		}
		recvStr := string(buf[:n])
		fmt.Printf("收到Client %s 端发来的数据: %s\n", client.srcIP, recvStr)

		commandResponse := client.handleCommand(recvStr)
		if commandResponse != nil {
			if err = commandResponse.Response(client); err != nil {
				fmt.Printf("send Response data to client: %s failed\n", client.srcIP)
				continue
			}
		} else {
			client.mutex.Lock()
			client.writer.Write([]byte(recvStr)) // 发送数据
			client.writer.Flush()
			client.mutex.Unlock()
		}
	}
}

func (client *clientHandler) handleCommand(line string) (r *result) {
	client.parseLine(line)
	fmt.Printf("client.command: %s, client.param: %s\n", client.command, client.param)
	if client.middleware[client.command] != nil {
		if err := client.middleware[client.command](client.param); err != nil {
			client.param = ""
			return &result{
				code: 500,
				msg:  fmt.Sprintf("Internal error: %s", err),
			}
		}
		client.param = ""
	}

	return nil
}

func getCommand(line string) []string {
	return strings.SplitN(strings.Trim(line, "\r\n"), " ", 2)
}

func (client *clientHandler) parseLine(line string) {
	params := getCommand(line)
	client.command = strings.ToUpper(params[0])
	if len(params) > 1 {
		client.param = params[1]
	}
}

func (client *clientHandler) writeMessage(code int, message string) error {
	line := fmt.Sprintf("%d %s", code, message)
	return client.writeLine(line)
}

func (client *clientHandler) writeLine(line string) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	if _, err := client.writer.WriteString(line + "\r\n"); err != nil {
		return err
	}
	if err := client.writer.Flush(); err != nil {
		return err
	}

	return nil
}
