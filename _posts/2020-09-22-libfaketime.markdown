---
title: libfaketime
layout: post
category: linux
author: 夏泽民
---
https://github.com/wolfcw/libfaketime

user@host> date
Tue Nov 23 12:01:05 CEST 2016

user@host> LD_PRELOAD=/usr/local/lib/libfaketime.so.1 FAKETIME="-15d" date
Mon Nov  8 12:01:12 CEST 2016

https://unix.stackexchange.com/questions/30444/libfaketime-and-mac-os-ld-preload

https://blog.csdn.net/cunxiedian8614/article/details/105694624

http://inorz.net/2018/03/26/modifies-the-system-time-for-a-single-application/
<!-- more -->
我们有业务程序在测试时，需要把时间设置到指定时间去做一些测试，像这些程序，我们之前的做法都是分配一台低配置的虚拟机，然后可以独立修改系统时间。来解决这类需求
问题
经常有很多人需要在同时测试，这时候没办法，有多少人有需求就要建多少台虚拟机。很多时候部分虚拟机属于空闲状态的，比较浪费资源。
考虑过用Docker，Docker的时间是跟宿主机一起的，也可以用Docker通过调整时区来修改时间，但只能修改24小时内，不能跨多天调整。但也算是一种方案了
Faketime
放狗搜索时找到个神器 libfaketime
GitHub上的介绍：libfaketime modifies the system time for a single application
使用
安装，非常简单。

git clone https://github.com/wolfcw/libfaketime.git
cd libfaketime
make && make install
GitHub页面上有了非常详细的使用说明了。举个很简单的粟子。

指定动态链接库

# 正常时间
[root@inorz.net ~]# date
Mon Mar 26 21:01:46 CST 2018
        
[root@inorz.net ~]# LD_PRELOAD=/usr/local/lib/faketime/libfaketime.so.1 FAKETIME="-1d" date
Sun Mar 25 21:01:48 CST 2018
faketime 命令

[root@inorz.net ~]# date
Mon Mar 26 21:04:42 CST 2018
[root@inorz.net ~]# faketime '2018-03-27 21:04:52' date
Tue Mar 27 21:04:52 CST 2018


