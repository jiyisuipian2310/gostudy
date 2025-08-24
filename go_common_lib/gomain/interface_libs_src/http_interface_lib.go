package interface_libs_src 

/*
#include <stdlib.h>
#include <string.h>

//begin http related defination
typedef enum {
    ERROR_TYPE_NONE = 0,            // 无错误
    ERROR_TYPE_CONNECTION_FAILED,   // 连接失败
    ERROR_TYPE_TIMEOUT,             // 超时
    ERROR_TYPE_TLS,                 // TLS错误
    ERROR_TYPE_DNS,                 // DNS解析错误
    ERROR_TYPE_READ_RESPONSE,       // 读取响应失败
    ERROR_TYPE_OTHER                // 其他错误
} HttpErrorType;

typedef void (*HttpCallback)(const char* response, int status, int error_type, void* user_data);

// 辅助函数，用于调用回调
static void invoke_http_callback(HttpCallback callback, const char* response, int status, int error_type, void* user_data) {
    if (callback != NULL) {
        callback(response, status, error_type, user_data);
    }
}
//end http related defination
*/
import "C"
import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	"unsafe"
)

// Go端的错误类型常量（与C端保持一致）
const (
	ErrorTypeNone             = C.ERROR_TYPE_NONE
	ErrorTypeConnectionFailed = C.ERROR_TYPE_CONNECTION_FAILED
	ErrorTypeTimeout          = C.ERROR_TYPE_TIMEOUT
	ErrorTypeTLS              = C.ERROR_TYPE_TLS
	ErrorTypeDNS              = C.ERROR_TYPE_DNS
	ErrorTypeReadResponse     = C.ERROR_TYPE_READ_RESPONSE
	ErrorTypeOther            = C.ERROR_TYPE_OTHER
)

var handleStore = struct {
	sync.Mutex
	m map[C.long]uintptr
}{m: make(map[C.long]uintptr)}

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     10 * time.Second,
	},
}

//export SendHttpRequestSync
func SendHttpRequestSync(cURL *C.char, cMethod *C.char, cHeaders *C.char, cBody *C.char, callback C.HttpCallback, userData unsafe.Pointer) {
	url := C.GoString(cURL)
	method := C.GoString(cMethod)
	headers := C.GoString(cHeaders)
	body := C.GoString(cBody)

	var headerMap map[string]string
	if headers != "" {
		if err := json.Unmarshal([]byte(headers), &headerMap); err != nil {
			cResp := C.CString("Failed to parse head data: " + err.Error())
			defer C.free(unsafe.Pointer(cResp))
			C.invoke_http_callback(callback, cResp, C.int(0), C.int(ErrorTypeOther), userData)
			return
		}
	}

	req, _ := http.NewRequestWithContext(
		context.Background(),
		method,
		url,
		strings.NewReader(body),
	)

	for k, v := range headerMap {
		req.Header.Set(k, v)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		errorType := determineErrorType(err)
		errorMsg := err.Error()

		cResp := C.CString(errorMsg)
		defer C.free(unsafe.Pointer(cResp))

		C.invoke_http_callback(callback, cResp, C.int(0), C.int(errorType), userData)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		cResp := C.CString("Failed to read response: " + err.Error())
		defer C.free(unsafe.Pointer(cResp))
		C.invoke_http_callback(callback, cResp, C.int(resp.StatusCode), C.int(ErrorTypeReadResponse), userData)
		return
	}

	cResp := C.CString(string(respBody))
	defer C.free(unsafe.Pointer(cResp))

	C.invoke_http_callback(callback, cResp, C.int(resp.StatusCode), C.int(ErrorTypeNone), userData)
}

//export SendHttpRequestAsync
func SendHttpRequestAsync(cURL *C.char, cMethod *C.char, cHeaders *C.char, cBody *C.char, callback C.HttpCallback, userData unsafe.Pointer) {
	url := C.GoString(cURL)
	method := C.GoString(cMethod)
	headers := C.GoString(cHeaders)
	body := C.GoString(cBody)

	h := C.long(uintptr(userData))

	handleStore.Lock()
	handleStore.m[h] = uintptr(userData)
	handleStore.Unlock()

	go func() {
		defer func() {
			handleStore.Lock()
			delete(handleStore.m, h)
			handleStore.Unlock()
		}()

		var headerMap map[string]string
		if headers != "" {
			if err := json.Unmarshal([]byte(headers), &headerMap); err != nil {
				cResp := C.CString("Failed to parse head data: " + err.Error())
				defer C.free(unsafe.Pointer(cResp))
				C.invoke_http_callback(callback, cResp, C.int(0), C.int(ErrorTypeOther), userData)
				return
			}
		}

		req, _ := http.NewRequestWithContext(
			context.Background(),
			method,
			url,
			strings.NewReader(body),
		)

		for k, v := range headerMap {
			req.Header.Set(k, v)
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			errorType := determineErrorType(err)
			errorMsg := err.Error()

			cResp := C.CString(errorMsg)
			defer C.free(unsafe.Pointer(cResp))

			C.invoke_http_callback(callback, cResp, C.int(0), C.int(errorType), userData)
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			cResp := C.CString("Failed to read response: " + err.Error())
			defer C.free(unsafe.Pointer(cResp))
			C.invoke_http_callback(callback, cResp, C.int(resp.StatusCode), C.int(ErrorTypeReadResponse), userData)
			return
		}

		cResp := C.CString(string(respBody))
		defer C.free(unsafe.Pointer(cResp))

		C.invoke_http_callback(callback, cResp, C.int(resp.StatusCode), C.int(ErrorTypeNone), userData)
	}()
}

// determineErrorType 确定错误类型
func determineErrorType(err error) int {
	switch {
	case isTimeoutError(err):
		return ErrorTypeTimeout
	case isConnectionError(err):
		return ErrorTypeConnectionFailed
	case isDNSError(err):
		return ErrorTypeDNS
	case isTLSError(err):
		return ErrorTypeTLS
	default:
		return ErrorTypeOther
	}
}

// isTimeoutError 检查是否为超时错误
func isTimeoutError(err error) bool {
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}
	return strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "Timeout") ||
		strings.Contains(err.Error(), "TIMEOUT")
}

// isConnectionError 检查是否为连接错误
func isConnectionError(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		return opErr.Op == "dial"
	}
	return strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "no such host") ||
		strings.Contains(err.Error(), "network is unreachable")
}

// isDNSError 检查是否为DNS错误
func isDNSError(err error) bool {
	if _, ok := err.(*net.DNSError); ok {
		return true
	}
	return strings.Contains(err.Error(), "DNS") ||
		strings.Contains(err.Error(), "dns") ||
		strings.Contains(err.Error(), "no such host")
}

// isTLSError 检查是否为TLS错误
func isTLSError(err error) bool {
	return strings.Contains(err.Error(), "tls") ||
		strings.Contains(err.Error(), "TLS") ||
		strings.Contains(err.Error(), "certificate") ||
		strings.Contains(err.Error(), "x509")
}

