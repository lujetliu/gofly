package main

import (
	"fmt"
	"log"
)

// 在 ./rule.go 中定义了接口规范后, 客户端就可以根据规范编写调用 rpc 的代码

// func main() {
// 	client, err := rpc.Dial("tcp", "localhost:8080")
// 	if err != nil {
// 		log.Fatal("dialing:", err)
// 	}

// 	var reply string
// 	err = client.Call(HelloServieName+".Hello", "world", &reply)
// 	// 相对于 ./client.go 唯一的变化

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(reply)
// }
// 但是 client.Call() 函数调用 rpc 方法比较繁琐, 同时参数的类型无法得到编译器
// 提供的安全保障, 因此为了简化客户端调用 rpc 函数, 在接口规范部分增加对客户
// 端的简单包装(./rule.go, 32), 在接口规范中针对客户端新增了 HelloServieClient
// 函数, 该类型也必须满足 HelloServieInterface 接口, 这样客户端用户就可以通过
// 接口对应的方法调用 rpc 函数, 同时提供了一个 DialHelloService 函数, 直接
// 拨号 HelloServie 服务; 此时简化客户端用户代码:
func main() {
	client, err := DialHelloService("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	err = client.Hello("world", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}

// 现在避免了 rpc 服务空间(HelloServieName), 方法名字(Hello)和参数类型不匹配
// 等低级错误的发生
// go run rule.go client1.go
