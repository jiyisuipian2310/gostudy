package main

import (
	"fmt"
	"time"
)

/*
每个 case 都必须是一个通道操作，即 <- 或 <-chan。
如果有多个 case 可执行，则 select 语句将随机选择一个可执行的 case。
如果有多个 case 都准备好接收，则 select 语句将随机选择一个执行。
如果没有 case 可执行，且有 default 分支，则执行 default 分支。
如果没有 case 可执行，且没有 default 分支，则 select 语句将阻塞。
*/

func main() {

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		//time.Sleep(1 * time.Second)
		//c1 <- "one"
		for {
			c1 <- "one"
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() {
		//time.Sleep(2 * time.Second)
		//c2 <- "two"
		for {
			c1 <- "two"
			time.Sleep(100 * time.Millisecond)
		}
	}()

	//一直循环遍历2个通道是否有数据，有数据则打印出来
	//疑问：如果c1和c2一直都有数据，c1对应的代码是不是一直会被执行，而c2对应的代码永远不会被执行？
	for {
		select {
		case msg1 := <-c1:
			fmt.Println("received ", msg1)
		case msg2 := <-c2:
			fmt.Println("received ", msg2)
		default:
			time.Sleep(1 * time.Millisecond)
			break
		}
	}
}
