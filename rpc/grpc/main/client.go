package main

import (
	context "context"
	"fmt"
	"log"

	grpc "google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	// grpc.Dial 负责和 gRPC 服务建立链接, 然后 NewHelloServiceClient()
	// 函数基于已经建立的链接构造 HelloServiceClient 接口对象, 通过接口
	// 定义的方法就可以调用服务器端对应的 gRPC 服务提供的方法
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &String{Value: "grpc"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.GetValue())
}

/*
	TODO: 模拟异步调用
	grpc 和标准库的 rpc 框架有一个区别, 即 grpc 生成的接口并不支持异步调用,
	不过可以在多个 goroutine 之间安全的共享 grpc 底层的 http/2 链接, 因此可
	以通过在另一个 goroutine 阻塞调用的方式模拟异步调用
*/
