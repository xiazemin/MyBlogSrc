---
title: golang热编译工具
layout: post
category: golang
author: 夏泽民
---
<!-- more -->
fswatch是一个工具, 通过检测文件的变化,并触发指定的命令
bee只适合Go语言; 而fswatch可以适用各种语言, 甚至是文件的远程同步
bee可以通过配置指定监控的文件夹; fswatch可以配置监控的文件夹并指定监控的深度(0代表当前目录)
bee可以指定监控文件的后缀; fswatch支持使用正则表达式, 来过滤监控到的文件.
bee.json需要重其他地方拷贝;但是.fsw.yml可以自动生成出来.
fswatch支持group kill. 这样可以确保fswatch停止后,不会有垃圾进程的存在.
fswatch会根据程序的运行时间自动判断, 是否为服务端程序, 并适当的修改重启策略.
bee功能庞大; fswatch的代码精简.

 fswatch init 直接运行这个命令就可以. 然后你会在目录下面发下一个.fswatch.json文件
 
 https://github.com/codeskyblue/fswatch
 
 Dogo是一个自动编译go项目的工具, 他的工作原理很简单, 当dogo监控的目录中go源代码文件发生修改,删除,增加的时候, 自动调用编译命令重新编译项目并重新启动新编译的可执行文件.
项目主页: https://github.com/liudng/dogo

运行bee run
https://github.com/beego/bee