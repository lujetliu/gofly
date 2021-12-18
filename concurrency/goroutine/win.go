package main

/* 赢者为王并发思想 */

import "fmt"

/*
 * 采用并发编程的动机很多: 并发编程可以简化问题, 例如一类问题对应一个处理线程
 * 会更简单; 并发编程可以可以提升性能, 其实对性能而言, 并不是程序运行速度快
 * 就表示用户体验好, 很多时候程序能快速响应用户请求才最重要; 可以采用类似的
 * 策略构建程序.
 *
 * 通过适当的开启一些冗余的线程, 尝试用不用的途径去解决同样的问题, 最终以赢者
 * 为王的方式提升程序的性能
 *
 */

// 当在浏览器里同时打开必应, 谷歌和百度搜索 "golang" 时, 当某个搜索最先返回
// 结果, 就可以关闭其他页面了
func main() {
	ch := make(chan string, 32)

	go func() {
		ch <- searchByBing("golang")
	}()
	go func() {
		ch <- searchByGoogle("golang")
	}()
	go func() {
		ch <- searchByBaidu("golang")
	}()

	// 此处, 当任意一个搜索返回结果时, 程序都会最先返回该结果(因为带缓存通道,
	// 所以不会阻塞)
	fmt.Println(<-ch)
}

func searchByBing(str string) string {
	return ""
}
func searchByGoogle(str string) string {
	return ""
}
func searchByBaidu(str string) string {
	return ""
}
