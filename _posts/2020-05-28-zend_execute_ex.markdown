---
title: zend_execute_ex
layout: post
category: php
author: 夏泽民
---
解释器引擎最终执行op的函数是zend_execute，实际上zend_execute是一个函数指针，在引擎初始化的时候zend_execute默认指向了execute,这个execute定义在{PHPSRC}/Zend/zend_vm_execute.h

ZEND_API void execute(zend_op_array *op_array TSRMLS_DC)
<!-- more -->
http://www.nowamagic.net/librarys/veda/detail/1580

   解释器引擎最终执行op的函数是zend_execute，实际上zend_execute是一个函数指针，在引擎初始化的时候zend_execute默认指向了execute,这个execute定义在
   {PHPSRC}/Zend/zend_vm_execute.h：
ZEND_API void execute(zend_op_array *op_array TSRMLS_DC)

zend_op_array简介

此类型的定义在{PHPSRC}/Zend/zend_compile.h:

 execute函数开始的时候是一些基础变量的申明，其中zend_execute_data *execute_data;是执行期的数据结构，此变量在进行一定的初始化之后将会被传递给每个op的handler函数作为参数，op在执行过程中随时有可能改变execute_data中的内容。

https://type.so/c/php-extension-in-action-get-arguments-after-zend-execute-ex.html
https://blog.csdn.net/phpkernel/article/details/5721384

http://www.nowamagic.net/librarys/veda/detail/1580