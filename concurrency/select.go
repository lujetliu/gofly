package main

import (
	"fmt"
	"sync"
	"time"
)

/* 并发的安全退出 */

/*
 * 当 goroutine 工作在错误的方向时, go 语言并没有提供直接终止 goroutine 的方法,
 * 因为这样会导致 goroutine 之间的共享变量处在未定义的状态;
 * 不同 goroutine 之间主要依靠通道进行通信和同步; 要同时处理多个通道的发送和接
 * 收操作, 需要使用 select 关键字, 当 select{} 有多个分支时, 会随机选择一个可
 * 用的通道分支, 如果没有可用的通道分支, 则选择 default 分支, 否则一直保持阻塞
 * 状态.
 *
 * 1. 可以使用 select{} 实现通道的超时判断
 */

// 当有多个通道均可操作时, select 会随机选择一个通道, 基于该特性实现一个生成
// 随机数
func main1() {
	ch := make(chan int)
	go func() {
		for {
			select {
			case ch <- 0:
			case ch <- 1:
			}
		}
	}()

	// for v := range ch {
	// 	fmt.Println(v)
	// }
	fmt.Println(<-ch)
}

// 通道的发送和接收操作是一一对应的, 如果要停止多个 goroutine, 可能需要创建
// 同样数量的通道; 但也可以通过 close() 关闭一个通道来实现广播的效果, 从关闭的
// 通道(假设此时通道为空, 即len(channel)=0)接收的操作均会收到一个零值和一个
// 可选的失败标志.
func worker1(cancel chan bool) {
	for {
		select {
		case <-cancel:
			return
		default:
			fmt.Println("hello")
			// 正常工作
		}

	}
}

func main2() {
	cancel := make(chan bool)

	for i := 0; i < 10; i++ {
		go worker1(cancel)
	}

	time.Sleep(time.Second)
	close(cancel)
}

// 上例通过 close() 关闭通道向多个 goroutine 广播退出的指令, 不过此程序仍然不够
// 健壮; 当每个 goroutine 收到退出指令退出时一般会进行一定的清理工作, 但是退出
// 的清理工作并不能保证被完成, main goroutine 并没有等待各个工作 goroutine
// 退出工作完成的机制, 现在使用 sync.WaitGroup 管理 goroutine.
func worker2(wg *sync.WaitGroup, cancel chan bool) {
	defer wg.Done()

	for {
		select {
		case <-cancel:
			return
		default:
			fmt.Println("hello")
			// 正常工作
		}
	}
}

func main() {
	cancel := make(chan bool)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker2(&wg, cancel)
	}
	time.Sleep(time.Second)
	close(cancel)
	wg.Wait()
}

// 现在每个工作者的 goroutine 的创建, 运行, 暂停和退出都是在 main() 函数
// 的安全控制之下了.
