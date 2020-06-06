---
title: stacktrace
layout: post
category: php
author: 夏泽民
---
https://github.com/rbspy/rbspy
https://github.com/oraoto/php-stacktrace
php-stacktrace: PHP进程外查看函数调用堆栈

生产环境多多少少会遇到CPU占用很高或者卡住的PHP进程，这时怎样才能知道这个进程在干啥呢？

一个方法是strace跟踪系统调用和参数，这样能大概知道PHP进程在干啥。要看到具体的PHP函数就需要用PHP扩展（xdebug、xhprof）或者用GDB调试，高级点还可以用DTrace。

上周发现了ruby-stacktrace，它直接读取ruby进程的内存来获取堆栈信息，不用GDB和扩展，所以性能很好，于是我也照着写了一个php-stacktrace，算是勉强能用的玩具。
<!-- more -->
使用
使用比较简单，下载解压即可：

$ ./php-stacktrace --help
php-stacktrace 0.1
Sampling profiler for PHP programs

USAGE:
    php-stacktrace <COMMAND> <DEBUGINFO> <PID>

FLAGS:
    -h, --help       Prints help information
    -V, --version    Prints version information

ARGS:
    <COMMAND>      trace or top or oneshot
    <DEBUGINFO>    Path to php debuginfo
    <PID>          PID of the PHP process you want to profile
三个参数都是必填的。

COMMAND可以是trace、top、oneshot。oneshot只查看一次就退出，trace和top会一直跟踪，trace的输出可以用来生成火焰图，top统计函数耗时。

DEBUGINFO是调试信息文件的路径，Linux通常要独立安装debuginfo包，因为不会从elf里解析路径，所以要通过这个参数指定，通常的路径是/usr/lib/debug/.dwz/php....（在隐藏目录里，是个小坑）。

PID就是要跟踪的PHP进程ID。

顺带一提，只支持非线程安全的PHP 7.1。

原理
众所周知，Zend VM是用C写的，而各种PHP函数调用的信息都会用C语言的struct/union来表示，所以只要两步就能拿到堆栈信息：

读取PHP进程的内存
在内存里找到函数调用堆栈信息
第一步可以通过ptrace或process_vm_readv实现。ptrace就是调试器所用的方法，它可以暂停PHP进程然后读取内存。process_vm_readv可以不暂停进程，性能可能更好，但是不可靠，因为PHP还在执行，堆栈信息不断变化，很容易读到错误的内存。

第二步就需要DWARF调试信息了，调试信息里记录了结构体大小、字段偏移信息，通过这些信息我们就可以准确地去读内存然后做解析。

原理还是很简单的。

下一步
复制vm_stack，尽量在一次process_vm_readv拿到主要的堆栈信息
增加作用域信息，现在只有函数名

http://www.suoniao.com/article/3999
http://www.mamicode.com/info-detail-2175996.html

