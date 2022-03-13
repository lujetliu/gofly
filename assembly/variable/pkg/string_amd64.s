#include "textflag.h"

/*
GLOBL ·NameData(SB), NOPTR, $8
DATA ·NameData(SB)/8,$"gopher"

GLOBL ·Name(SB), NOPTR, $16
DATA ·Name+0(SB)/8, $·NameData(SB)
DATA ·Name+8(SB)/8, $6

其中 $·NameData(SB) 对应的是包变量 NameData 的地址, 可以将其看做一个常量
*/

GLOBL ·Name(SB), NOPTR, $24

DATA ·Name+0(SB)/8, $·Name+16(SB)
DATA ·Name+8(SB)/8, $6
DATA ·Name+16(SB)/8, $"gopher"

