package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/saintfish/chardet"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func init() {
	// 创建日志文件
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("无法创建日志文件: %v", err)
	}

	// 设置日志输出到文件和控制台
	logrus.SetOutput(file)

	// 设置日志级别
	logrus.SetLevel(logrus.InfoLevel)

	// 设置日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

// 持续监控目录的版本
func continuousFileReader(dirPath string, fileChan chan<- string, wg *sync.WaitGroup, stopChan <-chan bool) {
	defer wg.Done()

	logrus.Info(fmt.Sprintf("开始持续监控目录: %s", dirPath))

	// 用于记录已经处理过的文件，避免重复处理
	processedFiles := make(map[string]bool)

	for {
		select {
		case <-stopChan:
			logrus.Info("收到停止信号，结束监控")
			close(fileChan)
			return
		default:
			// 遍历目录
			filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return nil
				}

				// 只处理普通文件，且未处理过的
				if !info.IsDir() && !processedFiles[path] {
					if strings.HasSuffix(strings.ToLower(info.Name()), ".txt") {
						fileChan <- path
						processedFiles[path] = true
						logrus.Info(fmt.Sprintf("正在处理文件: %s", info.Name()))
					}
				}
				return nil
			})

			time.Sleep(1 * time.Second)
		}
	}
}

func detectEncodingWithChardet(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(content)
	if err != nil {
		return "", err
	}

	logrus.Info(fmt.Sprintf("文件: %s, 检测到的编码: %s (语言: %s, 置信度: %d)",
		filename, result.Charset, result.Language, result.Confidence))
	return result.Charset, nil
}

func fileOutput(fileChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for filename := range fileChan {
		encoding, err := detectEncodingWithChardet(filename)
		if err != nil {
			logrus.Info(fmt.Sprintf("检测编码失败: %v", err))
			continue
		}

		if encoding == "UTF-8" {
			newFilename := strings.Replace(filename, ".txt", ".json", -1)
			os.Rename(filename, newFilename)
			continue
		}

		if encoding == "GB-18030" || encoding == "EUC-KR" {
			outputFile := strings.Replace(filename, ".txt", ".json", -1)
			output, err := os.Create(outputFile) // 创建输出文件（UTF-8编码）
			if err != nil {
				logrus.Info(fmt.Sprintf("无法创建输出文件: %v", err))
				continue
			}
			defer output.Close()

			input, err := os.Open(filename)
			if err != nil {
				logrus.Info(fmt.Sprintf("无法打开输入文件: %v", err))
				continue
			}

			// 创建GB18030到UTF-8的转换器
			decoder := simplifiedchinese.GB18030.NewDecoder()
			reader := transform.NewReader(input, decoder)

			// 复制并转换内容
			_, err = io.Copy(output, reader)
			if err != nil {
				logrus.Info(fmt.Sprintf("转换失败: %v", err))
				continue
			}

			input.Close()
			//os.Remove(filename)
			logrus.Info(fmt.Sprintf("成功转换: %s -> %s", filename, outputFile))
		}
	}

	logrus.Info("文件输出协程结束")
}

func main() {
	fileChan := make(chan string, 20)
	stopChan := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)

	dirPath := "." // 可以修改为你要监控的目录路径

	// 启动持续监控的读取协程
	go continuousFileReader(dirPath, fileChan, &wg, stopChan)

	// 启动输出协程
	go fileOutput(fileChan, &wg)

	for {
		// 运行10秒后停止
		time.Sleep(10 * time.Second)
	}

	// 发送停止信号
	close(stopChan)

	// 等待所有协程完成
	wg.Wait()

	logrus.Info("程序执行完毕")
}
