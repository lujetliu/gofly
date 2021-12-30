package main

import "net/rpc"

/*
 * 在涉及 rpc 的应用中, 作为开发人员至少有3种角色:
 * - 服务器端实现 rpc 方法的开发人员
 * - 客户端调用 rpc 方法的开发人员
 * - 制定服务端和客户端 rpc 接口规范的设计人员
 *
 * 为了利于后期的维护和工作的切割, 把以上几种角色的工作分成不同的部分
 *
 */

// 重构 HelloService 服务, 首先明确服务的名字和接口
// 将 rpc 服务的接口分为3部分:
// 1, 服务的名字
const HelloServieName = "path/to/pkg.HelloService" // 为了避免名字冲突, 在 rpc
// 服务的名字中增加了包路径前缀(rpc服务抽象的包路径, 并不等价 go 的包路径)

// 2, 服务要实现的详细方法列表
type HelloServiceInterface = interface {
	Hello(request string, reply *string) error
}

// 3, 注册该类型服务的函数
func RegisterHelloService(svc HelloServiceInterface) error {
	// 注册服务时, 编译器会要求传入的对象满足 HelloServiceInterface 接口
	return rpc.RegisterName(HelloServieName, svc)
}

type HelloServieClient struct {
	*rpc.Client
}

var _ HelloServiceInterface = (*HelloServieClient)(nil) // TODO:?

func DialHelloService(network, address string) (*HelloServieClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HelloServieClient{Client: c}, nil

}

func (h *HelloServieClient) Hello(request string, reply *string) error {
	return h.Client.Call(HelloServieName+",Hello", request, reply)
}
