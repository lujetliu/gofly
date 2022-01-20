package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

/*
 TODO: 进程调度, goroutine 调度
 在操作系统中, 每个进程都会有一个唯一的进程编号, 每个线程也有唯一的线程编号,
 在 go 中, 每个 goroutine 也有自己唯一的编号, 虽然 goroutine 有内在的编号,
 但是 go 却刻意没有提供获取该编号的接口;
 go 没有提供 goid 的原因是为了避免滥用, 因为大部分用户在轻松拿到 goid 后,
 就会在编程中不自觉的编写出强依赖 goid 的代码, 强依赖 goid 将导致这些代码
 不好移植, 同时也会导致并发模型复杂化, 同时,  go 语言中可能同时存在海量的
 goroutine, 但是每个 goroutine 何时被销毁并不好实时监控, 这会导致依赖 goid
 的资源无法自动回收(需要手工回收); 如果使用 go 汇编则可以忽略这些问题.
*/

// 纯 go 方式获取 goid
func GetGoid() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		// runtime.Stack() 获取当前函数的调用栈帧信息, TODO: 源码
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine ")
	)

	// TODO: strings 库的用法, 源码
	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}

	return int64(id)
}

func main() {
	fmt.Println(GetGoid())
}
