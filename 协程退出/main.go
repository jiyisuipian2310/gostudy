package main

import (
	"fmt"
	"time"
)

/*
此代码演示了 通过停止通道 通知协程退出
广播机制：close(stopCh) 会通知所有监听该通道的协程
*/

func main() {
	// 创建停止通道
	stopCh := make(chan struct{})

	// 启动3个协议处理协程
	for i := 1; i <= 2; i++ {
		go worker(i, stopCh)
	}

	// 运行2秒后, 关闭通道，通知所有协程退出
	time.Sleep(2 * time.Second)
	fmt.Println("\n主线程发送停止信号...")
	close(stopCh)

	// 等待一秒让协程退出
	time.Sleep(1 * time.Second)
	fmt.Println("主程序退出")
}

func worker(id int, stopCh <-chan struct{}) {
	num := 0
	for {
		select {
		case <-stopCh:
			// 检测到退出信号
			fmt.Printf("协程 %d: 收到退出信号，停止处理\n", id)
			return
		default:
			num++
			fmt.Printf("协程 %d: 正在处理数据, num: %d\n", id, num)
			time.Sleep(100 * time.Microsecond)
		}
	}
}
