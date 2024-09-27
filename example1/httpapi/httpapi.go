package httpapi

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"unsafe"
)

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

/*发送https post消息, 并获取响应数据, 不校验对端证书*/
func SendAndRecvHttpsPostMsg(targetUrl string, body string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //Do not verify peer certificate
	}

	payload := strings.NewReader(body)
	client := &http.Client{Transport: tr}
	resp, err := client.Post(targetUrl, "text/plain", payload)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return BytesToString(respBody), err
}

/*发送http post消息, 并获取响应数据*/
func SendAndRecvHttpPostMsg(targetUrl string, body string) (string, error) {
	payload := strings.NewReader(body)
	resp, err := http.Post(targetUrl, "text/plain", payload)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return BytesToString(respBody), err
}
