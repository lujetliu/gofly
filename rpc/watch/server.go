package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
	"time"
)

/*
   在很多系统中都提供了监控(watch)功能的接口, 当系统满足某种条件时 Watch() 方法
   返回监控的结果, 可以通过 rpc 框架实现一个基本的监控功能; 在 ./call.go 知
   client.send 是线程安全的, 所有也可以通过在不同的 goroutine 中同时并发阻塞调
   用 rpc 方法, 通过在一个独立的 goroutine 中调用 Watch() 方法进行监控.
*/

// 为了便于演示, 通过 rpc 构造一个简单的内存键值数据库
type KVStoreService struct {
	m      map[string]string
	fliter map[string]func(key string) // 对应每个 Watch() 调用时定义的过滤器函数列表
	mu     sync.Mutex
}

func NewKVStoreService() *KVStoreService {
	return &KVStoreService{
		m:      make(map[string]string),
		fliter: make(map[string]func(key string)),
	}
}

func (k *KVStoreService) Get(key string, value *string) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if v, ok := k.m[key]; ok {
		*value = v
		return nil
	}

	return fmt.Errorf("not found")
}

func (k *KVStoreService) Set(kv [2]string, reply *struct{}) error {
	// 用一个匿名的空结构体表示忽略输出参数 TODO: 为了满足 rpc 的函数规则?
	k.mu.Lock()
	defer k.mu.Unlock()

	key, value := kv[0], kv[1]

	if oldValue := k.m[key]; oldValue != value {
		for _, fn := range k.fliter {
			// 当修改某个键对应的值时会调用每一个过滤器函数
			fn(key)
		}
	}

	k.m[key] = value
	return nil
}

func (k *KVStoreService) Watch(timeoutSecond int, keyChanged *string) error {
	id := fmt.Sprintf("watch-%s-%03d", time.Now(), rand.Int())
	ch := make(chan string, 10) // buffered

	k.mu.Lock()
	k.fliter[id] = func(key string) { ch <- key }
	k.mu.Unlock()

	select {
	case <-time.After(time.Duration(timeoutSecond) * time.Second):
		return fmt.Errorf("timeout")
	case key := <-ch:
		*keyChanged = key
		return nil
	}
}

func main() {
	rpc.RegisterName("KVStoreService", NewKVStoreService())

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("listentcp error:", err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("accept error:", err)
	}
	rpc.ServeConn(conn)
}
