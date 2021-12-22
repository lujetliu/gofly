// hello1.c

#include "hello.h"
#include <stdio.h>

void SayHello(const char* s) {
	puts(s);
}

// 本源文件通过 #include "hello.h" 语句包含了 SayHello() 函数的声明, 这样可以
// 保证函数的实现满足模块对外公开的接口;
// 接口文件 ./hello.h 是 hello 模块的实现者和使用者共同的约定, 但是该约定并没
// 有要求必须使用 c 语言来实现 SayHello() 函数, 也可以使用 c++ 语言重新实现
// SayHello() 函数(./hello.cpp)
