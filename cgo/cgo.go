package main

/*
 * C 语言作为一个通用语言, 很多库会选择提供一个 c 兼容的 API, 然后用其他不同的
 * 编程语言实现; go 语言通过自带的叫 cgo 的工具支持 c 语言函数调用, 同时可以
 * 用 go 语言导出 c 动态库接口给其他语言使用.
 */

//#include <stdio.h>  // 引入 c 语言的 <stdio.h> 文件
import "C"

// import "fmt"

// func main() {
// 	fmt.Println("hello, cgo")
// }

// 通过 import "C" 语句启用 cgo 特性, 主函数中虽然没有调用 cgo 的相关函数, 但是
// go build 命令会在编译和链接阶段启动 gcc 编译器, 这已经是一个完整的 cgo 程序.

// 基于 c 标准库函数输出字符串
func main() {
	C.puts(C.CString("hello, cgo"))
	// 使用 C.CString 创建 c 语言字符串, 而且使用 puts() 函数向标准输出打印
	/*
	 * TODO: C 语言, GC
	 * 没有释放使用的 C.CString 创建的 C 语言字符串会导致内存泄露, 但是对这个
	 * 小程序来说足够安全, 程序退出后操作系统会自动回收程序的所有资源.
	 */
}

// 除了以上的调用方式, 还有 ./call/call1.go, ./call/call2.go 等其他的调用方式

/*
 * 在采用面向 c 语言 API 编程中, 可以彻底解放模块实现者的加锁, 实现者可以用任
 * 何编程语言实现模块, 只要满足公开的 API 约定, 可以用 c 语言(./call/hello.c)
 * 实现 SayHello() 函数, 也可以使用更复杂的 c++ 语句实现(./call/hello.cpp),
 * 当然也可以用汇编语言实现, 甚至 go 语言来重新实现 SayHello() 函数.
 *
 */
