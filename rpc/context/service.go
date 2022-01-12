package main

import (
	"fmt"
	"log"
	"net"
)

type HelloService struct {
	conn    net.Conn
	isLogin bool
}

func (h *HelloService) Login(request string, reply *string) error {
	if request != "user:password" {
		return fmt.Errorf("auth failed")
	}

	log.Println("login ok")
	h.isLogin = true // 有必要加锁?
	return nil
}

func (h *HelloService) Hello(request string, reply *string) error {
	if !h.isLogin {
		return fmt.Errorf("please login")
	}

	*reply = "hello: " + request + ", from" + h.conn.RemoteAddr().String()
	return nil
}
