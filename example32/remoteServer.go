package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type RemoteHttpHandler struct {
}

func newRemoteHttpHandler() *RemoteHttpHandler {
	return &RemoteHttpHandler{}
}

func (h *RemoteHttpHandler) DisConnectionCallback(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logrus.Info(fmt.Sprintf("remoteServer Received Message: %s", string(body)))

	/*
		logrus.Info(fmt.Sprintf("remoteServer Received Message: %s", string(body)))
		for name, port := range g_mapDBPort {
			targetURL := fmt.Sprintf("http://localhost:%s/disConnection", port)
			_, err := http.Post(targetURL, "application/json", bytes.NewBuffer(body))
			if err != nil {
				logrus.Warning(fmt.Sprintf("SendMessage to %s Failed, Url: %s", name, targetURL))
			} else {
				logrus.Info(fmt.Sprintf("SendMessage to %s Success, Url: %s", name, targetURL))
			}
		}
	*/

	var decodedBody = new(HostInfo)
	json.Unmarshal([]byte(string(body)), &decodedBody)
	if decodedBody.MasterAccount == "" {
		errMsg := "{\"resultCode\":1, resultDesc:\"No Found masterAccount or masterAccount is empty !\"}"
		logrus.Warning(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	bFind := KillPidByMasterAccount(decodedBody.MasterAccount)
	if !bFind {
		errMsg := fmt.Sprintf("{\"resultCode\":1, resultDesc:\"masterAccount[%s] isn't exist !\"}", decodedBody.MasterAccount)
		logrus.Warning(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// 返回目标程序的响应给客户端
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"resultCode\":0}"))
}

func (h *RemoteHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/p/platformstatus/breakUserAndProcess" {
		h.DisConnectionCallback(w, r)
		return
	}

	strErrMsg := fmt.Sprintf("Don't support url: %s", r.URL.Path)
	logrus.Warning(strErrMsg)
	http.Error(w, strErrMsg, http.StatusBadRequest)
}
