---
title: php execute_ex
layout: post
category: php
author: 夏泽民
---
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

static void (*old_zend_execute_ex) (zend_execute_data *execute_data TSRMLS_DC);

ZEND_API void learn_execute_ex (zend_execute_data *execute_data TSRMLS_DC)
{
    php_printf("====== extension debug start ======\n");
    php_printf("function name: %s\n", get_active_function_name(TSRMLS_C));

    old_zend_execute_ex(execute_data TSRMLS_CC);

    int stacked = 0;
    void **top;
    void **bottom;
    zval *arguments;
    smart_str buf = {0};

    array_init(arguments);
    
    top = zend_vm_stack_top(TSRMLS_C) - 1;
    if (top) {
        stacked = (int)(zend_uintptr_t) *top; // argc
        if (stacked) {
            bottom = zend_vm_stack_top(TSRMLS_C);
            EG(argument_stack)->top = top + 1;
            if (zend_copy_parameters_array(stacked, arguments TSRMLS_CC) == SUCCESS) {
                php_json_encode(&buf, arguments, 0 TSRMLS_CC);
            }
            EG(argument_stack)->top = bottom;
        }
    }

    smart_str_0(&buf);

    php_printf("%s\n", buf.c);

    smart_str_free(&buf);
    zval_dtor(arguments);

    php_printf("====== extension debug end ======\n");
}

PHP_MINIT_FUNCTION(learn)
{
    old_zend_execute_ex = zend_execute_ex;
    zend_execute_ex = learn_execute_ex;

    return SUCCESS;
}

PHP_MSHUTDOWN_FUNCTION(learn)
{
    zend_execute_ex = old_zend_execute_ex;

    return SUCCESS;
}
2015-11-04 00:38 更新
后来看到，其实不用上面这中方法就可以实现, php 5.5之后要从 prev 里面去取

/**
 * php_var_dump defined in this head file.
 */
#include "ext/standard/php_var.h"

zend_execute_data *real_execute_data = execute_data->prev_execute_data;

void **p = real_execute_data->function_state.arguments;
int arg_count = (int) (zend_uintptr_t) * p;
zval *argument_element;
int i;
// zval *obj = real_execute_data->object;
unsigned long start = mach_absolute_time();
for (i = 0; i < arg_count; i++) {
    argument_element = *(p - (arg_count - i));
    php_var_dump(&argument_element, 1);
}
<!-- more -->
字节码在 Zend 虚拟机中的解释执行 之 概述

execute_ex
我们来看看执行一个简单的脚本 test.php 的调用栈

execute_ex @ zend_vm_execute.h : 411
zend_execute @ zend_vm_execute.h : 474
php_execute_script @ zend.c : 1474
do_cli @ php_cli.c : 993
main @ php_cli.c : 1381 
由于是执行脚本文件，所以 do_cli 调用了 php_execute_script 函数，最终调用 execute_ex 函数：

 
ZEND_API void execute_ex(zend_execute_data *ex)
{
    DCL_OPLINE

#ifdef ZEND_VM_IP_GLOBAL_REG
    const zend_op *orig_opline = opline;
#endif
#ifdef ZEND_VM_FP_GLOBAL_REG
    zend_execute_data *orig_execute_data = execute_data;
    execute_data = ex;
#else
    zend_execute_data *execute_data = ex;
#endif


    LOAD_OPLINE();
    ZEND_VM_LOOP_INTERRUPT_CHECK();

    while (1) {
#if !defined(ZEND_VM_FP_GLOBAL_REG) || !defined(ZEND_VM_IP_GLOBAL_REG)
            int ret;
#endif
#if defined(ZEND_VM_FP_GLOBAL_REG) && defined(ZEND_VM_IP_GLOBAL_REG)
        ((opcode_handler_t)OPLINE->handler)(ZEND_OPCODE_HANDLER_ARGS_PASSTHRU);
        if (UNEXPECTED(!OPLINE)) {
#else
        if (UNEXPECTED((ret = ((opcode_handler_t)OPLINE->handler)www.90168.org(ZEND_OPCODE_HANDLER_ARGS_PASSTHRU)) != 0)) {
#endif
#ifdef ZEND_VM_FP_GLOBAL_REG
            execute_data = orig_execute_data;
# ifdef ZEND_VM_IP_GLOBAL_REG
            opline = orig_opline;
# endif
            return;
#else
            if (EXPECTED(ret > 0)) {
                execute_data = EG(current_execute_data);
                ZEND_VM_LOOP_INTERRUPT_CHECK();
            } else {
# ifdef ZEND_VM_IP_GLOBAL_REG
                opline = orig_opline;
# endif
                return;
            }
#endif
        }

    }
    zend_error_noreturn(E_CORE_ERROR, "Arrived at end of main loop which shouldn't happen");
}
和其它 C 语言编写的系统软件类似，函数中使用了大量的宏定义，通过宏定义的名字还是能大概看出其用途

DCL_OPLINE，变量声明

LOAD_OPLINE()，加载指令字节码

ZEND_VM_LOOP_INTERRUPT_CHECK()，interrupt 检测


解释器引擎最终执行op的函数是zend_execute，实际上zend_execute是一个函数指针，在引擎初始化的时候zend_execute默认指向了execute,这个execute定义在{PHPSRC}/Zend/zend_vm_execute.h：

 
  Zend引擎主要包含两个核心部分：编译、执行：

                           

    执行阶段主要用到的数据结构：

          opcode： php代码编译产生的zend虚拟机可识别的指令，php7有173个opcode，定义在 zend_vm_opcodes.hPHP中的所有语法实现都是由这些opcode组成的。

         

复制代码
struct _zend_op {
    const void *handler; //对应执行的C语言function，即每条opcode都有一个C function处理
    znode_op op1;   //操作数1
    znode_op op2;   //操作数2
    znode_op result; //返回值
    uint32_t extended_value; 
    uint32_t lineno; 
    zend_uchar opcode;  //opcode指令
    zend_uchar op1_type; //操作数1类型
    zend_uchar op2_type; //操作数2类型
    zend_uchar result_type; //返回值类型
};
复制代码
         zend_op_array : zend引擎执行阶段的输入数据结构，整个执行阶段都是操作这个数据结构。

             

                            

 

 

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
                              
PHP xhprof 扩展原理
由于公司项目，最近需要分析后端PHP接口的性能数据，就采用了FACEBOOK之前开源的一个扩展，现在市面上很多分支都是基于FB最开始的线开发的，但是由于FB已经停止维护，所以现在其他线都是自己个人在维护。

今天我分析的这个分支是兼容PHP7+的版本收集性能数据，先贴出GITHUB的链接  https://github.com/longxinH/xhprof

扩展的安装方式和PHP调用的API用法在github上readme.md上有详细说明，可以参考，我现在做的项目就是用这个来收集的，因为是私人维护的，不敢保证以后会不会启用，所以先自己了解下XHPROF扩展源码，防止后面维护人不维护，可以自己在继续维护起来。

 

分析：

起始外部扩展就是相当于内核一个模块，都是zend_module_entry结构体, 收集性能数据的原理，就是在模块初始化的时候代理了这个，代理了这个编译和执行OPCODE 的函数，覆盖了一层，加了自己的处理，自己的处理就是C代码的执行时间和PHP申请堆内存的计算。

/* Replace zend_compile with our proxy */
 
_zend_compile_file = zend_compile_file;
 
zend_compile_file = hp_compile_file;
 
 
 
/* Replace zend_compile_string with our proxy */
 
_zend_compile_string = zend_compile_string;
 
zend_compile_string = hp_compile_string;
 
 
 
/* Replace zend_execute with our proxy */
 
_zend_execute_ex = zend_execute_ex;
 
zend_execute_ex = hp_execute_ex;
这是zend_module_entry，第三方扩展都需要实现这个结构体，用来保存扩展的基本信息，类似设计模式里面的实现接口
/* Callback functions for the xhprof extension */
zend_module_entry xhprof_module_entry = {
#if ZEND_MODULE_API_NO >= 20010901
        STANDARD_MODULE_HEADER,
#endif
        "xhprof",                        /* Name of the extension */
        xhprof_functions,                /* List of functions exposed */
        PHP_MINIT(xhprof),               /* Module init callback */
        PHP_MSHUTDOWN(xhprof),           /* Module shutdown callback */
        PHP_RINIT(xhprof),               /* Request init callback */
        PHP_RSHUTDOWN(xhprof),           /* Request shutdown callback */
        PHP_MINFO(xhprof),               /* Module info callback */
#if ZEND_MODULE_API_NO >= 20010901
        XHPROF_VERSION,
#endif
        STANDARD_MODULE_PROPERTIES
};
 

/* Xhprof's global state.
 *
 * This structure is instantiated once.  Initialize defaults for attributes in
 * hp_init_profiler_state() Cleanup/free attributes in
 * hp_clean_profiler_state() */
存储扩展性能数据，是否开启，等信息的结构体
ZEND_BEGIN_MODULE_GLOBALS(xhprof)
 
    /*       ----------   Global attributes:  -----------       */
 
    /* Indicates if xhprof is currently enabled */
    int              enabled;
 
    /* Indicates if xhprof was ever enabled during this request */
    int              ever_enabled;
 
    /* Holds all the xhprof statistics */
    zval            stats_count;
 
    /* Indicates the current xhprof mode or level */
    int              profiler_level;
 
    /* Top of the profile stack */
    hp_entry_t      *entries;
 
    /* freelist of hp_entry_t chunks for reuse... */
    hp_entry_t      *entry_free_list;
 
    /* Callbacks for various xhprof modes */
    hp_mode_cb       mode_cb;
 
    /*       ----------   Mode specific attributes:  -----------       */
 
    /* Global to track the time of the last sample in time and ticks */
    struct timeval   last_sample_time;
    uint64           last_sample_tsc;
    /* XHPROF_SAMPLING_INTERVAL in ticks */
    long             sampling_interval;
    uint64           sampling_interval_tsc;
    int              sampling_depth;
    /* XHProf flags */
    uint32 xhprof_flags;
 
    char *root;
 
    /* counter table indexed by hash value of function names. */
    uint8  func_hash_counters[256];
 
    HashTable *trace_callbacks;
 
    /* Table of ignored function names and their filter */
    hp_ignored_functions *ignored_functions;
 
ZEND_END_MODULE_GLOBALS(xhprof)
/**
static const zend_ini_entry_def ini_entries[] = {
    { xhprof.output_dir, NULL, arg1, arg2, arg3, '', displayer, 7, sizeof(name)-1, sizeof('')-1 },
    { xhprof.collect_additional_info, NULL, arg1, arg2, arg3, '0', displayer, 7, sizeof(name)-1, sizeof('')-1 },
    { xhprof.sampling_interval, NULL, arg1, arg2, arg3, '100000', displayer, 7, sizeof(name)-1, sizeof('')-1 },
    { xhprof.sampling_depth, NULL, arg1, arg2, arg3, '100000', displayer, 7, sizeof(name)-1, sizeof('')-1 },
    { NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, 0, 0} 
 
}
 */
这是注册xhprof扩展所有的php.ini配置项
PHP_INI_BEGIN()
 
/* output directory:
 * Currently this is not used by the extension itself.
 * But some implementations of iXHProfRuns interface might
 * choose to save/restore XHProf profiler runs in the
 * directory specified by this ini setting.
 */
PHP_INI_ENTRY("xhprof.output_dir", "", PHP_INI_ALL, NULL)
 
/*
 * collect_additional_info
 * Collect mysql_query, curl_exec internal info. The default is 0.
 */
PHP_INI_ENTRY("xhprof.collect_additional_info", "0", PHP_INI_ALL, NULL)
 
/* sampling_interval:
 * Sampling interval to be used by the sampling profiler, in microseconds.
 */
#define STRINGIFY_(X) #X
#define STRINGIFY(X) STRINGIFY_(X)
 
STD_PHP_INI_ENTRY("xhprof.sampling_interval", STRINGIFY(XHPROF_DEFAULT_SAMPLING_INTERVAL), PHP_INI_ALL, OnUpdateLong, sampling_interval, zend_xhprof_globals, xhprof_globals)
 
/* sampling_depth:
 * Depth to trace call-chain by the sampling profiler
 */
STD_PHP_INI_ENTRY("xhprof.sampling_depth", STRINGIFY(INT_MAX), PHP_INI_ALL, OnUpdateLong, sampling_depth, zend_xhprof_globals, xhprof_globals)
PHP_INI_END()
/**
 * Module init callback.
 *
 * @author cjiang
 */
模块初始化的调用的函数
//int zm_startup_xhprof(xhprof)(int type, int module_number)
PHP_MINIT_FUNCTION(xhprof)
{
    ZEND_INIT_MODULE_GLOBALS(xhprof, php_xhprof_init_globals, NULL);
 
    REGISTER_INI_ENTRIES();
 
    hp_register_constants(INIT_FUNC_ARGS_PASSTHRU);
 
    /* Replace zend_compile with our proxy */
    _zend_compile_file = zend_compile_file;
    zend_compile_file  = hp_compile_file;
 
    /* Replace zend_compile_string with our proxy */
    _zend_compile_string = zend_compile_string;
    zend_compile_string = hp_compile_string;
 
    /* Replace zend_execute with our proxy */
    _zend_execute_ex = zend_execute_ex;
    zend_execute_ex  = hp_execute_ex;
 
    /* Replace zend_execute_internal with our proxy */
    _zend_execute_internal = zend_execute_internal;
    zend_execute_internal = hp_execute_internal;
 
#if defined(DEBUG)
    /* To make it random number generator repeatable to ease testing. */
    srand(0);
#endif
    return SUCCESS;
}
 
/**
 * Module shutdown callback.
 */
模块关闭的调用的函数
//int zm_shutdown_xhprof(xhprof)(int type, int module_number)
PHP_MSHUTDOWN_FUNCTION(xhprof)
{
    /* free any remaining items in the free list */
    hp_free_the_free_list();
 
    /* Remove proxies, restore the originals */
    zend_execute_ex       = _zend_execute_ex;
    zend_execute_internal = _zend_execute_internal;
    zend_compile_file     = _zend_compile_file;
    zend_compile_string   = _zend_compile_string;
 
    UNREGISTER_INI_ENTRIES();
 
    return SUCCESS;
}
 
/**
 * Request init callback. Nothing to do yet!
 */
请求初始化的调用的函数
//int zm_activate_xhprof(xhprof)(int type, int module_number)
PHP_RINIT_FUNCTION(xhprof)
{
#if defined(ZTS) && defined(COMPILE_DL_XHPROF)
    ZEND_TSRMLS_CACHE_UPDATE();
#endif
 
    return SUCCESS;
}
 
/**
 * Request shutdown callback. Stop profiling and return.
 */
//int zm_deactivate_xhprof(xhprof)(int type, int module_number)
请求完成的调用的函数
PHP_RSHUTDOWN_FUNCTION(xhprof)
{
    hp_end();
    return SUCCESS;
}
 
/**
 * Module info callback. Returns the xhprof version.
 */
phpinfo显示模块基本信息调用的函数
//int zm_info_xhprof(xhprof)(zend_module_entry *zend_module)
PHP_MINFO_FUNCTION(xhprof)
{
    php_info_print_table_start();
    php_info_print_table_header(2, "xhprof support", "enabled");
    php_info_print_table_row(2, "Version", XHPROF_VERSION);
    php_info_print_table_end();
    DISPLAY_INI_ENTRIES();
}
 
/**
 * Start XHProf profiling in hierarchical mode.
 *
 * @param  long $flags  flags for hierarchical mode
 * @return void
 * @author kannan
 */
PHP语言中 xhprof_enable函数
//zif_xhprof_enable(zend_execute_data *execute_data, zval *return_value)
PHP_FUNCTION(xhprof_enable)
{
    long  xhprof_flags = 0;              /* XHProf flags */
    zval *optional_array = NULL;         /* optional array arg: for future use */
 
    if (zend_parse_parameters(ZEND_NUM_ARGS(), "|lz", &xhprof_flags, &optional_array) == FAILURE) {
        return;
    }
 
    hp_get_ignored_functions_from_arg(optional_array);
 
    hp_begin(XHPROF_MODE_HIERARCHICAL, xhprof_flags);
}
 
/**
 * Stops XHProf from profiling in hierarchical mode anymore and returns the
 * profile info.
 *
 * @param  void
 * @return array  hash-array of XHProf's profile info
 * @author kannan, hzhao
 */
PHP语言中 xhprof_disable函数
//zif_xhprof_disable(zend_execute_data *execute_data, zval *return_value)
PHP_FUNCTION(xhprof_disable)
{
    if (XHPROF_G(enabled)) {
        hp_stop();
        RETURN_ZVAL(&XHPROF_G(stats_count), 1, 0);
    }
    /* else null is returned */
}
 
/**
 * Start XHProf profiling in sampling mode.
 *
 * @return void
 * @author cjiang
 */
PHP语言中 xhprof_sample_enable函数
//zif_xhprof_sample_enable(zend_execute_data *execute_data, zval *return_value)
PHP_FUNCTION(xhprof_sample_enable)
{
    long xhprof_flags = 0;    /* XHProf flags */
    hp_get_ignored_functions_from_arg(NULL);
    hp_begin(XHPROF_MODE_SAMPLED, xhprof_flags);
}
 
/**
 * Stops XHProf from profiling in sampling mode anymore and returns the profile
 * info.
 *
 * @param  void
 * @return array  hash-array of XHProf's profile info
 * @author cjiang
 */
PHP语言中 xhprof_sample_disable函数
//zif_xhprof_sample_disable(zend_execute_data *execute_data, zval *return_value)
PHP_FUNCTION(xhprof_sample_disable)
{
    if (XHPROF_G(enabled)) {
        hp_stop();
        RETURN_ZVAL(&XHPROF_G(stats_count), 1, 0);
    }
  /* else null is returned */
}
这个扩展的基本接口就是上面这些代码，接下来用一个例子来详细解释一下，怎么收集性能数据的。

贴出一段PHP代码

<?php
//phpinfo();
 
xhprof_enable(XHPROF_FLAGS_NO_BUILTINS  | XHPROF_FLAGS_CPU | XHPROF_FLAGS_MEMORY);
 
aaa();
 
$xhprof_data = xhprof_disable();
 
print_r($xhprof_data);
 
function aaa()
{
	bbb();
}
 
function bbb()
{
	
}
收集数据的过程其实针对函数嵌套调用递归收集的过程，xhprof_enable开启的这一行解析成 OPCODE，然后执行之前代理的hp_execute_ex函数，初始化main节点，下面就以图的形式来看下具体生成的结构，和调用链。



在xhprof_enable 和  xhprof_disable 之前的PHP代码，每一行都会被翻译成OPCODE，都会执行hp_execute_ex函数。

所以每次如果执行的是函数就会收集对应的性能数据
ZEND_DLEXPORT void hp_execute_ex (zend_execute_data *execute_data)
{
    if (!XHPROF_G(enabled)) {
        _zend_execute_ex(execute_data);
        return;
    }
 
    char *func = NULL;
    int hp_profile_flag = 1;
 
    func = hp_get_function_name(execute_data);
 
    if (!func) {
        _zend_execute_ex(execute_data);
        return;
    }
 
    zend_execute_data *real_execute_data = execute_data->prev_execute_data;
 
    BEGIN_PROFILING(&XHPROF_G(entries), func, hp_profile_flag, real_execute_data);
 
    _zend_execute_ex(execute_data);
 
    if (XHPROF_G(entries)) {
        END_PROFILING(&XHPROF_G(entries), hp_profile_flag);
    }
 
    efree(func);
}
输出：

Array
(
    [aaa==>bbb] => Array
        (
            [ct] => 1
            [wt] => 23
            [cpu] => 29
            [mu] => 832
            [pmu] => 0
        )
 
    [main()==>aaa] => Array
        (
            [ct] => 1
            [wt] => 88
            [cpu] => 89
            [mu] => 1408
            [pmu] => 0
        )
 
    [main()] => Array
        (
            [ct] => 1
            [wt] => 99
            [cpu] => 99
            [mu] => 1976
            [pmu] => 0
        )
 
)
xhprof_enable(XHPROF_FLAGS_NO_BUILTINS  | XHPROF_FLAGS_CPU | XHPROF_FLAGS_MEMORY);
这个函数相当于执行了扩展里面的，这个函数相当自己造了一个虚拟的 main()节点，
来存当前调用PHP函数的性能数据  ct wt mu pmu等数据
 
PHP_FUNCTION(xhprof_enable)
{
    long  xhprof_flags = 0;              /* XHProf flags */
    zval *optional_array = NULL;         /* optional array arg: for future use */
 
    if (zend_parse_parameters(ZEND_NUM_ARGS(), "|lz", &xhprof_flags, &optional_array) == FAILURE) {
        return;
    }
 
    hp_get_ignored_functions_from_arg(optional_array);
    //以main为根节点存储性能数据
    hp_begin(XHPROF_MODE_HIERARCHICAL, xhprof_flags);
}
 
static void hp_begin(long level, long xhprof_flags)
{
    if (!XHPROF_G(enabled)) {
        int hp_profile_flag = 1;
 
        XHPROF_G(enabled)      = 1;
        XHPROF_G(xhprof_flags) = (uint32)xhprof_flags;
        //先绑定每个函数开始，结束的回调函数，用来记录当时的性能点数据
        /* Initialize with the dummy mode first Having these dummy callbacks saves
         * us from checking if any of the callbacks are NULL everywhere. */
        XHPROF_G(mode_cb).init_cb     = hp_mode_dummy_init_cb;
        XHPROF_G(mode_cb).exit_cb     = hp_mode_dummy_exit_cb;
        XHPROF_G(mode_cb).begin_fn_cb = hp_mode_dummy_beginfn_cb;
        XHPROF_G(mode_cb).end_fn_cb   = hp_mode_dummy_endfn_cb;
 
        /* Register the appropriate callback functions Override just a subset of
        * all the callbacks is OK. */
        switch (level) {
            case XHPROF_MODE_HIERARCHICAL:
                XHPROF_G(mode_cb).begin_fn_cb = hp_mode_hier_beginfn_cb;
                XHPROF_G(mode_cb).end_fn_cb   = hp_mode_hier_endfn_cb;
                break;
            case XHPROF_MODE_SAMPLED:
                XHPROF_G(mode_cb).init_cb     = hp_mode_sampled_init_cb;
                XHPROF_G(mode_cb).begin_fn_cb = hp_mode_sampled_beginfn_cb;
                XHPROF_G(mode_cb).end_fn_cb   = hp_mode_sampled_endfn_cb;
                break;
        }
 
        /* one time initializations */
        hp_init_profiler_state(level);
 
        /* start profiling from fictitious main() */
        XHPROF_G(root) = estrdup(ROOT_SYMBOL);
 
        /* start profiling from fictitious main() */
        BEGIN_PROFILING(&XHPROF_G(entries), XHPROF_G(root), hp_profile_flag, NULL);
    }
}
 
 
#define BEGIN_PROFILING(entries, symbol, profile_curr, execute_data)        \
do {                                                                     \
    /* Use a hash code to filter most of the string comparisons. */     \
    uint8 hash_code  = hp_inline_hash(symbol);                          \
    profile_curr = !hp_ignore_entry_work(hash_code, symbol);                 \
    if (profile_curr) {                                                 \
        if (execute_data != NULL) {                                     \
            symbol = hp_get_trace_callback(symbol, execute_data); \
        }                                                               \
        hp_entry_t *cur_entry = hp_fast_alloc_hprof_entry();            \
        (cur_entry)->hash_code = hash_code;                             \
        (cur_entry)->name_hprof = symbol;                               \
        (cur_entry)->prev_hprof = (*(entries));                         \
        /* Call the universal callback */                               \
        hp_mode_common_beginfn((entries), (cur_entry));                 \
        /* Call the mode's beginfn callback */                          \
        XHPROF_G(mode_cb).begin_fn_cb((entries), (cur_entry));         \
        /* Update entries linked list */                                \
        (*(entries)) = (cur_entry);                                     \
    }                                                               \
} while (0)
 
begin_fn_cb
||   这个就是记录当前函数开始，记录的CPU和内存大小
void hp_mode_hier_beginfn_cb(hp_entry_t **entries, hp_entry_t  *current)
{
    /* Get start tsc counter */
    current->tsc_start = cycle_timer();
 
    /* Get CPU usage */
    if (XHPROF_G(xhprof_flags) & XHPROF_FLAGS_CPU) {
        current->cpu_start = cpu_timer();
    }
 
    /* Get memory usage */
    if (XHPROF_G(xhprof_flags) & XHPROF_FLAGS_MEMORY) {
        current->mu_start_hprof  = zend_memory_usage(0);
        current->pmu_start_hprof = zend_memory_peak_usage(0);
    }
}
 

$xhprof_data = xhprof_disable(); 结束收集数据
PHP这个函数对应的就是扩展里下面的这个函数，是用来停止收集，并且RETURN_ZVAL(&XHPROF_G(stats_count), 1, 0); 
返回已经收集到的数据，
数据格式在内核里面是一个二级HASHTABLE，对应PHP是一个二维数组
 
这个结束过程是对应这之前xhprof_enable的执行过程，是对main虚拟的节点做回源处理。
//zif_xhprof_disable(zend_execute_data *execute_data, zval *return_value)
PHP_FUNCTION(xhprof_disable)
{
    if (XHPROF_G(enabled)) {
        hp_stop();
        RETURN_ZVAL(&XHPROF_G(stats_count), 1, 0);
    }
    /* else null is returned */
}
 
static void hp_stop()
{
    int   hp_profile_flag = 1;
 
    /* End any unfinished calls */
    while (XHPROF_G(entries)) {
        END_PROFILING(&XHPROF_G(entries), hp_profile_flag);
    }
 
    if (XHPROF_G(root)) {
        efree(XHPROF_G(root));
        XHPROF_G(root) = NULL;
    }
 
    /* Stop profiling */
    XHPROF_G(enabled) = 0;
}
GDB调试的过程

(gdb) p *xhprof_globals.entries
$33 = {name_hprof = 0x7ffff1279230 "bbb", rlvl_hprof = 0, tsc_start = 594813346976, cpu_start = 56435, mu_start_hprof = 390296, pmu_start_hprof = 427056, prev_hprof = 0x1aad3b0, hash_code = 224 '\340'}
(gdb) p *xhprof_globals.entries->prev_hprof
$34 = {name_hprof = 0x7ffff12791e0 "aaa", rlvl_hprof = 0, tsc_start = 594308864680, cpu_start = 55591, mu_start_hprof = 390216, pmu_start_hprof = 427056, prev_hprof = 0x1aad360, hash_code = 104 'h'}
(gdb) p *xhprof_globals.entries->prev_hprof->prev_hprof
$35 = {name_hprof = 0x7ffff1202058 "main()", rlvl_hprof = 0, tsc_start = 594156065161, cpu_start = 54939, mu_start_hprof = 390136, pmu_start_hprof = 427056, prev_hprof = 0x0, hash_code = 32 ' '}


(gdb) p *xhprof_globals.entries->prev_hprof->prev_hprof->prev_hprof