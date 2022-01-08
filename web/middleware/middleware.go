package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

/* 使用中间件剥离非业务逻辑 */
// 对大多数的场景来说,非业务的需求都是在 HTTP  请求处理前做一些事情, 并且在响应
// 之后做一些事情

var logger = log.New(os.Stdout, "", 0)

// 对 helloHandler() 函数增加耗时统计, 可以使用一种叫函数适配器(function
// adapter) 的方法来对 helloHandler() 进行包装
func helloHandler(wr http.ResponseWriter, r *http.Request) {
	wr.Write([]byte("hello"))
}

func timeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(wr http.ResponseWriter, r *http.Request) {
			timeStart := time.Now()

			// next handler
			next.ServeHTTP(wr, r)

			timeElapsed := time.Since(timeStart)
			logger.Println(timeElapsed)
		},
	)
}

func main() {
	http.Handle("/", timeMiddleware(http.HandlerFunc(helloHandler)))
	http.ListenAndServe(":8080", nil)
	// ...
}

/* http 库的 Handle, HandlerFunc, 和 ServeHTTP 的关系

type handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

type (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}

即只要自己定义的 helloHandler() 函数的函数签名是:
func (ResponseWriter, *Request)
就和 http.HandlerFunc() 有了一致的签名, 可以将 helloHandler() 函数进行类型转换,
转换为 http.HandlerFunc(), 从而间接拥有 ServeHTTP() 方法实现 http.handler 接口
可见一个请求的调用链:
h = getHandler() => h.ServeHTTP(w, r) => h(w, r)
所以就可以使用中间件通过包装 handler 再返回一个新的 handler
*/

/*
 * 中间件要做的就是通过一个或多个函数对 handler 进行包装, 返回一个包括了各个
 * 中间件逻辑的函数链, 如:
 * customizeHandler = logger(timeout(ratelimit(helloHandler)))  TODO: 实现
 * 此流程在进行请求处理的时候不断的进行函数压栈再出栈, 类似与递归的执行流:

	 [exec of logger logic]				函数栈:[]
	 [exec of timeout logic]			函数栈:[logger]
	 [exec of ratelimit logic]			函数栈:[timeout/logger]
	 [exec of helloHandler logic]		函数栈:[ratelimit/timeout/logger]
	 [exec of ratelimit logic part2]    函数栈:[timeout/logger]
	 [exec of timeout logic part2]      函数栈:[logger]
	 [exec of logger logic part2]       函数栈:[]

*/
