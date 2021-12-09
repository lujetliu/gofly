package main

import (
	"sync"
	"sync/atomic"
	"time"
)

/*
 * 原子操作配合互斥锁可以实现高效的单件模式, 互斥锁的代价比普通整数的原子读写
 * 高很多, 在性能敏感的地方可以增加一个数字型的标志位, 通过原子检测标志位状态
 * 降低互斥锁的使用次数来提高性能.
 *
 *
 *  TODO:
 * 单件模式(又称单例模式)确保一个类只有一个实例，并提供一个全局访问点;
 * 一种用于确定整个应用程序中只有一个类实例且这个实例所占资源在整个应用
 * 程序中是共享时的程序设计方法.
 * 使用场景:
 * - 当类只有一个实例, 且客户可以从一个众所周知的访问点访问它时.
 * - 当这个唯一实例应该是通过子类化可扩展的, 并且客户应该无需更改代码就能
 * 使用一个扩展的实例时.
 */

type singleton struct{}

// var (
// 	instance    *singleton
// 	initialized uint32
// 	mu          sync.Mutex
// )

// func Instance() *singleton {
// 	if atomic.LoadUint32(&initialized) == 1 {
// 		return instance
// 	}

//  // 对读写 instance 操作加锁
// 	mu.Lock()
// 	defer mu.Unlock()

// 	if instance == nil {
// 		defer atomic.StoreUint32(&initialized, 1)
// 		instance = &singleton{}
// 	}
// 	return instance
// }

// 将通用的代码提取出来就成了标准库中 sync.Once 的实现:
type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}

	// 对读和写 o.done 这一前后过程加锁, 保证这个过程是一个整体(原子性)绑定执
	// 行, 防止在读写的中间其他 goroutine 对 o.done 的写操作.
	o.m.Lock()
	defer o.m.Unlock()

	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

// 基于 sync.Once 重新实现单件(singleton) 模式:
var (
	instance *singleton
	once     sync.Once
)

func Instance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})

	return instance
}

/*
 * sync/atomic 包对基本数值类型以及复杂对象的读写都提供了原子操作的支持,
 * atomic.Value 原子对象提供了 Load() 和 Store() 两个原子方法, 分别用于
 * 加载和保存数据, 返回值和参数都是 interface{} 类型, 因此可以用于任意的
 * 自定义复杂类型.
 */

var config atomic.Value

// 这是一个简化的生产者消费者模型, 后台 goroutine 生成最新的配置信息, 前台多个
// goroutine 获取最新的配置信息, 所有 goroutine 共享配置信息资源.
func main() {
	// 初始化配置信息
	config.Store(loadConfig())

	// 启动一个后台 goroutine, 加载更新后的配置信息
	go func() {
		for {
			time.Sleep(time.Second) // 一秒更新以此配置
			config.Store(loadConfig())
		}
	}()
	for i := 0; i < 10; i++ {
		go func() {
			// for r := range requests() {
			// 	c := config.Load()
			// ...
			// }
		}()
	}

}

func loadConfig() struct{} {
	// ...
	return struct{}{}
}

func requests() chan<- struct{} {
	// ...
	return nil
}
