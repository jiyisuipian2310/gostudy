package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	logrus_stack "github.com/Gurpartap/logrus-stack"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(false)
	stackLevels := []logrus.Level{logrus.PanicLevel, logrus.FatalLevel}
	logrus.AddHook(logrus_stack.NewHook(stackLevels, stackLevels))
}

func main() {
	config := AppConfig{}
	if !ParseConfig(&config) {
		return
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go processPidAndMasterAccount(&config)
	go startServer(&config)
	wg.Wait()
	logrus.Info(fmt.Sprintf("httpServer Quit!"))
}

func startServer(config *AppConfig) {
	if !config.HttpsFlag {
		StartHttpServer(config)
	} else {
		StartHttpsServer(config)
	}
}

func StartHttpServer(config *AppConfig) {
	logrus.Info(fmt.Sprintf("Remote httpServer started on http://0.0.0.0:%s", config.ListenPort))

	server := &http.Server{
		Addr:    ":" + config.ListenPort,
		Handler: newRemoteHttpHandler(),
	}

	err := server.ListenAndServe()
	if err != nil {
		logrus.Warning(fmt.Sprintf("Failed to start HTTP server: %v", err.Error()))
	}
}

func VerifyCallback(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	// 基本验证已经由TLS库完成，这里可以做额外验证
	if len(verifiedChains) == 0 || len(verifiedChains[0]) == 0 {
		return fmt.Errorf("no valid certificate chains\n")
	}

	// 获取客户端证书
	clientCert := verifiedChains[0][0]
	logrus.Info(fmt.Sprintf("clientCert CN: %v", clientCert.Subject.CommonName))

	// 示例1: 检查证书CN是否在允许列表中
	allowedCNs := []string{"client.example.com", "admin.example.com"}
	validCN := false
	for _, cn := range allowedCNs {
		if clientCert.Subject.CommonName == cn {
			validCN = true
			break
		}
	}
	if !validCN {
		return fmt.Errorf("certificate CN %s is not allowed\n", clientCert.Subject.CommonName)
	}

	// 示例2: 检查证书扩展或OU
	logrus.Info(fmt.Sprintf("OrganizationalUnit CN: %v", clientCert.Subject.OrganizationalUnit[0]))
	if !strings.Contains(clientCert.Subject.OrganizationalUnit[0], "China") {
		return fmt.Errorf("certificate OU is not from China department\n")
	}

	// 示例3: 检查证书有效期额外限制
	// if time.Now().AddDate(0, 0, 30).After(clientCert.NotAfter) {
	// 	return fmt.Errorf("certificate expires too soon")
	// }

	return nil
}

func StartHttpsServer(config *AppConfig) {
	logrus.Info(fmt.Sprintf("Remote httpsServer started on https://0.0.0.0:%s", config.ListenPort))

	var tlsConfig *tls.Config

	caFile := config.CAFile
	certFile := config.CertFile
	keyFile := config.KeyFile

	if config.VerifyClient {
		// 加载服务端证书和私钥
		serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("Failed to load server certificate: %v", err)
		}

		// 加载用于验证客户端证书的CA证书
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatalf("Failed to load CA certificate: %v", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			log.Fatalf("Failed to parse CA certificate")
		}

		tlsConfig = &tls.Config{
			ClientCAs:             caCertPool,                     // 用于验证客户端证书的CA
			ClientAuth:            tls.RequireAndVerifyClientCert, // 要求并验证客户端证书
			Certificates:          []tls.Certificate{serverCert},  // 服务端证书
			MinVersion:            tls.VersionTLS12,               // 设置最低TLS版本
			VerifyPeerCertificate: VerifyCallback,
		}
		certFile = ""
		keyFile = ""
	} else {
		tlsConfig = &tls.Config{
			ClientAuth: tls.NoClientCert,
		}
	}

	server := &http.Server{
		Addr:      ":" + config.ListenPort,
		TLSConfig: tlsConfig,
		Handler:   newRemoteHttpHandler(),
	}

	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		logrus.Warning(fmt.Sprintf("Failed to start HTTPS server: %v", err.Error()))
	}
}

func processPidAndMasterAccount(config *AppConfig) {
	var localAddr string
	if config.Debug {
		localAddr = ":62020"
		logrus.Info(fmt.Sprintf("local httpServer started on http://0.0.0.0:62020"))
	} else {
		localAddr = "127.0.0.1:62020"
		logrus.Info(fmt.Sprintf("local httpServer started on http://127.0.0.1:62020"))
	}

	server := &http.Server{
		Addr:    localAddr,
		Handler: newLocalHttpHandler(),
	}

	err := server.ListenAndServe()
	if err != nil {
		logrus.Warning(fmt.Sprintf("Failed to start HTTP server: %v", err.Error()))
	}
}
