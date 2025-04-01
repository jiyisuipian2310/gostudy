package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type UserCommand struct {
	LoginAddr   string `json:"address"`
	LoginUser   string `json:"user"`
	LoginPasswd string `json:"password"`
	Command     string `json:"command"`
}

type RemoteHttpHandler struct {
}

func newRemoteHttpHandler() *RemoteHttpHandler {
	return &RemoteHttpHandler{}
}

func (h *RemoteHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"resultCode\":0}"))

	go DealMessage(string(body))
}

func DealMessage(body string) {
	logrus.Info(fmt.Sprintf("remoteServer Received Message: %s", string(body)))

	cmd := UserCommand{}
	json.Unmarshal([]byte(body), &cmd)

	logrus.Info(fmt.Sprintf("LoginAddr: %s, LoginUser: %s, LoginPasswd: %s, Command: %s",
		cmd.LoginAddr, cmd.LoginUser, cmd.LoginPasswd, cmd.Command))

	// 建立SSH客户端连接
	client, err := ssh.Dial("tcp", cmd.LoginAddr, &ssh.ClientConfig{
		User:            cmd.LoginUser,
		Auth:            []ssh.AuthMethod{ssh.Password(cmd.LoginPasswd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		logrus.Warning(fmt.Sprintf("SSH(%s) dial error: %s", cmd.LoginAddr, err.Error()))
		return
	}

	session, err := client.NewSession()
	if err != nil {
		logrus.Warning(fmt.Sprintf("new session error, SSH(%s) dial error: %s", cmd.LoginAddr, err.Error()))
		return
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmd.Command); err != nil {
		logrus.Warning(fmt.Sprintf("Failed to run: %s", err.Error()))
	}
	logrus.Info(fmt.Sprintf("CommandResponse: \n%s", b.String()))
}
