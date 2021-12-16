package main

/*
 * 因为 go 强大的并发性能, 我们倾向于编写最大并发的程序, 以提供最高的性能;
 * 有时候需要适当的控制并发的程度, 这样不仅可以给其他任务让出一定的硬件资源,
 * 也可以适当降低功耗.
 *
 * TODO: 虚拟文件系统, gatefs 源码
 * go 语言自带的 godoc 程序中有一个 vfs 包, 其实现了一个虚拟的文件系统, 在 vfs
 * 包下的子包 gatefs 就是为了控制访问该虚拟文件系统的最大并发数.
 */

// import (
// 	"golang.org/x/tools/godoc/vfs"
// 	"golang.org/x/tools/godoc/vfs/gatefs"
// )

// func main() {
// 	fs := gatefs.New(vfs.OS("/path"), make(chan bool, 8))
//  // ...
// }

// 其中 vfs.OS("/path") 基于本地文件系统构造一个虚拟的文件系统, 然后 gatefs.New
// 基于现有的虚拟文件系统构造一个并发受控的虚拟文件系统, 其原理就是通过带缓存的
// 通道的发送和接收规则来实现最大并发阻塞, 如 ../channel.go 中 98 行的函数示例.

/*
 * 不过 gatefs 对此做了一抽象类型 gate, 增加了 enter() 和 leave() 方法分别对应
 * 并发代码的进入和离开; 当超出并发数目限制的时候, enter() 方法会阻塞直到并发数
 * 降下来为止.
 */

// type gate chan bool

// func (g gate) enter() { g <- true }
// func (g gate) leave() { <-g }

// gatefs 包装的新的虚拟文件系统就是将需要控制并发的方法增加了对 enter() 和
// leave() 的调用而已

// type gatefs struct {
// 	fs vfs.FileSystem
// 	gate
// }

// func (fs gatefs) Lstat(p string) (os.FileInfo, error) {
// 	fs.enter()
// 	defer fs.leave()

// 	return fs.fs.Lstat(p)
// }

// 因此, 不仅可以控制最大的并发数目, 而且可以通过带缓存通道的使用量和最大容量
// 比例来判断程序运行的并发率; 当通道为空时可以认为是空闲状态, 当通道满时可以
// 认为是繁忙状态, 这对于后台一些低级任务的运行是有参考价值的.
