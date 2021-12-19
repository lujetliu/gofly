package main

/* 剖析异常, TODO: panic, recover 源码 */

import (
	"errors"
	"fmt"
	"log"
)

// panic() 支持抛出任意类型的异常(而不仅是 error 类型的错误), recover() 函数
// 调用的返回值和 panic() 函数的输入参数类型一致, 函数签名如下:
// func panic(interface{})
// func recover() interface{}

// func main() {
// 	if r := recover(); r != nil {
// 		log.Fatal(r)
// 	}

// 	panic(123)

// 	if r := recover(); r != nil { // Warning: Unreachable code
// 		log.Fatal(r)
// 	}
// }

// 上面程序中的两个 recover() 调用都不能捕获任何异常, 在第一个 recover() 调用
// 执行时, 函数必然是在正常的非异常执行流程中, 这时 recover() 调用将返回 nil;
// 发生异常时, 第二个 recover() 调用不会被执行到, 因为 panic() 调用会导致马上
// 执行该函数体中 panic() 语句前已经注册 defer 的函数后返回
// 在非 defer 语句中执行 recover() 调用无意义

/*
 * 对 recover() 的调用有严格的要求: 必须在 defer() 函数中直接调用 recover(),
 * 如果在defer() 中调用的是 recover() 函数的包装函数, 异常将不会被捕获
 */

// 有时候希望包装自己的 MyRecover() 函数, 在内部增加必要的日志信息然后再调用
// recover(), 这是错误的做法
func f1() {
	defer func() {
		// 无法捕获异常
		if r := MyRecover(); r != nil {
			fmt.Println(r)
		}
	}()
	panic(1)
}

func MyRecover() interface{} {
	log.Println("trace...")
	return recover()
}

// 在嵌套的 defer() 函数中调用 recover(), 也会导致无法捕获异常
func f2() {
	defer func() {
		defer func() {
			// 无法捕获异常
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
	}()
	panic(1)
}

// 直接在 defer 语句中调用 MyRecover() 函数, 可以正常工作
func f3() {
	defer MyRecover() // 可捕获异常
	panic(1)
}

// 但是如果 defer 语句直接调用 recover() 函数, 则依然不能正常捕获异常
func f4() {
	defer recover()
	panic(1)
}

/*
 * TODO: 函数调用, 栈帧
 * 两层嵌套的 defer() 函数中直接调用 recover() 和一层 defer() 函数中调用
 * 包装的 MyRecover() 函数一样, 都是经过了两个函数帧才到达真正的 recover()
 * 函数, 此时, goroutine 对应的上一级栈帧中已经没有异常信息了.
 *
 *
 * 必须要和有异常的栈帧只隔一个栈帧, recover() 函数才能正常捕获异常, 即
 * recover() 函数捕获的是祖父一级调用函数栈帧的异常(刚好可以跨越一层
 * defer() 函数)
 */

// 当希望捕获到的异常转为错误时, 可针对不同类型分别进行处理, 从而返回
// 原始的信息
func foo() (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("Unknown panic:%v", r)
			}
		}
	}()
	panic("TODO")
}
