package main

import (
	"image/color"
	"sync"
)

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

func (f *File) ReadFile(offset int64, data []byte) int {
	// ...
	return 0
}

/*
 * 此时, 函数 CloseFile() 和 ReadFile() 就成了 File 类型独有的方法了(而不
 * 是File对象的方法), 它们也不再占用包级空间中的名字资源, 同时 File 类型
 * 已经明确了它们的操作对象, 因此方法名可简化为 Close() 和 Read().
 */
func (f *File) Close() error {
	// ...
	return nil
}

func (f *File) Read(offset int64, data []byte) int {
	// ...
	return 0
}

/*
 * 将第一个函数参数移动到函数前面, 从代码角度看很小的改动, 却从编程哲学角度上
 * 来说, go 语言已经进入面向对象语言的行列. 可以给任何自定义类型添加一个或
 * 多个方法, 每种类型对应的方法和类型的定义在同一个包中, 因此无法给 int 这种
 * 内置类型添加方法(方法的定义和类型的定义不在同一个包中); 对于给定的类型,
 * 每个方法的名字必须是唯一的, 同时方法和函数一样也不支持重载.
 *
 * 重载: 函数或者方法有相同的名称, 但是参数列表不相同的情形,
 * 这样的同名不同参数的函数或者方法之间，互相称之为重载函数或者重载方法
 */

// go 语言不支持传统面向对象中的继承特性, 而是以特有的组合方式支持了方法的
// 继承, 通过在结构体内置匿名的成员来实现继承:
type Point struct {
	X, Y float64
}

// 通过嵌入 Point 来提供 X 和 Y 两个字段
type ColoredPoint struct {
	Point // 匿名成员
	Color color.RGBA
}

func CallTest() {
	var cp ColoredPoint
	cp.X = 1 // cp.Point.X
}

// 通过嵌入匿名的成员, 不仅可以继承匿名成员的内部成员, 还可以继承匿名成员
// 所对应的方法; 一般称 Point 为基类, 把 ColoredPoint 看作 Point 的继承类
// 或子类, 不过这种方式继承的方法并不能实现 c++ 中虚函数的多态特性(?), 所
// 有继承来的方法的接收者参数依然是那个匿名成员本身, 而不是当前的变量.

type Cache struct {
	m          map[string]string
	sync.Mutex // 匿名成员
}

func (c *Cache) Lockup(key string) string {
	c.Lock()         // c.Mutex.Lock()
	defer c.Unlock() // c.Mutex.Unlock()

	return c.m[key]
}

// 在调用 c.Lock() 和 c.Unlock() 时, c 并不是方法 Lock() 和 Unlock() 真正的
// 接收者, 而是会将其展开为 c.Mutex.Lock() 和 c.Mutex.Unlock() 调用, 这种展开
// 是在编译阶段完成的, 没有运行时代价.

/*
 * TODO: 语言设计思想的理解, 继承, 多态, 接口
 * 在传统的面向对象语言(如c++或java)的继承中, 子类的方法是在运行时动态绑定到
 * 对象的, 因此基类实现的某些方法看到的 this 可能不是基类类型对应的对象, 这个
 * 特性会导致基类方法运行的不确定性; 在 go 语言中, 通过嵌入匿名的成员来"继承"
 * 的基类方法, this 就是实现该方法的类型的对象, go 语言中方法是编译时静态绑定
 * 的; 如果需要虚函数的多态特性, 则需要借助接口来实现.
 */
