package main

/*
	哈希表: 底层元素为链表的一维数组
	哈希碰撞(Hash Collision): 不同的键通过哈希函数可能产生相同的哈希值, 哈希
	碰撞导致同一个桶中可能存在多个元素, 有多种方式解决哈希冲突(理想状态下:
	一个桶中只有一个元素)
	- 开放寻址法
	- 拉链法(常用)
	(TODO: 为什么哈希表比较快), map 源码, 手写实现
	拉链法将同一个桶中的元素通过链表的形式进行链接; 随着桶中元素的增加,
	可以不断连接新的元素, 同时不用预先为元素分配内存; 拉链法的不足之处在
	于需要存储额外的指针用于链接元素, 这增加了整个哈希表的大小, 同时由于
	链表存储的地址不连续,所以无法高效利用 cpu 高速缓存;

	go 语言运行对值为 nil 的 map 进行访问(TODO)

	如果不能比较 map 中的 key 是否相同, 那么这些 key 就不能作为 map 的 key;
	- 布尔值, 整数值, 浮点值, 复数值, 字符串值都是可比较的
	- 指针值是可比较的, 指向相同的变量, 或者值均为 nil, 则指针相等
	- 通道值是可比较的, 都由相同的 make 函数调用创建(长度), 或者值均为 nil, 则相等
	- 接口值是可比较的, 如果两个接口值具有相同的动态类型和相等的动态值,
		或者两个接口值都为 nil, 则相等
	- 如果结构体的所有字段都是可比较的, 则它们的值是可比较的
	- 如果数组元素类型的值可比较, 则数组值可比较; 如果两个数组对应的元素相等,
		则相等
	- 切片, 函数和 map 是不可比较的(TODO)


	map 不支持并发的读写, 只支持并发读
	map 不需要从多个 goroutine 安全访问, 在实际情况下, map 可能是某些已经同步
	的较大数据结构或计算的一部分; 因此, 要求所有 map 操作都互斥将减慢大多数
	程序的速度, 而只是为了增加少数程序的安全性; 即支持并发读的原因是为了保证
	大多数场景下的查找效率;
	TODO: sync.Map 的实现和使用
*/

func main() {
	aa := make(map[int]int)
	go func() {
		for {
			aa[0] = 5
		}
	}()

	go func() {
		for {
			_ = aa[1]
		}
	}()
	select {}
}

// fatal error: concurrent map read and map write
