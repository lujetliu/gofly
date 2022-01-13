package main

// TODO: http 1/2 协议的区别
// gRPC 是谷歌公司基于 protobuf 开发的跨语言的开源 rpc 框架, grpc 基于 http/2
// 协议设计, 可以基于一个 http/2 连接提供多个服务, 对移动设备更加友好.

/*
 TODO: 底层源码, http/2 协议, tcp, unix 套接字
  ----------------------------------------------
 |                     应用程序                |
 |                 生成的 Stub 代码            |
 |                 gRPC 内核+解释器            |
 |					  HTTP/2                   |
 |               安全(TLS/SSL或ALTS等)         |
 |              Unix 套接字  |     TCP         |
 |                                             |
  ----------------------------------------------

  最底层为 TCP  或 Unix 套接字协议, 在此之上是 HTTP/2 协议的实现, 然后构建了
  针对 go 语言的  gRPC 核心库(grpc 内核+解释器); 应用程序通过 gRPC 插件生成
  的 Stub 代码和 gRPC 核心库通信, 也可以直接和 gRPC 核心库通信.

*/

// 从 protobuf 的角度看, gRPC 就是一个针对服务接口生成代码的生成器;
// 在 ../protobuf/protoc-gen-go-netrpc.go 中实现了一个简单的 protobuf 代码
// 生成器插件, 但这只是适配标准库的 rpc 框架的, 以下开始使用 gRPC
// ./hello.proto
// ./main/server.go, ./main/client.go, ./main/service.go, ./main/hello.pb.go
// TODO: 以上代码功能虽已实现, 但还需要深入理解原理
