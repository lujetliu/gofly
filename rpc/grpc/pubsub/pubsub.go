package main

/*
 在 ../../watch 中基于 go 内置的 rpc 实现了一个简化版的 Watch() 方法, 基于
 Watch() 的思路虽然也可以构造发布和订阅系统, 但是因为 rpc 缺乏流机制导致
 每次只能返回一个结果; 在发布订阅模式中, 由调用者发起的发布行为类似一个普通
 函数调用, 而被动的订阅者则类似 grpc 客户端单向流中的接收者; 本节中尝试基于
 grpc 的流特性构造一个发布和订阅系统.

 发布和订阅是一个常见的设计模式(TODO), 开源社区中已经存在很多实现, 其中 Docker
 项目提供了一个 pubsub 的极简实现; 参考本文件中的代码, 开始尝试基于 grpc 和
 pubsub 包提供一个跨网络的发布和订阅系统. ./pubsub.proto, ./main/pubsub.pb.go

 因此 Subscribe 是服务器端的单向流, 所以生成的 PubsubService_SubscribeServer
 接口只有 Send() 方法
	type PubsubService_SubscribeServer interface {
		Send(*String) error
		grpc.ServerStream
	}
*/

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/pkg/pubsub"
)

func main() {
	p := pubsub.NewPublisher(100*time.Millisecond, 10)

	//  以下函数可封装
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok { // TODO: 断言
			if strings.HasPrefix(key, "golang:") {
				return true
			}
		}
		return false
	})

	docker := p.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, "docker:") {
				return true
			}
		}
		return false
	})

	go p.Publish("hi")
	go p.Publish("golang: https://golang.org")
	go p.Publish("docker: https://www.docker.com/")
	time.Sleep(1)

	go func() {
		fmt.Println("golang topic: ", <-golang)
	}()

	go func() {
		fmt.Println("docker topic: ", <-docker)
	}()

	<-time.After(4 * time.Second)
}
