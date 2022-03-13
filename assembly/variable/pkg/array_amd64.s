
GLOBL ·Num(SB), $16
DATA ·num+0(SB)/8, $1
DATA ·num+8(SB)/8, $0

/*
	go 中的数组是值类型, 底层不含指针;
	此处的汇编代码定义变量时, 不需要 NOPTR 标志, 因为 go 编译器在 ./pkg.go 中
	声明 Num 类型中推导出该变量内部没有指针数据;
*/
