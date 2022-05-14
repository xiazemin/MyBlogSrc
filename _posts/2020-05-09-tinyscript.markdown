---
title: 写一个语言 编译器 虚拟机
layout: post
category: golang
author: 夏泽民
---
https://gocn.vip/topics/10361
https://github.com/elvin-du/tinyscript
整个项目包括三个东西：

创建了一个自己的语言
编译器
虚拟机
golang 实现的一个编译器，用来编译一个自己创建的语言（用来玩的），最后写了一个自定义虚拟机用来运行自定义语言。
<!-- more -->
语言介绍
为了跨平台（其实是为了方便开发 ^ ^），所以这个语言没有静态编译成硬件指令集，最后的机器码是我自己的定义的，和 MIPS 类似的（其实就是一个 mips 子集）虚拟指令集。为了运行这些指令集，我写了一个虚拟机。

语言和 golang 和 javascript 类似，实现了函数，类型声明，函数调用等最基本的一些语言元素，没有实现类，结构体，接口等复杂数据结构。 下面是用这个语言编程的例子：

func fact(int n)  int {
    if(n == 0) {
        return 1
    }
    return fact(n-1) * n
}
func main() void {
    return fact(2)
}
每个函数都实现了相应的 UnitTest