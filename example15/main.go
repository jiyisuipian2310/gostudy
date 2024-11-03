package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int, 10000) // 缓冲大小为10000
	go func() {
		defer close(c)
		for i := 0; i < 100; i++ {
			time.Sleep(1 * time.Second)
			fmt.Println("Sent:", i)
			c <- i
		}
	}()

	for {
		select {
		case val, ok := <-c:
			if !ok {
				fmt.Println("Channel closed!")
				return
			}
			fmt.Println("Received:", val)
		case <-time.After(200 * time.Millisecond):
			//fmt.Println("Timeout")
			//return
		}
	}
}
