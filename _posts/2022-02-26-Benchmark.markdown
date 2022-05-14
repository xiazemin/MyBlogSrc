---
title: Benchmark
layout: post
category: golang
author: 夏泽民
---
// 以BenchmarkXXX类似命名，并传入b *testing.B 参数
func BenchmarkLoopSum(b *testing.B) {
    for i := 0; i < b.N; i++ {
        total := 0
        for j := 0; j <= maxLoop; j++ {
            total += j
        }
    }
}

-bench grep
通过正则表达式过滤出需要进行benchtest的用例
-count n
跑n次benchmark，n默认为1
-benchmem
打印内存分配的信息
-benchtime=5s
自定义测试时间，默认为1s
<!-- more -->
5000表示测试次数，即test.B提供的N, ns/op表示每一个操作耗费多少时间(纳秒)。B/op表示每次调用需要分配16个字节。allocs/op表示每次调用有多少次分配
基准测试框架对一个测试用例的默认测试时间是 1 秒。开始测试时，当以 Benchmark 开头的基准测试用例函数返回时还不到 1 秒，那么 testing.B 中的 N 值将按 1、2、5、10、20、50……递增，同时以递增后的值重新调用基准测试用例函数。

https://www.cnblogs.com/linyihai/p/10977267.html

结合 pprof
go test -bench=. -benchmem -cpuprofile profile.out

go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out


覆盖测试
测试覆盖率就是测试运行到的被测试代码的代码数目。其中以语句的覆盖率最为简单和广泛，语句的覆盖率指的是在测试中至少被运行一次的代码占总代码数的比例。

测试整个包: go test -cover=true pkg_name

测试单个测试函数: go test -cover=true pkg_name -run TestSwap。

生成 HTML 报告

go test -cover=true pkg_name -coverprofile=out.out 将在当前目录生成覆盖率数据
配合 go tool cover -html=out.out 在浏览器中打开 HTML 报告。
或者使用 go tool cover -html=out.out -o=out.html 生成 HTML 文件。

https://www.cnblogs.com/bergus/articles/go-benchmark-xing-neng-ce-shi.html

BDD
https://github.com/smartystreets/goconvey

https://www.jianshu.com/p/da201abc1c6e
https://github.com/smartystreets/goconvey

B/op表示每次操作需要分配多少个字节，allocs/op表示每次操作有多少次alloc(内存分配)。

Go testing 库 testing.T 和 testing.B 简介
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

https://studygolang.com/articles/5494

第一列数字表示测试次数，一个测试用例的默认测试时间是 1 秒，所以表示1s内的执行次数，这个时间可以指定，使用**-benchtime=5s**；
ns/op表示每一次操作耗费多少时间(纳秒)；
因此从测试结果来看后者的耗时远小于前者。

我们知道，前者性能低下的原因可能与内存使用有关，因为golang中的string是字节切片，用“+”往里面拼接需要新建一个字节切片然后再往里面复制，可以使用benchmark来观察，使用命令：go test -bench=. -benchmem

https://blog.csdn.net/somanlee/article/details/107564151
