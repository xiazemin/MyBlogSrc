---
title: 序列化性能比较
layout: post
category: golang
author: 夏泽民
---
https://github.com/alecthomas/go_serialization_benchmarks

https://github.com/smallnest/gosercomp

 MessagePack,gogo/protobuf,和flatbuffers差不多，这三个优秀的库在序列化和反序列化上各有千秋，而且都是跨语言的。 从便利性上来讲，你可以选择MessagePack和gogo/protobuf都可以，两者都有大厂在用。 flatbuffers有点反人类，因为它的操作很底层，而且从结果上来看，序列化的性能要差一点。但是它有一个好处，那就是如果你只需要特定的字段， 你无须将所有的字段都反序列化。从结果上看，不反序列化字段每个调用只用了9.54纳秒，这是因为字段只有在被访问的时候才从byte数组转化为相应的类型。 因此在特殊的场景下，它可以提高N被的性能。但是序列化的代码的面相太难看了。

新增加了gencode的测试，它表现相当出色，而且生成的字节也非常的小。

https://studygolang.com/articles/9312
<!-- more -->
方式	优点	缺点
binary	性能高	不支持不确定大小类型 int、slice、string
gob	支持多种类型	性能低
json	支持多种类型	性能低于 binary 和 protobuf
protobuf	支持多种类型，性能高	需要单独存放结构，如果结构变动需要重新生成 .pb.go 文件

https://www.cnblogs.com/rsapaper/p/10237583.html


encoding/gob包(The way to go)
gob的优势
这里需要明确一点，gob只能用在golang中，所以在实际工程开发过程中，如果与其他端，或者其他语言打交道，那么gob是不可以的，我们就要使用json了。
Gob is much more preferred when communicating between Go programs. However, gob is currently supported only in Go and, well, C, so only ever use that when you’re sure no program written in any other programming language will try to decode the values.

gob的优势就是：发送方的结构和接受方的结构并不需要完全一致，例如定义一个结构体：

struct { A, B int }
下面的类型都是可以发送、接收的：

struct { A, B int } // the same
*struct { A, B int }    // extra indirection of the struct
struct { *A, **B int }  // extra indirection of the fields
struct { A, B int64 }   // different concrete value type; see below
以下可以接收：

struct { A, B int } // the same
struct { B, A int } // ordering doesn't matter; matching is by name
struct { A, B, C int }  // extra field (C) ignored
struct { B int }    // missing field (A) ignored; data will be dropped
struct { B, C int } // missing field (A) ignored; extra field (C) ignored.
下面的格式是有问题的：

struct { A int; B uint }    // change of signedness for B
struct { A int; B float }   // change of type for B
struct { }          // no field names in common
struct { C, D int }     // no field names in common

https://studygolang.com/articles/19855

https://www.jianshu.com/p/e8166a20935a?utm_campaign=maleskine&utm_content=note&utm_medium=seo_notes&utm_source=recommendation

包gob管理gobs流 - 二进制值在编码器（发送器）和解码器（接收器）之间交换
。 >
json包json实现了
RFC 4627中定义的JSON的编码和解码。

Gob更受欢迎在Go程序之间进行通信时。但是，目前仅在Go中支持gob，并且 C 只支持gob当你确定没有任何其他编程语言编写的程序会尝试解码这些值。

https://www.it1352.com/809206.html
