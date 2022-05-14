---
title: Linux的mmap内存映射机制
layout: post
category: linux
author: 夏泽民
---
<!-- more -->
一个进程应该包括一个mm_struct(memory manage struct), 该结构是进程虚拟地址空间的抽象描述,里面包括了进程虚拟空间的一些管理信息: start_code, end_code, start_data, end_data, start_brk, end_brk等等信息.另外,也有一个指向进程虚存区表(vm_area_struct: virtual memory area)的指针,该链是按照虚拟地址的增长顺序排列的.在Linux进程的地址空间被分作许多区(vma),每个区(vma)都对应虚拟地址空间上一段连续的区域, vma是可以被共享和保护的独立实体,这里的vma就是前面提到的内存对象. 

设备驱动的mmap实现主要是将一个物理设备的可操作区域（设备空间）映射到一个进程的虚拟地址空间。这样就可以直接采用指针的方式像访问内存的方式访问设备。在驱动中的mmap实现主要是完成一件事，就是实际物理设备的操作区域到进程虚拟空间地址的映射过程。同时也需要保证这段映射的虚拟存储器区域不会被进程当做一般的空间使用，因此需要添加一系列的保护方式。

Linux提供了内存映射函数mmap,它把文件内容映射到一段内存上(准确说是虚拟内存上),通过对这段内存的读取和修改,实现对文件的读取和修改 。普通文件被映射到进程地址空间后，进程可以向访问普通内存一样对文件进行访问，不必再调用read()，write（）等操作。
<img src="{{site.url}}{{site.baseurl}}/img/linuxMMap.jpeg"/>
mmap系统调用的实现过程是
1.先通过文件系统定位要映射的文件；
2.权限检查,映射的权限不会超过文件打开的方式,也就是说如果文件是以只读方式打开,那么则不允许建立一个可写映射； 
3.创建一个vma对象,并对之进行初始化； 
4.调用映射文件的mmap函数,其主要工作是给vm_ops向量表赋值；
5.把该vma链入该进程的vma链表中,如果可以和前后的vma合并则合并；
6.如果是要求VM_LOCKED(映射区不被换出)方式映射,则发出缺页请求,把映射页面读入内存中.

2、munmap函数
munmap(void * start, size_t length):
该调用可以看作是mmap的一个逆过程.它将进程中从start开始length长度的一段区域的映射关闭,如果该区域不是恰好对应一个vma,则有可能会分割几个或几个vma.
 msync(void * start, size_t length, int flags):
把映射区域的修改回写到后备存储中.因为munmap时并不保证页面回写,如果不调用msync,那么有可能在munmap后丢失对映射区的修改.其中flags可以是MS_SYNC, MS_ASYNC, MS_INVALIDATE, MS_SYNC要求回写完成后才返回, MS_ASYNC发出回写请求后立即返回, MS_INVALIDATE使用回写的内容更新该文件的其它映射.该系统调用是通过调用映射文件的sync函数来完成工作的.
brk(void * end_data_segement):
将进程的数据段扩展到end_data_segement指定的地址,该系统调用和mmap的实现方式十分相似,同样是产生一个vma,然后指定其属性.不过在此之前需要做一些合法性检查,比如该地址是否大于mm->end_code, end_data_segement和mm->brk之间是否还存在其它vma等等.通过brk产生的vma映射的文件为空,这和匿名映射产生的vma相似,关于匿名映射不做进一步介绍.库函数malloc就是通过brk实现的.