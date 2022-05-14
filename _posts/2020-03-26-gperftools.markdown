---
title: gperftools 定位CPU热点函数
layout: post
category: golang
author: 夏泽民
---
1、下载及安装libunwind 

下载地址： http://download.savannah.gnu.org/releases/libunwind/libunwind-1.1.tar.gz 

安装：./configure CFLAGS=-U_FORTRIFY_SOURCE ; make -j8 ; make install

2、下载及安装gperftools 

git地址：https://github.com/gperftools/gperftools 

安装： ./configure ; make -j8 ; make install

或者：yum install gperftools gperftools-devel

3、使用方法 

参考：https://dirtysalt.github.io/gperftools.html#orgheadline4 

举例：

LD_PRELOAD="/usr/local/lib/libprofiler.so" CPUPROFILE=cpu_perf.prof CPUPROFILE_FREQUENCY=100 ./exe

运行完成后会在当前目录下生成cpu_perf.prof文件，结合pprof工具（gperftools/bin）可生成txt或pdf等，例如：

pprof --pdf ./test cpu_perf.prof > cpu_perf.pdf

pprof --txt ./test cpu_perf.prof > cpu_perf.txt
<!-- more -->
谷歌的工具集，可查看CPU采样结果。pprof (google-perftool)，用于来分析程序，必须保证程序能正常退出。使用步骤：

1.准备工具，先安装工具包

libunwind-1.1.tar.gz

gperftools-2.1.tar.gz

解压后 configure到系统默认路径即可，之后直接-lprofiler

 

2.再安装图形工具

sudo yum install "graphviz*"

sudo yum install ghostscript

GPerfTools="/home/zhifeng.czf/tools/gperftools-bin"



3. 自己程序增加编译链接选项：

-lprofiler

 

附：

GPerlToos=/home/zhifeng.czf/tools/gperftools-bin

CCFLAGS=-fno-omit-frame-pointer-g

ALL_BINS=test_gperf

all:$(ALL_BINS)

test_gperf:test_gperf.o

         g++ $(CCFLAGS) -o $@ $^-L$(GPerfTools)/lib -Wl, -Bdynamic -lprofiler

.cpp.o:

         g++ $(CCFLAGS) -c -I./-I$(GPerfTools)/include -fPIC -o $@ $<

clean:

         rm -rf $(ALL_BINS) *.o



4. 运行：

CPUPROFILE=./profile CPUPROFILESIGNAL=12 ./xlongsrv

 

此外，对于不能正常结束的程序，如果CPUPROFILESIGNAL的方式不成功，可以自己调用：

ProfilerStart(), ProfilerStop();

 

5.查看结果

pprof --pdf ./xlongsrv ./profile.0  >  graph.pdf

pprof --gif ./xlongsrv ./profile.0  >  graph.gif

/home/zhifeng.czf/tools/gperf/bin/pprof--callgrind ./XLongSrv.astar.perf ./bs.prof > graph.callgrind

/home/zhifeng.czf/tools/gperf/bin/pprof-text ./XLongSrv.astar.perf ./bs.prof


https://blog.golang.org/profiling-go-programs
https://lrita.github.io/2017/05/26/golang-memory-pprof/
