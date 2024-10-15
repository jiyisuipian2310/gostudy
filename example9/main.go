package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	// SFTP服务器的信息
	host := "192.168.104.109:22"
	user := "root"
	password := "root100^YHN"
	remoteFilePath := "/home/yull/miniserver/miniserver.tar.gz" // 在服务器上的文件路径
	localFilePath := "./miniserver.tar.gz"                      // 本地文件保存路径

	// SSH配置
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         3 * time.Second,
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("currTime: %s\n", timeStr)

	// 连接SSH服务端
	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		log.Fatalf("Failed to connect to remote server: %v, currTime: %s", err, timeStr)
	}
	defer conn.Close()

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatalf("Failed to create sftp client: %v", err)
	}
	defer sftpClient.Close()

	// 打开远程文件
	remoteFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		log.Fatalf("Failed to open remote file: %v", err)
	}
	defer remoteFile.Close()

	// 读取远程文件内容
	fileBytes, err := ioutil.ReadAll(remoteFile)
	if err != nil {
		log.Fatalf("Failed to read remote file: %v", err)
	}

	// 创建本地文件以保存下载的文件
	localFile, err := os.Create(localFilePath)
	if err != nil {
		log.Fatalf("Failed to create local file: %v", err)
	}
	defer localFile.Close()

	// 将文件内容写入本地文件
	_, err = localFile.Write(fileBytes)
	if err != nil {
		log.Fatalf("Failed to write to local file: %v", err)
	}

	log.Println("File downloaded successfully")
}
