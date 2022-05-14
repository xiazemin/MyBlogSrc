---
title: execute_data
layout: post
category: php
author: 夏泽民
---
1.EG(executor_globals/zend_executor_globals)

PHP整个生命周期中最主要的一个结构，是一个全局变量，在main执行前分配(非ZTS下)，直到PHP退出，它记录着当前请求全部的信息

2.EX(execute_data/zend_execute_data)

在执行过程中最核心的一个结构，每次函数的调用、include/require、eval等都会生成一个新的结构，它表示当前的作用域、代码的执行位置以及局部变量的分配等等，
<!-- more -->
#define EX(element)             ((execute_data)->element)

struct _zend_execute_data {
    const zend_op       *opline;  //指向当前执行的opcode，初始时指向zend_op_array起始位置
    zend_execute_data   *call;             /* current call                   */
    zval                *return_value;  //返回值指针
    zend_function       *func;          //当前执行的函数（非函数调用时为空）
    zval                 This;          //这个值并不仅仅是面向对象的this，还有另外两个值也通过这个记录：call_info + num_args，分别存在zval.u1.reserved、zval.u2.num_args
    zend_class_entry    *called_scope;  //当前call的类
    zend_execute_data   *prev_execute_data; //函数调用时指向调用位置作用空间
    zend_array          *symbol_table; //全局变量符号表
#if ZEND_EX_USE_RUN_TIME_CACHE
    void               **run_time_cache;   /* cache op_array->run_time_cache */
#endif
#if ZEND_EX_USE_LITERALS
    zval                *literals;  //字面量数组，与func.op_array->literals相同
#endif
};

.Zend的执行流程

在Zend VM中zend_execute_data的zend_execute_data.current_execute_data,zend_execute_data.prev_execute_data保存的是调用方的信息，实现了call/ret，zend_execute_data后面会分配额外的内存空间用于局部变量的存储，实现了ebp/esp的作用。
step1: 为当前作用域分配一块内存，充当运行栈，zend_execute_data结构、所有局部变量、中间变量等等都在此内存上分配
step2: 初始化全局变量符号表，然后将全局执行位置指针EG(current_execute_data)指向step1新分配的zend_execute_data，然后将zend_execute_data.opline指向op_array的起始位置
step3: 从EX(opline)开始调用各opcode的C处理handler(即_zend_op.handler)，每执行完一条opcode将EX(opline)++继续执行下一条，直到执行完全部opcode，函数/类成员方法调用、if的执行过程：
step3.1: if语句将根据条件的成立与否决定EX(opline) + offset所加的偏移量，实现跳转
step3.2: 如果是函数调用，则首先从EG(function_table)中根据function_name取出此function对应的编译完成的zend_op_array，然后像step1一样新分配一个zend_execute_data结构，将EG(current_execute_data)赋值给新结构的prev_execute_data，再将EG(current_execute_data)指向新的zend_execute_data，最后从新的zend_execute_data.opline开始执行，切换到函数内部，函数执行完以后将EG(current_execute_data)重新指向EX(prev_execute_data)，释放分配的运行栈，销毁局部变量，继续从原来函数调用的位置执行
step3.3: 类方法的调用与函数基本相同，后面分析对象实现的时候再详细分析
step4: 全部opcode执行完成后将step1分配的内存释放，这个过程会将所有的局部变量"销毁"，执行阶段结束




4.函数的执行流程

【初始化阶段】 这个阶段首先查找到函数的zend_function，普通function就是到EG(function_table)中查找，成员方法则先从EG(class_table)中找到zend_class_entry，然后再进一步在其function_table找到zend_function，接着就是根据zend_op_array新分配 zend_execute_data 结构并设置上下文切换的指针
【参数传递阶段】 如果函数没有参数则跳过此步骤，有的话则会将函数所需参数传递到 初始化阶段 新分配的 zend_execute_data动态变量区
【函数调用阶段】 这个步骤主要是做上下文切换，将执行器切换到调用的函数上，可以理解会在这个阶段__递归调用zend_execute_ex__函数实现call的过程(实际并一定是递归，默认是在while(1){...}中切换执行空间的，但如果我们在扩展中重定义了zend_execute_ex用来介入执行流程则就是递归调用)
【函数执行阶段】 被调用函数内部的执行过程，首先是接收参数，然后开始执行opcode
【函数返回阶段】 被调用函数执行完毕返回过程，将返回值传递给调用方的zend_execute_data变量区，然后释放zend_execute_data以及分配的局部变量，将上下文切换到调用前，回到调用的位置继续执行，这个实际是函数执行中的一部分，不算是独立的一个过程

https://www.debug8.com/php/t_9433.html

在PHP源码中通常是一个execute_data变量，与该变量相关的宏是#define EX(element) execute_data.element。

查看了下opcode的执行，发现_zend_execute_data对于PHP代码的执行也是关键的数据结构。所以本文主要围绕该数据结构展开来讲~

PHP代码的执行
总所周知，PHP代码的执行会经历：词法分析->语法分析->opcode生成->执行opcode。

对“opcode生成->执行opcode”展开来就是：

zend_compile_file()将PHP代码文件编译为op_array（opcode的数组）
初始化一个execute_data，之后赋值：EX(op_array) = op_array;EX(opline) = op_array->opcodes;
执行opcode：EX(opline)->handler(execute_data)，这是在一个while语句中执行的。
execute_data可以理解为一个执行上下文，通过gdb单步调试，发现：

当我们执行php hello.php时，会创建一个execute_data，当出现include/require或函数调用时，则会新建一个execute_data后，切换到新的execute_data执行，当一个execute_data执行结束，会切换到上层的execute_data继续执行。
http://yangxikun.github.io/php/2016/03/24/php-execute_data.html

https://www.php.net/manual/zh/migration55.internals.php
Zend/zend_vm_execute.h文件中的execute函数实现中，zend_execute_data类型的execute_data变量贯穿整个中间代码的执行过程， 其在调用时并没有直接使用execute_data，而是使用EX宏代替，其定义在Zend/zend_compile.h文件中，如下：

#define EX(element) execute_data.element
因此我们在execute函数或在opcode的实现函数中会看到EX(fbc)，EX(object)等宏调用， 它们是调用函数局部变量execute_data的元素：execute_data.fbc和execute_data.object。 execute_data不仅仅只有fbc、object等元素，它包含了执行过程中的中间代码，上一次执行的函数，函数执行的当前作用域，类等信息。

http://www.phppan.com/2012/02/php-execute-data/

https://type.so/c/php-extension-in-action-get-arguments-after-zend-execute-ex.html

https://blog.csdn.net/weixin_34405354/article/details/90652400

https://www.cnblogs.com/hellohell/p/9101803.html

https://blog.csdn.net/xiaolei1982/article/details/52140544
