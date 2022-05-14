---
title: benchmark
layout: post
category: golang
author: 夏泽民
---
网关服务本身没有业务逻辑处理，仅作为统一入口进行请求转发，因此我们主要关注下列指标

吞吐量：每秒钟可以处理的请求数
响应时间：从客户端发出请求，到收到回包的总耗时
定位瓶颈
一般后台服务的瓶颈主要为 CPU，内存，IO 操作中的一个或多个。若这三者的负载都不高，但系统吞吐量低，基本就是代码逻辑出问题了。

在代码正常运行的情况下，我们要针对某个方面的高负载进行优化，才能提高系统的性能。golang 可通过 benchmark 加 pprof 来定位具体的性能瓶颈。

benchmark 简介
go test -v gate_test.go -run=none -bench=. -benchtime=3s -cpuprofile cpu.prof -memprofile mem.prof
-run 知道单次测试，一般用于代码逻辑验证
-bench=. 执行所有 Benchmark，也可以通过用例函数名来指定部分测试用例
-benchtime 指定测试执行时长
-cpuprofile 输出 cpu 的 pprof 信息文件
-memprofile 输出 heap 的 pprof 信息文件。
-blockprofile 阻塞分析，记录 goroutine 阻塞等待同步（包括定时器通道）的位置
-mutexprofile 互斥锁分析，报告互斥锁的竞争情况
benchmark 测试用例常用函数

b.ReportAllocs() 输出单次循环使用的内存数量和对象 allocs 信息
b.RunParallel() 使用协程并发测试
b.SetBytes(n int64) 设置单次循环使用的内存数量
<!-- more -->

http://team.jiunile.com/blog/2020/05/go-performance.html

https://mp.weixin.qq.com/s?__biz=MzAxMTA4Njc0OQ==&mid=2651439020&idx=1&sn=c2094f4dccb53385dc207958e7f42f9e&chksm=80bb615eb7cce8481eb7a8f09d4a13e2974b3785c241dd31245647cd7540dde414d64f2b3719&scene=21#wechat_redirect