package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go func() {
		time.Sleep(10 * time.Second) // 10秒后停止定时器
		cancel()
	}()

	for {
		select {
		case <-ticker.C:
			fmt.Println("定时输出：", time.Now().Format("15:04:05"))
		case <-ctx.Done():
			fmt.Println("定时器停止")
			return
		}
	}
}
