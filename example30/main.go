package main

import (
	"fmt"
	"sync"
)

type ChannelMap struct {
	cmd chan command
	m   map[string]int
}

type command struct {
	action string // "get", "set", "delete"
	key    string
	value  int
	result chan<- result
}

type result struct {
	value int
	ok    bool
}

func NewChannelMap() *ChannelMap {
	sm := &ChannelMap{
		cmd: make(chan command),
		m:   make(map[string]int),
	}
	go sm.run()
	return sm
}

func (m *ChannelMap) run() {
	for cmd := range m.cmd {
		switch cmd.action {
		case "get":
			value, ok := m.m[cmd.key]
			cmd.result <- result{value, ok}
		case "set":
			m.m[cmd.key] = cmd.value
		case "delete":
			delete(m.m, cmd.key)
		}
	}
}

func (m *ChannelMap) Set(key string, value int) {
	m.cmd <- command{action: "set", key: key, value: value}
}

func (m *ChannelMap) Get(key string) (int, bool) {
	res := make(chan result)
	m.cmd <- command{action: "get", key: key, result: res}
	r := <-res
	return r.value, r.ok
}

func (m *ChannelMap) Del(key string) {
	m.cmd <- command{action: "delete", key: key}
}

// 通过 channel 实现一个并发安全的 map
func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	cm := NewChannelMap()

	go func() {
		defer wg.Done()
		for i := 0; i <= 100000; i++ {
			cm.Set(fmt.Sprintf("key_%d", i), i+100)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i <= 100000; i++ {
			v, ok := cm.Get(fmt.Sprintf("key_%d", i))
			if !ok {
				fmt.Printf("key_%d isn't exist !\n")
			} else {
				fmt.Printf("key_%d, value: %d\n", i, v)
			}
		}
	}()

	wg.Wait()
}
