---
title: benchmark
layout: post
category: golang
author: 夏泽民
---
testing包提供了对Go包的自动测试支持。 这是和go test 命令相呼应的功能， go test 命令会自动执行所以符合格式
func TestXXX(t *testing.T) 

当带着 -bench=“.” （ 参数必须有！）来执行*go test命令的时候性能测试程序就会被顺序执行。符合下面格式的函数被认为是一个性能测试程序，
func BenchmarkXxx(b *testing.B)

执行 go test -bench=”.” 后 结果 ：
// 表示测试全部通过
>PASS                       
// Benchmark 名字 - CPU            循环次数             平均每次执行时间 
BenchmarkLoops-2                  100000             20628 ns/op     
BenchmarkLoopsParallel-2          100000             10412 ns/op   
//      哪个目录下执行go test         累计耗时
ok      swap/lib                   2.279s
<!-- more -->
testing.T
判定失败接口
Fail 失败继续
FailNow 失败终止
打印信息接口
Log 数据流 （cout　类似）
Logf format (printf 类似）
SkipNow 跳过当前测试
Skiped 检测是否跳过
综合接口产生：

Error / Errorf 报告出错继续 [ Log / Logf + Fail ]
Fatel / Fatelf 报告出错终止 [ Log / Logf + FailNow ]
Skip / Skipf 报告并跳过 [ Log / Logf + SkipNow ]
testing.B
首先 ， testing.B 拥有testing.T 的全部接口。
SetBytes( i uint64) 统计内存消耗， 如果你需要的话。
SetParallelism(p int) 制定并行数目。
StartTimer / StopTimer / ResertTimer 操作计时器
testing.PB
Next() 接口 。 判断是否继续循环

b.N是如何自动调整的？
内存统计是如何实现的？
SetBytes()其使用场景是什么？

数据结构
源码包src/testing/benchmark.go:B定义了性能测试的数据结构，我们提取其比较重要的一些成员进行分析：

type B struct {
	common                         // 与testing.T共享的testing.common，负责记录日志、状态等
	importPath       string // import path of the package containing the benchmark
	context          *benchContext
	N                int            // 目标代码执行次数，不需要用户了解具体值，会自动调整
	previousN        int           // number of iterations in the previous run
	previousDuration time.Duration // total duration of the previous run
	benchFunc        func(b *B)   // 性能测试函数
	benchTime        time.Duration // 性能测试函数最少执行的时间，默认为1s，可以通过参数'-benchtime 10s'指定
	bytes            int64         // 每次迭代处理的字节数
	missingBytes     bool // one of the subbenchmarks does not have bytes set.
	timerOn          bool // 是否已开始计时
	showAllocResult  bool
	result           BenchmarkResult // 测试结果
	parallelism      int // RunParallel creates parallelism*GOMAXPROCS goroutines
	// The initial states of memStats.Mallocs and memStats.TotalAlloc.
	startAllocs uint64  // 计时开始时堆中分配的对象总数
	startBytes  uint64  // 计时开始时时堆中分配的字节总数
	// The net total of this test after being run.
	netAllocs uint64 // 计时结束时，堆中增加的对象总数
	netBytes  uint64 // 计时结束时，堆中增加的字节总数
}
其主要成员如下：

common： 与testing.T共享的testing.common，管理着日志、状态等；
N：每个测试中用户代码执行次数
benchFunc：测试函数
benchTime：性能测试最少执行时间，默认为1s，可以通过能数-benchtime 2s指定
bytes：每次迭代处理的字节数
timerOn：计时启动标志，默认为false，启动计时为true
startAllocs：测试启动时记录堆中分配的对象数
startBytes：测试启动时记录堆中分配的字节数
netAllocs：测试结束后记录堆中新增加的对象数，公式：结束时堆中分配的对象数-
netBytes：测试对事后记录堆中新增加的字节数
关键函数
启动计时：B.StartTimer()
StartTimer()负责启动计时并初始化内存相关计数，测试执行时会自动调用，一般不需要用户启动。

func (b *B) StartTimer() {
	if !b.timerOn {
		runtime.ReadMemStats(&memStats)     // 读取当前堆内存分配信息
		b.startAllocs = memStats.Mallocs    // 记录当前堆内存分配的对象数
		b.startBytes = memStats.TotalAlloc  // 记录当前堆内存分配的字节数
		b.start = time.Now()                // 记录测试启动时间
		b.timerOn = true                   // 标记计时标志
	}
}
StartTimer()负责启动计时，并记录当前内存分配情况，不管是否有“-benchmem”参数，内存都会被统计，参数只决定是否要在结果中输出。

停止计时：B.StopTimer()
StopTimer()负责停止计时，并累加相应的统计值。

func (b *B) StopTimer() {
	if b.timerOn {
		b.duration += time.Since(b.start)                   // 累加测试耗时
		runtime.ReadMemStats(&memStats)                     // 读取当前堆内存分配信息
		b.netAllocs += memStats.Mallocs - b.startAllocs     // 累加堆内存分配的对象数
		b.netBytes += memStats.TotalAlloc - b.startBytes    // 累加堆内存分配的字节数
		b.timerOn = false                                  // 标记计时标志
	}
}
需要注意的是，StopTimer()并不一定是测试结束，一个测试中有可能有多个统计阶段，所以其统计值是累加的。

重置计时：B.ResetTimer()
ResetTimer()用于重置计时器，相应的也会把其他统计值也重置。

func (b *B) ResetTimer() {
	if b.timerOn {
		runtime.ReadMemStats(&memStats)     // 读取当前堆内存分配信息
		b.startAllocs = memStats.Mallocs    // 记录当前堆内存分配的对象数
		b.startBytes = memStats.TotalAlloc  // 记录当前堆内存分配的字节数
		b.start = time.Now()                // 记录测试启动时间
	}
	b.duration = 0                          // 清空耗时
	b.netAllocs = 0                         // 清空内存分配对象数
	b.netBytes = 0                          // 清空内存分配字节数
}
ResetTimer()比较常用，典型使用场景是一个测试中，初始化部分耗时较长，初始化后再开始计时。

设置处理字节数：B.SetBytes(n int64)
// SetBytes records the number of bytes processed in a single operation.
// If this is called, the benchmark will report ns/op and MB/s.
func (b *B) SetBytes(n int64) {
	b.bytes = n
}
这是一个比较含糊的函数，通过其函数说明很难明白其作用。

其实它是用来设置单次迭代处理的字节数，一旦设置了这个字节数，那么输出报告中将会呈现“xxx MB/s”的信息，用来表示待测函数处理字节的性能。待测函数每次处理多少字节数只有用户清楚，所以需要用户设置。

举个例子，待测函数每次执行处理1M数据，如果我们想看待测函数处理数据的性能，那么我们在测试中设置SetByte(1024 *1024)，假如待测函数需要执行1s的话，那么结果中将会出现 “1 MB/s”（约等于）的信息。示例代码如下所示：

func BenchmarkSetBytes(b *testing.B) {
    b.SetBytes(1024 * 1024)
    for i := 0; i < b.N; i++ {
        time.Sleep(1 * time.Second) // 模拟待测函数
    }
}

打印结果：

E:\OpenSource\GitHub\RainbowMango\GoExpertProgrammingSourceCode\GoExpert\src\gotest>go test -bench SetBytes benchmark_test.go
BenchmarkSetBytes-4            1        1010392800 ns/op           1.04 MB/s
PASS
ok      command-line-arguments  1.412s
 
可以看到测试执行了一次，花费时间约1S，数据处理能力约为1MB/s。

报告内存信息：
func (b *B) ReportAllocs() {
	b.showAllocResult = true
}
ReportAllocs() 用于设置是否打印内存统计信息，与命令行参数“-benchmem”一致，但本方法只作用于单个测试函数。

性能测试是如何启动的
性能测试要经过多次迭代，每次迭代可能会有不同的b.N值，每次迭代执行测试函数一次，跟据此次迭代的测试结果来分析要不要继续下一次迭代。

我们先看一下每次迭代时所用到的方法，runN():

func (b *B) runN(n int) {
	b.N = n                       // 指定B.N
	b.ResetTimer()                // 清空统计数据
	b.StartTimer()                // 开始计时
	b.benchFunc(b)                // 执行测试
	b.StopTimer()                 // 停止计时
}
该方法指定b.N的值，执行一次测试函数。

与T.Run()类似，B.Run()也用于启动一个子测试，实际上用户编写的任何一个测试都是使用Run()方法启动的，我们看下B.Run()的伪代码：

func (b *B) Run(name string, f func(b *B)) bool {
	sub := &B{                          // 新建子测试数据结构
		common: common{
			signal:  make(chan bool),
			name:    name,
			parent:  &b.common,
		},
		benchFunc:  f,
	}
	if sub.run1() { // 先执行一次子测试，如果子测试不出错且子测试没有子测试的话继续执行sub.run()
	sub.run()       // run()里决定要执行多少次runN()
	}
	b.add(sub.result) // 累加统计结果到父测试中
	return !sub.failed
}
所有的测试都是先使用run1()方法执行一次测试,run1()方法中实际上调用了runN(1)，执行一次后再决定要不要继续迭代。

测试结果实际上以最后一次迭代的数据为准，当然，最后一次迭代往往意味着b.N更大，测试准确性相对更高。

B.N是如何调整的？
B.launch()方法里最终决定B.N的值。我们看下伪代码：

func (b *B) launch() { // 此方法自动测算执行次数，但调用前必须调用run1以便自动计算次数
	d := b.benchTime
	for n := 1; !b.failed && b.duration < d && n < 1e9; { // 最少执行b.benchTime（默认为1s）时间，最多执行1e9次
		last := n
		n = int(d.Nanoseconds()) // 预测接下来要执行多少次，b.benchTime/每个操作耗时
		if nsop := b.nsPerOp(); nsop != 0 {
			n /= int(nsop)
		}
		n = max(min(n+n/5, 100*last), last+1) // 避免增长较快，先增长20%，至少增长1次
		n = roundUp(n) // 下次迭代次数向上取整到10的指数，方便阅读
		b.runN(n)
	}
}
不考虑程序出错，而且用户没有主动停止测试的场景下，每个性能测试至少要执行b.benchTime长的秒数，默认为1s。先执行一遍的意义在于看用户代码执行一次要花费多长时间，如果时间较短，那么b.N值要足够大才可以测得更精确，如果时间较长，b.N值相应的会减少，否则会影响测试效率。

最终的b.N会被定格在某个10的指数级，是为了方便阅读测试报告。

内存是如何统计的？
我们知道在测试开始时，会把当前内存值记入到b.startAllocs和b.startBytes中，测试结束时，会用最终内存值与开始时的内存值相减，得到净增加的内存值，并记入到b.netAllocs和b.netBytes中。

每个测试结束，会把结果保存到BenchmarkResult对象里，该对象里保存了输出报告所必需的统计信息：

type BenchmarkResult struct {
	N         int           // 用户代码执行的次数
	T         time.Duration // 测试耗时
	Bytes     int64         // 用户代码每次处理的字节数，SetBytes()设置的值
	MemAllocs uint64        // 内存对象净增加值
	MemBytes  uint64        // 内存字节净增加值
}
其中MemAllocs和MemBytes分别对应b.netAllocs和b.netBytes。

那么最终统计时只需要把净增加值除以b.N即可得到每次新增多少内存了。

每个操作内存对象新增值：

func (r BenchmarkResult) AllocsPerOp() int64 {
	return int64(r.MemAllocs) / int64(r.N)
}
每个操作内存字节数新增值：

func (r BenchmarkResult) AllocedBytesPerOp() int64 {
	if r.N <= 0 {
		return 0
	}
	return int64(r.MemBytes) / int64(r.N)
}

https://blog.csdn.net/weixin_33834137/article/details/91699646
https://github.com/RainbowMango/GoExpertProgramming
https://studygolang.com/articles/5494

