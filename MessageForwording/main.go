package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
}

func NewProxyInfo(clientconn net.Conn) *ProxyInfo {
	return &ProxyInfo{
		ClientConn:   clientconn,
		ClientReader: bufio.NewReader(clientconn),
		ClientWriter: bufio.NewWriter(clientconn),
		ClientAddr:   clientconn.RemoteAddr().String(),
	}
}

func (p *ProxyInfo) ClientToRemote(ctx context.Context) error {
	defer func() {
		//fmt.Printf("goroutine ClientToRemote exit, client %s -> remote %s\n", p.ClientAddr, p.RemoteAddress)
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
					errmsg = fmt.Sprintf("客户端 %s 断开连接\n", p.ClientAddr)
				} else if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					//fmt.Println("Client 读取超时")
					continue
				} else {
					errmsg = fmt.Sprintf("读取客户端 %s 数据失败, 原因：%s\n", p.ClientAddr, err.Error())
				}
				log.Println(errmsg)
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
		//fmt.Printf("goroutine RemoteToClient exit, remote %s -> client %s\n", p.RemoteAddress, p.ClientAddr)
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
					errmsg = fmt.Sprintf("远端 %s 断开连接\n", p.RemoteAddress)
				} else if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					//fmt.Println("remote 读取超时")
					continue
				} else {
					errmsg = fmt.Sprintf("读取远端 %s 数据失败, 原因：%s\n", p.RemoteAddress, err.Error())
				}
				log.Println(errmsg)
				return fmt.Errorf(errmsg)
			}

			p.ClientWriter.Write(buffer[:len])
			p.ClientWriter.Flush()
		}
	}
}

func handleConn(clientconn net.Conn, targetAddr string) {
	defer func() {
		//log.Printf("handleConn exit\n")
		defer clientconn.Close()
	}()

	remoteconn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Println("Dial error: ", err)
		return
	}

	proxyInfo := NewProxyInfo(clientconn)
	proxyInfo.RemoteConn = remoteconn
	proxyInfo.RemoteAddress = targetAddr
	proxyInfo.RemoteReader = bufio.NewReader(remoteconn)
	proxyInfo.RemoteWriter = bufio.NewWriter(remoteconn)

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error { return proxyInfo.ClientToRemote(ctx) })
	g.Go(func() error { return proxyInfo.RemoteToClient(ctx) })

	if err = g.Wait(); err != nil {
		//log.Println("Error happened: ", err)
	}
}

// 配置结构体
type Config struct {
	ListenPort int
	TargetAddr string
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
	flag.Parse() // 解析命令行参数

	if config.ListenPort == 0 || config.TargetAddr == "" {
		flag.Usage()
		os.Exit(1)
	}

	listenAddr := fmt.Sprintf("0.0.0.0:%d", config.ListenPort)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("无法监听端口 %d: %v", config.ListenPort, err)
	}

	log.Printf("服务启动，本地监听端口 %d, 转发到远端地址 %s", config.ListenPort, config.TargetAddr)
	defer listener.Close()

	for {
		netConn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}

		remoteAddr := netConn.RemoteAddr()
		log.Printf("收到连接请求，客户端地址: %s, 转发地址: %s\n", remoteAddr.String(), config.TargetAddr)
		go handleConn(netConn, config.TargetAddr)
	}
}
