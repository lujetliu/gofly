package main

//#include <hello.h>
import "C"

func main() {
	C.SayHello(C.CString("hello, world\n"))
}

// TODO: 如何运行
