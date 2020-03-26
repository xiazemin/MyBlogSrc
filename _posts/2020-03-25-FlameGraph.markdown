---
title: FlameGraph
layout: post
category: golang
author: 夏泽民
---
安装
go get github.com/uber/go-torch
# 再安装 brendangregg/FlameGraph 
export PATH=$PATH:/absolute/path/FlameGraph-master
# 还需要安装一个graphviz用来画内存图
yum install graphviz

https://studygolang.com/articles/11556
<!-- more -->
1.安装组件
安装go-torch
go get github.com/uber/go-torch
安装 FlameGraph
cd $WORK_PATH && git clone https://github.com/brendangregg/FlameGraph.git
export PATH=$PATH:$WORK_PATH/FlameGraph
安装graphviz
yum install graphviz(CentOS, Redhat)


import "net/http"
import _ "net/http/pprof"
func main() {
    // 主函数中添加
    go func() {
		http.HandleFunc("/program/html", htmlHandler) // 用来查看自定义的内容
		log.Println(http.ListenAndServe("0.0.0.0:8080", nil))
	}()
}

使用
# 用 -u 分析CPU使用情况
./go-torch -u http://127.0.0.1:8080
# 用 -alloc_space 来分析内存的临时分配情况
./go-torch -alloc_space http://127.0.0.1:8080/debug/pprof/heap --colors=mem
# 用 -inuse_space 来分析程序常驻内存的占用情况；
./go-torch -inuse_space http://127.0.0.1:8080/debug/pprof/heap --colors=mem
# 画出内存分配图
go tool pprof -alloc_space -cum -svg http://127.0.0.1:8080/debug/pprof/heap > heap.svg


查看
使用浏览器查看svg文件，程序运行中，可以登录 http://127.0.0.1:10086/debug/pprof/ 查看程序实时状态 在此基础上，可以通过配置handle来实现自定义的内容查看，可以添加Html格式的输出，优化显示效果

func writeBuf(buffer *bytes.Buffer, format string, a ...interface{}) {
	(*buffer).WriteString(fmt.Sprintf(format, a...))
}
func htmlHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, statusHtml())
}
// 访问 localhost:8080/program/html 可以看到一个表格，一秒钟刷新一次
func statusHtml() string {
	var buf bytes.Buffer
	buf.WriteString("<html><meta http-equiv=\"refresh\" content=\"1\">" +
		"<body><h2>netflow-decoder status count</h2>" +
		"<table width=\"500px\" border=\"1\" cellpadding=\"5\" cellspacing=\"1\">" +
		"<tr><th>NAME</th><th>TOTAL</th><th>SPEED</th></tr>")
	writeBuf(&buf, "<tr><td>UDP</td><td>%d</td><td>%d</td></tr>",
		lastRecord.RecvUDP, currSpeed.RecvUDP)
	...
	writeBuf(&buf, "</table><p>Count time: %s</p><p>Time now: %s</p>",
		countTime.Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
	buf.WriteString("</body></html>")
	return buf.String()
}


#!/bin/bash
go build -o m  main.go
./m

$go-torch -u http://127.0.0.1:8080
INFO[13:48:29] Run pprof command: go tool pprof -raw -seconds 30 http://127.0.0.1:8080/debug/pprof/profile
ERROR: No stack counts found

$go-torch -alloc_space http://127.0.0.1:8080/debug/pprof/heap --colors=mem
INFO[13:51:37] Run pprof command: go tool pprof -raw -seconds 30 -alloc_space http://127.0.0.1:8080/debug/pprof/heap
INFO[13:51:37] Writing svg to torch.svg


$go-torch -inuse_space http://127.0.0.1:8080/debug/pprof/heap --colors=mem
INFO[13:53:28] Run pprof command: go tool pprof -raw -seconds 30 -inuse_space http://127.0.0.1:8080/debug/pprof/heap
INFO[13:53:28] Writing svg to torch.svg

$go tool pprof -alloc_space -cum -svg http://127.0.0.1:8080/debug/pprof/heap > heap.svg
Fetching profile over HTTP from http://127.0.0.1:8080/debug/pprof/heap
Saved profile in /Users/didi/pprof/pprof.alloc_objects.alloc_space.inuse_objects.inuse_space.006.pb.gz

启动服务
使用工具进行打压，压至程序问题暴露
使用go tool pprof http://localhost:9999/pprof/profile收集数据，收集30s后在进入命令行模式后输入quit退出，会生成名为【pprof.二进制名.samples.cpu.00x.pb.gz】的文件，默认放到$HOME/pprof/下
使用go tool pprof -http=:8080 [步骤5中生成的文件]解释数据并生成调用栈graph图示，点击view切换成火焰图，

火焰图(或者叫冰柱图)会默认从左到右按处理时间从大到小排列，方便定位问题；
根据火焰图的结果，分析程序响应变慢时，哪个函数占据了更多的处理时间，可以更直观的定位问题

Golang的性能调优手段
Go语言内置的CPU和Heap profiler
Go强大之处是它已经在语言层面集成了profile采样工具,并且允许我们在程序的运行时使用它们，
使用Go的profiler我们能获取以下的样本信息：
CPU profiles
Heap profiles
block profile、traces等

Go语言常见的profiling使用场景
基准测试文件：例如使用命令go test . -bench . -cpuprofile prof.cpu 生成采样文件后，再通过命令 go tool pprof [binary] prof.cpu 来进行分析。

import _ net/http/pprof：如果我们的应用是一个web服务，我们可以在http服务启动的代码文件(eg: main.go)添加 import _ net/http/pprof，这样我们的服务 便能自动开启profile功能，有助于我们直接分析采样结果。

通过在代码里面调用 runtime.StartCPUProfile或者runtime.WriteHeapProfile


go-torch是Uber
公司开源的一款针对Go语言程序的火焰图生成工具，能收集 stack traces,并把它们整理成火焰图，直观地程序给开发人员。
go-torch是基于使用BrendanGregg创建的火焰图工具生成直观的图像，很方便地分析Go的各个方法所占用的CPU的时间， 火焰图是一个新的方法来可视化CPU的使用情况，本文中我会展示如何使用它辅助我们排查问题。

安装
1.首先，我们要配置FlameGraph
的脚本
FlameGraph 是profile数据的可视化层工具，已被广泛用于Python和Node
git clone https://github.com/brendangregg/FlameGraph.git

2.检出完成后，把flamegraph.pl
拷到我们机器环境变量$PATH的路径中去，例如：
cp flamegraph.pl /usr/local/bin

3.在终端输入 flamegraph.pl -h
是否安装FlameGraph成功
$ flamegraph.pl -hOption h is ambiguous (hash, height, help)USAGE: /usr/local/bin/flamegraph.pl [options] infile > outfile.svg --title # change title text --width # width of image (default 1200) --height # height of each frame (default 16) --minwidth # omit smaller functions (default 0.1 pixels) --fonttype # font type (default "Verdana") --fontsize # font size (default 12) --countname # count type label (default "samples") --nametype # name type label (default "Function:") --colors # set color palette. choices are: hot (default), mem, io, # wakeup, chain, java, js, perl, red, green, blue, aqua, # yellow, purple, orange --hash # colors are keyed by function name hash --cp # use consistent palette (palette.map) --reverse # generate stack-reversed flame graph --inverted # icicle graph --negate # switch differential hues (blue<->red) --help # this message eg, /usr/local/bin/flamegraph.pl --title="Flame Graph: malloc()" trace.txt > graph.svg

4.安装go-torch
有了flamegraph的支持，我们接下来要使用go-torch展示profile的输出，而安装go-torch很简单，我们使用下面的命令即可完成安装
go get -v github.com/uber/go-torch

5.使用go-torch命令
$ go-torch -hUsage: go-torch [options] [binary] <profile source>pprof Options: -u, --url= Base URL of your Go program (default: http://localhost:8080) -s, --suffix= URL path of pprof profile (default: /debug/pprof/profile) -b, --binaryinput= File path of previously saved binary profile. (binary profile is anything accepted by https://golang.org/cmd/pprof) --binaryname= File path of the binary that the binaryinput is for, used for pprof inputs -t, --seconds= Number of seconds to profile for (default: 30) --pprofArgs= Extra arguments for pprofOutput Options: -f, --file= Output file name (must be .svg) (default: torch.svg) -p, --print Print the generated svg to stdout instead of writing to file -r, --raw Print the raw call graph output to stdout instead of creating a flame graph; use with Brendan Gregg's flame graph perl script (see https://github.com/brendangregg/FlameGraph) --title= Graph title to display in the output file (default: Flame Graph) --width= Generated graph width (default: 1200) --hash Colors are keyed by function name hash --colors= set color palette. choices are: hot (default), mem, io, wakeup, chain, java, js, perl, red, green, blue, aqua, yellow, purple, orange --cp Use consistent palette (palette.map) --reverse Generate stack-reversed flame graph --inverted icicle graphHelp Options: -h, --help Show this help message


https://github.com/uber-archive/go-torch
https://blog.golang.org/profiling-go-programs

内存使用调优
内存调优主要是使用上面那个pprof图，观察流程是否合理，是否可以简化，以及每个函数的内存分配情况，具体过程不像上面那么清洗，都是小修小补，故直接总结一些可能不够可靠的经验：

减少不必要的临时变量，函数的参数如果比较长则应该传递指针
在字节流处理中，原来经常出现使用 bytes.NewBuffer(buffer) 作为参数的情况，这种用法是为了使用 bytes.Buffer 的一系列函数，但是需要重新申请一次空间，其实这样会多申请一个bytes.Buffer对象，如果操作比较简单，可以直接对buffer数组进行，不用转换。还有就是 string 的转换也会申请空间，比如把 []byte 转 string ，做个简单的处理又转成 []byte 发送出去 ，可以尽量去掉中间的过程
如果已知切片大小，直接make出指定长度，否则频繁的 grow 占用资源

google/pprof是一个性能可视化和分析工具，由Google的工程师开发。虽然自称不是Google官方的工具，但是项目挂在google的team下，而且还在Google其它项目中得到应用，是非常好的一个性能剖析工具。

go tool pprof 复制了一份google/pprof的代码， 封装了一个golang的工具，用来分析Go pprof包产生的剖析数据,也就是最终数据的处理和分析还是通过gogole/pprof来实现的。

这样，你至少就用两种方式来分析Go程序的 pprof数据：

go tool pprof : Go封装的pprof的工具
pprof: 原始的pprof工具
pprof读写一组profile.proto格式的数据，产生可视化的数据分析报告，数据是protocol buffer格式的数据，具体格式可以参考: profile.proto。因此，它可以分析可以任意产生这种格式的程序，不管程序是什么语言开发的。

它可以读取本地的剖析数据，或者通过http访问线上的实时的剖析数据，具体使用方法可以参考官方的说明。

今天8月份的时候，pprof发布了新的UI。 新的UI提供了顶部菜单(工具栏)， 可以提供各种不同的功能的切换，非常的方便。 同时，展示也提供了新的样式，更加的好看，SVG图中的展示也更加醒目。

现在, 另一个很重要的功能火焰图也被合并到主分支，这样，我们不用再利用第三方的工具go-torch等来查看火焰图。 这也意味着， 明年二月份发布的Go 1.10中我们可以直接通过go tool pprof查看火焰图了。

如果你不想等待到明年二月份，你可以下载最新的pprof来查看。

go get -u github.com/google/pprof


一、perf 命令
让我们从 perf 命令（performance 的缩写）讲起，它是 Linux 系统原生提供的性能分析工具，会返回 CPU 正在执行的函数名以及调用栈（stack）。

通常，它的执行频率是 99Hz（每秒99次），如果99次都返回同一个函数名，那就说明 CPU 这一秒钟都在执行同一个函数，可能存在性能问题。


$ sudo perf record -F 99 -p 13204 -g -- sleep 30
上面的代码中，perf record表示记录，-F 99表示每秒99次，-p 13204是进程号，即对哪个进程进行分析，-g表示记录调用栈，sleep 30则是持续30秒。

运行后会产生一个庞大的文本文件。如果一台服务器有16个 CPU，每秒抽样99次，持续30秒，就得到 47,520 个调用栈，长达几十万甚至上百万行。

为了便于阅读，perf record命令可以统计每个调用栈出现的百分比，然后从高到低排列。


$ sudo perf report -n --stdio


这个结果还是不易读，所以才有了火焰图。

二、火焰图的含义
火焰图是基于 perf 结果产生的 SVG 图片，用来展示 CPU 的调用栈。



y 轴表示调用栈，每一层都是一个函数。调用栈越深，火焰就越高，顶部就是正在执行的函数，下方都是它的父函数。

x 轴表示抽样数，如果一个函数在 x 轴占据的宽度越宽，就表示它被抽到的次数多，即执行的时间长。注意，x 轴不代表时间，而是所有的调用栈合并后，按字母顺序排列的。

火焰图就是看顶层的哪个函数占据的宽度最大。只要有"平顶"（plateaus），就表示该函数可能存在性能问题。

颜色没有特殊含义，因为火焰图表示的是 CPU 的繁忙程度，所以一般选择暖色调。

三、互动性
火焰图是 SVG 图片，可以与用户互动。

（1）鼠标悬浮

火焰的每一层都会标注函数名，鼠标悬浮时会显示完整的函数名、抽样抽中的次数、占据总抽样次数的百分比。下面是一个例子。


mysqld'JOIN::exec (272,959 samples, 78.34 percent)
（2）点击放大

在某一层点击，火焰图会水平放大，该层会占据所有宽度，显示详细信息。



左上角会同时显示"Reset Zoom"，点击该链接，图片就会恢复原样。

（3）搜索

按下 Ctrl + F 会显示一个搜索框，用户可以输入关键词或正则表达式，所有符合条件的函数名会高亮显示。

四、火焰图示例
下面是一个简化的火焰图例子。

首先，CPU 抽样得到了三个调用栈。


func_c 
func_b 
func_a 
start_thread 

func_d 
func_a 
start_thread 

func_d 
func_a 
start_thread
上面代码中，start_thread是启动线程，调用了func_a。后者又调用了func_b和func_d，而func_b又调用了func_c。

经过合并处理后，得到了下面的结果，即存在两个调用栈，第一个调用栈抽中1次，第二个抽中2次。


start_thread;func_a;func_b;func_c 1 
start_thread;func_a;func_d 2
有了这个调用栈统计，火焰图工具就能生成 SVG 图片。



上面图片中，最顶层的函数g()占用 CPU 时间最多。d()的宽度最大，但是它直接耗用 CPU 的部分很少。b()和c()没有直接消耗 CPU。因此，如果要调查性能问题，首先应该调查g()，其次是i()。

另外，从图中可知a()有两个分支b()和h()，这表明a()里面可能有一个条件语句，而b()分支消耗的 CPU 大大高于h()。

五、局限
两种情况下，无法画出火焰图，需要修正系统行为。

（1）调用栈不完整

当调用栈过深时，某些系统只返回前面的一部分（比如前10层）。

（2）函数名缺失

有些函数没有名字，编译器只用内存地址来表示（比如匿名函数）。

六、Node 应用的火焰图
Node 应用的火焰图就是对 Node 进程进行性能抽样，与其他应用的操作是一样的。


$ perf record -F 99 -p `pgrep -n node` -g -- sleep 30
详细的操作可以看这篇文章。

七、浏览器的火焰图
Chrome 浏览器可以生成页面脚本的火焰图，用来进行 CPU 分析。

打开开发者工具，切换到 Performance 面板。然后，点击"录制"按钮，开始记录数据。这时，可以在页面进行各种操作，然后停止"录制"。

这时，开发者工具会显示一个时间轴。它的下方就是火焰图。