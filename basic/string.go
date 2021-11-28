package main

import (
	"fmt"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

/*
 * 一个字符串是一个不可改变的字节序列, 和数组不同, 字符串的元素不可修改,
 * 是一个只读的字节数组, 每个字符串的长度虽然也是固定的, 但是字符串的长度
 * 并不是字符串类型的一部分.
 *
 */

// rune 类型占用四个字节, 称为字符类型, 它表示的是一个 Unicode字符
// (Unicode是一个可以表示世界范围内的绝大部分字符的编码规范), 字符类型用
// 单引号, 字符串类型用双引号;
// rune 只是 int32 类型的别名, 并不是重新定义的类型, rune 用于表示每个
// Unicode 码点

// 字符串的底层结构在 reflect.StringHeader 中定义:
type StringHeader struct {
	Data uintptr // 字符串指向的底层字节数组
	Len  int     // 字符串的字节长度
}

/*
 *
 * 字符串其实是一个结构体, 因此字符串的赋值操作就是 reflect.StringHeader 结构
 * 体的复制过程, 并不会涉及底层字节数组的复制
 *
 */

func main() {
	s := "hello, world"
	// hello := s[:5]
	// world := s[7:]
	s1 := "hello, world"[:5]
	s2 := "hello, world"[7:]
	/*
	 * TODO: 理解
	 * 字符串虽然不是切片, 但是支持切片操作, 不同位置的切片底层访问的是同一块内存
	 * 数据(因为字符串是只读的, 所有相同的字符串面值常量通常对应同一个字符串常量)
	 */
	// 字符串的数组类似, 内置的 len() 函数返回字符串的长度, 也可以通过
	// reflect.StringHeader 结构访问字符串的长度(不推荐使用)

	fmt.Println("len(s):", (*reflect.StringHeader)(unsafe.Pointer(&s)).Len)
	// fmt.Println(len(s1))
	fmt.Println("len(s1):", (*reflect.StringHeader)(unsafe.Pointer(&s1)).Len)
	fmt.Println("len(s2):", (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len)
}

func forOnString(s string, forBody func(i int, r rune)) {
	for i := 0; len(s) > 0; {
		r, size := utf8.DecodeRuneInString(s)
		forBody(i, r)
		s = s[size:]
		i += size
	}
}
