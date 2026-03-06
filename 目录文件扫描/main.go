package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type AppManager struct {
	wg       sync.WaitGroup
	workerCh chan string
	mainCh   chan string
	stopCh   chan struct{}
	stopOnce sync.Once // 确保只关闭一次stopCh
}

/* main.go 存在问题
通道关闭顺序问题：main() 中先关闭了 workerCh 和 mainCh，但工作协程可能还在尝试向 mainCh 发送数据

发送操作阻塞：当 mainCh 已关闭，向关闭的通道发送数据会导致 panic

停止信号处理不完整：工作协程收到 stopCh 信号后立即返回，但可能还有未完成的数据发送
*/

func (app *AppManager) subWorker(id int) {
	defer app.wg.Done()

	for {
		select {
		case value, ok := <-app.workerCh:
			if !ok {
				log.Printf("worker %d detect workerCh closed, exit\n", id)
				return
			}

			files, err := os.ReadDir(value)
			if err != nil {
				app.stopOnce.Do(func() { close(app.stopCh) })
				log.Printf("worker %d ReadDir Failed: %s, exit!\n", id, err)
				return
			}

			for _, file := range files {
				name := file.Name() // 获取文件/目录名称
				if file.IsDir() {
					app.mainCh <- (value + "/" + name)
				} else {
					fmt.Printf("worker %d, file: %s\n", id, value+"/"+name)
				}
			}
		case <-app.stopCh:
			log.Printf("worker %d receive quit notifation, exit!\n", id)
			return
		}
	}
}

func (app *AppManager) mainWorker(timeout int) {
	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	defer timer.Stop()

	for {
		select {
		case value := <-app.mainCh:
			app.workerCh <- value
			timer.Reset(time.Duration(timeout) * time.Second)
		case <-timer.C:
			log.Println("No data for 5 seconds, exiting")
			app.stopOnce.Do(func() { close(app.stopCh) })
			return
		case <-app.stopCh:
			log.Printf("mainWorker receive quit notifation, exit!\n")
			return
		}
	}
}

func main() {
	app := &AppManager{
		workerCh: make(chan string, 1000),
		mainCh:   make(chan string, 1000),
		stopCh:   make(chan struct{}),
	}

	for i := 0; i < 10; i++ {
		app.wg.Add(1)
		go app.subWorker(i + 1)
	}

	app.workerCh <- "/root"

	go app.mainWorker(5)

	<-app.stopCh

	close(app.workerCh)
	close(app.mainCh)

	app.wg.Wait()
}
