---
title: interface
layout: post
category: golang
author: 夏泽民
---
https://xargin.com/go-and-interface/
https://github.com/cch123/go-internals
<!-- more -->
http://blog.comprehend.in/2017/05/07/golang-type-assert.html

Go语言中的类型断言，语法上是这样的:

x.(T)

其中，x是interface接口的表达式，T是类型，称为被断言类型。

补充一下，接口有接口值的概念，其包括动态类型和动态值两部分。

类型断言根据T的不同可以分为两种情况:

1. T是具体类型
类型断言首先检查x的动态类型是否是T.
如果是，类型断言的结果就是x的动态值.
如果否，操作就会panic
例如，

var w io.Writer
w = os.Stdout
f := w.(*os.File) // success: f == os.Stdout
c := w.(*bytes.Buffer) // panic: interface holds *os.File, not *bytes.Buffer

其中，os.Stdout的类型就是*os.File

2. T是接口类型
类型断言首先检查x的动态类型是否满足T。
如果断言成功，x的动态值不会被提取，结果仍然是以前的动态类型和动态值。但结果的类型是接口类型T.

换句话说，对接口类型的断言，结果的类型是T，不是x的类型，从而改变了可访问的方法集合(通常更大)，但会保留x接口值中的动态类型和动态值。

举个例子，下面的语句中，
第一个类型断言后，
w、rw两个接口的动态值都是os.Stout，动态类型都是*os.File。

w的类型是io.Writer，只暴露Write方法，
rw的类型是io.ReadWriter,暴露Read和Write方法。

var w io.Writer
w = os.Stdout
rw := w.(io.ReadWriter) // success: *os.File has both Read and Write

w = new(ByteCounter)
rw = w.(io.ReadWriter) // panic: *ByteCounter has no Read method

另外，无论被断言的类型是什么，如果操作数x是一个nil 接口类型，断言就会失败。

我们几乎不需要对一个更少限制性接口类型(或者说拥有更少的方法)进行断言，
这就跟赋值操作是一样的。

w = rw // io.ReadWriter is assignable to io.Writer
w = rw.(io.Writer) // fails only if rw == nil

https://studygolang.com/articles/2157
https://www.cnblogs.com/susufufu/p/7353290.html