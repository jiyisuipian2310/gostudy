package tcpserver

import (
	"example8/config"
	"fmt"
	"net"
	"time"
)

var cfgFile = "./config.yaml"

type TcpServer struct {
	listener net.Listener
	shutdown bool
	cfg      *config.Config
}

func NewTcpServer() (*TcpServer, error) {
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	server := &TcpServer{
		cfg: cfg,
	}

	return server, nil
}

func (s *TcpServer) DisplayCfgInfo() {
	fmt.Printf("ListenAddr: %s\n", s.cfg.ListenAddr)
}

func (s *TcpServer) listen() (err error) {
	l, err := net.Listen("tcp", s.cfg.ListenAddr)
	if err != nil {
		return err
	}
	s.listener = l
	return nil
}

func (s *TcpServer) serve() error {
	for {
		netConn, err := s.listener.Accept()
		if err != nil {
			if s.listener != nil {
				return err
			}
		}

		// set linger 0 and tcp keepalive setting between client connection
		conn := netConn.(*net.TCPConn)
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(time.Duration(600) * time.Second)
		conn.SetLinger(0)

		if s.cfg.IdleTimeout > 0 {
			conn.SetDeadline(time.Now().Add(time.Duration(s.cfg.IdleTimeout) * time.Second))
		}

		c := newClientHandler(conn, s.cfg)
		c.DisplayClientInfo()
		c.ProcessMessage()
	}
}

func (s *TcpServer) StartServer() (err error) {
	if err := s.listen(); err != nil {
		return err
	}

	errChannel := make(chan error)
	fmt.Printf("Listening address %s\n", s.listener.Addr())

	done := make(chan struct{})

	go func() {
		if err := s.serve(); err != nil {
			if !s.shutdown {
				errChannel <- err
			}
		}
		done <- struct{}{}
	}()

	<-done

	return nil
}
