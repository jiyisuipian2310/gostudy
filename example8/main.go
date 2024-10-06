package main

import (
	"example8/tcpserver"
	"fmt"
)

func main() {
	tcpServer, err := tcpserver.NewTcpServer()
	if err != nil {
		fmt.Printf("NewTcpServer error:%v\n", err)
		return
	}

	tcpServer.AddMiddleware("user", ProcessUserCommand)

	tcpServer.DisplayCfgInfo()

	tcpServer.StartServer()
}

type CustomError struct {
	ErrorCode int
	Message   string
}

// 实现error接口的Error方法
func (e *CustomError) Error() string {
	return fmt.Sprintf("Error Code: %d, Message: %s", e.ErrorCode, e.Message)
}

func ProcessUserCommand(param string) error {
	if param == "" {
		return &CustomError{
			ErrorCode: 1001,
			Message:   "param is not valid",
		}
	}

	fmt.Printf("UserParam:%v\n", param)
	return nil
}
