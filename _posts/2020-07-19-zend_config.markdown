---
title: zend_config
layout: post
category: php
author: 夏泽民
---
 ./configure --with-php-config=/usr/local/bin/php-config ./configure --enable-debug

checking for bison version... ./configure: line 5574: -z: command not found
2.3
configure: error: bison 3.0.0 or later is required to generate PHP parsers (excluded versions: none).
localhost:php-src didi$ git branch

brew install bison

checking for bison version... ./configure: line 5574: -z: command not found
3.6.2 (ok)
checking for re2c... no
configure: error: re2c 0.13.4 is required to generate PHP lexers.

brew install re2c


生成了Zend/zend_config.h

/Users/didi/PhpstormProjects/c/php-src/Zend/zend_config.w32.h:23:10: fatal error: '../main/config.w32.h' file not found
#include <../main/config.w32.h>

In file included from /usr/include/php/main/php.h:35:

/usr/include/php/Zend/zend.h:51:11: fatal error: 'zend_config.h' file not found

# include <zend_config.h>

          ^

1 error generated.

make: *** [redis.lo] Error 1
两个问题都解决了

cd ext/alae/cmake-build-debug/

make

/Users/didi/PhpstormProjects/c/php-src/Zend/zend_virtual_cwd.h:22:10: fatal error: 'TSRM.h' file not found
#include "TSRM.h"
<!-- more -->
https://l1905.github.io/php/2020/02/28/macos-pecl-xdebug-mongodb/

CMakefileList.txt加上

include_directories(${PHP_SOURCE}/TSRM)
问题解决

ld: symbol(s) not found for architecture x86_64
clang: error: linker command failed with exit code 1 (use -v to see invocation)
make[2]: *** [alae] Error 1
make[1]: *** [CMakeFiles/alae.dir/all] Error 2
make: *** [all] Error 2

https://blog.51cto.com/peterxu/1795036

SET(CMAKE_C_COMPILER g++)

SET(CMAKE_MODULE_LINKER_FLAGS "-lstdc++.8")

SET(CMAKE_SHARED_LINKER_FLAGS "-lstdc++.8")
SET(CMAKE_EXE_LINKER_FLAGS "-lstdc++.8")
SET(CMAKE_STATIC_LINKER_FLAGS "-lstdc++.8")

https://blog.csdn.net/esrrhs/article/details/52700332
https://blog.csdn.net/m0_38130105/article/details/84234774

32 warnings and 5 errors generated.
make[2]: *** [CMakeFiles/alae.dir/alae.c.o] Error 1
make[1]: *** [CMakeFiles/alae.dir/all] Error 2
make: *** [all] Error 2

/Users/didi/PhpstormProjects/c/php-src/ext/alae/alae.c:624:21: error: use of undeclared identifier 'EX_CONSTANT'

7.2不再支持7.1的这个宏
{% raw %}
/* constant in currently executed function */
#define EX_CONSTANT(node) \
	RT_CONSTANT_EX(EX_LITERALS(), node)

/* run-time constant */
# define RT_CONSTANT_EX(base, node) \
	((zval*)(((char*)(base)) + (node).constant))
	

# define EX_LITERALS() \
	EX(literals)
	
 error: no member named 'literals' in '_zend_execute_data'
    dim = EX_CONSTANT(opline->op2);

	
1、zend_execute_data:opcode执行期间非常重要的一个结构，记录着当前执行的zend_op、返回值、所属函数/对象指针、符号表等
struct _zend_execute_data {
    const zend_op       *opline;           /* executed opline 指向第一条opcode */
    zend_execute_data   *call;             /* current call                   */
    zval                *return_value;
    zend_function       *func;             /* executed op_array              */
    zval                 This;
#if ZEND_EX_USE_RUN_TIME_CACHE
    void               **run_time_cache;
#endif
#if ZEND_EX_USE_LITERALS
    zval                *literals;
#endif
    zend_class_entry    *called_scope;
    zend_execute_data   *prev_execute_data;
    zend_array          *symbol_table;
};

2、zend_op:zend指令
//zend.compile.h
struct _zend_op {
    const void *handler;  //该指令调用的处理函数
    znode_op op1; //操作数1
    znode_op op2; //操作数2
    znode_op result; 
    uint32_t extended_value;
    uint32_t lineno;
    zend_uchar opcode; //opcode指令编号
    zend_uchar op1_type; //操作数1类型
    zend_uchar op2_type; 
    zend_uchar result_type;
};

7.1有这个结构7.2没有了

# define EX_LITERALS()

error: cannot initialize a parameter of type 'zend_object *'
      (aka '_zend_object *') with an lvalue of type 'zval *' (aka '_zval_struct *')

https://segmentfault.com/a/1190000004340427
https://wiki.jikexueyuan.com/project/extending-embedding-php/2.html


在PHP源码中，我们可以见到诸如PHPAPI ZEND_API TSRM_API等xxx_API(当然还有其他格式的)这样的宏

关于它们的定义都是类似于

#if defined(__GNUC__) && __GNUC__ >= 4
# define ZEND_API __attribute__ ((visibility("default")))
#else
# define ZEND_API
#endif
 

一、预定义__GNUC__宏

    1 __GNUC__ 是gcc编译器编译代码时预定义的一个宏。需要针对gcc编写代码时， 可以使用该宏进行条件编译。

    2 __GNUC__ 的值表示gcc的版本。需要针对gcc特定版本编写代码时，也可以使用该宏进行条件编译。

    3 __GNUC__ 的类型是“int”，该宏被扩展后， 得到的是整数字面值。可以通过仅预处理，查看宏扩展后的文本。

所以我们知道ZEND_API定义为：

如果编译器使用的是gcc且GNUC的版本大于等于4,则定义ZEND_API为 __attribute__ ((visibility("default")))

那__attribute__到底是干嘛的，有什么作用呢？

查阅关于C的相关资料得出结论:

__attribute__ ((visibility("default")))定义的函数都是可见的

GCC 有个visibility属性, 该属性是说, 启用这个属性:

1. 当-fvisibility=hidden时

动态库中的函数默认是被隐藏的即 hidden. 除非显示声明为__attribute__((visibility("default"))).

2. 当-fvisibility=default时

 动态库中的函数默认是可见的.除非显示声明为__attribute__((visibility("hidden"))).

$ /Library/Developer/CommandLineTools/usr/bin/ld -v
@(#)PROGRAM:ld  PROJECT:ld64-274.2
configured to support archs: armv6 armv7 armv7s arm64 i386 x86_64 x86_64h armv6m armv7k armv7m armv7em (tvOS)
LTO support using: LLVM version 8.0.0, (clang-800.0.42.1)
TAPI support using: Apple TAPI version 1.30
{% endraw %}


