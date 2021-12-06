package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
 * 一般静态编程语言都有严格的类型系统, 使编译器可以深入检查出程序员的出格举动,
 * 但是, 过于严格的类型系统会使得编程过于繁琐; go 语言试图让程序员在安全和灵活
 * 之间取得一个平衡, 它在提供严格的类型检查同时, 通过接口实现了对鸭子类型的支
 * 持, 使得安全动态的编程变的相对容易.
 *
 * go 的接口类型是对其他类型行为的抽象和概括, 因为接口类型不会和特定的实现
 * 细节绑定在一起, 通过这种抽象的方式可以让对象更加灵活和具有适应能力; 很多
 * 面向对象的语言都有相似的接口概念, 在 go 语言中接口类型的独特之处在于它满足
 * 隐式实现的鸭子类型.
 *
 * 鸭子类型: "只要走起路来像鸭子, 叫起来也像鸭子, 那么就可以把它当做鸭子"
 * (duck type) 是动态类型的一种风格; 在使用鸭子类型的语言中, 一个可以接受
 * 任意类型的对象的函数, 可以调用对象的 "走" 和 "叫"方法, 如果这些需要被
 * 调用的方法不存在, 将引发一个运行时错误; go 语言中的面向对象就是如此,
 * 如果一个对象只要看起来像是某种接口类型的实现, 那么它就可以作为该接口类型
 * 使用, 这种设计可以创建一个新的接口类型满足已经存在的具体类型(实现其方法),
 * 却不用去破坏这些类型原有的定义; 当使用的类型来自不受控制的包时这种设计尤其
 * 灵活有用;
 * go 语言的接口类型是延迟绑定(?), 可以实现类似虚函数的多态功能.(TODO)
 *
 */

// fmt.Printf() 函数的设计就是基于接口的, 它的正真功能由 fmt.Fprinf() 提供,
// 任意隐式满足 fmt.Stringer 接口的对象都可以打印, 不满足的依然可以通过
// 反射的技术打印, fmt.Fprint() 函数的签名如下:
// func Fprintf(w io.Writer, format string, arsgs ...interface{}) (int, error)

// 其中 io.Writer 是用于输出的接口, error 是内置的错误接口:
// type Writer interface {
// 	Write(p []byte) (n int, err error)
// }

// type error interface {
// 	Error() string
// }

// 可以通过定制自己的输出对象, 将每个字符转换为大写字符后输出:
type UpperWriter struct {
	io.Writer
}

// 继承了 io.Writer 的 Write, 并重写 Write 方法
func (p *UpperWriter) Write(data []byte) (n int, err error) {
	return p.Writer.Write(bytes.ToUpper(data))
	/*
	 * 在其他地方可以调用 p.Write(bytes.ToUpper(data)), 但是在重写的
	 * Write 函数中调用会出现循环调用, 发生栈溢出; 使用 htop 观察
	 * 内存使用突增1G
	 *
	 * runtime: goroutine stack exceeds(超过) 1000000000-byte(0.9G) limit
	 * runtime: sp=0xc0201603a8 stack=[0xc020160000, 0xc040160000]
	 * fatal error: stack overflow

	 * runtime stack:
	 * runtime.throw(0x4be4df, 0xe)
	 * /usr/local/go/src/runtime/panic.go:1116 +0x72
	 * runtime.newstack()
	 * /usr/local/go/src/runtime/stack.go:1067 +0x78d
	 * runtime.morestack()
	 * /usr/local/go/src/runtime/asm_amd64.s:449 +0x8f
	 */
}

func main() {
	fmt.Fprintln(&UpperWriter{os.Stdout}, "hello, world")
}

type UpperString string

func (s UpperString) String() string {
	return strings.ToUpper(string(s))
}
