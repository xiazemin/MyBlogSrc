---
title: copylocks vet
layout: post
category: golang
author: 夏泽民
---
go tool vet -copylocks 命令简介

Go 语言安装包附带 vet 命令行工具。该工具能对程序源码运行一套启发式算法以发现可疑的程序结构，如无法执行的代码或对 fmt.Printf 函数的错误调用（指 arguments 没有对齐 format 参数）


<!-- more -->
{% raw %}
package main
import "fmt"

func f() {
    fmt.Printf("%d\n")
    return
    fmt.Println("Done")
}
> go tool vet vet.go
vet.go:8: unreachable code
vet.go:6: missing argument for Printf("%d"): format reads arg 1, have only 0 args
本文专讲该工具的 copylocks 选项。让我们看看它能做什么以及如何在实际的程序中发挥作用。

假设程序使用互斥锁进行同步：

package main
import "sync"

type T struct {
    lock sync.Mutex
}
func (t *T) Lock() {
    t.lock.Lock()
}
func (t T) Unlock() {
    t.lock.Unlock()
}

func main() {
    t := T{lock: sync.Mutex{}}
    t.Lock()
    t.Unlock()
    t.Lock()
}
如果变量 v 是可寻址的，并且 &v 的方法集合包含 m，那么 v.m() 是 (&v).m() 的简写。
想一想上述程序运行的结果可能是什么...

程序会进入死锁状态：

fatal error: all goroutines are asleep — deadlock!
goroutine 1 [semacquire]:
sync.runtime_Semacquire(0x4201162ac)
    /usr/local/go/src/runtime/sema.go:47 +0x30
sync.(*Mutex).Lock(0x4201162a8)
    /usr/local/go/src/sync/mutex.go:85 +0xd0
main.(*T).Lock(0x4201162a8)
...
运行上述程序得到了糟糕的结果，根本原因是把 receiver 按值传递给 Unlock 方法，所以 t.lock.Unlock() 实际上是由 lock 的副本调用的。我们很容易忽视这点，特别在更大型的程序中。Go 编译器不会检测这方面，因为这可能是程序员有意为之。该 vet 工具登场啦...

> go tool vet vet.go
vet.go:13: Unlock passes lock by value: main.T
选项 copylocks (默认启用) 会检测拥有 Lock 方法 (实际需要 pointer receiver) 的 type 是否按值传递。如果是这种情况，则会发出警告。

sync 包有使用该机制的例子，它有一个命名为 noCopy 的特殊 type。为了避免某 type 按值拷贝 (实际上通过 vet 工具进行检测)，需要往 struct 定义中添加一个 field(如 WaitGroup):

package main
import "sync"
type T struct {
    wg sync.WaitGroup
}
func fun(T) {}
func main() {
    t := T{sync.WaitGroup{}}
    fun(t)
}
> go tool vet lab.go
lab.go:9: fun passes lock by value: main.T contains sync.WaitGroup contains sync.noCopy
lab.go:13: function call copies lock value: main.T contains sync.WaitGroup contains sync.noCopy
深入理解该机制

under-the-hood

vet 工具的源文件放在 /src/cmd/vet 路径下。vet 的每个选项都利用 register 函数进行注册，该函数其中两个参数分别是一个可变参数 (类型是该选项所关注的 AST 结点类型) 和一个回调函数。该回调函数将因特定类型的结点事件触发。对于 copylocks 选项，需要检测的结点包含 return 语句。最终都会转到 lockPath，它验证传递的值是否属于某个 type(拥有一个需要 pointer receiver 的 Lock 方法)。在整个处理过程中，go/ast 包被广泛使用。可以在 Go 源码可测试的示例中找到对该包的简单介绍。
{% endraw %}

https://studygolang.com/articles/14479?from=timeline