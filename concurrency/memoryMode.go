package main

import (
	"fmt"
	"sync"
)

/*
 * TODO: cache, 页表, 储存器
 * 顺序一致性内存模型
 * 知乎: https://zhuanlan.zhihu.com/p/422848235
 */

// 如果只是简单的在线程之间进行数据同步, 原子操作已经提供了一些同步的保障,
// 不过这个保障有一个前提: 顺序一致性内存模型
var a string
var done bool

func setup() {
	a = "hello, world"
	done = true
}

func main() {
	go setup()
	for !done {
	}
	fmt.Println(a)
	// TODO: 模拟
	// go 语言并不保证在 main() 函数中观测到的对 done 的写入操作发生在对字符串 a 的
	// 写操作之后, 因此程序可能打印出一个空字符串; 更有可能, 因为两个 grountine
	// 之间没有同步事件, setup 对 done 的写入操作甚至无法被 main 看到, main() 函
	// 数有可能陷入死循环中.
	// 一致性内存模型的问题比较系统, 需要考虑硬件和软件的实现; 以及 cache 和
	// cpu 的 store buffer 的问题;
	// 可能被 go 语言的编译器优化, 也或者 goroutine 执行过快, 与main主goroutine
	// 的运行在效率上出现了顺序, 所以导致模拟出上述情况

}

/*
 * go 中同一个 grountine 内部, 顺序一致性内存模型是得到保证的, 但是不同的
 * grountine 之间, 并不满足顺序一致性的内存模型, 需要通过明确定义的同步事件
 * 来作为同步的参考; 如果两个事件不可排序, 就说这两个事件是并发的; 为了最大化
 * 并行, go 语言的编译器和处理器在不影响上述规定的前提下可能会对执行语句重新
 * 排序(cpu也会对一些指令进行乱序执行)
 *
 */

// 如果一个并发程序无法确定事件的顺序关系, 那么程序的运行则会有不确定的结果:
func main1() {
	go fmt.Println("hello, world")
}

// TODO: goroutine 调度原理
// 根据 go 语言规范, main() 函数退出时程序结束, 不会等待任何后台 goroutine,
//  goroutine 的执行和 main1 函数的返回事件是并发的, 谁都有可能先发生, 所以
//  main1 的执行结果是未知的.
// 用 ./atomic.go 中描述的原子操作并不能解决问题, 因为无法确定两个原子操作
// 之间的顺序, 只有通过同步原语来给两个事件明确排序:
func main2() {
	done := make(chan int)

	go func() {
		fmt.Println("hello, world")
		done <- 1
	}()
	<-done
}

// 在同一个 goroutine 内满足顺序一致性原则, 现在使用 channel 规定了两个
// goroutine 执行完成的顺序, 此时程序可以正常打印结果.
// 同时, 也可以使用 sync.Mutex 互斥量实现同步:
func main3() {
	var mu sync.Mutex

	mu.Lock()
	go func() {
		fmt.Println("hello, world")
		mu.Unlock()
	}()

	// 后台 goroutine 的 mu.Unlock() 必然在 fmt.Println 完成后发生, main 函数
	// 的第二个 mu.Lock 必然在 mu.Unlock() 之后发生(sync.Mutex保证), 此时main
	// 函数便会在后台 goroutine 执行完成后退出
	mu.Lock()
}
