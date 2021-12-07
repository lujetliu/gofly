package main

/*
 * go 语言通过简单特性的组合, 可以轻易实现鸭子面向对象和虚拟继承等高级特性.
 * TODO: 语言设计思想, 虚拟继承
 */

import (
	"fmt"
	"testing"
)

/*
 * go 对基础类型的类型一致性要求严格, 但对于接口类型的转换则非常灵活.
 * 对于基础类型(非接口类型)不支持隐式的转换, 无法将一个 int 类型的值
 * 直接赋值给 int64 类型, 也无法将 int 类型的值赋值给底层是 int 类型
 * 的新定义命名类型的变量.
 * 对象和接口之间的转换, 接口和接口之间的转换都可能是隐式的转换
 */

// var a io.ReadCloser = (*os.File)(f) // 隐式转换, *os.File 满足 io.ReadCloser
// var b io.Reader = a // 隐式转换, io.ReadCloser 满足 io.Reader 接口
// var a io.Closer = a // 隐式转换, io.ReadCloser 满足 io.Closer 接口
// var a io.Reader = c.(io.Reader) // 显式转换, io.Closer 不满足 io.Reader 接口

// 有时候对象和接口之间太灵活了, 需要人为的限制这种无意之间的适配, 常见的做法
// 是定义一个特殊方法来区分接口, 如下:
// - runtime 包中的 Error 接口就定义了一个特有的 RuntimeError() 方法, 用于避免
// 其他类型无意适配该接口.
// The Error interface identifies a run time error.
type Error interface {
	error // 继承

	// RuntimeError is a no-op function but
	// serves to distinguish types that are run time
	// errors from ordinary errors: a type is a
	// run time error if it has a RuntimeError method.
	RuntimeError()
}

// - 在 Protobuf 中, Message 接口也采用了类似的方法, 也定义了一个特有的
// ProtoMessage 方法, 用于避免其他类型无意中适配了该接口.
type Message interface {
	Reset()
	String() string
	ProtoMessage()
}

// 但上述做法只能算是"君子协定", 如果有人故意伪造一个 proto.Message 接口也很
// 容易, 再严格一点的做法就是给接口定义一个私有方法, 只有满足了这个私有方法
// 的对象才能满足这个接口, 而私有方法的名字是包含在包的绝对路径名的, 因此
// 只有在包内部实现这个私有方法才能满足这个接口. 测试包中的 testing.TB 接口
// 就是采用类似的方法:
// type TB interface {
// 	Cleanup(func())
// 	Error(args ...interface{})
// 	Errorf(format string, args ...interface{})
// 	Fail()
// 	FailNow()
// 	Failed() bool
// 	Fatal(args ...interface{})
// 	Fatalf(format string, args ...interface{})
// 	Helper()
// 	Log(args ...interface{})
// 	Logf(format string, args ...interface{})
// 	Name() string
// 	Skip(args ...interface{})
// 	SkipNow()
// 	Skipf(format string, args ...interface{})
// 	Skipped() bool
// 	TempDir() string

// 	// A private method to prevent users implementing the
// 	// interface and so future additions to it will not
// 	// violate(违反) Go 1 compatibility(兼容性).
// 	private()
// }

/*
 * 通过私有方法禁止外部对象实现接口的做法的代价:
 * - 导致接口只能在包内部使用, 外部包在正常情况下无法创建满足该接口的对象.
 * - 这种防护措施非绝对, 恶意的用户依然可以绕开这种保护机制.
 *
 *
 *
 * 通过在结构体中嵌入匿名类型成员, 可以继承匿名类型的方法; 被嵌入的匿名成员
 * 不一定是普通类型, 也可以是接口类型;
 * 可以通过嵌入匿名的 testing.TB 接口来伪造私有方法, 因为接口方法是
 * 延迟绑定(TODO), 所以在编译时私有方法是否真的存在并不重要.
 */
type TB struct {
	testing.TB // 嵌入匿名接口
}

func (t *TB) Fatal(args ...interface{}) { // 重写继承自 testing.TB 的 Fatal方法
	fmt.Println("TB.Fatal disabled")
}

func main() {
	var tb testing.TB = new(TB) // 将对象隐式转换为 testing.TB 接口类型
	// 再通过 testing.TB 接口来调用自己的 Fatal() 方法.
	tb.Fatal("hello, playground")
}

/*
 * TODO
 * 这种通过嵌入匿名接口或嵌入匿名指针对象来实现继承的做法其实是一种纯需继承,
 * 继承的只是接口指定的规范, 真正的实现在运行的时候才被注入.
 */
