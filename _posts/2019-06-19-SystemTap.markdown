---
title: SystemTap
layout: post
category: linux
author: 夏泽民
---
SystemTap 是监控和跟踪运行中的 Linux 内核的操作的动态方法。这句话的关键词是动态，因为 SystemTap 没有使用工具构建一个特殊的内核，而是允许您在运行时动态地安装该工具。它通过一个名为Kprobes 的应用编程接口（API）来实现该目的
<!-- more -->
SystemTap 与一种名为 DTrace 的老技术相似，该技术源于 Sun Solaris 操作系统。在 DTrace 中，开发人员可以用 D 编程语言（C 语言的子集，但修改为支持跟踪行为）编写脚本。DTrace 脚本包含许多探针和相关联的操作，这些操作在探针 “触发” 时发生。例如，探针可以表示简单的系统调用，也可以表示更加复杂的交互，比如执行特定的代码行。清单 1 显示了 DTrace 脚本的一个简单例子，它计算每个进程发出的系统调用的数量（注意，使用字典将计数和进程关联起来）。该脚本的格式包含探针（在发出系统调用时触发）和操作（对应的操作脚本）。

清单 1. 计算每个进程的系统调用的简单 DTrace 脚本
1
2
3
4
5
6
syscall:::entry 
{ 
 
  @num[pid,execname] = count(); 
 
}
DTrace 是 Solaris 最引人注目的部分，所以在其他操作系统中开发它并不奇怪。DTrace 是在 Common Development and Distribution License (CDDL) 之下发行的，并且被移植到 FreeBSD 操作系统中。

另一个非常有用的内核跟踪工具是 ProbeVue，它是 IBM 为 IBM® AIX® 操作系统 6.1 开发的。您可以使用 ProbeVue 探查系统的行为和性能，以及提供特定进程的详细信息。这个工具使用一个标准的内核以动态的方式进行跟踪。清单 2 显示了 ProbeVue 脚本的一个例子，它指出发出 sync 系统调用的特定进程。

清单 2. 指出哪个进程调用 sync 的简单 ProbeVue 脚本
1
2
3
4
5
@@syscall:*:sync:entry
{
  printf( "sync() syscall invoked by process ID %d\n", __pid );
  exit();
}
考虑到 DTrace 和 ProbeVue 在各自的操作系统中的巨大作用，为 Linux 操作系统策划一个实现该功能的开源项目是势不可挡的。SystemTap 从 2005 年开始开发，它提供与 DTrace 和 ProbeVue 类似的功能。许多社区还进一步完善了它，包括 Red Hat、Intel、Hitachi 和 IBM 等。

这些解决方案在功能上都是类似的，在触发探针时使用探针和相关联的操作脚本。现在，我们看一下 SystemTap 的安装，然后探索它的架构和使用。

安装 SystemTap
您可能仅需一个 SystemTap 安装就可以支持 SystemTap，具体情况取决于您的分发版和内核。对于其他情况，需要使用一个调试内核映像。这个小节介绍在 Ubuntu version 8.10 (Intrepid Ibex) 上安装 SystemTap 的步骤，但这并不是一个具有代表性的 SystemTap 安装。在 参考资料 部分中，您可以找到在其他分发版和版本上安装 SystemTap 的更多信息。

对大部分用户而言，安装 SystemTap 都非常简单。对于 Ubuntu，使用 apt-get：

1
$ sudo apt-get install systemtap
在安装完成之后，您可以测试内核看它是否支持 SystemTap。为此，使用以下简单的命令行脚本：

1
$ sudo stap -ve 'probe begin { log("hello world") exit() }'
如果该脚本能够正常运行，您将在标准输出 [stdout] 中看到 “hello world”。如果没有看到这两个单词，则还需要其他工作。对于 Ubuntu 8.10，需要使用一个调试内核映像。应该使用 apt-get 获取包 linux-image-debug-generic 就可以获得它的。但这里不能直接使用 apt-get，因此您可以下载该包并使用 dpkg 安装它。您可以下载通用的调用映像包并按照以下的方式安装它：

1
2
3
$ wget http://ddebs.ubuntu.com/pool/main/l/linux/
          linux-image-debug-2.6.27-14-generic_2.6.27-14.39_i386.ddeb
$ sudo dpkg -i linux-image-debug-2.6.27-14-generic_2.6.27-14.39_i386.ddeb
现在，已经安装了通用的调试映像。对于 Ubuntu 8.10，还需要一个步骤：SystemTap 分发版有一个问题，但可以通过修改 SystemTap 源代码轻松解决。查看 参考资料 获得如何更新运行时 time.c 文件的信息。

如果您使用定制的内核，则需要确保启用内核选项 CONFIG_RELAY、CONFIG_DEBUG_FS、CONFIG_DEBUG_INFO 和 CONFIG_KPROBES。

SystemTap 的架构
让我们深入探索 SystemTap 的某些细节，理解它如何在运行的内核中提供动态探针。您还将看到 SystemTap 是如何工作的，从构建进程脚本到在运行的内核中激活脚本。

动态地检查内核
SystemTap 用于检查运行的内核的两种方法是 Kprobes 和 返回探针。但是理解任何内核的最关键要素是内核的映射，它提供符号信息（比如函数、变量以及它们的地址）。有了内核映射之后，就可以解决任何符号的地址，以及更改探针的行为。

Kprobes 从 2.6.9 版本开始就添加到主流的 Linux 内核中，并且为探测内核提供一般性服务。它提供一些不同的服务，但最重要的两种服务是 Kprobe 和 Kretprobe。Kprobe 特定于架构，它在需要检查的指令的第一个字节中插入一个断点指令。当调用该指令时，将执行针对探针的特定处理函数。执行完成之后，接着执行原始的指令（从断点开始）。

Kretprobes 有所不同，它操作调用函数的返回结果。注意，因为一个函数可能有多个返回点，所以听起来事情有些复杂。不过，它实际使用一种称为 trampoline 的简单技术。您将向函数条目添加一小段代码，而不是检查函数中的每个返回点。这段代码使用 trampoline 地址替换堆栈上的返回地址 —— Kretprobe 地址。当该函数存在时，它没有返回到调用方，而是调用 Kretprobe（执行它的功能），然后从 Kretprobe 返回到实际的调用方。

假如现在有这么一个需求：需要获取正在运行的 Linux 系统的信息，如我想知道系统什么时候发生系统调用，发生的是什么系统调用等这些信息，有什么解决方案呢？

最原始的方法是，找到内核系统调用的代码，加上我们需要获得信息的代码、重新编译内核、安装、选择我们新编译的内核重启。这种做法对于内核开发人员简直是梦魇，因为一遍做下来至少得需要1个多小时，不仅破坏了原有内核代码，而且如果换了一个需求又得重新做一遍上面的工作。所以，这种调试内核的方法效率是极其底下的。
之后内核引入了一种Kprobe机制，可以用来动态地收集调试和性能信息的工具，是一种非破坏性的工具，用户可以用它跟踪运行中内核任何函数或执行的指令等。相比之前的做法已经有了质的提高了，但Kprobe并没有提供一种易用的框架，用户需要自己去写模块，然后安装，对用户的要求还是蛮高的。
systemtap 是利用Kprobe 提供的API来实现动态地监控和跟踪运行中的Linux内核的工具，相比Kprobe，systemtap更加简单，提供给用户简单的命令行接口，以及编写内核指令的脚本语言。对于开发人员，systemtap是一款难得的工具。
下面将会介绍systemtap的安装、systemtap的工作原理以及几个简单的示例。

systemtap 的安装
我的主机 Linux 发行版是32位 Ubuntu13.04，内核版本 3.8.0-30。由于 systemtap 运行需要内核的调试信息支撑，默认发行版的内核在配置时这些调试开关没有打开，所以安装完systemtap也是无法去探测内核信息的。 下面我以两种方式安装并运行 systemtap：

方法一
编译内核以支持systemtap 
我们重新编译内核让其支持systemtap，首先你想让内核中有调试信息，编译内核时需要加上 -g 标志；其次，你还需要在配置内核时将 Kprobe 和 debugfs 开关打开。最终效果是，你能在内核 .config 文件中看到下面四个选项是设置的：

  CONFIG_DEBUG_INFO
  CONFIG_KPROBES
  CONFIG_DEBUG_FS
  CONFIG_RELAY
配置完之后，按照之前你编译内核的步骤编译即可。

获取systemtap源码 
从此地址 https://sourceware.org/systemtap/ftp/releases下载已经发布的systemtap的源代码，截至目前（2013.9.17）最新版本为systemtap-2.3。下载完之后解压。 当然你还可以使用 git 去克隆最新的版本（2.4），命令如下：

  git clone git://sources.redhat.com/git/systemtap.git
编译安装systemtap 
如果你下载的是最新版本的systemtap，那么你需要新版的 elfutils，可以从https://fedorahosted.org/releases/e/l/elfutils/ 下载elfutils-0.156 版本。下载之后解压缩到适合的目录（我放在~/Document/ 下），不需要安装，只要配置systemtap时指定其位置即可。 进入之前解压systemtap的目录，使用下面命令进行配置：

   ./configure --with-elfutils=~/Document/elfutils-0.156
以这里方法配置之后，你只需要再运行 make install 即完成systemtap的编译安装。如果需要卸载的话，运行 make uninstall。

方法二
由于发行版的内核默认无内核调试信息，所以我们还需要一个调试内核镜像，在http://ddebs.ubuntu.com/pool/main/l/linux/ 找到你的内核版本相对应的内核调试镜像（版本号包括后面的发布次数、硬件体系等都必须一致），如针对我上面的内核版本，就可以用如下命令下载安装内核调试镜像：

$ wget http://ddebs.ubuntu.com/pool/main/l/linux/linux-image-debug-3.8.0-30-generic_dbgsym_3.8.0-30.43_i386.ddeb
$ sudo dpkg -i linux-image-debug-3.8.0-30-generic_dbgsym_3.8.0-30.43_i386.ddeb
一般这种方法下，你只需要使用apt在线安装systemtap即可：

$sudo apt-get install systemtap
当然方法二仅限于Ubuntu发行版，至于其他的发行版并不能照搬，网上也有很多相关的资料。

systemtap 测试示例
安装完systemtap之后，我们需要测试一下systemtap是否能正确运行：

示例一：打印hello systemtap
以root用户或者具有sudo权限的用户运行以下命令：

$stap -ve 'probe begin { log("hello systemtap!") exit() }'
如果安装正确，会得到如下类似的输出结果：

Pass 1: parsed user script and 96 library script(s) using 55100virt/26224res/2076shr/25172data kb, in 120usr/0sys/119real ms.
Pass 2: analyzed script: 1 probe(s), 2 function(s), 0 embed(s), 0 global(s) using 55496virt/27016res/2172shr/25568data kb, in 0usr/0sys/4real ms.
Pass 3: translated to C into "/tmp/stapYqNuF9/stap_e2d1c1c9962c809ee9477018c642b661_939_src.c" using 55624virt/27380res/2488shr/25696data kb, in 0usr/0sys/0real ms.
Pass 4: compiled C into "stap_e2d1c1c9962c809ee9477018c642b661_939.ko" in 1230usr/160sys/1600real ms.
Pass 5: starting run.
hello systemtap!
Pass 5: run completed in 0usr/10sys/332real ms.
示例二：打印4s内所有open系统调用的信息
创建systemtap脚本文件test2.stp:

#!/usr/bin/stap

probe begin 
{
    log("begin to probe")
}

probe syscall.open
{
    printf ("%s(%d) open (%s)\n", execname(), pid(), argstr)
}

probe timer.ms(4000) # after 4 seconds
{
    exit ()
}

probe end
{
    log("end to probe")
}
将该脚本添加可执行的权限 chmod +x test2.stp ，使用./test2.stp 运行该脚本，即可打印4s内所有open系统调用的信息，打印格式为：进程名（进程号）打开什么文件。 大家可以自行去测试，如果两个示例都能正确运行，基本上算是安装成功了！

systemtap 工作原理
systemtap 的核心思想是定义一个事件（event），以及给出处理该事件的句柄（Handler）。当一个特定的事件发生时，内核运行该处理句柄，就像快速调用一个子函数一样，处理完之后恢复到内核原始状态。这里有两个概念：

事件（Event）：systemtap 定义了很多种事件，例如进入或退出某个内核函数、定时器时间到、整个systemtap会话启动或退出等等。
句柄（Handler）：就是一些脚本语句，描述了当事件发生时要完成的工作，通常是从事件的上下文提取数据，将它们存入内部变量中，或者打印出来。
Systemtap 工作原理是通过将脚本语句翻译成C语句，编译成内核模块。模块加载之后，将所有探测的事件以钩子的方式挂到内核上，当任何处理器上的某个事件发生时，相应钩子上句柄就会被执行。最后，当systemtap会话结束之后，钩子从内核上取下，移除模块。整个过程用一个命令 stap 就可以完成。
