package main

import context "context"

// 基于./hello.pb.go 中 服务端的 HelloServiceServer 接口可以重新实
// 现 HelloService 服务
type HelloServiceImpl struct{}

func (h *HelloServiceImpl) Hello(
	ctx context.Context, args *String,
) (*String, error) {
	reply := &String{Value: "hello: " + args.GetValue()}
	return reply, nil
}
