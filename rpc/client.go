package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	// 通过 rpc.Dial() 拨号 rpc 服务, client.Call() 调用具体的 rpc 方法
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	// 参数1: 用点号连接的 rpc 服务名字和方法名字
	// 参数2, 3: Hello 方法的两个参数
	err = client.Call("HelloService.Hello", "world", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}
