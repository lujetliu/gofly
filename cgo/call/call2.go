package main

//void SayHello(const char* s);
import "C"

// TODO: 运行报错, 如何调用 c 源文件中的函数
func main() {
	C.SayHello(C.CString("hello, cgo"))
}

// TODO: c 源文件的编译, 静态库, 动态库
// 可以将对应的 c 文件 ./hello.c 编译打包为静态库或动态库文件供使用; 如果是以
// 静态库或动态库方式引用 SayHello() 函数, 需要将对应的 c 源文件移出当前目录
// (CGO构建程序会自动构建当前目录下的 c 源文件, 从而导致 c 函数名冲突)
