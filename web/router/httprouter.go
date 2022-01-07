package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
 * go 的 net/http 包提供了基础的路由函数组合与丰富的功能函数, 如果项目的路由
 * 在个位数, URI 固定且不通过 URI 来传递参数的情况下使用 http 提供的默认路由
 * 就可以满足功能, 但在复杂的场景下, 就需要使用路由器(软件层面)来处理路由
 *
 * go 的 web 框架大致分为:
 * - Router 框架
 * - MVC 框架
 */

// go 应用最广泛的路由器是 httprouter(TODO: 源码), 很多开源的路由器框架都是
// 基于 httprouter 进行一定程度的改造和封装; gin 就是 httprouter 的变种.

// 请求路由
// 在常见的 web 框架中, 路由器是必备的组件, go 的路由器也称为 http 的多路复用器,

// REST(TODO), 除了GET和POST之外还使用了HTTP协议定义的几种其他标准化语义
const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" // RFC 5789
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

// REST 风格的 API 重度依赖请求路径, 会将很多参数放在请求 URI

// httprouter
// 流行的开源 go web 框架大多使用 httprouter, 或基于 httprouter 的变种对路由
// 进行支持; httprouter 的路由使用显式匹配, 在设计路由时应规避会导致路由冲突
// 的情况:(TODO:源码)
// 冲突:
// GET /user/info/:name
// GET /user/:id
// 不冲突:
// GET /user/info/:name
// POST /user/:id
// 如果两个路由拥有一致的 http 方法(get/post/put/delete)和请求路径前缀, 且在某
// 个位置出现了 A 路由是 wildcard(指:id这种形式)参数, B 路由则是普通字符串,
// 则会发送路由冲突, 在初始化阶段会直接 panic

// 同时, httprouter 考虑到字典树的深度, 在初始化时会对参数的数量进行限制,
// 所以在路由中的参数数目不能超过 255, 否则会导致 httprouter 无法识别后续
// 参数;

/*
 除支持路径中的 wildcard 参数外, httprouter 还可以支持 * 号进行通配, 不过 * 号
 开头的参数只能放在路由的结尾, 如下:
 Pattern: /src/*filepath
 /src/                       filepath = ""
 /src/somefile.go            filepath = "somefile.go"
 /src/subdir/somefile.go     filepath = "subdir/somefile.go"

 这种设计在 REST 风格中不太常见, 主要是为了能使 httprouter 做简单的 http 静态
 文件服务器
*/

/*
  httprouter 原理
  httprouter 和很多衍生路由库使用的数据结构被称为压缩动态检索树(Compressing
  Dynamic Trie), 是检索树(Trie Tree) 的一种.

  如下是一个典型的检索树结构, (TODO: vim绘图插件)


							     O
								/|\
							   / | \
						      /  |  \
						     /   |   \
						    /    |    \
						  |/    \|/    \|
						  N		 K      U -----------> T
						 /\
						/  \
			      	   /    \
					  /      \|
					 /        L ---------------> /(URI中的/)
				   |/                          /
			       /(URI中的/)                /
                                            |/
                                            N
											 \
											  \
											   \
											    \|
												 E

  检索树常用来进行字符串检索, 例如为给定的字符串序列建立检索树, 对于目标字符
  串, 只要从根结点开始深度优先搜索, 即可判断出该字符串是否出现过, 时间复杂度
  为 O(n), n 可以认为是目标字符串的长度; 因为字符串本身不像数值类型可以进行
  数值比较, 两个字符串对比的时间复杂度取决于字符串长度, 如果不用检索树实现,
  就要对历史字符串进行排序, 再利用二分搜索类的算法进行搜索, 时间复杂度只高不
  低, 可认为检索树是一种空间换时间的典型做法;
  普通的检索树明显的缺点是每个字母都需要建立一个子结点, 这样会导致检索树的
  层级比较深, 压缩检索树相对平衡了检索树的优点和缺点.

  如下是一个典型的压缩检索树结构:

							     /user

						       /  |   \
							  /   |    \
						     /    |     \|
						   |/     \    /search/
						 /info     \|       \
                          /       /addr      \
                         /         |          \|
						 \         |          :city_id
                         \|        \|            /
						 :id       :id          /
                                                |
                                                |
                                               \|/
                                              :type



  现在每个结点上可以存储字符串了, 这也是"压缩"的含义, 使用压缩检索树可以减少
  树的层数, 同时因为每个结点上的数据存储也比通常的检索树更多, 所以程序的局部
  性更好(一个结点的路径加载到缓存即可进行多个字符的对比), 从而对 cpu 友好.
*/

// 压缩动态检索树创建过程
// 路由设定如下:
// - PUT "/user/installations/:installation_id/repositories/:reposit"
//
// - GET /marketplace_listing/plans
// - GET /marketplace_listing/plans/:id/accounts
// - GET /search
// - GET /status
// - GET /support
//

/* 根结点创建 ---------------------------------------------------------------*/
/*
	httprouter 的 Router 结构体中存储压缩检索树使用数据结构:
	type Router struct {
		trees map[string]*node
		// ...
	}

	trees 字段的键为 HTTP1.1 的 RFC 中定义的各种方法:
	GET
	HEAD
	OPTIONS
	POST
	PUT
	PATCH
	DELETE

	每一种方法对应的都是一棵独立的压缩检索树, 这些树彼此之间不共享数据

*/

func main() {
	r := httprouter.New()
	r.PUT("/user/installations/:installation_id/repositories/:reposit",
		func(http.ResponseWriter, *http.Request, httprouter.Params) {})
	// ...
	// r.tress 字段未导出, 如何获取压缩检测树的结点? TODO
}

// PUT 对应的根结点会被创建出来, 其中:
// radix(计数根) 的结点类型为 *httprouter.node
// path: 当前结点对应的路径中的字符串
// wildChild: 子结点是否为参数结点, 即 :id 这种类型的结点
// nType: 当前结点的类型, 4个枚举值
// - static: 非根结点的普通字符串结点
// - root: 根结点
// - param: 参数结点, 如 :id
// - catchAll: 通配符结点, 如 *anyway
// indices: 子结点索引, 当子结点为非参数类型, 即本结点的 wildChild 为 false
// 时, 会将每个子结点的首字母放入该索引数组(字符串)
/*
    PUT 树

	path: "/user/installations/"
	wildChild: true(子结点为参数结点)
	nType: root
	indices: ""



	path: ":installation_id"
	wildChild: false
	nType: param
	indices: ""



	path: "/repositories/"
	wildChild: true
	nType: default (TODO:?)
	indices: ""



	path: ":reposit"
	wildChild: false
	nType: param
	indices: ""

*/

/* 子结点插入 ---------------------------------------------------------------*/
// - GET /marketplace_listing/plans, 此路由没有参数, 所以 path 都被存储到根结点
/*

    GET 树

	path: "/marketplace_listing/plans/" (TODO:不需要最后的 / ?)
	wildChild: false
	nType: root
	indices: ""

*/

// - GET /marketplace_listing/plans/:id/accounts, 在 GET 树中和上条路径有相同
// 的前缀, 所以作为上条路径的子节点进行插入
/*

    GET 树

	path: "/marketplace_listing/plans/"
	wildChild: true(注意子结点插入前后的变化)
	nType: root
	indices: ""

	path: ":id"
	wildChild: false
	nType: param
	indices: ""

	path: "/accounts"
	wildChild: false
	nType: default
	indices: ""

	// 由于 :id 这个结点是只有一个字符串的普通子结点, 因此 indices 不需要处理

*/

/* 边分裂 ------------------------------------------------------------------*/
// - GET /search, 会导致树的边分裂

/*
    GET 树

	path: "/"
	wildChild: false
	nType: root
	indices: "ms"

		|             \
		|              \
		|               \
		|                \
		|                 \|
		|
		|                path: "search"
		|                wildChild: false
		|                nType: default
		|                indices: ""
		|
	   \|/

	path: "marketplace_listing/plans/"
	wildChild: true
	nType: root
	indices: ""
		|
		|
	   \|/
	path: ":id"
	wildChild: false
	nType: param
	indices: ""
		|
		|
	   \|/
	path: "/accounts"
	wildChild: false
	nType: default
	indices: ""

	原有路径和新的路径在初始的 / 位置发生分裂, 这样就需要把原有的根结点内容下
	移, 再将新路由 search 同样作为子结点挂在根结点之后; 此时因为出现多个子结点,
	根结点的 indices 提供子结点索引, 其值 "ms" 代表子结点的首字母, 分别为
	m(marketplace_listing), s(search)

*/

// 继续插入以下子结点
// - GET /status
// - GET /support
/*
    GET 树

	path: "/"
	wildChild: false
	nType: root
	indices: "ms"

		|             \
		|              \                            ------->  path: "earch"
		|               \                          |          wildChild: false
		|                \                         |          nType: default
		|                 \|                       |          indices: ""
		|                                          |
		|                path: "s"                 |
	 	|                wildChild: false          |         path: "tatus"
		|                nType: default   ---------|----->   wildChild: false
		|                indices: "etu"            |         nType: default
		|                                          |         indices: ""
	   \|/                                          \
										             \
	path: "marketplace_listing/plans/"                \
	wildChild: true                                   |
	nType: root                                      \|/
	indices: ""                                 path: "upport"
	    |                                       wildChild: false
		|                                       nType: default
	   \|/                                      indices:""
	path: ":id"
	wildChild: false
	nType: param
	indices: ""
		|
		|
	   \|/
	path: "/accounts"
	wildChild: false
	nType: default
	indices: ""

	原有路径和新的路径在初始的 / 位置发生分裂, 这样就需要把原有的根结点内容下
	移, 再将新路由 search 同样作为子结点挂在根结点之后; 此时因为出现多个子结点,
	根结点的 indices 提供子结点索引, 其值 "ms" 代表子结点的首字母, 分别为
	m(marketplace_listing), s(search)

*/

/* 子结点冲突处理 -----------------------------------------------------------*/
// 只要发生冲突, 在初始化阶段就会 panic
/*
   在路由本身只有字符串的情况下, 不会发生任何冲突; 只有当路由中含有通配符(如:id)
   或 catchAll 的情况下才可能发生冲突, 冲突处理分几种情况:
   - 在插入 wildChild 时父结点的 children 数组非空且 wildChild 被设置为 false.
	   如: GET /user/getAll 和 GET /user/:id/getAddr

   - 在插入 wildChild 结点时, 父结点的 children 数组非空且 wildChild 被设置为
	 true, 但该父结点的已存在的 wildChild 子结点与即将插入的 wildChild 名字不
	 一样. 如: GET /user/:id/info 和 GET /user/:name/info

  - 在插入 catchAll 结点时, 父结点的 children 非空.
		如: GET /src/abc 和 GET /src/*filename

  - 在插入 static 结点时, 父结点的 wildChild 字段被设置为 true.

  - 在插入 static 结点时, 父结点的 children 非空, 同时子结点 nType 为 catchAll
*/
