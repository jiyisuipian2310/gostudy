package main

import (
	"context"
	"encoding/json"
	"fmt"
	"grpcdemo/proto"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

const (
	ADD_STUDENT_INFO    = 1
	GET_STUDENT_INFO    = 2
	DEL_STUDENT_INFO    = 3
	UPDATE_STUDENT_INFO = 4
	SHOW_STUDENT_INFO   = 5
)

type StudentInfo struct {
	Name    string `json:"name"`
	Age     int32  `json:"age"`
	Address string `json:"address"`
}

type server struct {
	proto.UnimplementedMessageServiceServer

	MapStuInfo map[string]*StudentInfo
}

func (s *server) SendMessage(ctx context.Context, req *proto.RequestMessage) (*proto.ResponseMessage, error) {
	var err error
	if req.ReqCode == ADD_STUDENT_INFO {
		err = s.AddStudentInfo(req.ReqData)
	} else if req.ReqCode == SHOW_STUDENT_INFO {
		_, err = s.ShowStudentInfo(req.ReqData)
	} else if req.ReqCode == DEL_STUDENT_INFO {
		err = s.DeleteStudentInfo(req.ReqData)
	} else if req.ReqCode == UPDATE_STUDENT_INFO {
		err = s.UpdateStudentInfo(req.ReqData)
	}

	respdata := fmt.Sprintf("Hello %s, Server received message", req.ReqData)
	return &proto.ResponseMessage{RespCode: 0, RespData: respdata}, err
}

func (s *server) GetStudentInfo(req *proto.RequestMessage, stream grpc.ServerStreamingServer[proto.ResponseMessage]) error {
	for i := 0; i < 10; i++ {
		result := "Hello number " + strconv.Itoa(i)
		res := &proto.ResponseMessage{
			RespCode: 0,
			RespData: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func (s *server) AddStudentInfo(reqdata string) error {
	var sinfo StudentInfo
	if err := json.Unmarshal([]byte(reqdata), &sinfo); err != nil {
		return err
	}

	s.MapStuInfo[sinfo.Name] = &sinfo
	return nil
}

func (s *server) DeleteStudentInfo(reqdata string) error {
	var sinfo StudentInfo
	if err := json.Unmarshal([]byte(reqdata), &sinfo); err != nil {
		return err
	}

	delete(s.MapStuInfo, sinfo.Name)
	return nil
}

func (s *server) UpdateStudentInfo(reqdata string) error {
	var sinfo StudentInfo
	if err := json.Unmarshal([]byte(reqdata), &sinfo); err != nil {
		return err
	}

	_, exist := s.MapStuInfo[sinfo.Name]
	if exist {
		s.MapStuInfo[sinfo.Name] = &sinfo
	}
	return nil
}

func (s *server) ShowStudentInfo(reqdata string) (string, error) {
	for k, v := range s.MapStuInfo {
		stuinfo, _ := json.Marshal(v)
		fmt.Printf("Name: %s, Info: %s\n", k, string(stuinfo))
	}

	return "", nil
}

func main() {
	listen, err := net.Listen("tcp", ":8001")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	proto.RegisterMessageServiceServer(s, &server{
		MapStuInfo: make(map[string]*StudentInfo),
	})

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
