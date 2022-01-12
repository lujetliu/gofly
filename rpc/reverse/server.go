package main

/*
 * TODO: 深入研究tcp的连接, 以彻底理解客户端和服务端角色
 *
 * 通常的rpc 基于客户/服务端结构, 但如果在内网提供的 rpc 服务, 在外网无法链接
 * 服务器, 这时可参考类似反向代理的技术, 首先从内网主动链接到外网的 tcp 服务器,
 * 然后基于 tcp 连接向外网提供 rpc 服务, TODO: 实验
 * TODO: http 是否可以用这种方式?
 */

import (
	"net"
	"net/rpc"
	"time"
)

type HelloService struct{}

func (h *HelloService) Hello(request string, reply *string) error {
	*reply = "hello, " + request
	return nil
}

func main() {
	rpc.Register(new(HelloService))
	// TODO: 与 ../server.go 中的 RegisterName 区别?
	// rpc.RegisterName("HelloService", new(HelloService))

	for {
		// 方向 rpc 的内网服务将不再主动提供 tpc 监听服务, 而是首先主动链接
		// 到对方的 tcp 服务器, 然后基于每个建立的 tcp 链接向对方提供 rpc 服务
		conn, _ := net.Dial("tcp", "localhost:8080")
		if conn == nil {
			time.Sleep(time.Second)
			continue
		}

		rpc.ServeConn(conn)
		conn.Close()
	}
}
