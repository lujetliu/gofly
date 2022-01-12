package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

// TODO: net/rpc 源码

// 客户端 rpc 的实现原理
// go 的 rpc 最简单的使用方式是通过 client.Call() 方法进行同步阻塞调用

/*
	TODO: 相关代码

	// Call invokes the named function, waits for it to complete, and returns its error status.
	func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
	}
	通过 Client.Go() 方法进行一次异步调用, 返回表示这次调用的 Call 结构体, 然后
	等待 Call 的 Done 通道返回调用结果



	// Go invokes the function asynchronously. It returns the Call structure representing
	// the invocation. The done channel will signal when the call is complete by returning
	// the same Call object. If done is nil, Go will allocate a new channel.
	// If non-nil, done must be buffered or Go will deliberately crash.
	func (client *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
		call := new(Call)
		call.ServiceMethod = serviceMethod
		call.Args = args
		call.Reply = reply
		if done == nil {
			done = make(chan *Call, 10) // buffered.
		} else {
			// If caller passes done != nil, it must arrange that
			// done has enough buffer for the number of simultaneous
			// RPCs that will be using that channel. If the channel
			// is totally unbuffered, it's best not to run at all.
			if cap(done) == 0 {
				log.Panic("rpc: done channel is unbuffered")
			}
		}
		call.Done = done
		client.send(call)
		return call
	}
	client.send() 调用是线程安全的, 因此可以从多个 goroutine 同时向同一个 rpc
	链接发送调用指令




	func (call *Call) done() {
		select {
		case call.Done <- call:
			// ok
		default:
			// We don't want to block here. It is the caller's responsibility to make
			// sure the channel has enough buffer space. See comment in Go().
			if debugLog {
				log.Println("rpc: discarding Call reply due to insufficient Done chan capacity")
			}
		}
	}
	// 当调用完成或发送错误时, 将调用 call.done() 方法通知完成


*/

// 使用 Client.Go() 调用 ../server.go 中的 HelloService 服务
// func doClientWork(client *rpc.Client) {
// 	helloCall := client.Go("HelloService.Hello", "hello", new(string), nil)

// 	// ...
// 	helloCall = <-helloCall.Done
// 	if err := helloCall.Error; err != nil {
// 		log.Fatal(err)
// 	}

// 	args := helloCall.Args.(string)
// 	reply := helloCall.Reply.(string)
// 	fmt.Println(args, reply)
// }

func doClientWork(client *rpc.Client) {
	go func() {
		var keyChanged string
		err := client.Call("KVStoreService.Watch", 30, &keyChanged)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("watch:", keyChanged)
	}()

	err := client.Call("KVStoreService.Set", [2]string{"abc", "abc-value"},
		new(struct{}))
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 3)
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	doClientWork(client)
}
