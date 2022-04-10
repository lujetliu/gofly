package main

/*
	堆栈信息
	当程序遇到致命错误(如并发读写哈希表), 会 panic 输出一系列堆栈信息, 了解
	程序输出的堆栈信息, 有助于快速了解并定位问题

	堆栈信息是非常有用的排查问题的方式, 同时也可得知函数调用时传递的参数,
	加深理解 go 内部类型的结构, 以及值传递和指针传递的区别;
	go 可以配置 GOTRACEBACK 环境变量在程序异常终止时生成 core dump 文件,
	生成的文件可以有 dlv 或者 gdb 等高级调试工具进行分析调试(TODO)

*/

func trace(arr []int, a int, b int) int {
	panic("test trace")
	return 0
}

func main() {
	arr := []int{1, 2, 3}
	trace(arr, 5, 6)
}

// go run -gcflags="-l" stack_trace.go
// go build -gcflags="-l" stack_trace.go
// 使用  -gcflags="-l" 禁止函数的内联优化, 否则内联函数中不会打印函数的参数

// panic: test trace

// goroutine 1 [running]:
// main.trace({0xc000046770?, 0x404739?, 0x60?}, 0x0?, 0x0?)
//         /home/lucas/gogo/gofly/stack/stack_trace.go:11 +0x27
// main.main()
//         /home/lucas/gogo/gofly/stack/stack_trace.go:17 +0x57
// exit status 2

// 其中 0xc000046770, 0x404739, 0x60 对应 slice 的结构
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

// 文件名和行号后的 0x27/0x57 代表当前函数中的下一个要执行的指令位置, 其是距离
// 当前函数起始位置的偏移量

// 使用 go tool objdump -S -s "main.trace" ./stack_trace 命令反汇编执行文件

// TEXT main.trace(SB) /home/lucas/gogo/gofly/stack/stack_trace.go
// func trace(arr []int, a int, b int) int {
//   0x4551e0		493b6610		CMPQ 0x10(R14), SP
//   0x4551e4		7622			JBE 0x455208
//   0x4551e6		4883ec18		SUBQ $0x18, SP
//   0x4551ea		48896c2410		MOVQ BP, 0x10(SP)
//   0x4551ef		488d6c2410		LEAQ 0x10(SP), BP
// 	panic("test trace")
//   0x4551f4		488d05254f0000		LEAQ 0x4f25(IP), AX
//   0x4551fb		488d1dae260200		LEAQ 0x226ae(IP), BX
//   0x455202		e8f962fdff		CALL runtime.gopanic(SB)
//   0x455207		90			NOPL
// func trace(arr []int, a int, b int) int {
//   0x455208		e853cdffff		CALL runtime.morestack_noctxt.abi0(SB)
//   0x45520d		ebd1			JMP main.trace(SB)

// 可以看出当前函数的起始位置 0x4551e0  + 0x27 为 NOPL 指令(TODO), 对应了
// 当前 panic 函数的下一条指令
