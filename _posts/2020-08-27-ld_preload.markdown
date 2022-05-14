---
title: ld_preload
layout: post
category: linux
author: 夏泽民
---
https://www.cnblogs.com/net66/p/5609026.html
https://www.zcfy.cc/article/dynamic-linker-tricks-using-ld-preload-to-cheat-inject-features-and-investigate-programs
<!-- more -->
https://stackoverflow.com/questions/426230/what-is-the-ld-preload-trick

If you set LD_PRELOAD to the path of a shared object, that file will be loaded before any other library (including the C runtime, libc.so). So to run ls with your special malloc() implementation, do this:

$ LD_PRELOAD=/path/to/my/malloc.so /bin/ls

在开始讲述为什么要当心LD_PRELOAD环境变量之前，请让我先说明一下程序的链接。所谓链接，也就是说编译器找到程序中所引用的函数或全局变量所存在的位置。一般来说，程序的链接分为静态链接和动态链接，静态链接就是把所有所引用到的函数或变量全部地编译到可执行文件中。动态链接则不会把函数编译到可执行文件中，而是在程序运行时动态地载入函数库，也就是运行链接。所以，对于动态链接来说，必然需要一个动态链接库。动态链接库的好处在于，一旦动态库中的函数发生变化，对于可执行程序来说是透明的，可执行程序无需重新编译。这对于程序的发布、维护、更新起到了积极的作用。对于静态链接的程序来说，函数库中一个小小的改动需要整个程序的重新编译、发布，对于程序的维护产生了比较大的工作量。


 

当然，世界上没有什么东西都是完美的，有好就有坏，有得就有失。动态链接所带来的坏处和其好处一样同样是巨大的。因为程序在运行时动态加载函数，这也就为他人创造了可以影响你的主程序的机会。试想，一旦，你的程序动态载入的函数不是你自己写的，而是载入了别人的有企图的代码，通过函数的返回值来控制你的程序的执行流程，那么，你的程序也就被人“劫持”了。


 

LD_PRELOAD简介



 

在UNIX的动态链接库的世界中，LD_PRELOAD就是这样一个环境变量，它可以影响程序的运行时的链接（Runtime linker），它允许你定义在程序运行前优先加载的动态链接库。这个功能主要就是用来有选择性的载入不同动态链接库中的相同函数。通过这个环境变量，我们可以在主程序和其动态链接库的中间加载别的动态链接库，甚至覆盖正常的函数库。一方面，我们可以以此功能来使用自己的或是更好的函数（无需别人的源码），而另一方面，我们也可以以向别人的程序注入恶意程序，从而达到那不可告人的罪恶的目的。


 

我们知道，Linux的用的都是glibc，有一个叫libc.so.6的文件，这是几乎所有Linux下命令的动态链接中，其中有标准C的各种函数。对于GCC而言，默认情况下，所编译的程序中对标准C函数的链接，都是通过动态链接方式来链接libc.so.6这个函数库的

https://blog.csdn.net/haoel/article/details/1602108

