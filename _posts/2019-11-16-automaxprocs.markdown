---
title: automaxprocs
layout: post
category: golang
author: 夏泽民
---
GC停顿经常超过100ms
现象
有同事反馈说, 最近开始试用公司的k8s, 部署在docker里的go进程有问题, 接口耗时很长, 而且还有超时. 逻辑很简单, 只是调用了kv存储, kv存储一般响应时间<5ms, 而且量很少, 小于40qps, 该容器分配了0.5个核的配额, 日常运行CPU不足0.1个核. 
<!-- more -->
GC停顿经常超过100ms
现象
有同事反馈说, 最近开始试用公司的k8s, 部署在docker里的go进程有问题, 接口耗时很长, 而且还有超时. 逻辑很简单, 只是调用了kv存储, kv存储一般响应时间<5ms, 而且量很少, 小于40qps, 该容器分配了0.5个核的配额, 日常运行CPU不足0.1个核. 



复现
我找了个容器, 踢掉访问流量. 用ab 50并发构造些请求看看. 网络来回延时60ms, 但是平均处理耗时200多ms, 99%耗时到了679ms.



用ab处理时, 看了下CPU和内存信息, 都没啥问题. docker分配的是0.5个核. 这里也没有用到那么多.



看了下监控, GC STW(stop the world)超过10ms, 50-100ms的都很多, 还有不少超过100ms的. Go不是声称在1.8后GC停顿基本是小于1ms的吗?





gc信息及trace
看看该进程的runtime信息, 发现内存很少，gc-pause很大，GOMAXPROCS为76，是机器的核数。 





export GODEBUG=gctrace=1, 重启进程看看. 可以看出gc停顿的确很严重.

gc 
111
 
@
97.209s
 
1
%:
 
82
+
7.6
+
0.036
 ms clock
,
 
6297
+
0.66
/
6.0
/
0
+
2.7
 ms cpu
,
 
9
->
12
->
6
 MB
,
 
11
 MB goal
,
 
76
 P

gc 
112
 
@
97.798s
 
1
%:
 
0.040
+
93
+
0.14
 ms clock
,
 
3.0
+
0.55
/
7.1
/
0
+
10
 ms cpu
,
 
10
->
11
->
5
 MB
,
 
12
 MB goal
,
 
76
 P

gc 
113
 
@
99.294s
 
1
%:
 
0.041
+
298
+
100
 ms clock
,
 
3.1
+
0.34
/
181
/
0
+
7605
 ms cpu
,
 
10
->
13
->
6
 MB
,
 
11
 MB goal
,
 
76
 P

gc 
114
 
@
100.892s
 
1
%:
 
99
+
200
+
99
 ms clock
,
 
7597
+
0
/
5.6
/
0
+
7553
 ms cpu
,
 
11
->
13
->
6
 MB
,
 
13
 MB goal
,
 
76
 P

在一台有go sdk的服务器上对服务跑一下trace, 再把trace文件下载到本地看看

curl 
-
o trace
.
out
 
'http://ip:port/debug/pprof/trace?seconds=20'

sz 
./
trace
.
out

下图可见有一个GC的wall time为172ms,而本次gc的两个stw阶段,sweep termination和mark termination都在80多ms的样子, 几乎占满了整个GC时间, 这当然很不科学. 





原因及解决方法
原因
这个服务是运行在容器里的, 容器和母机共享一个内核, 容器里的进程看到的CPU核数也是母机CPU核数, 对于Go应用来说, 会默认设置P(为GOMAXPROCS)的个数为CPU核数. 我们从前面的图也可以看到, GOMAXPROCS为76, 每个使用中的P都有一个m与其绑定, 所以线程数也不少, 上图中的为171.然而分配给该容器的CPU配额其实不多, 仅为0.5个核, 而线程数又不少.

猜测: 对于linux的cfs(完全公平调度器)来说, 当前容器内所有的线程(轻量级进程)都在一个调度组内. 为了保证效率, 对于每个被运行的task, 除非因为阻塞等原因主动切换, 那么至少保证其运行/proc/sys/kernel/schedmingranularity_ns的时间, 可以看到为4ms.

容器中Go进程没有正确的设置GOMAXPROCS的个数, 导致可运行的线程过多, 可能出现调度延迟的问题. 正好出现进入gc发起stw的线程把其他线程停止后, 其被调度器切换出去, 很久没有调度该线程, 实质上造成了stw时间变得很长(正常情况0.1ms的处理过程因为调度延迟变成了100ms级别).

解决方法
解决的方法, 因为可运行的P太多, 导致占用了发起stw的线程的虚拟运行时间, 且CPU配额也不多. 那么我们需要做的是使得P与CPU配额进行匹配. 我们可以选择:

增加容器的CPU配额.

容器层让容器内的进程看到CPU核数为配额数

根据配额设置正确的GOMAXPROCS

第1个方法: 没太大效果, 把配额从0.5变成1, 没本质的区别(尝试后, 问题依旧). 

第2点方法: 因为我对k8s不是很熟, 待我调研后再来补充. 

第3个方法: 设置GOMAXPROCS最简单的方法就是启动脚本添加环境变量 

GOMAXPROCS=2 ./svr_bin 这种是有效的, 但也有不足, 如果部署配额大一点的容器, 那么脚本没法跟着变.

uber的库automaxprocs
uber有一个库, go.uber.org/automaxprocs, 容器中go进程启动时, 会正确设置GOMAXPROCS. 修改了代码模板. 我们在go.mod中引用该库

go
.
uber
.
org
/
automaxprocs v1
.
2.0

并在main.go中import

import
 
(

    _ 
"go.uber.org/automaxprocs"

)

效果
automaxprocs库的提示
使用automaxprocs库, 会有如下日志:

对于虚拟机或者实体机

8核的情况下: 2019/11/07 17:29:47 maxprocs: Leaving GOMAXPROCS=8: CPU quota undefined

对于设置了超过1核quota的容器

2019/11/08 19:30:50 maxprocs: Updating GOMAXPROCS=8: determined from CPU quota

对于设置小于1核quota的容器

2019/11/08 19:19:30 maxprocs: Updating GOMAXPROCS=1: using minimum allowed GOMAXPROCS

如果docker中没有设置quota

2019/11/07 19:38:34 maxprocs: Leaving GOMAXPROCS=79: CPU quota undefined

此时建议在启动脚本中显式设置 GOMAXPROCS

请求响应时间
设置完后, 再用ab请求看看，网络往返时间为60ms, 99%请求在200ms以下了, 之前是600ms. 同等CPU消耗下, qps从差不多提升了一倍. 





runtime及gc trace信息
因为分配的是0.5核, GOMAXPROC识别为1. gc-pause也很低了, 几十us的样子. 同时也可以看到线程数从170多降到了11. 





gc 
97
 
@
54.102s
 
1
%:
 
0.017
+
3.3
+
0.003
 ms clock
,
 
0.017
+
0.51
/
0.80
/
0.75
+
0.003
 ms cpu
,
 
9
->
9
->
4
 MB
,
 
10
 MB goal
,
 
1
 P

gc 
98
 
@
54.294s
 
1
%:
 
0.020
+
5.9
+
0.003
 ms clock
,
 
0.020
+
0.51
/
1.6
/
0
+
0.003
 ms cpu
,
 
8
->
9
->
4
 MB
,
 
9
 MB goal
,
 
1
 P

gc 
99
 
@
54.406s
 
1
%:
 
0.011
+
4.4
+
0.003
 ms clock
,
 
0.011
+
0.62
/
1.2
/
0.17
+
0.003
 ms cpu
,
 
9
->
9
->
4
 MB
,
 
10
 MB goal
,
 
1
 P

gc 
100
 
@
54.597s
 
1
%:
 
0.009
+
5.6
+
0.002
 ms clock
,
 
0.009
+
0.69
/
1.4
/
0
+
0.002
 ms cpu
,
 
9
->
9
->
5
 MB
,
 
10
 MB goal
,
 
1
 P

gc 
101
 
@
54.715s
 
1
%:
 
0.026
+
2.7
+
0.004
 ms clock
,
 
0.026
+
0.42
/
0.35
/
1.4
+
0.004
 ms cpu
,
 
9
->
9
->
4
 MB
,
 
10
 MB goal
,
 
1
 P

上下文切换
以下为并发50, 一共处理8000个请求的perf stat结果对比. 默认CPU核数76个P, 上下文切换13万多次, pidstat查看system cpu消耗9%个核.  而按照quota数设置P的数量后, 上下文切换仅为2万多次, cpu消耗3%个核. 







automaxprocs原理解析
这个库如何根据quota设置GOMAXPROCS呢, 代码有点绕, 看完后, 其实原理不复杂. docker使用cgroup来限制容器CPU使用, 使用该容器配置的cpu.cfsquotaus/cpu.cfsperiodus即可获得CPU配额. 所以关键是找到容器的这两个值.

获取cgroup挂载信息
cat /proc/self/mountinfo

....

1070
 
1060
 
0
:
17
 
/ /
sys
/
fs
/
cgroup ro
,
nosuid
,
nodev
,
noexec 
-
 tmpfs tmpfs ro
,
mode
=
755

1074
 
1070
 
0
:
21
 
/ /
sys
/
fs
/
cgroup
/
memory rw
,
nosuid
,
nodev
,
noexec
,
relatime 
-
 cgroup cgroup rw
,
memory

1075
 
1070
 
0
:
22
 
/ /
sys
/
fs
/
cgroup
/
devices rw
,
nosuid
,
nodev
,
noexec
,
relatime 
-
 cgroup cgroup rw
,
devices

1076
 
1070
 
0
:
23
 
/ /
sys
/
fs
/
cgroup
/
blkio rw
,
nosuid
,
nodev
,
noexec
,
relatime 
-
 cgroup cgroup rw
,
blkio

1077
 
1070
 
0
:
24
 
/ /
sys
/
fs
/
cgroup
/
hugetlb rw
,
nosuid
,
nodev
,
noexec
,
relatime 
-
 cgroup cgroup rw
,
hugetlb

1078
 
1070
 
0
:
25
 
/ /
sys
/
fs
/
cgroup
/
cpu
,
cpuacct rw
,
nosuid
,
nodev
,
noexec
,
relatime 
-
 cgroup cgroup rw
,
cpuacct
,
cpu

1079
 
1070
 
0
:
26
 
/ /
sys
/
fs
/
cgroup
/
cpuset rw
,
nosuid
,
nodev
,
noexec
,
relatime 
-
 cgroup cgroup rw
,
cpuset

1081
 
1070
 
0
:
27
 
/ /
sys
/
fs
/
cgroup
/
net_cls rw
,
nosuid
,
nodev
,
noexec
,
relatime 
-
 cgroup cgroup rw
,
net_cls

....

cpuacct,cpu在/sys/fs/cgroup/cpu,cpuacct这个目录下.

获取该容器cgroup子目录
cat /proc/self/cgroup

10
:
net_cls
:
/kubepods/
burstable
/
pod62f81b5d
-
xxxx
/
xxxx92521d65bff8

9
:
cpuset
:
/kubepods/
burstable
/
pod62f81b5d
-
xxxx
/
xxxx92521d65bff8

8
:
cpuacct
,
cpu
:
/kubepods/
burstable
/
pod62f81b5d
-
xxxx
/
xxxx92521d65bff8

7
:
hugetlb
:
/kubepods/
burstable
/
pod62f81b5d
-
5ce0
-
xxxx
/
xxxx92521d65bff8

6
:
blkio
:
/kubepods/
burstable
/
pod62f81b5d
-
5ce0
-
xxxx
/
xxxx92521d65bff8

5
:
devices
:
/kubepods/
burstable
/
pod62f81b5d
-
5ce0
-
xxxx
/
xxxx92521d65bff8

4
:
memory
:
/kubepods/
burstable
/
pod62f81b5d
-
5ce0
-
xxxx
/
xxxx92521d65bff8

....

该容器的cpuacct,cpu具体在/kubepods/burstable/pod62f81b5d-xxxx/xxxx92521d65bff8子目录下

计算quota
cat 
/
sys
/
fs
/
cgroup
/
cpu
,
cpuacct
/
kubepods
/
burstable
/
pod62f81b5d
-
5ce0
-
xxxx
/
xxxx92521d65bff8
/
cpu
.
cfs_quota_us

50000



cat 
/
sys
/
fs
/
cgroup
/
cpu
,
cpuacct
/
kubepods
/
burstable
/
pod62f81b5d
-
5ce0
-
xxxx
/
xxxx92521d65bff8
/
cpu
.
cfs_period_us

100000

两者相除得到0.5, 小于1的话，GOMAXPROCS设置为1，大于1则设置为计算出来的数。

核心函数
automaxprocs库中核心函数如下所示, 其中cg为解析出来的cgroup的所有配置路径. 分别读取cpu.cfs_quota_us和cpu.cfs_period_us, 然后计算. 





官方issue
谷歌搜了下, 也有人提过这个问题 

runtime: long GC STW pauses (≥80ms) #19378 https://github.com/golang/go/issues/19378

总结
容器中进程看到的核数为母机CPU核数，一般这个值比较大>32, 导致go进程把P设置成较大的数，开启了很多P及线程

一般容器的quota都不大，0.5-4，linux调度器以该容器为一个组，里面的线程的调度是公平，且每个可运行的线程会保证一定的运行时间，因为线程多, 配额小, 虽然请求量很小, 但上下文切换多, 也可能导致发起stw的线程的调度延迟, 引起stw时间升到100ms的级别，极大的影响了请求

通过使用automaxprocs库, 可根据分配给容器的cpu quota, 正确设置GOMAXPROCS以及P的数量, 减少线程数，使得GC停顿稳定在<1ms了. 且同等CPU消耗情况下, QPS可增大一倍，平均响应时间由200ms减少到100ms. 线程上下文切换减少为原来的1/6

同时还简单分析了该库的原理. 找到容器的cgroup目录, 计算cpuacct,cpu下cpu.cfs_quota_us/cpu.cfs_period_us, 即为分配的cpu核数.

当然如果容器中进程看到CPU核数就是分配的配额的话, 也可以解决这个问题. 这方面我就不太了解了.