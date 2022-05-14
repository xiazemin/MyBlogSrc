---
title: php graphviz 可视化工具
layout: post
category: php
author: 夏泽民
---
PEAR（扩展与应用库：PHP Extension and Application Repository），是PHP官方开源类库，可以使用pear list列出所有已经安装的包。通过pear install可以安装需要的包。

PECL，是PHP的扩展库，可以通过 PEAR 的 Package Manager 的管理方式来下载和安装扩展代码。

以安装yaconf为例：

$ ./pecl install yaconf

phpdbg的其他功能可以通过使用phpdbg --help查看。

phpize 命令用来动态安装扩展，如果在安装PHP时没有安装某个扩展，可以通过这个命令随时安装。

1.2.2 使用GDB调试PHP7
GDB是一个由GNU开源组织发布的、UNIX/LINUX操作系统下的、基于命令行的、功能强大的程序调试工具。当我们的程序发生coredump，通过GDB可以从 core文件中复现场景，定位问题。

这里，我们演示下如何通过GDB来调试PHP 程序。首先我们编写一段简单的代码test.php：

<?php

$a = '1';

echo $a;

下面我们开始进行gdb调试，运行gdb php:

$ gdb php

 (gdb)

使用b命令在main函数入口增加断点：

(gdb) b main

Breakpoint 1 at 0x797df0: file /home/vagrant/php7/php-7.1.0/sapi/cli/php_cli.c, line 1181.

使用r命令运行test.php

(gdb) r test.php

对于在php-fpm下运行的 PHP 程序如何调试呢？
output/conf/php-fpm.conf

// 添加以下配置项

[www.local]

pm=static

pm.max_children=1

pm.start_servers=1

pm.min_spare_servers=1

pm.max_spare_servers=1

$ gdb php

(gdb) attach 4458
<!-- more -->
vld 扩展
PHP代码的执行实际是在执行代码解析后的各种 opcode。通过 vld 扩展可以很方便地看到执行过程中的 opcode。扩展可以从 https://github.com/derickr/vld 下载安装，下边是安装示例：

$ git clone https://github.com/derickr/vld.git

命令行执行：

$ php -dvld.active=1 vld.php

vld 扩展有下边几个参数。

1）vld.active：是否在执行 PHP 的同时激活vld：1激活，0不激活（默认不激活）；

2）vld.execute：是否输出程序的执行结果：1输出，0不输出（默认输出）；

3）vld.verbosity：显示更详细的opcode信息，开启后可以显示每个opcode的操作数的类型等信息；

dot 是一种描述图形的语言，可以由Graphviz工具包来绘制 dot 描述的图形。vld 扩展可以直接通过命令来生成dot 脚本

$ php -dvld.active=1 -dvld.save_paths=1 vld.php

$ dot -Tpng /tmp/paths.dot -o paths.png


windows下的source insight、Mac下的Understand以及Linux下的Vim+Ctags

在某个函数上点击右键，选择Graphical Views→Declaration命令，可以看到该函数的调用关系

IntelliJ IDEA 有 Call Hierarchy的功能

首先在偏好设置中搜索关键字 Call Hierarchy

找到快捷键 ctrl+alt+h

选中一个函数array_slice，按快捷键，在右边生成了该函数所有的调用关系

与查看用例find usage不同的是，Call Hierarchy功能会递归的寻找用例的用例，直到找到没有入口函数为止


https://www.jianshu.com/p/8cc511b83710

https://www.zhihu.com/question/34495043/answer/244410441

静态分析又有两种方法，一是分析源码，二是分析编译后的目标文件。

分析源码获得的调用图的质量取决于分析工具对编程语言的理解程度，比如能不能找出正确的C++重载函数。Doxygen是源码文档化工具，也能绘制调用图，它似乎是自己分析源码获得函数调用关系的。GNU cflow也是类似的工具，不过它似乎偏重分析流程图（flowchart）。

对编程语言的理解程度最好的当然是编译器了，所以有人想出给编译器打补丁，让它在编译时顺便记录函数调用关系。CodeViz<!--StartFragment -->（其灵感来自Martin Devera (Devik) 的工具）就属于此类，它（1.0.9版）给GCC 3.4.1打了个补丁。另外一个工具egypt的思路更巧妙，不用大动干戈地给编译器打补丁，而是让编译器自己dump出调用关系，然后分析分析，交给Graphviz去绘图。不过也有人另起炉灶，自己写个C语言编译器（ncc），专门分析调用图，勇气可嘉。不如要是对C++语言也这么干，成本不免太高了。分析C++的调用图，还是借助编译器比较实在。

https://blog.csdn.net/happmaoo/article/details/83201305
http://students.ceid.upatras.gr/~sxanth/ncc/
http://www.gson.org/egypt/egypt.html

http://luxik.cdi.cz/~devik/mm.htm
http://www.gnu.org/software/cflow/
https://blog.csdn.net/Solstice/article/details/488865
