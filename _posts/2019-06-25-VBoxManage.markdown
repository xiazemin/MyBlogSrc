---
title: VBoxManage
layout: post
category: docker
author: 夏泽民
---
VBoxManage是VirtualBox的命令行接口。利用他，你可以在主机操作系统的命令行中完全地控制VirtualBox。VBoxManage支持GUI可访问的全部功能，而且更多。VBoxManage展示了虚拟化引擎的全部特征，包括GUI无法访问的。
<!-- more -->
使用VBoxManage时要记住两件事：
第一，VBoxManage必须和一个具体和“子命令”一起使用，比如“list”或“createvm“或“startvm”。
第二，大多数子命令需要在其后指定特定的虚拟机。有两种方式：
指定虚拟机的名称，和其在GUI中显示的一样。注意，如果名称包含空格，必须将全部名称包含在双引号中（和命令行参数包含空格时的要求一样）。
例如：
VBoxManage startvm "Windows XP"
指定其UUID，VirtualBox用来引用虚拟机的内部唯一标识符。设上述名称为“Windows XP”的虚拟机有如下UUID，下面的命令有同样的效果：
 
VBoxManage startvm 670e746d-abea-4ba6-ad02-2a3b043810a5
使用VBoxManage list vms可列出当前注册的所有虚拟机的名称及其对应的UUID。
通过命令行控制VirtualBox的典型用法如下：
使用命令新建虚拟机并立即在VirtualBox中注册，使VBoxManage createvm的--register选项：
 
$ VBoxManage createvm --name "SUSE 10.2" --register
VirtualBox Command Line Management Interface Version 3.1.6
(C) 2005-2010 Sun Microsystems, Inc.
All rights reserved.
 
Virtual machine 'SUSE 10.2' is created.
UUID: c89fc351-8ec6-4f02-a048-57f4d25288e5
Settings file: '/home/username/.VirtualBox/Machines/SUSE 10.2/SUSE 10.2.xml'
从上面的输出可以看到，一个新的虚拟机被创建，带有一个新的UUID和新的XML的设置文件。
显示虚拟机的配置，使用VBoxManage showvminfo；详见“VBoxManage showvminfo”。
修改虚拟机的设置，使用VBoxManage modifyvm，例如：
 
VBoxManage modifyvm "Windows XP" --memory "512MB"
详见“VBoxManage modifyvm”。
控制虚拟机的运行，使用下列其中一个：
启动当前关闭的虚拟机，使用VBoxManage startvm；详见“VBoxManage startvm”。
暂停或保存当前运行的虚拟机，使用VBoxManage controlvm；详见“VboxManage controlvm”。
命令概述
不带参数运行VBoxManage或使用了无效的参数，将显示下面的语法图。注意，根据主机平台，输出可能会稍有不同；如有疑问，请检查VBoxManage在您的特定主机的可用命令输出。
（译者注：没翻译语法图，请运行VBoxManage查看输出，原文见http://www.virtualbox.org/manual/ch08.html#id2535703）。
每次调用VBoxManage，只能执行一个命令。但是，一个命令可能支持几个子命令在同一行被调用。接下来的部分是每个命令的详细参考。
VBoxManage list
list命令提供你的系统和VirtualBox当前设置的相关信息。
VboxManage list有如下可用子命令：
vms 列出当前在VirtualBox注册的所有虚拟机。默认显示包含每个虚拟机的名字和UUID的紧凑列表。如果指定了--long或--l参数，将显示和showvminfo命令一样的详细列表。
runningvms 用和vms相同的格式列出当前正在运行的虚拟机的唯一标识符（UUID）。
hdds，dvds，floppies 显示当前所有在VirtualBox注册的虚拟磁盘镜像的信息，包括其所有设置，在VirtualBox中的UUID和与其关联的所有文件。
ostypes 列出VirtualBox目前支持的所有客户机操作系统，及其在modifyvm命令中引用它的标识符。
hostdvds，hostfloppies，hostifs 相应地，列出主机上的DVD，软驱和网络接口，及用来在VirtualBox中访问他们的名字。
hostusb 提供主机上的USB设备的信息，特别是用来建立USB筛选器的信息和当前是否被主机使用。
usbfilters 列出所有在VirtualBox中注册的全局USB筛选器——即，所有虚拟机都可能访问的设备的筛选器——及其参数。
systemproperties 显示部分VirutalBox的全局设置，比如客户机内存和虚拟硬盘尺寸的最大和最小值，文件夹设置和当前使用的验证库。
hddbackends 列出所有VirtualBox已知的硬盘驱动器后端。除了后端本身的名字，还显示了功能说明、配置和其他有用信息。
VBoxManage showvminfo
showvminfo命令显示特定虚拟机的信息。这和VBoxManage list vms --long为所有虚拟机显示的内容相同。
你将得到类似下面的信息：
$ VBoxManage showvminfo "Windows XP"
VirtualBox Command Line Management Interface Version 3.1.6
(C) 2005-2010 Sun Microsystems, Inc.
