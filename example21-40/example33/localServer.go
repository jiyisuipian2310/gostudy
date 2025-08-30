package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

var g_PidAndMasterMutex sync.Mutex
var g_mapPidAndMasterAccount map[string]string

type HostInfo struct {
	Opcode        string `json:"opcode"`
	MasterAccount string `json:"masterAccount"`
	Pid           string `json:"pid"`
}

func init() {
	g_mapPidAndMasterAccount = make(map[string]string)
}

/************************** LocalHttpHandler ******************************/
type LocalHttpHandler struct {
}

func newLocalHttpHandler() *LocalHttpHandler {
	return &LocalHttpHandler{}
}

func (h *LocalHttpHandler) operatePidAndMasterAccount(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "local server failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logrus.Info(fmt.Sprintf("localServer Received Message: %s", string(body)))

	var decodedBody = new(HostInfo)
	json.Unmarshal([]byte(body), &decodedBody)
	logrus.Info(fmt.Sprintf("opcode: %s, masterAccount: %s, pid: %s", decodedBody.Opcode, decodedBody.MasterAccount, decodedBody.Pid))
	if decodedBody.Opcode == "" {
		http.Error(w, "opcode or masterAccount or pid is empty", http.StatusBadRequest)
		logrus.Warning("opcode or masterAccount or pid is empty")
		return
	}

	if decodedBody.Opcode == "add" {
		if decodedBody.MasterAccount == "" || decodedBody.Pid == "" {
			http.Error(w, "masterAccount or pid is empty", http.StatusBadRequest)
			logrus.Warning("opcode is add, masterAccount or pid is empty")
			return
		}
		AddPidAndMaster(decodedBody.Pid, decodedBody.MasterAccount)
	} else if decodedBody.Opcode == "delete_by_pid" {
		if decodedBody.Pid == "" {
			http.Error(w, "pid is empty", http.StatusBadRequest)
			logrus.Warning("opcode is delete_by_pid, pid is empty")
			return
		}
		DeletePidAndMasterByPid(decodedBody.Pid)
	} else if decodedBody.Opcode == "delete_by_master_account" {
		if decodedBody.MasterAccount == "" {
			http.Error(w, "pid is empty", http.StatusBadRequest)
			logrus.Warning("opcode is delete_by_master_account, MasterAccount is empty")
			return
		}
		DeletePidAndMasterByMasterAccount(decodedBody.MasterAccount)
	} else if decodedBody.Opcode == "display" {
		DisplayPidAndMaster()
	} else if decodedBody.Opcode == "kill" {
		if decodedBody.MasterAccount == "" {
			http.Error(w, "pid is empty", http.StatusBadRequest)
			logrus.Warning("opcode is kill, MasterAccount is empty")
			return
		}
		KillPidByMasterAccount(decodedBody.MasterAccount)
	}

	// 返回目标程序的响应给客户端
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *LocalHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/operatePidAndMasterAccount" {
		h.operatePidAndMasterAccount(w, r)
		return
	}

	strErrMsg := fmt.Sprintf("Local Server don't support url: %s", r.URL.Path)
	logrus.Warning(strErrMsg)
	http.Error(w, strErrMsg, http.StatusBadRequest)
}

func AddPidAndMaster(pid, masterAccount string) {
	g_PidAndMasterMutex.Lock()
	g_mapPidAndMasterAccount[pid] = masterAccount
	g_PidAndMasterMutex.Unlock()
}

func DeletePidAndMasterByPid(pid string) {
	g_PidAndMasterMutex.Lock()
	delete(g_mapPidAndMasterAccount, pid)
	g_PidAndMasterMutex.Unlock()
}

func DeletePidAndMasterByMasterAccount(masterAccount string) {
	g_PidAndMasterMutex.Lock()
	for key, value := range g_mapPidAndMasterAccount {
		if value == masterAccount {
			delete(g_mapPidAndMasterAccount, key)
		}
	}
	g_PidAndMasterMutex.Unlock()
}

func KillPidByMasterAccount(masterAccount string) bool {
	bfind := false

	g_PidAndMasterMutex.Lock()
	for key, value := range g_mapPidAndMasterAccount {
		if value == masterAccount {
			pid, err := strconv.Atoi(key)
			if err == nil {
				bfind = true
				syscall.Kill(pid, syscall.SIGKILL)
				delete(g_mapPidAndMasterAccount, key)
				logrus.Info(fmt.Sprintf("Have find pid[%d] by masterAccount[%s], kill it", pid, masterAccount))
			}
		}
	}
	g_PidAndMasterMutex.Unlock()

	if !bfind {
		logrus.Warning(fmt.Sprintf("Haven't find pid by masterAccount[%s]", masterAccount))
		return false
	}

	return true
}

func DisplayPidAndMaster() {
	g_PidAndMasterMutex.Lock()
	logrus.Info(fmt.Sprintf("g_mapPidAndMasterAccount size: %d", len(g_mapPidAndMasterAccount)))
	for k, v := range g_mapPidAndMasterAccount {
		logrus.Info(fmt.Sprintf("pid: %s, masterAccount:%s", k, v))
	}
	g_PidAndMasterMutex.Unlock()
}
