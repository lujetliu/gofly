// +build go1.10

package main

//void SayHello(_GoString_ s);
import "C"
import "fmt"

func main() {
	C.SayHello("hello, world\n")
}

//export SayHello
func SayHello(s string) {
	fmt.Print(s)
}

// 以上虽然看起来全是 go 代码, 但是执行的时候先是从 go 语言的 main() 函数到
// cgo 自动生成的 c 语言版本的 SayHello() 桥接函数, 最后又回到语言环境的
// SayHello() 函数; 这段代码包含了 cgo 的精华, TODO: 加深理解
