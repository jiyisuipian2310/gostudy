package tcpserver

import (
	"example8/config"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	fmt.Printf("cfginfo - ListenAddr: %s\n", s.cfg.ListenAddr)
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
		fmt.Printf("before accept\n")
		netConn, err := s.listener.Accept()
		fmt.Printf("after accept\n")
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
		go c.ProcessMessage()
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

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

L:
	for {
		switch <-ch {
		case syscall.SIGHUP:
			fmt.Printf("syscall.SIGHUP received, stopping server\n")
		case syscall.SIGTERM:
			fmt.Printf("syscall.SIGTERM received, stopping server\n")
		case syscall.SIGINT:
			fmt.Printf("syscall.SIGINT received, stopping server\n")
			if err := s.stop(); err != nil {
				errChannel <- err
			}
			break L
		}
	}

	<-done

	return nil
}

func (s *TcpServer) stop() error {
	s.shutdown = true
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return err
		}
	}
	return nil
}
