package main

// 在开始的例子中, 全部的 cgo 代码都在一个 go 文件中(../cgo.go), 然后通过面向
// c 接口的编程技术将 SayHello() 分别拆分到不同的 c 文件(./hello.c, ./hello.h),
// 而 main 依然是 go 文件; 而用 go 函数重新实现了 c 语言接口的 SayHello()函数;
// 现在将例子中的几个文件重新合并到一个 go 文件中.

//void SayHello(char* s)
import "C"
import "fmt"

func main() {
	C.SayHello(C.CString("hello, world\n"))
}

//export SayHello
func SayHello(s *C.char) {
	fmt.Println(C.GoString(s))
}

// TODO: 如何运行
// 现在版本的 cgo 代码中 c 语言代码的比例已经很少了, 但是依然可以进一步以 go
// 语言思维来提炼 cgo 代码, 分析发现 SayHello() 的参数如果可以直接使用 go
// 字符串是最直接的, 在 go1.10 中增加了一个 _GoString_ 预定义的 c 语言类型,
// 用于表示 go 语言字符串: ./main1.go
