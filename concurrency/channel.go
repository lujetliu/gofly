package main

import "fmt"

/*
 * 通道(channel)是在 goroutine 之间进行同步的主要方法, 在无缓冲的通道上的每一
 * 次发送都有与其对应的接收操作相匹配, 发送和接收操作通常发生在不同的 goroutine
 * 上(在同一个 goroutine 上执行两个操作很容易导致死锁);
 * TODO: channel 源码, 内部原理
 */

/*
 * 如果试图通过在主 main 中加入 time.Sleep() 语句使主 goroutine 休眠, 从而等待
 * 子 goroutine 的执行完成; 但是这样构建的程序是不稳健的, 依然有失败的可能性.
 *
 * 严谨的并发程序的正确性不应该依赖于 cpu 的执行速度和休眠时间等不可靠的因素.
 * 严谨的并发也应该是可以静态推导出结果的; 根据 goroutine 内顺序一致性, 结合
 * 通道和 sync 事件的可排序性来推导, 完成各个 goroutine 各段代码的偏序关系
 * 排序(TODO); 如果两个事件无法根据此规则来排序, 则它们是并发的, 也就是执行
 * 先后顺序不可靠的.
 *
 * 解决同步问题的思路: 使用显式的同步
 */

// 在无缓冲的通道上的发送操作总在对应的接收操作完成前发生
var done = make(chan bool)
var msg string

func aGoroutine() {
	msg = "hello, world"
	done <- true // close(done)
	// 若在关闭通道后继续从中接收数据, 接收者会收到该通道返回的零值;
	// 所以在此例中可使用 close(done) 关闭通道代替 done <- true, 作用相同
}

func main1() {
	go aGoroutine()
	<-done
	fmt.Println(msg)
	// 这里因为对通道的发送操作是在子 goroutine 中进行的, 所以使用有缓冲
	// 通道也是等效的
}

/* 基于以上规则, 交换两个 goroutine 中的接收和发送操作也是可以的(很危险)   */
// 对于从无缓冲通道进行的接收, 发生在对该通道进行的发送完成之前
func bGoroutine() {
	msg = "hello, world"
	<-done
}

func main2() {
	go bGoroutine()
	done <- true
	fmt.Println(msg)
}

// 虽然也可以保证打印出 "hello, world", 因为在 main 中 done <- true 发送完成前,
// 后台 goroutine 的接收已经开始, 这保证了 msg 的赋值操作被执行; 也就是说
// 对无缓冲通道, 接收方和发送方都准备好后才会接收和发送; 但若该通道为带缓冲
// 的 done = make(bool, 1), 则 main 中的 done <- true 发送操作将不会被后台
// goroutine 的 <-done 接收操作阻塞, 该程序将无法打印出期望的结果.

/*
 * 对于带缓冲(缓冲大小为C)的通道, 对于通道中的第K个接收操作发生在第K+C个发送
 * 操作完成之前, 如果将C设置为0自然就对应无缓冲的通道, 也就是第K个接收完成在
 * 第K个发送完成之前, 因为无缓冲的通道只能同步发1个, 所以就简化为前面无缓冲
 * 通道的规则: 对于从无缓冲通道进行的接收, 发生在对该通道进行的发送完成之前.
 *
 */

// 基于带缓存通道, 可以将打印操作扩展到 N 个:
func main3() {
	done := make(chan int, 10) // 带10个缓冲区

	// 开 N 个后台 goroutine 进行打印工作
	for i := 0; i < cap(done); i++ {
		// len(channel) 内部的数据长度
		// cap(channel) channel 的容量
		go func() {
			fmt.Println("hello, world")
			done <- 1
		}()
	}

	// 等待 N 个后台 goroutine 执行完成
	for i := 0; i < cap(done); i++ { // 可以用 range
		<-done
	}
}

/*
 * TODO: sync 源码, WaitGroup 的实现原理, 进程调度, 信号量
 * 对于这种要等待 N 个线程完成后再进行下一步的同步操作, sync 提供了
 * sync.WaitGroup 来等待一组 goroutine 完成
 *
 */

// 可以通过控制通道的缓冲区的大小控制并发执行的 goroutine 的最大数目
var limit = make(chan int, 3)
var work []func()

func main4() {
	for _, w := range work {
		go func(w func()) {
			limit <- 1
			w()
			<-limit
			// 在 goroutine 中对 limit 的发送和接收操作限制了同时执行的
			// goroutine 的数量
		}(w)
	}

	// <-make(chan int)
	// select {}
	for {
	}
	// select {} 是一个空的通道选择语句, 会导致死锁, <- make(chan int) 也一样;
	// for{} 会出现死循环

}
