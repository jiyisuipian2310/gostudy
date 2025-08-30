package main

import (
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

const addr = "192.168.104.100:22"
const user = "root"
const password = "root100^YHN"

func main() {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	//设置输入输出
	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	if err := session.Run("ls -alh"); err != nil { //执行命令
		log.Fatal("Failed to run command: ", err)
	}
}
