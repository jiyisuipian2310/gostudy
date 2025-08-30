package main

import "fmt"

/*
go 通道：有缓冲的通道和无缓冲的通道
有缓冲的通道：ch := make(chan string, 10)
	向有缓冲通道写入数据时，如果缓冲区未满，则写操作将成功，程序将继续执行
	缓冲区已满，则写操作将阻塞，直到有空闲缓冲区可用时，方可写入
无缓冲的通道：ch := make(chan string)
	通道缓冲区的大小为 0
	无缓冲通道的发送和接收操作都是阻塞的，因此必须有接收者准备好接收才能进行发送操作，反之亦然
*/

func writeChannel(ch chan string, done chan struct{}) {
	count := 0
	for {
		ch <- fmt.Sprintf("hello_%d", count)
		count++
		if count == 10000 {
			break
		}
	}

	done <- struct{}{}
	close(ch)
}

//向一个已经关闭的通道读取数据，第二个返回值是false
func readChannel(ch chan string, done chan struct{}) {
	for {
		str, ok := <-ch
		fmt.Println(fmt.Sprintf("---> %s, ok=%v", str, ok))
		if !ok {
			break
		}
	}

	done <- struct{}{}
}

func main() {
	done := make(chan struct{})
	ch := make(chan string)
	go writeChannel(ch, done)
	go readChannel(ch, done)

	//等待2个协程都退出
	for i := 0; i < 2; i++ {
		fmt.Println(<-done)
	}
}
