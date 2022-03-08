#include "textflag.h"

/*
  command-line-arguments
  runtime.gcdata: missing Go type information for global symbol var/pkg.Id: size 8
 
  其实Go汇编语言中定义的数据并没有所谓的类型, 
  每个符号只不过是对应一块内存而已，因此 Id 符号也是没有类型的;
  但是Go语言是再带垃圾回收器的语言，而Go汇编语言是工作在自动垃圾回收体系
  框架内的, 当Go语言的垃圾回收器在扫描到 Id 变量的时候，无法知晓该变量内部
  是否包含指针，因此就出现了上述错误, 错误的根本原因并不是 Id 没有类型, 而
  是 Id 变量没有标注是否会含有指针信息, 使用 "textflag.h" 的 NOPTR 标志其
  不包含指针数据;
*/

GLOBL ·Id(SB), NOPTR, $8

DATA ·Id+0(SB)/1, $0x37
DATA ·Id+1(SB)/1, $0x25
DATA ·Id+2(SB)/1, $0x00
DATA ·Id+3(SB)/1, $0x00
DATA ·Id+4(SB)/1, $0x00
DATA ·Id+5(SB)/1, $0x00
DATA ·Id+6(SB)/1, $0x00
DATA ·Id+7(SB)/1, $0x00

