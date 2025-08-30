package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type ProxyInfo struct {
	ClientConn   net.Conn
	ClientReader *bufio.Reader
	ClientWriter *bufio.Writer
	ClientAddr   string

	RemoteConn    net.Conn
	RemoteReader  *bufio.Reader
	RemoteWriter  *bufio.Writer
	RemoteAddress string

	DstIp     string `json:"dstip"`
	DstPort   int    `json:"dstport"`
	DstDomain string `json:"dstdomain"`

	SyncWG sync.WaitGroup
}

func NewProxyInfo(clientconn net.Conn) *ProxyInfo {
	return &ProxyInfo{
		ClientConn:   clientconn,
		ClientReader: bufio.NewReader(clientconn),
		ClientWriter: bufio.NewWriter(clientconn),
		ClientAddr:   clientconn.RemoteAddr().String(),
	}
}

func (p *ProxyInfo) ReadProxyData(reader *bufio.Reader) (err error) {
	var buf [7]byte
	err = binary.Read(reader, binary.LittleEndian, &buf)
	if err != nil {
		return err
	}

	proxy_flag := buf[:5]
	if string(proxy_flag) != "proxy" {
		return fmt.Errorf("my_proxy_handler refuse, reason: magic[%s] != proxy", proxy_flag)
	}

	length := binary.BigEndian.Uint16(buf[5:])
	if int(length) > 256 {
		return fmt.Errorf("my_proxy_handler refuse, reason: proxylen[%d] > 256\n", length)
	}

	fmt.Printf("proxylen[%d]\n", length)

	proxydata := make([]byte, length)

	//从 reader 中最少读取 length 个字节
	_, err = io.ReadFull(reader, proxydata)
	if err != nil {
		return fmt.Errorf("ReadFull error: %s", err.Error())
	}

	fmt.Printf("proxydata: %s\n", proxydata)

	err = json.Unmarshal([]byte(proxydata), p)
	if err != nil {
		return fmt.Errorf("json.Unmarshal error: ", err.Error())
	}

	return nil
}

func (p *ProxyInfo) ClientToRemote(chClientClose chan<- string, chRemoteClose <-chan string) error {
	defer func() {
		fmt.Printf("goroutine ClientToRemote exit, client %s -> remote %s\n", p.ClientAddr, p.RemoteAddress)
	}()

	dataCh := make(chan []byte)
	errCh := make(chan string, 1)

	go func() {
		buffer := make([]byte, 40960)
		for {
			len, err := p.ClientReader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					errCh <- fmt.Sprintf("Client %s closed the connection", p.ClientAddr)
				} else {
					errCh <- fmt.Sprintf("Read client %s data failed %s", p.ClientAddr, err.Error())
				}
				return
			}

			dataCh <- buffer[:len]
		}
	}()

	for {
		select {
		case data := <-dataCh:
			p.RemoteWriter.Write(data)
			p.RemoteWriter.Flush()
		case errmsg := <-errCh:
			chClientClose <- errmsg
			return fmt.Errorf("%s", errmsg)
		case errmsg := <-chRemoteClose:
			fmt.Printf("Detect remote quit in ClientToRemote: %s\n", errmsg)
			return fmt.Errorf("%s", errmsg)
		default:
			time.Sleep(50 * time.Millisecond) // 模拟其他工作
		}
	}
	return nil
}

func (p *ProxyInfo) RemoteToClient(chClientClose <-chan string, chRemoteClose chan<- string) {
	defer func() {
		defer p.RemoteConn.Close()
		fmt.Printf("goroutine RemoteToClient exit, remote %s -> client %s\n", p.RemoteAddress, p.ClientAddr)
		p.SyncWG.Done()
	}()

	dataCh := make(chan []byte)
	errCh := make(chan string, 1)

	go func() {
		buffer := make([]byte, 40960)
		for {
			len, err := p.RemoteReader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					errCh <- fmt.Sprintf("remote %s closed the connection", p.RemoteAddress)
				} else {
					errCh <- fmt.Sprintf("read remote %s data failed %s", p.RemoteAddress, err.Error())
				}
				return
			}

			dataCh <- buffer[:len]
		}
	}()

loop:
	for {
		select {
		case data := <-dataCh:
			p.ClientWriter.Write(data)
			p.ClientWriter.Flush()
		case errmsg := <-chClientClose:
			fmt.Printf("Detect client quit in RemoteToClient: %s\n", errmsg)
			break loop
		case errmsg := <-errCh:
			chRemoteClose <- errmsg
			break loop
		default:
			time.Sleep(50 * time.Millisecond) // 模拟其他工作
		}
	}
}

func handleConn(clientconn net.Conn) {
	defer func() {
		fmt.Printf("handleConn exit\n")
		defer clientconn.Close()
	}()

	proxyInfo := NewProxyInfo(clientconn)
	fmt.Printf("client %s connected\n", proxyInfo.ClientAddr)

	err := proxyInfo.ReadProxyData(proxyInfo.ClientReader)
	if err != nil {
		fmt.Println("ReadProxyData error: ", err)
		return
	}

	fmt.Printf("dstip: %s, dstport: %d, dstdomain: %s\n", proxyInfo.DstIp, proxyInfo.DstPort, proxyInfo.DstDomain)

	proxyInfo.RemoteAddress = fmt.Sprintf("%s:%d", proxyInfo.DstIp, proxyInfo.DstPort)
	remoteconn, err := net.Dial("tcp", proxyInfo.RemoteAddress)
	if err != nil {
		fmt.Println("Dial error: ", err)
		return
	}

	proxyInfo.RemoteConn = remoteconn
	proxyInfo.RemoteReader = bufio.NewReader(remoteconn)
	proxyInfo.RemoteWriter = bufio.NewWriter(remoteconn)

	chRemoteClose := make(chan string, 1)
	chClientClose := make(chan string, 1)

	proxyInfo.SyncWG.Add(1)
	go proxyInfo.RemoteToClient(chClientClose, chRemoteClose)

	err = proxyInfo.ClientToRemote(chClientClose, chRemoteClose)
	proxyInfo.SyncWG.Wait()
	if err != nil {
		return
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error: ", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}

		go handleConn(conn)
	}
}
