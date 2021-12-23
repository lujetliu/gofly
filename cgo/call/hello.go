package main

// cgo 不仅用于 go 语言中调用 c 函数, 还可以导出 go 语言函数给 c 语言函数调用

import "C"
import "fmt"

//export SayHello
func SayHello(s *C.char) {
	fmt.Print(C.GoString(s))
}

// 通过 cgo 的 //export SayHello 指令将 go 语言实现的函数 SayHello() 导出为 C
// 语言函数, 为了适配 cgo 导出的 c 语言函数, 禁止了在函数的声明语句中的 const
// 修饰符;
// 这里其实有两个版本的  SayHello() 函数, 一个是 go 语言环境的, 另一个是 c 语言
// 环境的; cgo 生成的 c 语言版本的 SayHello() 函数最终会通过桥接代码调用 go
// 语言版本的 SayHello() 函数; TODO: 理解

/*
 * 通过面向 c 语言接口的编程技术, 不仅解放了函数的实现者, 同时也简化了函数的
 * 使用, 现在可以将 SayHello() 当做一个标准库的函数使用: ./hello1.go
 */
