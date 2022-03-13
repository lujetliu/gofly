#include "textflag.h"


/*

type reflect.SliceHeader struct {
		Data uintptr
		Len int
		Cap int
}

//  切片类型和字符串类型的头结构体比较相似, 切片类型头结构体只是比字符串头
// 结构体多了 Cap 字段

*/

GLOBL ·SliceData(SB), NOPTR, $24 // var SliceData []byte

DATA ·SliceData+0(SB)/8, $text<>(SB)  // SliceHeader.Data
DATA ·SliceData+8(SB)/8, $12 // SliceHeader.Len
DATA ·SliceData+16(SB)/8, $16 // SliceHeader.Cap

GLOBL text<>(SB), NOPTR, $16 
DATA text<>+0(SB)/8, $"Hello Wo"
DATA text<>+8(SB)/8, $"rld!"
