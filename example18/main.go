package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"golang.org/x/sync/errgroup"
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

func (p *ProxyInfo) ClientToRemote(ctx context.Context) error {
	defer func() {
		fmt.Printf("goroutine ClientToRemote exit, client %s -> remote %s\n", p.ClientAddr, p.RemoteAddress)
	}()

	buffer := make([]byte, 40960)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			p.ClientConn.SetReadDeadline(time.Now().Add(1 * time.Second))
			len, err := p.ClientReader.Read(buffer)
			if err != nil {
				errmsg := ""
				if err == io.EOF {
					errmsg = fmt.Sprintf("Client %s closed the connection\n", p.ClientAddr)
				} else if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					//fmt.Println("Client 读取超时")
					continue
				} else {
					errmsg = fmt.Sprintf("Read client %s data failed %s\n", p.ClientAddr, err.Error())
				}
				fmt.Println(errmsg)
				return fmt.Errorf(errmsg)
			}

			p.RemoteWriter.Write(buffer[:len])
			p.RemoteWriter.Flush()
		}
	}
}

func (p *ProxyInfo) RemoteToClient(ctx context.Context) error {
	defer func() {
		defer p.RemoteConn.Close()
		fmt.Printf("goroutine RemoteToClient exit, remote %s -> client %s\n", p.RemoteAddress, p.ClientAddr)
	}()

	buffer := make([]byte, 40960)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			p.RemoteConn.SetReadDeadline(time.Now().Add(1 * time.Second))
			len, err := p.RemoteReader.Read(buffer)
			if err != nil {
				errmsg := ""
				if err == io.EOF {
					errmsg = fmt.Sprintf("remote %s closed the connection\n", p.RemoteAddress)
				} else if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					//fmt.Println("remote 读取超时")
					continue
				} else {
					errmsg = fmt.Sprintf("read remote %s data failed %s\n", p.RemoteAddress, err.Error())
				}
				fmt.Println(errmsg)
				return fmt.Errorf(errmsg)
			}

			p.ClientWriter.Write(buffer[:len])
			p.ClientWriter.Flush()
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

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error { return proxyInfo.ClientToRemote(ctx) })
	g.Go(func() error { return proxyInfo.RemoteToClient(ctx) })

	if err = g.Wait(); err != nil {
		fmt.Println("Error happened: ", err)
	}
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
