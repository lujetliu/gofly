package main

/*
 * 用互斥锁保保护一个数值型的共享资源繁琐且效率低下; 标准库的 sync/atomic 包
 * 对原子操作提供了丰富的支持, 重新实现 ./mutex.go
 *
 */

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var total uint64

func worker(wg *sync.WaitGroup) {
	defer wg.Done()

	var i uint64
	for i = 0; i <= 1000; i++ {
		// for 循环中使用外部变量, 避免短变量声(:=) 捕获 i 的类型(int)
		atomic.AddUint64(&total, i)
		// AddUint64 函数保证了 total 的读取, 更新和保存是一个原子操作,
		// 因此在多线程中访问也是安全的
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go worker(&wg)
	go worker(&wg)
	wg.Wait()

	fmt.Println(total)
}
