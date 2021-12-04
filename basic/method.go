package main

/*
 * 方法是面向对象编程(Object-Oriented Programming, OOP)的一个特性, 在 c++ 中
 * 方法对应一个类对象的成员函数, 关联到某个具体对象的虚表; 而 go 的方法却是
 * 关联到类型的, 这样可以在编译阶段完成方法的静态绑定, 一个面向对象的程序会
 * 用方法来表达其属性对应的操作, 这样就不需要直接去操作对象, 而是借助方法来
 * 进行.
 *
 */

/*
 * 面向对象编程更多的只是一种思想, go 语言的祖先 c 语言虽然不支持面向对象, 但
 * c 语言的标准库中的 File 相关的函数也用到了面向对象编程的思想.
 */
// 实现一组 c 语言风格的 File 函数

// 文件对象
type File struct {
	fd int
}

// 打开文件
func OpenFile(name string) (f *File, err error) {
	// ...
	return nil, nil
}

// 关闭文件
func CloseFile(f *File) error {
	// ...
	return nil
}

// 读文件数据
func ReadFile(f *File, offset int64, data []byte) int {
	// ...
	return 0
}

// OpenFile 类似于构造函数, 用于打开文件对象; CloseFile 类似于析构函数(TODO),
// 用于关闭文件对象, ReadFile 则类似于普通的成员函数, 这3个函数都是普通函数.
// CloseFile 和 ReadFile 作为普通函数, 需要占用包级空间中的名字资源(?), 它们
// 只是针对 File 类型对象的操作; 如果希望这类函数和操作对象的类型紧密绑定在一
// 起时, go 语言的做法是将函数 CloseFile 和 ReadFile 的第一个参数移动到函数名
// 的开头:
func (f *File) CloseFile() error {
	// ...
	return nil
}

func (c *File) ReadFile(offset int64, data []byte) int {
	// ...
	return 0
}

/*
 * 此时, 函数 CloseFile 和 ReadFile 就成了 File 类型独有的方法了(而不是File对象
 * 的方法), 它们也不再占用包级空间中的名字资源, 同时 File 类型已经明确了它们的
 * 操作对象, 因此方法名可简化为 Close 和 Read.
 *
 * 将第一个函数参数移动到函数前面, 从代码角度看很小的改动, 却从编程哲学角度上
 * 来说, go 语言已经进入面向对象语言的行列.
 *
 */
