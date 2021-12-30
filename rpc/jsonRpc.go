package main

/* 跨语言的 rpc */

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

/*
 * go 标准库的 rpc 默认采用 go 特有的 gob 编码, 因此从其他语言调用 go 语言实现
 * 的 rpc 服务比较困难; go 语言的 rpc 框架有两个特色的设计:
 * - rpc 数据打包时可以通过插件实现自定义的编码和解码
 * - rpc 建立在抽象的 io.ReadWriteCloser 接口之上, 可以将 rpc 架设在不同的
 * 通信协议上.
 *
 */

//---------------------------------------------------------------server

// 通过官方自带的 net/rpc/jsonrpc 扩展实现一个跨语言的 rpc
type HelloService struct{}

func (h *HelloService) Hello(request string, reply *string) error {
	*reply = "hello, " + request
	return nil
}

// 基于 json 重新实现 rpc 服务
func main() {
	// TODO: RegisterName 源码
	// RegisterName 会检查传入的第二个参数(interface{}) 的方法集, 如果方法
	// 集没有任何方法会返回 error
	rpc.RegisterName("HelloService", new(HelloService))

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
		// jsonrpc.NewServerCodec(conn) 是针对服务器端的 json 编解码器
	}
}

//---------------------------------------------------------------client
func client() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("net.Dial:", err)
	}

	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	// 基于 conn 建立针对客户端的 json 编解码器

	var reply string

	err = client.Call("HelloService.Hello", "world", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}

// TODO: nc 工具的使用, golang 实现一个简单的 nc 工具
// 用一个普通的 tcp 服务代替 go 语言的 rpc 服务器获取客户端调用时发送的数据
// 格式, 如: nc -l 8080
//> go run jsonRpc.go(client函数改为main, 作为客户端), nc 作为服务端
//> {"method":"HelloService.Hello","params":["world"],"id":0}
// 客户端发送的数据是 json 编码的数据
// 请求的 json 数据对象在内部对应两个结构体:
// 客户端: clientRequest, 服务器端: serverRequest; 内容类似(TODO)
type clientRequest struct {
	Method string         `json:"method"`
	Params [1]interface{} `json:"params"`
	Id     uint64         `json:"id"`
}

type serverRequest struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
	Id     *json.RawMessage `json:"id"`
}

// 在获取到 rpc 调用对应的 json 数据后, 可以向 rpc 服务器发送 json 数据模拟
// rpc 方法调用()
//> go run jsonRpc.go(server函数改为main, 作为服务端), nc 作为客户端
//> echo -e '{"method":"HelloService.Hello","params":["world"],"id":0}' | nc localhost 8080
//> {"id":0,"result":"hello, world","error":null}
// 其中 id 对应输入的 id 参数, 对于顺序调用, id 不是必需的; 但 go 的 rpc 框架
// 支持异步调用, 当返回结果的顺序和调用的顺序不一致时, 可以通过 id 识别对应的
// 调用.
// 返回的 json 数据也对应内部的两个结构体:
// 客户端: clientResponse, 服务器端: serverResponse; 内容类似(TODO)
type clientResponse struct {
	Id     uint64           `json:"id"`
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
}

type serverResponse struct {
	Id     *json.RawMessage `json:"id"`
	Result interface{}      `json:"result"`
	Error  interface{}      `json:"error"`
}

/*
因此无论采用何种语言, 只要遵循同样的 json 结构, 以同样的流程就可以和 go
语言编写的 rpc 服务进行通信, 即实现了跨语言的 rpc.
*/
