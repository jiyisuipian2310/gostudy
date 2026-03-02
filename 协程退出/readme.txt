在 Go 语言中，协程（goroutine）之间通常通过 channel 来进行通信和同步
要通知其他协程退出，最常用、最推荐的方式是使用 context 包 或 关闭 channel 的方式来广播退出信号