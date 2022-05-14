---
title: Interface
layout: post
category: golang
author: 夏泽民
---
https://mp.weixin.qq.com/s/9CciL8ifi_Q_qfyNEPLNlg
<!-- more -->
一个 interface{} 可以包含任何数据，同时他也是一个非常有用的参数，因为他可以是任何的 type。要了解 interface{} 是如何工作的，它是怎样可以符合任何类型的，我们必须要先知道它名字背后的概念。



Interfaces（接口）

这里有一个对空接口很好的定义，by Jordan Oreilli:



interface 包含两种功能：一坨方法的集合，同时他也是一个类型 (type)



interface{} 的类型是一个没有方法的接口。由于它没有需要实现的方法，所以所有的类型都至少实现了零个方法，自然满足了 interface 类型的条件，所有的类型都符合空的 interface



一个使用 Interface{} 类型做参数的方法，可以接收任何类型的参数。Go 会为方法把参数转化为 interface 类型。

Russ Cox 曾写过一篇内部描述 interfaces 的文章，描述了一个 interface 由两部分组成：

一个指针指向存储类型信息

一个指针指向具体的数据



Russ  用 C 语言描述了 interface 的定义







虽然现在的 runtime 使用 Go 写的，但是表述依然是相同的。我们可以通过打印出 interface{} 的指针地址来证明他：





func main() {
  var i int8 = 1
  read(i)
}

//go:noinline
func read(i interface{}) {
  println(i)
}

print:
(0x10591e0,0x10be5c6)



两个指针地址一个指向了类型信息，另一个指向了值。



Underlying structure



空接口的底层描述是在 reflection package 的文档里：

type emptyInterface struct {
   typ  *rtype            // word 1 with type description
   word unsafe.Pointer    // word 2 with the value
}
和上面说的一样，我们可以看到空 interface 有一个描述 type 的指针，以及包含具体数据的 word 





rtype 的结构体包含了具体的 type 的描述：

type rtype struct {
   size       uintptr
   ptrdata    uintptr
   hash       uint32
   tflag      tflag
   align      uint8
   fieldAlign uint8
   kind       uint8
   alg        *typeAlg
   gcdata     *byte
   str        nameOff
   ptrToThis  typeOff
}



在这些字段中，有一些我们已知的内容：

size 是字节大小

kind 包含了 int8, int16, bool 等类型

align 是这个类型的变量的对齐方式



根据嵌入到空接口的类型，我们可以映射出字段或者方法：

type structType struct {
   rtype
   pkgPath name
   fields  []structField
}


该结构体还有两个包含字段列表的映射关系。它清楚地展示了将内建的类型转换成空接口将会导致一次水平转换，在字段描述和它的值被存在内存中的地方



这是我们看到的空接口的描述：





interface is composed by two words



现在让我们看下，从空接口实际上可以转换到那种类型上。



Conversions



让我们试试用空接口来转换错误的类型

func main() {
  var i int8 = 1
  read(i)
}

//go:noinline
func read(i interface{}) {
  n := i.(int16)
  println(n)
}


虽然从 int8 转到 int16 是可以的，但是程序会 panic:



panic: interface conversion: interface {} is int8, not int16

goroutine 1 [running]:
main.read(0x10592e0, 0x10be5c1)
main.go:10 +0x7d
main.main()
main.go:5 +0x39
exit status 2


我们生产 asm 代码来看看 Go 到底检查了什么





code generated while checking the type of an empty interface



步骤如下



步骤 1： 用 int16 (指令 LEAQ：Load Effective Address)的类型与空接口的内部类型（指令 MOVQ：读取空接口 48字节偏移的内存）作比较（指令 CMPQ）

步骤 2：指令 JNE，如果不等于就跳转，跳转到在第 3 步创建的指令，来处理错误信息

步骤 3：代码抛出异常，并生成上一步的错误信息

步骤 4：这是错误指令的结尾。这个特定的指令由显示该指令（main.go:10+0x7d）错误信息所引用



任何从空接口内部类型的转换都是在原始类型转换之后进行。转换成空接口然后再转换回原始类型对你的程序耗时有一些影响。让我们运行一些基准测试来大致了解一下。



性能

这里有两个基准测试。一个是复制一个结构体，另外一个是复制一个空接口。

package main_test

import (
  "testing"
)

var x MultipleFieldStructure

type MultipleFieldStructure struct {
  a int
  b string
  c float32
  d float64
  e int32
  f bool
  g uint64
  h *string
  i uint16
}

//go:noinline
func emptyInterface(i interface {}) {
  s := i.(MultipleFieldStructure)
  x = s
}

//go:noinline
func typed(s MultipleFieldStructure) {
  x = s
}

func BenchmarkWithType(b *testing.B) {
  s := MultipleFieldStructure{a: 1, h: new(string)}
  for i := 0; i < b.N; i++ {
    typed(s)
  }
}

func BenchmarkWithEmptyInterface(b *testing.B) {
  s := MultipleFieldStructure{a: 1, h: new(string)}
  for i := 0; i < b.N; i++ {
    emptyInterface(s)
  }
}


结果如下



BenchmarkWithType-8               300000000           4.24 ns/op
BenchmarkWithEmptyInterface-8      20000000           60.4 ns/op


将类型转换为空接口，再从空接口转换回类型的两次过程消耗量超过 55 纳秒。而且时长会随着结构体内部字段的增加而增加：



BenchmarkWithType-8             100000000         17 ns/op
BenchmarkWithEmptyInterface-8    10000000        153 ns/op


但是，用指针转换回相同的结构体指针是一个好办法。转换看起来是下面这个样子



func emptyInterface(i interface {}) {
  s := i.(*MultipleFieldStructure)
  y = s
}


现在结果有显著的不同



BenchmarkWithType-8                 2000000000          2.16 ns/op
BenchmarkWithEmptyInterface-8       2000000000          2.02 ns/op


对于像 int 或者 string 这样的基本类型，性能测试略微不同



int:

BenchmarkWithTypeInt-8              2000000000          1.42 ns/op
BenchmarkWithEmptyInterfaceInt-8    1000000000          2.02 ns/op


string:

BenchmarkWithTypeString-8           1000000000          2.19 ns/op
BenchmarkWithEmptyInterfaceString-8  50000000           30.7 ns/op




合理并节制使用空接口，在大多数情况下，空接口会对你的程序性能造成一些真正的影响。
