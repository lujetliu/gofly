package main

/*
   rpc 是远程过程调用, 因此每次调用的参数和返回值不能太大, 否则将影响每次调用的
   响应时间, 因此传统的 rpc 方法调用对上传和下载较大数据的场景并不适合; 同时传
   统 rpc 模式也不适用于时间不确定的订阅和发布模式(TODO), 为此, grpc 框架针对
   服务器端和客户端分别提供了流特性.

   TODO: 流特性底层原理, 与一个 tpc 连接的关系
*/

// 服务器端或客户端的单向流是双向流的特例, 在 HelloService 增加一个支持双向流的
// Channel() 方法(./hello.proto)
// service HelloService {
// 	rpc Hello (String) returns (String);
// 	rpc Channel (stream String) returns (stream String);
// }
// 关键字 stream 指定启用流特性, 参数部分是接收客户端参数的流, 返回值是返回给
// 客户端的流

/*
	// 重新生成代码(./main/hello.pb.go), 接口中新增了 Channel() 方法的定义

	// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
	type HelloServiceClient interface {
		Hello(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
		Channel(ctx context.Context, opts ...grpc.CallOption) (HelloService_ChannelClient, error)
	}

	// HelloServiceServer is the server API for HelloService service.
	type HelloServiceServer interface {
		Hello(context.Context, *String) (*String, error)
		Channel(HelloService_ChannelServer) error
	}


	// 客户端返回的 HelloService_ChannelClient 可用于和服务器端进行双向通信
	type HelloService_ChannelClient interface {
		Send(*String) error
		Recv() (*String, error)
		grpc.ClientStream
	}
	// 服务器端的参数是 HelloService_ChannelServer , 用于和客户端的双向通信
	type HelloService_ChannelServer interface {
		Send(*String) error
		Recv() (*String, error)
		grpc.ServerStream
	}
	// 可以看出客户端和服务器端的流辅助接口均定义了 Send() 和 Recv(), 用于流
	// 数据的双向通信

*/
