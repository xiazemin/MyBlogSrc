---
title: execute_data
layout: post
category: php
author: 夏泽民
---
php的op_array与execute_data的关系
php分为几个阶段包括生成opcode阶段和执行opcode阶段，其实分别对应的就是上面两个数据结构，
并且两个数据结构都是在解析到新的函数时分配新的空间，然后层层嵌套，最外层总是有个大的op_array与execute_data,具体点说就是这两个数据结构存储的是当前函数下的变量环境。
然后就是上面两个不同阶段存储该阶段应该存储的数据，然后可供下一层调用。
<!-- more -->
http://www.voidcn.com/article/p-qjlxtflt-mc.html

https://www.php.cn/php-weizijiaocheng-392486.html
https://blog.csdn.net/xiaolei1982/article/details/52140544

zend_op_array : zend引擎执行阶段的输入数据结构，整个执行阶段都是操作这个数据结构。
{% raw %}
      　　　　　　 zend_op_array有三个核心部分：opcode指令(对应c的指令)

                                                   字面量存储(变量初始值、调用的函数名称、类名称、常量名称等等称之为字面量)

                                                   变量分配的情况 (当前array定义的变量 临时变量的数量 编号，执行初始化一次性分配zval，使用时完全按照标号索引不是根据变量名)

         

           zend_executor_globals     PHP整个生命周期中最主要的一个结构，是一个全局变量，在main执行前分配(非ZTS下)，直到PHP退出，它记录着当前请求全部的信息，经常见到的一个宏EG操作的就是这个结构。

                                定义在zend_globals.h中：

 

                                   

 

                

               zend_execute_data  是执行过程中最核心的一个结构，每次函数的调用、include/require、eval等都会生成一个新的结构，它表示当前的作用域、代码的执行位置以及局部变量的分配等等，等同于机器码执行过程中stack的角色，后面分析具体执行流程的时候会详细分析其作用。 

              zend_execute_data与zend_op_array的关联关系：

                                         

2.执行过程

        Zend的executor与linux二进制程序执行的过程是非常类似的。

        在C程序执行时有两个寄存器ebp、esp分别指向当前作用栈的栈顶、栈底，局部变量全部分配在当前栈，函数调用、返回通过call、ret指令完成，调用时call将当前执行位置压入栈中，返回时ret将之前执行位置出栈，跳回旧的位置继续执行。

        Zend VM中zend_execute_data就扮演了这两个角色，zend_execute_data.prev_execute_data保存的是调用方的信息，实现了call/ret，zend_execute_data后面会分配额外的内存空间用于局部变量的存储，实现了ebp/esp的作用。

                    a. 为当前作用域分配一块内存，充当运行栈，zend_execute_data结构、所有局部变量、中间变量等等都在此内存上分配

                    b.初始化全局变量符号表，然后将全局执行位置指针EG(current_execute_data)指向步骤a新分配的zend_execute_data，然后将zend_execute_data.opline指向op_array的起始位置

                    c.从EX(opline)开始调用各opcode的C处理handler(即_zend_op.handler)，每执行完一条opcode将EX(opline)++继续执行下一条，直到执行完全部opcode

                                if语句将根据条件的成立与否决定EX(opline) + offset所加的偏移量，实现跳转

                                如果是函数调用，则首先从EG(function_table)中根据function_name取出此function对应的编译完成的zend_op_array，然后像步骤a一样新分配一个zend_execute_data结构，将EG(current_execute_data)赋值给新结构的prev_execute_data，再将EG(current_execute_data)指向新的zend_execute_data，最后从新的zend_execute_data.opline开始执行，切换到函数内部，函数执行完以后将EG(current_execute_data)重新指向EX(prev_execute_data)，释放分配的运行栈，销毁局部变量，继续从原来函数调用的位置执行

                                类方法的调用与函数基本相同

                    d.全部opcode执行完成后将步骤a分配的内存释放，这个过程会将所有的局部变量"销毁"，执行阶段结束

                                    

 

                              首先根据zend_execute_data、当前zend_op_array中局部/临时变量数计算需要的内存空间，编译阶段zend_op_array的结果，在编译过程中已经确定当前作用域下有多少个局部变量(func->op_array.last_var)、临时/中间/无用变量(func->op_array.T)，从而在执行之初就将他们全部分配完成。
{% endraw %}
https://www.cnblogs.com/hellohell/p/9101803.html
https://www.ucloud.cn/yun/28588.html
