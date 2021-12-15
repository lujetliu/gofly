package main

/*
 *
 * 原子操作是并发编程中"最小的且不可并行化"的操作, 通常如果多个并发体对同一个
 * 共享资源进行的操作是原子操作, 那么同一时刻最多只能有一个并发体对该资源进行
 * 操作; 从线程角度看, 在当前线程修改共享资源期间, 其他线程是不能访问该资源的,
 * 原子操作对多线程并发编程模型来说, 不会发生有别于单线程的意外情况, 共享资源
 * 的完整性可以得到保证.
 *
 * 标准库的 sync/atomic 包对原子操作提供了丰富的支持;
 */

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var total1 uint64

func worker1(wg *sync.WaitGroup) {
	defer wg.Done()

	var i uint64
	for i = 0; i <= 1000; i++ {
		// for 循环中使用外部变量, 避免短变量声(:=) 捕获 i 的类型(int)
		atomic.AddUint64(&total1, i)
		// AddUint64 函数保证了 total 的读取, 更新和保存是一个原子操作,
		// 因此在多线程中访问也是安全的
	}
}

func main1() {
	var wg sync.WaitGroup
	wg.Add(2)
	go worker1(&wg)
	go worker1(&wg)
	wg.Wait()

	fmt.Println(total1)
}

/*
 * 一般, 原子操作都是通过"互斥"访问来保证的, 通常是由特殊的 CPU 指令提供保护,
 * 如果仅仅想模拟粗粒度的原子操作, 可借助 sync.Mutex 实现, 但使用 sync.Mutex
 * 繁琐且效率底下, 推荐使用原子操作或通道.
 */

var total2 struct {
	sync.Mutex
	value int
}

func worker2(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i <= 1000; i++ {
		// time.Sleep(1000 * time.Microsecond)
		// 如果不加锁模拟数据竞态, 因为 goroutines 执行速度过快, 所有需要加
		// 延迟模拟数据竞态执行结果
		total2.Lock()     //  临界区:
		total2.value += i // 通过加锁和解锁保证语句在同一时刻只被一个线程访问
		total2.Unlock()   //
	}
	// 1001000
	/*
	 * 对多线程模型的程序, 进出临界区前后进行加锁和解锁都是必需的,如果没有锁的
	 * 保护, total 的最终值将由于多线程之间的竞争而可能不确定
	 */
}

func main2() {
	var wg sync.WaitGroup
	wg.Add(2)
	go worker2(&wg)
	go worker2(&wg)
	wg.Wait()

	fmt.Println(total2.value)
}
