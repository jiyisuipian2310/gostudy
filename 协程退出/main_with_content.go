package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Server struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewServer() *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		ctx:    ctx,
		cancel: cancel,
	}
}

// 启动多个工作goroutine
func (s *Server) Start(workerCount int) {
	for i := 0; i < workerCount; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
}

// worker工作函数
func (s *Server) worker(id int) {
	defer s.wg.Done()
	num := 0
	for {
		select {
		case <-s.ctx.Done():
			time.Sleep(time.Duration(id) * 100 * time.Millisecond)
			fmt.Printf("Worker %d 已退出\n", id)
			return
		default:
			num++
			fmt.Printf("Worker %d 处理任务, num: %d\n", id, num)
			time.Sleep(200 * time.Microsecond)
		}
	}
}

// 优雅关闭
func (s *Server) Shutdown() {
	fmt.Println("正在关闭服务...")
	s.cancel() // 发送取消信号

	// 等待所有goroutine退出，带超时控制
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("所有goroutine已退出")
	case <-time.After(5 * time.Second):
		fmt.Println("等待超时，强制退出")
	}
}

func main() {
	server := NewServer()

	// 启动5个worker
	server.Start(3)
	fmt.Println("服务已启动，运行中...")

	// 模拟服务运行8秒
	time.Sleep(5 * time.Second)

	// 优雅关闭
	server.Shutdown()
	fmt.Println("程序结束")
}
