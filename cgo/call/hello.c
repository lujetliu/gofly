// hello.c
#include <stdio.h>

// 为了运行外部引用, 函数前不需要 static 修饰符 
void SayHello(const char* s) {
	puts(s);
}

