package main

import (
	"fmt"
	"sync"
)

/*
 * 同步的方式有多种, 会使用 sync.Mutex 来实现同步通信; 根据 sync 包中的
 * Lock() 和 Unlock() 函数的注释可知, 不能直接对一个未加锁状态的 sync.Mutex
 * 进行解锁, 这会导致运行时异常; 对被占用的锁进行加锁会导致阻塞.
 */

/*
// TODO: sync 的源码, 进程信号量
// $GOROOT/src/sync/mutex.go
// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
//
// A locked Mutex is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a Mutex and then
// arrange for another goroutine to unlock it.
func (m *Mutex) Unlock() {
	if race.Enabled {
		_ = m.state
		race.Release(unsafe.Pointer(m))
	}
}

// Lock locks m.
// If the lock is already in use, the calling goroutine
// blocks until the mutex is available.
func (m *Mutex) Lock() {
	// Fast path: grab unlocked mutex.
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	// Slow path (outlined so that the fast path can be inlined)
	m.lockSlow()
}
*/

func main() {
	var mu sync.Mutex

	go func() {
		fmt.Println("hello, world")
		mu.Lock()
	}()

	mu.Unlock()

}

// 因此 mu.Lock() 和 mu.Unlock() 并不在同一个 goroutine 中, 在它们的执行顺序
// 上不满足顺序一致性内存模型, 所以它们是不可排序的(并发), 所以在 mu.Lock()
// 中的 mu.Unlock() 可能先执行导致出现运行时异常
// 输出:
// fatal error: sync: unlock of unlocked mutex

// 修复的方式是在 main() 函数所在 goroutine 中执行两次 mu.Lock(); 当第二次加锁
// 时会因为锁已经被占用(不是递归锁, TODO)而阻塞, 等待另一个 goroutine 执行解锁,
// 解锁会导致 main() 中的 mu.Lock() 阻塞状态取消, 此时两个 goroutine 无其他同步
// 事件参考, 退出的事件是并发的, 虽无法确定两个 goroutine 退出的先后顺序, 但是
// 可以保证打印工作的顺利完成.

func main2() {
	var mu sync.Mutex

	mu.Lock()
	go func() {
		fmt.Println("hello, world")
		mu.Unlock()
	}()

	mu.Lock() // 锁被占用, 阻塞
}
