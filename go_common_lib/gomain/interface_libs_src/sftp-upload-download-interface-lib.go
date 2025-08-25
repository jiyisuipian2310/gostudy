package interface_libs_src 

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"io"
	"os"
	"unsafe"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//export UploadFileByPassword
func UploadFileByPassword(address *C.char, username *C.char, password *C.char, remotePath *C.char, localPath *C.char, errMsg *C.char, errMsgLen C.int) C.int {
	if errMsg == nil || errMsgLen < 128 {
		return -1
	}

	// 转换C字符串到Go字符串
	goAddress := C.GoString(address)
	goUserName := C.GoString(username)
	goPassword := C.GoString(password)
	goLocalPath := C.GoString(localPath)
	goRemotePath := C.GoString(remotePath)

	config := &ssh.ClientConfig{
		User: goUserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(goPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

    setErrorMessage := func(message string) {
        cMsg := C.CString(message)
        defer C.free(unsafe.Pointer(cMsg)) // 确保释放内存
        C.strncpy(errMsg, cMsg, C.size_t(errMsgLen)-1)
    }

	// 连接SSH服务器
	conn, err := ssh.Dial("tcp", goAddress, config)
	if err != nil {
		setErrorMessage(fmt.Sprintf("Failed to connect to %s: %v", address, err))
		return -1
	}
	defer conn.Close()

	// 创建SFTP客户端
	client, err := sftp.NewClient(conn)
	if err != nil {
		setErrorMessage(fmt.Sprintf("Failed to create SFTP client: %v", err))
		return -1
	}
	defer client.Close()

	// 打开本地文件
	localFile, err := os.Open(goLocalPath)
	if err != nil {
		setErrorMessage(fmt.Sprintf("Failed to open local file '%s': %v", goLocalPath, err))
		return -1
	}
	defer localFile.Close()

	// 创建远程文件
	remoteFile, err := client.Create(goRemotePath)
	if err != nil {
		setErrorMessage(fmt.Sprintf("Failed to create remote file '%s': %v", goRemotePath, err))
		return -1
	}
	defer remoteFile.Close()

	// 复制文件内容（上传）
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		setErrorMessage(fmt.Sprintf("Failed to upload file: %v", err))
		return -1
	}

	return 0
}

//export DownloadFileByPassword
func DownloadFileByPassword(address *C.char, username *C.char, password *C.char, remotePath *C.char, localPath *C.char, errMsg *C.char, errMsgLen C.int) C.int {
	if errMsg == nil || errMsgLen < 128 {
		return -1
	}

	// 转换C字符串到Go字符串
	goAddress := C.GoString(address)
	goUserName := C.GoString(username)
	goPassword := C.GoString(password)
	goRemotePath := C.GoString(remotePath)
	goLocalPath := C.GoString(localPath)

	// 建立SSH连接
	config := &ssh.ClientConfig{
		User: goUserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(goPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 5*time.Second,
	}

    setErrorMessage := func(message string) {
        cMsg := C.CString(message)
        defer C.free(unsafe.Pointer(cMsg)) // 确保释放内存
        C.strncpy(errMsg, cMsg, C.size_t(errMsgLen)-1)
    }

	// 连接SSH服务器
	conn, err := ssh.Dial("tcp", goAddress, config)
	if err != nil {
		setErrorMessage(fmt.Sprintf("Failed to connect: %v", err))
		return -1
	}
	defer conn.Close()

	// 创建SFTP客户端
	client, err := sftp.NewClient(conn)
	if err != nil {
		setErrorMessage(fmt.Sprintf("Failed to create SFTP client: %v", err))
		return -1
	}
	defer client.Close()

	// 打开远程文件
	remoteFile, err := client.Open(goRemotePath)
	if err != nil {
		setErrorMessage(fmt.Sprintf("Failed to open remote file: %v", err))
		return -1
	}
	defer remoteFile.Close()

	// 创建本地文件
	localFile, err := os.Create(goLocalPath)
	if err != nil {
		setErrorMessage(fmt.Sprintf("Failed to create local file: %v", err))
		return -1
	}
	defer localFile.Close()

	// 复制文件内容
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		setErrorMessage(fmt.Sprintf("Failed to copy file: %v", err))
		return -1
	}

	return 0
}
