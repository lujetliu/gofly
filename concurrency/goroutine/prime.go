package main

import "fmt"

/* 素数筛 */
// 素数又叫质数（prime number), 有无限个; 质数定义为在大于1的自然数中,
// 除了1和它本身以外不再有其他因数

func GenerateNatural() <-chan int {
	ch := make(chan int)
	go func() {
		for i := 2; ; i++ {
			ch <- i
		}
	}()
	return ch
}

func PrimeFilter(in <-chan int, prime int) chan int {
	out := make(chan int)
	go func() {
		for i := range in {
			if i%prime != 0 {
				out <- i
			}
		}
	}()
	return out
}

func main() {
	ch := GenerateNatural()
	for i := 0; i < 10; i++ {
		prime := <-ch // 新出现的素数
		// 每次循环迭代时, 通道中的第一个数必是素数, 先读取并打印; 然后基于通道
		// 中剩余的数列, 并以当前取出的素数为筛子过滤后面的素数, 不同的素数筛
		// 对应的通道是串联在一起的
		// TODO: 对通道中的数据流转加深理解
		// TODO: 对通道的数据流转以动态图的方式实时展示
		fmt.Printf("%v: %v\n", i+1, prime)
		ch = PrimeFilter(ch, prime) // 基于新素数构造的过滤器
	}
}

/*
 * TODO: 加深理解和验证
 * 素数筛展示了一种优雅的并发程序结构, 但是因为每个并发体处理的任务粒度太细微,
 * 程序整体的性能并不理想; 对于细粒度的并发程序, csp 模型中固有的消息传递的代价
 * 过高(多线程并发模型同样要面临线程启动的代价)
 *
 */
