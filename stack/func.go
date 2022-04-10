package main

/*
	go 的函数调用栈
	TODO:golang 汇编
	栈在不同场景下的定义不同, 有时候指一种先入后出的数据结构, 有时指操作系统
	组织内存的形式; 在大多数现代计算机系统中, 每个线程都有一个被称为栈的内存
	区域, 遵循一种后进先出(LIFO)的形式, 增长方向为从高地址到低地址;

	当函数执行时, 函数的参数, 返回值, 局部变量会被压入栈中, 当函数退出时, 这
	些数据会被回收; 当函数还没有退出就调用另一个函数时, 形成了一条函数调用链

	每个函数在执行过程中都使用一块栈内存来保存返回地址, 局部变量, 函数参数等,
	这块区域称为函数的栈帧(stack frame); 栈的大小会随着函数调用层级的增加而
	扩大, 随函数的返回而缩小, 即函数的调用层级越深, 消耗的栈空间越大;

	因为数据是以后进先出的方式添加和删除的, 所以基于堆栈的内存分配较简单,
	并且通常比基于堆的动态内存分配快的多, 当函数退出时, 堆栈上的内存会自动
	高效的回收(垃圾回收的最初形式)

	// TODO: 搭配 go 设计与实现 Page101 理解
	go 语言函数的参数和返回值存储在栈中, 很多主流的编程语言将参数和返回值存储
	在寄存器中(c); 存储在栈中的好处在于所有平台都可以使用相同的约定, 从而容易
	开发出可移植, 跨平台的代码, 同时这种方式简化了协程, 延迟调用和反射调用的
	实现;(寄存器的值不能跨函数调用, 存活), (TODO: 跨函数调用就是闭包?), 简化了
	垃圾回收器件的栈扫描和对栈扩容的处理;

	c语言同时使用寄存器(前6个)和栈传递参数, 使用 eax 寄存器传递返回值(c语言的
	函数只能有一个返回值); 而 go 语言使用栈传递参数和返回值; 对比:
	- C 语言的方式能狗极大的减少函数调用的额外开销, 但是也增加了实现的复杂度
		- cpu 访问寄存器的开销比访问栈的开销低几十倍
		- 需要单独处理函数参数过多的情况
	- Go 语言的方式能够降低实现的复杂度并支持多返回值, 但是牺牲了函数的性能
		- 不需要考虑超过寄存器数量的参数如何传递
		- 不需要考虑不同架构上的寄存器差异
		- 函数入参和出参的内存空间需要在栈上进行分配


	go 语言支持匿名函数和闭包, 闭包(closure) 是实现词法绑定的技术(TODO), 闭包
	包含了函数的入口地址和其关联的环境, 闭包和普通函数最大的区别在于闭包函数
	中可以引用闭包外的变量;

*/

func mul(a, b int) int {
	return a * b
}

func main() {
	mul(3, 4)
}

// TODO: golang 汇编
// "".mul STEXT nosplit size=29 args=0x18 locals=0x0
// 	0x0000 00000 (stack.go:3)	TEXT	"".mul(SB), NOSPLIT|ABIInternal, $0-24
// 	0x0000 00000 (stack.go:3)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
// 	0x0000 00000 (stack.go:3)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
// 	0x0000 00000 (stack.go:3)	MOVQ	$0, "".~r2+24(SP)
// 	0x0009 00009 (stack.go:4)	MOVQ	"".b+16(SP), AX
// 	0x000e 00014 (stack.go:4)	MOVQ	"".a+8(SP), CX
// 	0x0013 00019 (stack.go:4)	IMULQ	AX, CX
// 	0x0017 00023 (stack.go:4)	MOVQ	CX, "".~r2+24(SP)
// 	0x001c 00028 (stack.go:4)	RET
// 	0x0000 48 c7 44 24 18 00 00 00 00 48 8b 44 24 10 48 8b  H.D$.....H.D$.H.
// 	0x0010 4c 24 08 48 0f af c8 48 89 4c 24 18 c3           L$.H...H.L$..
// "".main STEXT size=71 args=0x0 locals=0x20
// 	0x0000 00000 (stack.go:7)	TEXT	"".main(SB), ABIInternal, $32-0
// 	0x0000 00000 (stack.go:7)	MOVQ	(TLS), CX
// 	0x0009 00009 (stack.go:7)	CMPQ	SP, 16(CX)
// 	0x000d 00013 (stack.go:7)	PCDATA	$0, $-2
// 	0x000d 00013 (stack.go:7)	JLS	64
// 	0x000f 00015 (stack.go:7)	PCDATA	$0, $-1
// 	0x000f 00015 (stack.go:7)	SUBQ	$32, SP
// 	0x0013 00019 (stack.go:7)	MOVQ	BP, 24(SP)
// 	0x0018 00024 (stack.go:7)	LEAQ	24(SP), BP
// 	0x001d 00029 (stack.go:7)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
// 	0x001d 00029 (stack.go:7)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
// 	0x001d 00029 (stack.go:8)	MOVQ	$3, (SP)
// 	0x0025 00037 (stack.go:8)	MOVQ	$4, 8(SP)
// 	0x002e 00046 (stack.go:8)	PCDATA	$1, $0
// 	0x002e 00046 (stack.go:8)	CALL	"".mul(SB)  (执行对 mul 函数的调用)
// 	0x0033 00051 (stack.go:9)	MOVQ	24(SP), BP
// 	0x0038 00056 (stack.go:9)	ADDQ	$32, SP
// 	0x003c 00060 (stack.go:9)	RET
// 	0x003d 00061 (stack.go:9)	NOP
// 	0x003d 00061 (stack.go:7)	PCDATA	$1, $-1
// 	0x003d 00061 (stack.go:7)	PCDATA	$0, $-2
// 	0x003d 00061 (stack.go:7)	NOP
// 	0x0040 00064 (stack.go:7)	CALL	runtime.morestack_noctxt(SB)
// 	0x0045 00069 (stack.go:7)	PCDATA	$0, $-1
// 	0x0045 00069 (stack.go:7)	JMP	0
// 	0x0000 64 48 8b 0c 25 00 00 00 00 48 3b 61 10 76 31 48  dH..%....H;a.v1H
// 	0x0010 83 ec 20 48 89 6c 24 18 48 8d 6c 24 18 48 c7 04  .. H.l$.H.l$.H..
// 	0x0020 24 03 00 00 00 48 c7 44 24 08 04 00 00 00 e8 00  $....H.D$.......
// 	0x0030 00 00 00 48 8b 6c 24 18 48 83 c4 20 c3 0f 1f 00  ...H.l$.H.. ....
// 	0x0040 e8 00 00 00 00 eb b9                             .......
// 	rel 5+4 t=17 TLS+0
// 	rel 47+4 t=8 "".mul+0
// 	rel 65+4 t=8 runtime.morestack_noctxt+0
// go.cuinfo.packagename. SDWARFINFO dupok size=0
// 	0x0000 6d 61 69 6e                                      main
// ""..inittask SNOPTRDATA size=24
// 	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
// 	0x0010 00 00 00 00 00 00 00 00                          ........
// gclocals·33cdeccccebe80329f1fdbee7f5874cb SRODATA dupok size=8
// 	0x0000 01 00 00 00 00 00 00 00                          ........
