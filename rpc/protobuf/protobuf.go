package main

// TODO: Protobuf 基本用法

/*
 * Protobuf 是 Protocol Buffers 的简称, 谷歌开发的一种数据描述语言,
 * Protobuf 刚开源时的定位类似于 XML, JSON 等数据描述语言, 通过附带工具
 * 生成代码并实现将结构化数据序列化的功能; 但在 rpc 的应用上, 最主要的
 * 是 Protobuf 作为接口规范的描述语言, 可以作为设计安全的跨语言的 rpc 接口
 * 的基础工具.
 */

// 通过 Protobuf 保证 rpc 的接口规范和安全, Protobuf 中最基本的数据单元是
// message, 是类似 go 语言中结构体的存在; 在 message 中可以嵌套 message 或
// 其他的基础数据类型的成员.
