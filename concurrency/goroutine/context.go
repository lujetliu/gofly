package main

/* 上下文 */

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
 * 在 go1.7 发布时, 标准库增加了一个 context 包, 用来简化对于处理单个请求的多个
 * goroutine 之间与请求域的数据, 超时和退出等操作, 可以用 context 包来重新实现
 * ./select.go 中的 goroutine 安全退出和超时控制.
 *
 * TODO: context 的使用, context 源码, 理解上下文原理
 *
 * 在一些 goroutine 使用无限循环时:
 * - 使用 close(channel) 会导致发送无限数据到 channel 的 goroutine 抛出 panic
 * - 使用 sync.WaitGroup 则会等待 goroutine 无限执行下去
 * 此时需要引入新的 channel, 显示接收停止执行的信号, 通知 goroutine 退出
 */

func worker(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	for {
		select {
		default:
			fmt.Println("hello")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// 当并发体超时或 main  主动停止工作者 goroutine 时, 每个工作者都可以安全退出
func main1() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(ctx, &wg)
	}

	time.Sleep(1 * time.Second)
	cancel()

	wg.Wait()
}

// TODO: 模拟/观察 goroutine 泄露
// go 语言是具有内存自动回收特性, 因此内存一般不会泄露; 在 ./prime.go
// 素数筛例子中,  GenerateNatural() 和 PrimeFilter() 函数内部都启动了新的
// goroutine, 当 main() 函数不再使用通道时, 后台 goroutine 有泄露的风险(无限
// for循环), 可以通过 context 包来避免这个问题, 以下是改进的素数筛实现:

// 返回生成自然数序列的通道 2, 3, 4, ......
func GenerateNatural(ctx context.Context) chan int {
	ch := make(chan int)
	go func() {
		for i := 2; ; i++ {
			// 无限循环是为了使得只需要在 main 中做简单的修改就可以得到
			// 需要数量的素数
			select {
			case <-ctx.Done():
				return
			case ch <- i:
				// 注意: 这里使用 for 生成无穷的自然数序列, 所以不能在 main
				// 中通过 close(ch) 通知该 goroutine 完成工作, 否则会因为往
				// 已经关闭的 channel 发送数据出现 panic; 因为这里的 for 是
				// 无限循环, 因此必须引入其他的 channel 来通知该 goroutine
			}
		}
	}()
	return ch
}

//  通道过滤器: 删除能被素数整除的数
func PrimeFilter(in <-chan int, prime int, ctx context.Context) chan int {
	out := make(chan int)
	go func() {
		for {
			// 对应 line58 同理
			if i := <-in; i%prime != 0 {
				select {
				case <-ctx.Done():
					// 对这里进行拓展, 也可以使用一个通道接收系统中断信号,
					// 由用户使用 ctrl+c 控制输出的素数的数量
					return
				case out <- i:
					// 对已经关闭的 channel 的接收会立即返回其类型的零值,
					// 即使 close(out) 或 close(in) 也不能使该 goroutine 结束
					// 工作, 必须引入其他的 channel 来通知该 goroutine
				}
			}
		}
	}()
	return out
}

func main() {
	//  通过 Context 控制后台 goroutine 的状态
	ctx, cancel := context.WithCancel(context.Background())

	ch := GenerateNatural(ctx)
	for i := 0; i < 10; i++ {
		prime := <-ch // 新出现的素数
		fmt.Printf("%v: %v\n", i+1, prime)
		ch = PrimeFilter(ch, prime, ctx) // 基于新素数构造的过滤器
	}

	cancel()
}

// 当 main() 函数完成前, 通过调用 cancel() 通知后台 goroutine 退出, 这样就可
// 避免 goroutine 的泄露
