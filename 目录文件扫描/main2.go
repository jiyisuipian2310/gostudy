package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

//main2.go 代码不存在问题

/*
阻塞读取通道的值
<-timer.C  // 会阻塞直到通道有值可读

非阻塞读取通道的值
select {
	case <-timer.C:  // 如果有值，就读取
	default:         // 如果没有值，立即执行 default
}
*/

type AppManager struct {
	wg       sync.WaitGroup
	workerCh chan string
	mainCh   chan string
	stopCh   chan struct{}
	stopOnce sync.Once
}

func (app *AppManager) subWorker(id int) {
	defer app.wg.Done()

	for {
		select {
		case value, ok := <-app.workerCh:
			if !ok {
				log.Printf("worker %d: workerCh closed, exit\n", id)
				return
			}

			files, err := os.ReadDir(value)
			if err != nil {
				app.stopOnce.Do(func() { close(app.stopCh) })
				log.Printf("worker %d ReadDir Failed: %s, exit!\n", id, err)
				return
			}

			for _, file := range files {
				name := file.Name()
				if file.IsDir() {
					subdir := ""
					if value == "/" {
						subdir = "/" + name
					} else {
						subdir = value + "/" + name
					}

					// 非阻塞发送到 mainCh
					select {
					case app.mainCh <- subdir:
						// 发送成功
					case <-app.stopCh:
						// 收到停止信号，退出
						log.Printf("worker %d received stop signal while sending to mainCh\n", id)
						return
					default:
						// mainCh 已满，但还没收到停止信号，可以重试或继续
						// 这里选择等待一小段时间后重试
						select {
						case app.mainCh <- subdir:
						case <-app.stopCh:
							return
						case <-time.After(100 * time.Millisecond):
							// 超时后继续尝试
						}
					}
				} else {
					fmt.Printf("worker %d, file: %s\n", id, value+"/"+name)
				}
			}
		case <-app.stopCh:
			log.Printf("worker %d receive quit notification, exit!\n", id)
			return
		}
	}
}

func (app *AppManager) mainWorker(timeout int) {
	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	defer timer.Stop()

	// 确保定时器初始状态正确
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(time.Duration(timeout) * time.Second)

	for {
		select {
		case value, ok := <-app.mainCh:
			if !ok {
				log.Println("mainCh closed, exiting mainWorker")
				return
			}

			// 非阻塞发送到 workerCh
			select {
			case app.workerCh <- value:
				// 重置定时器
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(time.Duration(timeout) * time.Second)
			case <-app.stopCh:
				log.Println("mainWorker received stop signal while sending")
				return
			}

		case <-timer.C:
			log.Printf("No data for %d seconds, exiting", timeout)
			app.stopOnce.Do(func() { close(app.stopCh) })
			return

		case <-app.stopCh:
			log.Printf("mainWorker receive quit notification, exit!\n")
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

	// 启动工作协程
	for i := 0; i < 10; i++ {
		app.wg.Add(1)
		go app.subWorker(i + 1)
	}

	// 启动主工作协程
	go app.mainWorker(5)

	// 发送初始目录
	select {
	case app.workerCh <- "/root":
	case <-time.After(time.Second):
		log.Println("Timeout sending initial directory")
	}

	// 等待停止信号
	<-app.stopCh
	log.Println("Stop signal received, shutting down...")

	// 关闭 workerCh，让工作协程退出
	close(app.workerCh)

	// 等待所有工作协程完成
	app.wg.Wait()

	// 最后关闭 mainCh
	close(app.mainCh)

	log.Println("Application shutdown complete")
}
