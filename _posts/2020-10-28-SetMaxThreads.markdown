---
title: SetMaxThreads
layout: post
category: golang
author: 夏泽民
---
https://github.com/golang/go/issues/16076
https://studygolang.com/articles/18392

func setProgramLimits() {
	// Swarm runnable threads could be large when the number of nodes is large
	// or under request bursts. Most threads are occupied by network connections.
	// Increase max thread count from 10k default to 50k to accommodate it.
	const maxThreadCount int = 50 * 1000
	debug.SetMaxThreads(maxThreadCount)
}
https://golang.hotexamples.com/zh/examples/runtime.debug/-/SetMaxThreads/golang-setmaxthreads-function-examples.html
<!-- more -->
SetMaxStack设置该以被单个go程调用栈可使用的内存最大值。如果任何go程在增加其调用栈时超出了该限制，程序就会崩溃。SetMaxStack返回之前的设置。默认设置在32位系统是250MB，在64位系统是1GB。
SetMaxThreads设置go程序可以使用的最大操作系统线程数。如果程序试图使用超过该限制的线程数，就会导致程序崩溃。SetMaxThreads返回之前的设置，初始设置为10000个线程。


fmt.Println("runtime.NumCPU:", runtime.NumCPU())
fmt.Println("runtime.NumCgoCall:", runtime.NumCgoCall())
fmt.Println("runtime.NumGoroutine:", runtime.NumGoroutine())
fmt.Println("runtime.GOMAXPROCS:", runtime.GOMAXPROCS(0)) //GOMAXPROCS设置可同时执行的最大CPU数

https://www.jianshu.com/p/9337693037d8

https://juejin.im/entry/6844903597780516872
