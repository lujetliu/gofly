

TEXT ·Swap(SB), NOSPLIT, $0
// TEXT ·Swap(SB), NOSPLIT, $0-32 

// NOSPLIT 标志会禁止汇编器为汇编函数插入栈分裂的代码, 指示叶子函数不进行栈
// 分裂; 对应 go 程序的  //go:nosplit 注释[注意: //后没有空格] (TODO: 栈分裂)


 // func Swap(a, b, c, int) int
 // func Swap(a, b, c, d int) 
 // func Swap() (a, b, c, d int)
 // func Swap() (a []int, d int)
 // 对于汇编函数, 只要函数的名字和参数大小一致就是相同的函数了, 而且在 go 汇编
 // 语言中, 输入参数和返回值参数没有任何区别(都是倒序入栈)







