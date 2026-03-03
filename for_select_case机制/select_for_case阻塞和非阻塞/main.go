package main

import (
	"log"
	"sync"
	"time"
)

// 阻塞操作：如果没有 case 满足条件且没有 default，select 会阻塞
func coroutine_block1(wg *sync.WaitGroup) {
	ticker := time.NewTicker(1 * time.Second)
	defer wg.Done()
	defer ticker.Stop()

	number := 0

	for {
		select {
		case <-ticker.C:
			number++
			log.Println("coroutine_block1 timeout output number: ", number)
			if number == 5 {
				return
			}
		}
	}
}

// 阻塞操作：如果没有 case 满足条件且没有 default，select 会阻塞
// 如果 select 中只有 ticker 相关的case，可以改成如下写法：
func coroutine_block2(wg *sync.WaitGroup) {
	ticker := time.NewTicker(1 * time.Second)
	defer wg.Done()
	defer ticker.Stop()

	number := 0

	for range ticker.C {
		number++
		log.Println("coroutine_block2 timeout output number: ", number)
		if number == 5 {
			return
		}
	}
}

// 非阻塞操作：如果没有 case 满足条件，存在default，select 会执行default中的代码
func coroutine_nonblock(wg *sync.WaitGroup) {
	ticker := time.NewTicker(2 * time.Second)
	defer wg.Done()
	defer ticker.Stop()

	number := 0

	for {
		select {
		case <-ticker.C:
			number++
			log.Println("timeout output number: ", number)
			if number == 100 {
				return
			}
		default:
			number++
			log.Println("this is default case", number)
			time.Sleep(200 * time.Millisecond)
			if number == 100 {
				return
			}
		}
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go coroutine_block1(&wg)

	wg.Add(1)
	go coroutine_block2(&wg)

	//wg.Add(1)
	//go coroutine_nonblock(&wg)

	wg.Wait() // 等待所有协程结束

	log.Println("all coroutine quit, app exit !")
}
