package main

import (
	"fmt"
	"sync"
)

type RWMutexMap struct {
	rw sync.RWMutex
	m  map[string]int
}

func NewRWMutexMap() *RWMutexMap {
	return &RWMutexMap{
		m: make(map[string]int),
	}
}

func (m *RWMutexMap) Set(key string, v int) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.m[key] = v
}

func (m *RWMutexMap) Get(key string) (int, bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	v, ok := m.m[key]
	return v, ok
}

func (m *RWMutexMap) Del(key string) {
	m.rw.Lock()
	defer m.rw.Unlock()
	delete(m.m, key)
}

//写锁就是互斥锁，不可以被多个 goroutine 同时持有,
//比如在调用写锁 m.rw.Lock() 时，如果有读锁或写锁已经被其他 goroutine 持有，
//则当前 goroutine 会被阻塞

// 读锁类似可重入锁，多个 goroutine 可以同时持有读锁，
// 比如在调用读锁 m.rw.RLock() 时，如果有其他 goroutine 也持有读锁，
// 则当前 goroutine 还可以继续获取这把读锁，但此时如果有 goroutine 尝试获取写锁则会被阻塞
func main() {
	m := NewRWMutexMap()

	// 并发读写 map
	go func() {
		for {
			m.Set("k", 1)
			fmt.Println("set k:", 1)
		}
	}()

	go func() {
		for {
			v, _ := m.Get("k")
			fmt.Println("read k:", v)
		}
	}()

	select {}
}
