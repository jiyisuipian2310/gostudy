package main

import (
	"main/myserver"
)

func main() {
	serve, _ := myserver.NewServer()
	serve.Serve()
}
