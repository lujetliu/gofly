package main

/*
 TODO: 对上下文的理解(context); tcp 连接
 基于上下文可以针对不同客户端定制化的 rpc 服务, 可以通过为每个 tcp 连接提供独立
 的 rpc 服务来实现对上下文特性的支持
*/

import (
	"fmt"
	"log"
	"net"
)

type HelloService struct {
	conn    net.Conn
	isLogin bool // 属于一个单独的 tcp 连接, 对其的修改或读取不存在数据竞态
}

func (h *HelloService) Login(request string, reply *string) error {
	if request != "user:password" {
		return fmt.Errorf("auth failed")
	}

	log.Println("login ok")
	h.isLogin = true
	return nil
}

func (h *HelloService) Hello(request string, reply *string) error {
	if !h.isLogin {
		return fmt.Errorf("please login")
	}

	*reply = "hello: " + request + ", from " + h.conn.RemoteAddr().String()
	return nil
}
