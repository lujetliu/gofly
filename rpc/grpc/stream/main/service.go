package main

import (
	context "context"
	"io"
)

type HelloServiceImpl struct{}

func (h *HelloServiceImpl) Hello(
	ctx context.Context, args *String,
) (*String, error) {
	reply := &String{Value: "hello: " + args.GetValue()}
	return reply, nil
}

func (h *HelloServiceImpl) Channel(stream HelloService_ChannelServer) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &String{Value: "hello: " + args.GetValue()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}

// 服务器端在循环中接收客户端发来的数据, 如果遇到 io.EOF 表示客户端流关闭,
// 如果函数退出表示服务器端流关闭, 生成返回的数据通过流发送给客户端, 双向流
// 数据的发送和接收是完全独立的行为
