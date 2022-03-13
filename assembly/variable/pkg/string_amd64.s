#include "textflag.h"


/*
	其中以 · 表示是当前包的变量;
	在汇编中定义全局变量时, 无法为变量指定具体的类型, 需要在
	go 文件中声明, 在汇编中只关心变量的名字和内存大小;
	在 go 文件中声明对应的变量, 垃圾回收才会根据变量的类型来管理
	其中与指针相关的内存数据;

	type reflect.StringHeader struct {
			Data uintptr
			Len int
	}
*/

/*
GLOBL ·NameData(SB), NOPTR, $8
DATA ·NameData(SB)/8,$"gopher"

GLOBL ·Name(SB), $16
DATA ·Name+0(SB)/8, $·NameData(SB)
DATA ·Name+8(SB)/8, $6

其中 $·NameData(SB) 对应的是包变量 NameData 的地址, 可以将其看做一个常量
*/

// 因为字符串是只读的, 使用以下方法避免包外访问 ·NameData
GLOBL ·Name(SB), NOPTR, $24

DATA ·Name+0(SB)/8, $·Name+16(SB)
DATA ·Name+8(SB)/8, $6
DATA ·Name+16(SB)/8, $"gopher"


GLOBL ·Helloworld(SB), NOPTR, $16

// 局部变量
/* 当前文件内的私有变量 text(以<>为扩展名) */
GLOBL text<>(SB), NOPTR, $16 // 即使私有变量表示的字符串只有12个字符长, 
// 仍然需要将变量的长度扩展为2的指数倍
DATA text<>+0(SB)/8, $"Hello Wo"
DATA text<>+8(SB)/8, $"rld!"

// 用私有变量 text 对应的内存地址对应的常量初始化字符串头结构体中的
// Data, 手工指定 Len 部分为字符串的长度
DATA ·Helloworld+0(SB)/8, $text<>(SB)  // StringHeader.Data
DATA ·Helloworld+8(SB)/8, $12 // StringHeader.Len
// 字符串是只读类型, 要避免在汇编中直接修改字符串底层数据的内容
