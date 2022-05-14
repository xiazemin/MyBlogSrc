---
title: hypervisor
layout: post
category: docker
author: 夏泽民
---
Hypervisor，又称虚拟机监视器（英语：virtual machine monitor，缩写为 VMM），是用来建立与执行虚拟机器的软件、固件或硬件。
被Hypervisor用来执行一个或多个虚拟机器的电脑称为主体机器（host machine），这些虚拟机器则称为客体机器（guest machine）。hypervisor提供虚拟的作业平台来执行客体操作系统（guest operating systems），负责管理其他客体操作系统的执行阶段；这些客体操作系统，共同分享虚拟化后的硬件资源。
<!-- more -->
Hypervisor——一种运行在基础物理服务器和操作系统之间的中间软件层，可允许多个操作系统和应用共享硬件。也可叫做VMM（ virtual machine monitor ），即虚拟机监视器。
Hypervisors是一种在虚拟环境中的“元”操作系统。他们可以访问服务器上包括磁盘和内存在内的所有物理设备。Hypervisors不但协调着这些硬件资源的访问，也同时在各个虚拟机之间施加防护。当服务器启动并执行Hypervisor时，它会加载所有虚拟机客户端的操作系统同时会分配给每一台虚拟机适量的内存，CPU，网络和磁盘。

种类
目前市场上各种x86 管理程序(hypervisor)的架构存在差异，三个最主要的架构类别包括：
· I型：虚拟机直接运行在系统硬件上，创建硬件全仿真实例，被称为“裸机”型。
裸机型在虚拟化中Hypervisor直接管理调用硬件资源，不需要底层操作系统，也可以将Hypervisor看
作一个很薄的操作系统。这种方案的性能处于主机虚拟化与操作系统虚拟化之间。
II型：虚拟机运行在传统操作系统上，同样创建的是硬件全仿真实例，被称为“托管（宿主）”型。
托管型/主机型Hypervisor运行在基础操作系统上，构建出一整套虚拟硬件平台
（CPU/Memory/Storage/Adapter），使用者根据需要安装新的操作系统和应用软件，底层和上层的
操作系统可以完全无关化，如Windows运行Linux操作系统。主机虚拟化中VM的应用程序调用硬件资
源时需要经过：VM内核->Hypervisor->主机内核，因此相对来说，性能是三种虚拟化技术中最差的。
Ⅲ型：虚拟机运行在传统操作系统上，创建一个独立的虚拟化实例（容器），指向底层托管操作系统，被称为“操作系统虚拟化”。

目前市场主要厂商及产品：VMware vSphere、微软Hyper-V、Citrix XenServer 、IBM PowerVM、Red Hat Enterprise Virtulization、Huawei FusionSphere、开源的KVM、Xen、VirtualBSD等。


hypervisor 之于操作系统类似于操作系统之于进程。它们为执行提供独立的虚拟硬件平台，而虚拟硬件平台反过来又提供对底层机器的虚拟的完整访问。但并不是所有 hypervisor 都是一样的，这是件好事，因为 Linux 就是以灵活性和选择性著称。本文首先简要介绍虚拟化和 hypervisor，然后探索两个基于 Linux 的 hypervisor。

虚拟化 就是通过某种方式隐藏底层物理硬件的过程，从而让多个操作系统可以透明地使用和共享它。这种架构的另一个更常见的名称是平台虚拟化。在典型的分层架构中，提供平台虚拟化的层称为 hypervisor （有时称为虚拟机管理程序 或 VMM）。

首先，类似于将用户空间应用程序和内核函数连接起来的系统调用，一个通常可用的虚拟化调用（hapercall，hypervisor 对操作系统进行的系统调用）层允许来宾系统向宿主操作系统发出请求。可以在内核中虚拟化 I/O，或通过来宾操作系统的代码支持它。故障必须由 hypervisor 亲自处理，从而解决实际的故障，或将虚拟设备故障发送给来宾操作系统。hypervisor 还必须处理在来宾操作系统内部发生的异常。（毕竟，来宾操作系统发生的错误仅会停止该系统，而不会影响 hypervisor 或其他来宾操作系统）。hypervisor 的核心要素之一是页映射器，它将硬件指向特定操作系统（来宾或 hypervisor）的页。最后，需要使用一个高级别的调度器在hypervisor和来宾操作系统之间传输控制。

首先是 KVM，它是首个被集成到 Linux 内核的 hypervisor 解决方案，并且实现了完整的虚拟化。其次是 Lguest，这是一个实验 hypervisor，它通过少量的更改提高准虚拟化。

KVM 针对运行在 x86 硬件硬件上的、驻留在内核中的虚拟化基础结构。KVM 是第一个成为原生 Linux 内核（2.6.20）的一部分的 hypervisor，它是由 Avi Kivity 开发和维护的，现在归 Red Hat 所有。

这个 hypervisor 提供 x86 虚拟化，同时拥有到 PowerPC® 和 IA64 的通道。另外，KVM 最近还添加了对对称多处理（SMP）主机（和来宾）的支持，并且支持企业级特性，比如活动迁移（允许来宾操作系统在物理服务器之间迁移）。

KVM 是作为内核模块实现的，因此 Linux 只要加载该模块就会成为一个hypervisor。KVM 为支持 hypervisor 指令的硬件平台提供完整的虚拟化（比如 Intel® Virtualization Technology [Intel VT] 或 AMD Virtualization [AMD-V] 产品）。KVM 还支持准虚拟化来宾操作系统，包括 Linux 和 Windows®。

这种技术由两个组件实现。第一个是可加载的 KVM 模块，当在 Linux 内核安装该模块之后，它就可以管理虚拟化硬件，并通过 /proc 文件系统公开其功能。第二个组件用于 PC 平台模拟，它是由修改版 QEMU 提供的。QEMU 作为用户空间进程执行，并且在来宾操作系统请求方面与内核协调。


硬件仿真
在物理机上创建一个模拟硬件的程序，来仿真所有的硬件，在这个程序之上运行虚拟机，最典型的就是QEMU了。

优点：虚拟机操作系统（VM OS）不需要更改
缺点：由于所有的硬件都是软件模拟的，所以性能很差

全虚拟化
虚拟机的操作系统与底层硬件是完全隔离的，由Hypervisor捕捉并进行转化由VM OS对硬件的调用代码，比较典型的有KVM。

优点：无需更改虚拟机操作系统，兼容性好。
缺点：性能一般，特别是I/O性能

半虚拟化
与全虚拟化技术类似，利用Hypervisor来实现对底层硬件的共享访问，但VM OS中需要集成半虚拟化相关的代码，也就是让虚拟机自己知道是一个虚拟机，来配合Hypervisor。通过这种方式无需捕捉特权指令，所以性能非常好。最典型的的是Xen。

优点：性能好
缺点：需要对VM OS进行更改

这里曾经自己有一个疑惑，不太清楚硬件仿真和全虚拟化的区别在哪里，也可能是中文名字的诱导，现在还算清晰一些，硬件仿真的方式，虚拟机执行的指令都是由仿真程序模拟的，而全虚拟化中的虚拟机的指令是经过Hypervisor转给底层硬件的，后者如果还算是真正的调用了底层硬件的话，前者根本就是假货，都是仿真程序模拟的。是由根本区别的。
Hypervisor又是什么？
上面说到虚拟化技术相当于一个榨汁的过程，更准确点来说算是一个完整榨汁的方案，那Hypervisor就是一个榨汁机。

虚拟化 就是通过某种方式隐藏底层物理硬件的过程，从而让多个操作系统可以透明地使用和共享它。这种架构的另一个更常见的名称是平台虚拟化。在典型的分层架构中，提供平台虚拟化的层称为 hypervisor （有时称为虚拟机管理程序 或 VMM）。来宾操作系统称为虚拟机（VM），因为对这些 VM 而言，硬件是专门针对它们虚拟化的。

hypervisor 可以划分为两大类。首先是类型 1，这种 hypervisor 是直接运行在物理硬件之上的。其次是类型 2，这种 hypervisor 运行在另一个操作系统（运行在物理硬件之上）中。类型 1 hypervisor 的一个例子是基于内核的虚拟机（KVM —— 它本身是一个基于操作系统的 hypervisor）。类型 2 hypervisor 包括 QEMU 和 WINE。

SNIA的定义：Hypervisor – Software that provides virtual machine environments which are used by guest operating systems，翻译过来就是：为客户端操作系统提供虚拟机环境的软件。

Hypervisor 是一种运行在物理服务器和操作系统之间的中间软件层（可以是软件程序，也可以是固件程序），可允许多个操作系统和应用共享一套基础物理硬件，因此也可以看作是虚拟环境中的“元”操作系统，它可以协调访问服务器上的所有物理设备和虚拟机，也叫虚拟机监视器VMM（Virtual Machine Monitor）。

Hypervisor是所有虚拟化技术的核心。非中断地支持多工作负载迁移的能力是Hypervisor的基本功能。当服务器启动并执行Hypervisor时，它会给每一台虚拟机分配适量的内存、CPU、网络和磁盘，并加载所有虚拟机的客户操作系统。

Hypervisor 翻译过来就是超级监督者，被引申用为超级管理程序、超多功能管理器、虚拟机管理器、VMM。

两种类型的Hypervisor，可以看到：

一种是裸机型，直接运行在硬件设备上的，这种架构搭建的虚拟化环境称Bare-Metal Hardware Virtualization裸机虚拟化环境；

一种是主机托管型，运行在具有虚拟化功能的操作系统上的，构建的是Hosted Virtualization 主机虚拟化环境。

两种Hypervisor的各自代表产品

Type 1: Hyper-V，VMware ESX Server，Citrix XenServer，Oracle OVM for SPARC，KVM；

Type 2: VMware Server，VMware Fusion，Microsoft Virtual Server，Oracle VM Virtual Box，Oracle VM for x86， Solaris Zones，Parallels，Microsoft Virtual PC。
Hypervisor有如下优点：

提高主机硬件的使用效率。因为一个主机可以运行多个虚拟机，这样主机的硬件资源能被高效充分的利用起来。
虚拟机移动性强。传统软件强烈捆绑在硬件上，转移一个软件至另一个服务器上耗时耗力（比如重新安装）；然而，虚拟机与硬件是独立的，这样使得虚拟机可以在本地或远程虚拟服务器上低消耗转移。
虚拟机彼此独立。一个虚拟机的奔溃不会影响其他分享同一硬件资源的虚拟机，大大提升安全性。
易保护，易恢复。Snapshot技术可以记录下某一时间点下的虚拟机状态，这使得虚拟机在错误发生后能快速恢复。



