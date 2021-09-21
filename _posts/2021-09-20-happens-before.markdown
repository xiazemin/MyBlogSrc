---
title: happens-before
layout: post
category: golang
author: 夏泽民
---
happens-before是一个术语，并不仅仅是Go语言才有的。假设A和B表示一个多线程的程序执行的两个操作。如果A happens-before B，那么A操作对内存的影响 将对执行B的线程(且执行B之前)可见。

无论使用哪种编程语言，有一点是相同的：如果操作A和B在相同的线程中执行，并且A操作的声明在B之前，那么A happens-before B。

还有一点是，在每门语言中，无论你使用那种方式获得，happens-before关系都是可传递的：如果A happens-before B，同时B happens-before C，那么A happens-before C。当这些关系发生在不同的线程中，传递性将变得非常有用。

happens-before并不是指时序关系，并不是说A happens-before B就表示操作A在操作B之前发生。它就是一个术语，就像光年不是时间单位一样。happens-before 是一系列语言规范中定义的操作间的关系。它和时间的概念独立。这和我们通常说”A在B之前发生”时表达的真实世界中事件的时间顺序不同。

A happens-before B并不意味着A在B之前发生。
A在B之前发生并不意味着A happens-before B。

<!-- more -->
{% raw %}
int A = 0;
int B = 0;
void main()
{
    A = B + 1; // (1)
    B = 1; // (2)
}
{% endraw %}

(1) happens-before (2)。但是，如果我们使用gcc -O2编译这个代码，编译器将产生一些指令重排序。有可能执行顺序是这样子的：

将B的值取到寄存器
将B赋值为1
将寄存器值加1后赋值给A

也就是到第二条机器指令(对B的赋值)完成时，对A的赋值还没有完成。换句话说，(1)并没有在(2)之前发生!

根据定义，操作(1)对内存的影响必须在操作(2)执行之前对其可见。换句话说，对A的赋值必须有机会对B的赋值有影响.

但是在这个例子中，对A的赋值其实并没有对B的赋值有影响。即便(1)的影响真的可见，(2)的行为还是一样。所以，这并不能算是违背happens-before规则。

golang happen before 的保证
1) 单线程

2) Init 函数

如果包P1中导入了包P2，则P2中的init函数Happens Before 所有P1中的操作
main函数Happens After 所有的init函数
3) Goroutine

Goroutine的创建Happens Before所有此Goroutine中的操作
Goroutine的销毁Happens After所有此Goroutine中的操作
4) Channel

对一个元素的send操作Happens Before对应的receive 完成操作
对channel的close操作Happens Before receive 端的收到关闭通知操作
对于Unbuffered Channel，对一个元素的receive 操作Happens Before对应的send完成操作
对于Buffered Channel，假设Channel 的buffer 大小为C，那么对第k个元素的receive操作，Happens Before第k+C个send完成操作。可以看出上一条Unbuffered Channel规则就是这条规则C=0时的特例
5) Lock

Go里面有Mutex和RWMutex两种锁，RWMutex除了支持互斥的Lock/Unlock，还支持共享的RLock/RUnlock。

对于一个Mutex/RWMutex，设n < m，则第n个Unlock操作Happens Before第m个Lock操作。
对于一个RWMutex，存在数值n，RLock操作Happens After 第n个UnLock，其对应的RUnLock Happens Before 第n+1个Lock操作。
简单理解就是这一次的Lock总是Happens After上一次的Unlock，读写锁的RLock HappensAfter上一次的UnLock，其对应的RUnlock Happens Before 下一次的Lock。

var l sync.Mutex
var a string
func f() {
    a = "hello, world" // (1)
    l.Unlock() // (2)
}
func main() {
    l.Lock() // (3)
    go f()
    l.Lock() // (4)
    print(a) // (5)
}
(1) happens-before (2) happens-before (4) happens-before (5)

6) Once

once.Do中执行的操作，Happens Before 任何一个once.Do调用的返回。

https://www.liuvv.com/p/6196d525.html

https://zhuanlan.zhihu.com/p/33829310


https://blog.csdn.net/m0_37055174/article/details/104404066

在单一goroutine 中Happens Before所要表达的顺序就是程序执行的顺序，happens before原则指出在单一goroutine 中当满足下面条件时候，对一个变量的写操作w1对读操作r1可见：

读操作r1没有发生在写操作w1前
在读操作r1之前，写操作w1之后没有其他的写操作w2对变量进行了修改
在一个goroutine里面，不存在并发，所以对变量的读操作r1总是对最近的一个写操作w1的内容可见，但是在多goroutine下则需要满足下面条件才能保证写操作w1对读操作r1可见：

写操作w1先于读操作r1
任何对变量的写操作w2要先于写操作w1或者晚于读操作r1
这两条条件相比第一组的两个条件更加严格，因为它要求没有任何写操作与w1或者读操作r1并发的运行，而是要求在w1操作前或读操作r1后发生。

在一个goroutine时候，不存在与w1或者r1并发的写操作，所以前面两种定义是等价的：一个读操作r1总是对最近的一个对写操作w1的内容可见。但是当有多个goroutines并发访问变量时候，就需要引入同步机制来建立happen-before条件来确保读操作r1对写操作w1写的内容可见。

需要注意的是在go内存模型中将多个goroutine中用到的全局变量初始化为它的类型零值在内被视为一次写操作，另外当读取一个类型大小比机器字长大的变量的值时候表现为是对多个机器字的多次读取，这个行为是未知的，go中使用sync/atomic包中的Load和Store操作可以解决这个问题。

解决多goroutine下共享数据可见性问题的方法是在访问共享数据时候施加一定的同步措施，比如sync包下的锁或者通道。

https://studygolang.com/articles/18297
https://studygolang.com/articles/33719

https://zhuanlan.zhihu.com/p/409599399

