package main

/* 生产者/消费者模型 */

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

/*
 * 生产者和消费者问题是线程模型中的经典问题:生产者和消费者在同一时间段内共用
 * 同一个存储空间, 生产者往存储空间中添加产品, 消费者从存储空间中取走产品,
 * 当存储空间为空时, 消费者阻塞, 当存储空间满时, 生产者阻塞.
 *
 *
 * 并发编程中最常见的就是生产者/消费者模型, 该模型主要通过平衡生产线程和消费
 * 线程的工作能力来提高程序的整体处理数据的速度;
 * 生产者生产一些数据, 放到成果队列中, 同时消费者从成果队列中获取这些数据;
 * 这样生产和消费就变成了异步的两个过程, 等成果队列中没有数据时, 消费者就
 * 进入饥饿的等待中; 而当成果队列中数据已满时, 生产者则面临因产品积压导致
 * cpu 被剥夺的问题.
 * TODO: 进程调度, 实验剥夺空闲进程cpu资源
 */

// 生产者, 生成 factor 整倍的序列
func Producer(factor int, out chan<- int) {
	for i := 0; ; i++ {
		out <- i * factor
	}
}

// 消费者
func Consumer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

func main() {
	ch := make(chan int, 64) // 成果队列

	go Producer(3, ch) // 生成3的倍数的序列
	go Producer(5, ch) // 生成5的倍数的序列
	go Consumer(ch)    // 消费生成的队列

	// 使用 ctrl+c 退出
	sig := make(chan os.Signal, 1)
	// 接收系统中断信号, TODO: 常见的 signal, 系统信号
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <-sig)
}

// 此例中有两个生产者, 并且在两个生产者之间无同步事件参考, 它们是并发的
