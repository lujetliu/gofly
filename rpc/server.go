package main

import (
	"log"
	"net"
	"net/rpc"
)

// 基于 rpc 的打印例子
type HelloService struct{}

/*
 * TODO: rpc 规则, rpc 相关函数实现
 * Hello() 方法必须满足 go 语言的 RPC 规则: 方法只能有两个可序列化(TODO)的参数,
 * 其中第二个参数是指针类型, 并且返回一个 error 类型, 同时必须是公开的方法
 *
 */
func (h *HelloService) Hello(request string, reply *string) error {
	*reply = "hello, " + request
	return nil
}

func main() {
	rpc.RegisterName("HelloService", new(HelloService))
	// RegisterName() 函数调用会将对象类型中所有满足rpc规则的对象方法注册为
	// rpc函数, 所有注册的方法会放在 HelloService 服务的空间之下, 然后建立
	// 一个唯一的 TCP 链接, 并且通过 rpc.ServeConn() 函数在该 tcp 链接上为
	// 客户端提供 rpc 服务, TODO: 这其实是在不同的语言, 不用的网络实体间以
	// 面向对象的模型搭建一个沟通的桥梁; rpc 思想

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("Accept error:", err)
	}
	rpc.ServeConn(conn)
}
