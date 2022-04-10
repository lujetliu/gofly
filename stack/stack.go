package main

import (
	"os"
	"runtime/pprof"
)

/*
	栈转移和栈扩容(TODO)
	go 在线程上实现了用户态更加轻量的协程, 线程的大小一般是在创建时指定的,
	为了避免出现栈溢出(stack overflow) 的错误, 默认的栈大小会相对较大(2MB),
	而在 go 中每个协程都有一个栈, 每个栈的大小在初始化时候为2KB;
	go 的协程栈是可以扩容的, 最大的协程栈在64位操作系统中为1GB, 在32为系统中
	为250MB;
	栈的扩容和调整由运行时完成, 栈的管理主要为:(TODO)
	- 触发扩容的时机
	- 栈调整的方式


	扩容:
	为了应对频繁的栈调整, 对获取栈的内存进行了优化, 特别是小栈; 在Linux中,
	会对2KB/4KB/8KB/16KB 的小栈进行专门的优化, 即在全局及每个逻辑处理器(P)
	中预先分配这些小栈的缓存池, 从而避免频繁的申请堆内存;

	栈调试(TODO)
	有很多方式可以用于调试栈
	- 可以在源码级别进行调试
		go 语言在源码级别提供了栈相关的多种级别的调试, 用户调试栈的扩容和分配等;
		但这些静态变量没有暴露给用户, 要使用这些变量, 需要直接修改 go 的源码并
		重新进行编译:
		const (
			stackDebug = 0
			stackFromSystem = 0
			stackFaultOnFree = 0
			stackPoisonCopy = 0
			stackNoCache = 0
		)
	- 使用 runtime/debug.PrintStack 方法
	- 使用标准库 pprof 获取某一时刻的堆栈信息(pprof 是特征分析的强大工具, TODO)
		利用 pprof 的协程栈调试, 可以分析是否发生协程泄露, 当前程序使用最多的
		函数是什么, 并分析 cpu 的瓶颈, 可视化等特性;

*/

func main() {
	debugPprof(99)
}

// func debugStack(num int) {
// 	debug.PrintStack() // 打印当前的堆栈信息
// }

func debugPprof(num int) {
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1) // 获取当前时刻协程的栈信息
}

// goroutine profile: total 1
// 1 @ 0x459de5 0x4a59f5 0x4a580d 0x4a276b 0x4aee1d 0x4aedbe 0x435512 0x45e501
// #	0x459de4	runtime/pprof.runtime_goroutineProfileWithLabels+0x24	/usr/local/go/src/runtime/mprof.go:753
// #	0x4a59f4	runtime/pprof.writeRuntimeProfile+0xb4			/usr/local/go/src/runtime/pprof/pprof.go:725
// #	0x4a580c	runtime/pprof.writeGoroutine+0x4c			/usr/local/go/src/runtime/pprof/pprof.go:685
// #	0x4a276a	runtime/pprof.(*Profile).WriteTo+0x14a			/usr/local/go/src/runtime/pprof/pprof.go:332
// #	0x4aee1c	main.debugPprof+0x3c					/home/lucas/gogo/gofly/stack/stack.go:52
// #	0x4aedbd	main.main+0x1d						/home/lucas/gogo/gofly/stack/stack.go:45
// #	0x435511	runtime.main+0x211					/usr/local/go/src/runtime/proc.go:250
