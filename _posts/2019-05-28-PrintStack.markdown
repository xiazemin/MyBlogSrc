---
title: PrintStack
layout: post
category: golang
author: 夏泽民
---
 runtime/debug 库可以把调用堆栈打出来
 package main
import (
    "fmt"
    "runtime/debug"
)
func test3() {
    fmt.Printf("%s", debug.Stack())
    debug.PrintStack()
}
Stack trace是指堆栈回溯信息，在当前时间，以当前方法的执行点开始，回溯调用它的方法的方法的执行点，然后继续回溯，这样就可以跟踪整个方法的调用,大家比较熟悉的是JDK所带的jstack工具，可以把Java的所有线程的stack trace都打印出来。
<!-- more -->
异常退出情况下输出stacktrace
通过panic
如果应用中有没recover的panic,或者应用在运行的时候出现运行时的异常，那么程序自动会将当前的goroutine的stack trace打印出来。

输出下面的stack trace:
dump go run p.go
panic: panic from m3
goroutine 1 [running]:
panic(0x596a0, 0xc42000a1a0)
	/usr/local/Cellar/go/1.7.4/libexec/src/runtime/panic.go:500 +0x1a1
main.m3()
如果想让它把所有的goroutine信息都输出出来，可以设置 GOTRACEBACK=1:
GOTRACEBACK=1 go run p.go
panic: panic from m3

这个信息将两个goroutine的stack trace都打印出来了，而且goroutine 4的状态是sleep状态。

Go官方文档对这个环境变量有介绍：

The GOTRACEBACK variable controls the amount of output generated when a Go program fails due to an unrecovered panic or an unexpected runtime condition. By default, a failure prints a stack trace for the current goroutine, eliding functions internal to the run-time system, and then exits with exit code 2. The failure prints stack traces for all goroutines if there is no current goroutine or the failure is internal to the run-time. GOTRACEBACK=none omits the goroutine stack traces entirely. GOTRACEBACK=single (the default) behaves as described above. GOTRACEBACK=all adds stack traces for all user-created goroutines. GOTRACEBACK=system is like “all” but adds stack frames for run-time functions and shows goroutines created internally by the run-time. GOTRACEBACK=crash is like “system” but crashes in an operating system-specific manner instead of exiting. For example, on Unix systems, the crash raises SIGABRT to trigger a core dump. For historical reasons, the GOTRACEBACK settings 0, 1, and 2 are synonyms for none, all, and system, respectively. The runtime/debug package's SetTraceback function allows increasing the amount of output at run time, but it cannot reduce the amount below that specified by the environment variable. See https://golang.org/pkg/runtime/debug/#SetTraceback.

你可以设置 none、all、system、single、crash，历史原因， 你可以可是设置数字0、1、2，分别代表none、all、system。

通过SIGQUIT信号
如果程序没有发生panic,但是程序有问题，"假死“不工作，我们想看看哪儿出现了问题，可以给程序发送SIGQUIT信号，也可以输出stack trace信息。
func m3() {
	time.Sleep(time.Hour)
}
你可以运行kill -SIGQUIT <pid> 杀死这个程序，程序在退出的时候输出strack trace。

正常情况下输出stacktrace
上面的情况是必须要求程序退出才能打印出stack trace信息，但是有时候我们只是需要跟踪分析一下程序的问题，而不希望程序中断运行。所以我们需要其它的方法来执行。

你可以暴漏一个命令、一个API或者监听一个信号，然后调用相应的方法把stack trace打印出来。

打印出当前goroutine的 stacktrace
通过debug.PrintStack()方法可以将当前所在的goroutine的stack trace打印出来，如下面的程序。
打印出所有goroutine的 stacktrace
可以通过pprof.Lookup("goroutine").WriteTo将所有的goroutine的stack trace都打印出来，如下面的程序：
import (
	"os"
	"runtime/pprof"
	"time"
)
func m3() {
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	time.Sleep(time.Hour)
}

较完美的输出 stacktrace
你可以使用runtime.Stack得到所有的goroutine的stack trace信息，事实上前面debug.PrintStack()也是通过这个方法获得的。

为了更方便的随时的得到应用所有的goroutine的stack trace信息，我们可以监听SIGUSER1信号，当收到这个信号的时候就将stack trace打印出来。发送信号也很简单，通过kill -SIGUSER1 <pid>就可以，不必担心kill会将程序杀死，它只是发了一个信号而已。

配置 http/pprof
如果你的代码中配置了 http/pprof,你可以通过下面的地址访问所有的groutine的堆栈：
http://localhost:8888/debug/pprof/goroutine?debug=2.
