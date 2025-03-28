package main

import (
	"bufio"
	"context"
	"fmt"
	"grpcdemo/proto"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ADD_STUDENT_INFO    = 1
	GET_STUDENT_INFO    = 2
	DEL_STUDENT_INFO    = 3
	UPDATE_STUDENT_INFO = 4
	SHOW_STUDENT_INFO   = 5
)

func main() {
	var err error

	var serviceHost = "127.0.0.1:8001"

	conn, err := grpc.Dial(serviceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	client := proto.NewMessageServiceClient(conn)
	_, err = client.SendMessage(context.TODO(), &proto.RequestMessage{
		ReqCode: ADD_STUDENT_INFO,
		ReqData: "{\"name\":\"yull2\", \"age\": 20, \"address\": \"beijing\"}",
	})

	if err != nil {
		fmt.Printf("---> %s\n", err.Error())
		return
	}

	_, err = client.SendMessage(context.TODO(), &proto.RequestMessage{
		ReqCode: ADD_STUDENT_INFO,
		ReqData: "{\"name\": \"zhangsan2\", \"age\": 25, \"address\": \"shanghai\"}",
	})

	if err != nil {
		fmt.Println(err)
	}

	_, err = client.SendMessage(context.TODO(), &proto.RequestMessage{
		ReqCode: SHOW_STUDENT_INFO,
		ReqData: "",
	})

	resStream, err := client.GetStudentInfo(context.TODO(), &proto.RequestMessage{})
	if err != nil {
		log.Fatalf("error while calling GetStudentInfo RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GetStudentInfo: %v", msg.RespData)
	}

	fmt.Println("按回车键退出程序...")
	in := bufio.NewReader(os.Stdin)
	_, _, _ = in.ReadLine()
}
