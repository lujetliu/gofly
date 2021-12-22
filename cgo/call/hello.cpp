// hello.cpp


#include <iostream>

extern "C" {
	#include "hello.h"
}

void SayHello(const char* s) {
	std::cout << s;
}

// TODO: 如何调用
// 在 c++ 版本的 SayHello() 函数中, 通过 c++ 特有的 std::cout 输出流输出字符串,
// 同时为了保证 c++ 语言实现的 SayHello() 函数满足 c 语言头文件 hello.h 定义的
// 函数规范, 需要通过 extern "C" 语句指示该函数的链接符号遵循 C 语言的规则.
//


