package myserver

import "fmt"

type Server struct {
}

func NewServer() (*Server, error) {
	s := &Server{}
	return s, nil
}

func (s *Server) Serve() error {
	var err error
	var pktType uint8
	var pkt requestPacket

	pktType, err = getPacketType(3)
	if err != nil {
		fmt.Printf("getPacketType error: %v\n", err.Error())
		return err
	}

	fmt.Printf("pktType=%d\n", pktType)

	data := make([]byte, 100)
	pkt, err = makePacket(pktType, data)
	if err != nil {
		fmt.Printf("makePacket error: %v\n", err.Error())
		return err
	}

	err = s.dealRequestPacket(pkt)
	if err != nil {
		fmt.Printf("dealRequestPacket error: %v\n", err.Error())
		return err
	}

	{
		pkt, err = makePacket(2, data)
		if err != nil {
			fmt.Printf("makePacket error: %v\n", err.Error())
			return err
		}

		err = s.dealRequestPacket(pkt)
		if err != nil {
			fmt.Printf("dealRequestPacket error: %v\n", err.Error())
			return err
		}
	}

	return nil
}

func (s *Server) dealRequestPacket(p requestPacket) error {
	switch p.(type) {
	case *InitPacket:
		fmt.Printf("dealRequestPacket, this is InitPacket\n")
	case *LoginPacket:
		fmt.Printf("dealRequestPacket, this is LoginPacket\n")
	case defaultPacket:
		fmt.Printf("dealRequestPacket, this is defaultPacket\n")
	default:
		return fmt.Errorf("dealRequestPacket, unknown packet type %T", p)
	}

	return nil
}
