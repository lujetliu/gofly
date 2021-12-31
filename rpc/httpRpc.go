package main

// TODO: 应用场景

import (
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

/*
 * go 语言内在的 rpc 框架支持在 http 协议上提供 rpc 服务, 但是框架的 http 服务
 * 同样采用了内置的 Gob 协议(TODO), 并且没有提供采用协议的接口(?), 因此从其他
 * 语言也无法访问.
 *
 */

type HelloService struct{}

func (h *HelloService) Hello(request string, reply *string) error {
	*reply = "hello, " + request
	return nil
}

// 在 http 协议上提供 jsonrpc 服务
func main() {
	rpc.RegisterName("HelloService", new(HelloService)) // TODO: 此行作用?

	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}

		// 基于 conn 构建针对服务器端的 json 编码解码器, 最后通过
		// rpc.ServeRequest 函数为每次请求处理一次 rpc 方法调用
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})

	http.ListenAndServe(":8080", nil)
}

// 模拟 rpc 调用
//> curl localhost:8080/jsonrpc -X POST --data '{"method":"HelloService.Hello","params":["world"],"id":0}'
//> {"id":0,"result":"hello, world","error":null}
