package main

/*
	函数标识符通过 TEXT 汇编指令定义, 表示该行开始的指令定义在 TEXT 内存段,
	TEXT 后的指令一般对应函数的实现, 但是对 TEXT 指令本身来说并不关心后面
	是否有指令, 因此 TEXT 和 LABEL 定义的符号是类似的, 区别只是 LABEL 是用
	于跳转标号, 但是本质上都是通过标识符映射一个内存地址;
	语法如下:
	TETX symbol(SB), [flags, ] $framesize[-argsize]
	TEXT 指令, 函数名, 可选的 flags  标志, 函数帧大小, 和可选的函数参数大小
	- 函数名后面的(SB), 表示函数名符号相对于伪寄存器 SB 的偏移量, 二者组合在一起
		是绝对地址;
	- 标志部分用于指示函数的一些特殊行为
		- NOSPLIT: 指示叶子函数不进行栈分裂(TODO: 叶子函数, 栈分裂);
		- WRAPPER: 表示是一个包装函数, 在 panic 和 runtime.caller 等某些处理
			函数帧的地方不会增加函数帧计数 (TODO: 函数帧);
		- NEEDCTXT: 表示需要一个上下文参数, 一般用于闭包函数
	- framesize 表示函数的局部变量需要多大栈空间, 其中包含调用其他函数时准备
		调用参数的隐式空栈空间(TODO)
	- argsize 表示参数大小, 因为go汇编推荐在go代码中声明函数, 所以编译器可以从
		go 的函数声明中推导出函数参数的大小, 因此 argsize 可以省略

*/

//go:nosplit
func Wap(a, b int) (int, int) // 实现 ./swap.s
