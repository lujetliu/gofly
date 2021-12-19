package main

/* 错误处理  */

import (
	"fmt"
	"log"
	"syscall"
)

/*
 * TODO: syscall
 * 导致出现错误的原因不止一种, 很多时候用户希望了解更多的错误信息; 在 c 语言中
 * 默认采用一个整数类型的 errno 表达错误, 这样就可以根据需要定义多种错误类型;
 * 在 go  语言中, syscall.Errno 就是对应 c 语言中 errno 类型的错误, 在 syscall
 * 包中的接口, 如果有返回错误的话, 底层也是 syscall.Errno 错误类型
 */

// 通过 syscall 包的接口来修改文件的模式时, 如果遇到错误可以通过将 err 强制
// 断言为 syscall.Errno 错误类型来处理; 还可以进一步的通过类型查询或类型断言
// 获取底层真实的错误类型, 这样就可以获取更详细的错误信息
func main1() {
	err := syscall.Chmod(":invalid path:", 0666)
	if err != nil {
		log.Fatal(err.(syscall.Errno))
	}
}

// 在 go 语言中错误被认为是一种可以预期的结果, 而异常则是一种非预期的结果, 发
// 生异常可能表示程序中存在bug或发生了其他不可控的问题; go 语言推荐使用
// recover() 函数将内部异常转为错误处理, 这使得用户可以真正的关心业务相关
// 的错误处理;
// 以 json 解析器为例说明 recover() 的使用场景, 考虑到解析器的复杂性, 即使
// 某个语言解析器目前工作正常, 也无法肯定它没有漏洞; 因此, 当某个异常出现
// 时, 不会选择让解析器崩溃, 而是会将 panic 异常当做普通的解析错误, 并附加
// 额外信息提醒用户报告此错误.
// func ParseJson(input string) (s *Syntax, err error) {
// 	defer func() {
// 		if p := recover(); p != nil {
// 			err = fmt.Errorf("JSON: internal error: %v", p)
// 		}
// 	}()

// 	// ...
// }

// go 语言库的实现习惯是: 即使在包内部使用了 panic, 在导出函数
// 时也会被转化为明确的错误值.

/*
 * 有时为了方便上层用户理解, 底层实现者会将底层的错误信息重新包装为新的错误
 * 类型返回给上层用户; 上层用户在遇到错误时, 很容易从业务层面理解错误发生的
 * 原因, 但是在上层用户获得新的错误的同时, 也丢失了底层最原始的错误类型(只
 * 剩下错误描述信息了);
 * 为了记录这种错误类型在包装的变迁过程中的信息, 一般会定义一个辅助的
 * WrapError() 函数(wrap, 包装), 用于包装原始的错误, 同时保留完全的原始错误信息,
 * 为了方便问题的定位并且能记录错误发生时的函数调用状态, 需要在出现严重
 * 错误的时候保存完整的函数调用信息, 同时为了支持RPC等跨网络的传输, 还需要
 * 将错误序列化为类似 JSON 格式的数据, 然后在从这些数据中将错误解码恢复出来.
 */

// 可以定义以下所需类型
type Error interface {
	Caller() []CallerInfo
	Wraped() []error
	Code() int
	error

	private()
}

type CallerInfo struct {
	FuncName string
	FileName string
	FileLine int
}

// TODO: 设计思路, 理解(go语言高级编程, Page69)
// Error 是 error 类型的扩展, 用于给错误增加调用栈信息, 同时支持错误的多级
// 嵌套包装, 支持错误码格式, 为了方便使用可以定义以下辅助函数:
func New(msg string) error
func NewWithCode(code int, msg string) error

func Wrap(err error, msg string) error
func WrapWithCode(code int, err error, msg string) error

func FromJson(json string) (Error, error)
func ToJson(err error) string

// 示例使用
func f1() {
	err := NewWithCode(404, "http error code")
	fmt.Println(err.(Error).Code())
}

// 在 go 语言中, 错误处理有一套独特的编码风格, 检查某个子函数是否失败后, 通常
// 将处理失败的逻辑代码放在处理成功的代码之前, 如果某个错误会导致函数返回, 那么
// 成功时的逻辑代码不应该放在 else 语句块中, 而应该放在函数体中
// func f2() {
// 	f, err := os.Open("filename.ext")
// 	if err != nil {
// 		// 失败, 马上返回
// 	}
// 	// 正常的处理流程
// }
// go 语言中大部分函数的代码结构几乎相同, 首先是一系列的初始检查, 用于防止错误
// 发生, 之后是函数的实际逻辑.
