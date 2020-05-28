---
title: zend_execute_data
layout: post
category: web
author: 夏泽民
---
研究下PHP Generator的实现，发现对于Generator很重要的一个数据结构为_zend_execute_data，在PHP源码中通常是一个execute_data变量，与该变量相关的宏是#define EX(element) execute_data.element。

查看了下opcode的执行，发现_zend_execute_data对于PHP代码的执行也是关键的数据结构

https://www.laruence.com/2009/04/28/719.html
https://www.laruence.com/2011/09/13/2139.html

这里的EX宏对应全局变量execute_data，EG宏对应全局变量executor_globals，要区分开
  https://segmentfault.com/a/1190000019382316
  
  zend_executor_globals     PHP整个生命周期中最主要的一个结构，是一个全局变量，在main执行前分配(非ZTS下)，直到PHP退出，它记录着当前请求全部的信息，经常见到的一个宏EG操作的就是这个结构。
  
   zend_execute_data  是执行过程中最核心的一个结构，每次函数的调用、include/require、eval等都会生成一个新的结构，它表示当前的作用域、代码的执行位置以及局部变量的分配等等，等同于机器码执行过程中stack的角色
   https://www.cnblogs.com/hellohell/p/9101803.html
<!-- more -->
PHP代码的执行会经历：词法分析->语法分析->opcode生成->执行opcode。

对“opcode生成->执行opcode”展开来就是：

zend_compile_file()将PHP代码文件编译为op_array（opcode的数组）
初始化一个execute_data，之后赋值：EX(op_array) = op_array;EX(opline) = op_array->opcodes;
执行opcode：EX(opline)->handler(execute_data)，这是在一个while语句中执行的。
execute_data可以理解为一个执行上下文，通过gdb单步调试，发现：

当我们执行php hello.php时，会创建一个execute_data，当出现include/require或函数调用时，则会新建一个execute_data后，切换到新的execute_data执行，当一个execute_data执行结束，会切换到上层的execute_data继续执行。


每一个op_array的最后一条opcode都是RETURN，这也是为什么在类似PHP框架的配置文件中可以直接<?php return [配置信息]; ?>，而加载配置文件的代码为$config = require('配置文件');。

通过VLD扩展可以打印出一个PHP执行文件产生的opcode：

_zend_execute_data结构体说明
	<img src="{{site.url}}{{site.baseurl}}/img/zend_execuable.png"/>
_zend_execute_data结构体的初始化
在Zend/zend_execute.c中的函数zend_execute_data *i_create_execute_data_from_op_array(zend_op_array *op_array, zend_bool nested TSRMLS_DC)负责创建新的execute_data，其内存分配是在EG(argument_stack)（指向某一段内存，_zend_vm_stack数据结构，不足时会增长）上进行的。
typedef struct _zend_vm_stack *zend_vm_stack;
struct _zend_vm_stack {
    void **top;
    void **end;
    zend_vm_stack prev;
};

在一个初始化好的EG(argument_stack)上分配一个execute_data：
http://yangxikun.github.io/php/2016/03/24/php-execute_data.html
https://blog.csdn.net/xiaolei1982/article/details/52140544
https://segmentfault.com/a/1190000015748477

zend_execute_data.prev_execute_data保存的是调用方的信息，实现了call/ret，zend_execute_data后面会分配额外的内存空间用于局部变量的存储，实现了ebp/esp的作用。

                    a. 为当前作用域分配一块内存，充当运行栈，zend_execute_data结构、所有局部变量、中间变量等等都在此内存上分配

                    b.初始化全局变量符号表，然后将全局执行位置指针EG(current_execute_data)指向步骤a新分配的zend_execute_data，然后将zend_execute_data.opline指向op_array的起始位置

                    c.从EX(opline)开始调用各opcode的C处理handler(即_zend_op.handler)，每执行完一条opcode将EX(opline)++继续执行下一条，直到执行完全部opcode

                                if语句将根据条件的成立与否决定EX(opline) + offset所加的偏移量，实现跳转

                                如果是函数调用，则首先从EG(function_table)中根据function_name取出此function对应的编译完成的zend_op_array，然后像步骤a一样新分配一个zend_execute_data结构，将EG(current_execute_data)赋值给新结构的prev_execute_data，再将EG(current_execute_data)指向新的zend_execute_data，最后从新的zend_execute_data.opline开始执行，切换到函数内部，函数执行完以后将EG(current_execute_data)重新指向EX(prev_execute_data)，释放分配的运行栈，销毁局部变量，继续从原来函数调用的位置执行

                                类方法的调用与函数基本相同

                    d.全部opcode执行完成后将步骤a分配的内存释放，这个过程会将所有的局部变量"销毁"，执行阶段结束
   https://www.cnblogs.com/hellohell/p/9101803.html
   
 解释器引擎最终执行op的函数是zend_execute，实际上zend_execute是一个函数指针，在引擎初始化的时候zend_execute默认指向了execute,这个execute定义在{PHPSRC}/Zend/zend_vm_execute.h：
 
 执行过程详解
execute函数开始的时候是一些基础变量的申明，其中zend_execute_data *execute_data;是执行期的数据结构，此变量在进行一定的初始化之后将会被传递给每个op的handler函数作为参数，op在执行过程中随时有可能改变execute_data中的内容。
https://blog.csdn.net/weixin_34405354/article/details/90652400

php扩展实践zend_execute_ex层获取实参
其实在实现的 php 函数里面是很容易获取到的，参考 php 的 builtin 函数 func_get_args() 就可以知道了。

void **p;
int arg_count;
int i;
zend_execute_data *ex = EG(current_execute_data);

if (!ex || !ex->function_state.arguments) {
    RETURN_FALSE;
}

p = ex->function_state.arguments;
arg_count = (int)(zend_uintptr_t) *p;

for (i = 0; i < arg_count; i++) {
    zval *element, *arg;
    arg = *((zval **) (p - (arg_count - i)));
    php_var_dump(&arg, 1 TSRMLS_CC);
}
但是在 zend_execute_ex 中，是不能使用 function_state.arguments 来获取参数的，需要从 argument_stack 中获取调用函数的实参。                   
                    
https://type.so/c/php-extension-in-action-get-arguments-after-zend-execute-ex.html

https://www.ucloud.cn/yun/28588.html
Zend/zend_vm_execute.h文件中的execute函数实现中，zend_execute_data类型的execute_data变量贯穿整个中间代码的执行过程， 其在调用时并没有直接使用execute_data，而是使用EX宏代替，其定义在Zend/zend_compile.h文件中，如下：

#define EX(element) execute_data.element
因此我们在execute函数或在opcode的实现函数中会看到EX(fbc)，EX(object)等宏调用， 它们是调用函数局部变量execute_data的元素：execute_data.fbc和execute_data.object。 execute_data不仅仅只有fbc、object等元素，它包含了执行过程中的中间代码，上一次执行的函数，函数执行的当前作用域，类等信息

http://www.phppan.com/2012/02/php-execute-data/
https://www.php.cn/php-weizijiaocheng-392486.html

https://zhuanlan.zhihu.com/p/72162007


https://segmentfault.com/a/1190000019382316	