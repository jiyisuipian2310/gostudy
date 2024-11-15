package main

import (
	"context"
	"fmt"
	"grpcdemo/proto"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedSendMessageServer
}

func (s *server) Say(ctx context.Context, req *proto.SayRequest) (*proto.SayResponse, error) {
	fmt.Println("OpCode: ", req.OpCode)
	fmt.Println("OpMessage: ", req.OpMessage)
	return &proto.SayResponse{Name: "Hello " + req.OpMessage}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":8001")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	proto.RegisterSendMessageServer(s, &server{})
	//reflection.Register(s)

	defer func() {
		s.Stop()
		listen.Close()
	}()

	fmt.Println("Serving 8001...")
	err = s.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
