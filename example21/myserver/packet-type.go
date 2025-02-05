package myserver

import (
	"encoding"
	"fmt"
)

const (
	eInitPacket  = 1
	eLoginPacket = 2
	eDataPacket  = 3
)

type requestPacket interface {
	encoding.BinaryUnmarshaler
}

type responsePacket interface {
	encoding.BinaryMarshaler
}

/*********************************************************/
type InitPacket struct {
}

func (p *InitPacket) UnmarshalBinary(data []byte) error {
	return nil
}

func (p *InitPacket) MarshalBinary() ([]byte, error) {
	return nil, nil
}

/*********************************************************/
type LoginPacket struct {
}

func (p *LoginPacket) UnmarshalBinary(data []byte) error {
	return nil
}

func (p *LoginPacket) MarshalBinary() ([]byte, error) {
	return nil, nil
}

/*********************************************************/
type DataPacket struct {
}

func (p *DataPacket) UnmarshalBinary(data []byte) error {
	return nil
}

func (p *DataPacket) respond(svr *Server) responsePacket {
	return nil
}

// func (p *DataPacket) MarshalBinary() ([]byte, error) {
// 	return nil, nil
// }

/*********************************************************/

type defaultPacket interface {
	encoding.BinaryUnmarshaler
	respond(svr *Server) responsePacket
}

func getPacketType(t uint8) (uint8, error) {
	return t, nil
}

/*********************************************************/

func makePacket(pktType uint8, data []byte) (requestPacket, error) {
	var pkt requestPacket
	switch pktType {
	case eInitPacket:
		pkt = &InitPacket{}
	case eLoginPacket:
		pkt = &LoginPacket{}
	case eDataPacket:
		pkt = &DataPacket{}
	default:
		return nil, fmt.Errorf("makePacket, unknown packet type: %d", pktType)
	}

	if err := pkt.UnmarshalBinary(data); err != nil {
		return pkt, err
	}

	return pkt, nil
}
