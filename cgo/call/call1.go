package main

/* 在 go 中调用自己定义的 c 函数 */

/*
#include <stdio.h>

static void SayHello(const char* s) {
	puts(s);
}
*/
import "C"

func main() {
	C.SayHello(C.CString("heelo, cgo"))
}
