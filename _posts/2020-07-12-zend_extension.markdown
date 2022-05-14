---
title: zend_extension
layout: post
category: php
author: 夏泽民
---
http://blog.mallol.cn/2018/31c7372c.html
http://yangxikun.github.io/php/2016/07/10/php-zend-extension.html
根据 PHP 版本，zend_extension 指令可以是以下之一：

 

zend_extension (non ZTS, non debug build)
zend_extension_ts ( ZTS, non debug build)
zend_extension_debug (non ZTS, debug build)
zend_extension_debug_ts ( ZTS, debug build)

ZTS：ZEND Thread Safety

可通过phpinfo()查看ZTS是否启用，从而决定用
zend_extension还是zend_extension_ts，当然试一下怎么生效也可以。
https://blog.csdn.net/taft/article/details/596291

如果你要加载xdebug扩展，就需要用zend_extension，那么zend_extension和extension加载扩展有啥区别？

解释：

一般的模块，或者说大部分的模块，都是用extension加载的。
极少数和zend内核关系紧密的需要用zend_extension 进行加载，如xdebug、 opcache...等等
一般的mysql都属于extension 而 像xdebug却是zend_extension
错误提示

NOTICE: PHP message: PHP Warning: Xdebug MUST be loaded as a Zend extension

分析：
mysql扩展配置如下
extension=mysql.so  即可
xdebug 扩展如下  注意要添加红色的路径 否则就会出错
并且还有开启html_errors=on
[Xdebug]
zend_extension = /usr/lib64/php/modules/xdebug.so
<!-- more -->
通常在php.ini中，通过extension=*加载的扩展我们称为PHP扩展，通过zend_extension=*加载的扩展我们称为Zend扩展，但从源码的角度来讲，PHP扩展应该称为“模块”（源码中以module命名），而Zend扩展称为“扩展”（源码中以extension命名）。

两者最大的区别在于向引擎注册的钩子。少数的扩展，例如xdebug、opcache，既是PHP扩展，也是Zend扩展，但它们在php.ini中的加载方式得用zend_extension=*，具体原因下文会说明。

先来看看两种类型扩展主要的数据结构：

PHP扩展的结构体
{% raw %}
typedef struct _zend_module_entry zend_module_entry;
typedef struct _zend_module_dep zend_module_dep;
 
struct _zend_module_entry {
    unsigned short size;
    unsigned int zend_api;
    unsigned char zend_debug;
    unsigned char zts;
    const struct _zend_ini_entry *ini_entry;
    const struct _zend_module_dep *deps;
    const char *name;
    const struct _zend_function_entry *functions;          /* PHP Functions */
    int (*module_startup_func)(INIT_FUNC_ARGS);             /* MINIT */
    int (*module_shutdown_func)(SHUTDOWN_FUNC_ARGS);        /* MSHUTDOWN */
    int (*request_startup_func)(INIT_FUNC_ARGS);            /* RINIT */
    int (*request_shutdown_func)(SHUTDOWN_FUNC_ARGS);       /* RSHUTDOWN */
    void (*info_func)(ZEND_MODULE_INFO_FUNC_ARGS);
    const char *version;
    size_t globals_size;
#ifdef ZTS
    ts_rsrc_id* globals_id_ptr;
#else
    void* globals_ptr;
#endif
    void (*globals_ctor)(void *global TSRMLS_DC);
    void (*globals_dtor)(void *global TSRMLS_DC);
    int (*post_deactivate_func)(void);
    int module_started;
    unsigned char type;
    void *handle;
    int module_number;
    const char *build_id;
};
 
struct _zend_module_dep {
    const char *name;       /* module name */
    const char *rel;        /* version relationship: NULL (exists), lt|le|eq|ge|gt (to given version) */
    const char *version;    /* version */
    unsigned char type;     /* dependency type */
};
Zend扩展的结构体：

/* Typedef's for zend_extension function pointers */
typedef int (*startup_func_t)(zend_extension *extension);
typedef void (*shutdown_func_t)(zend_extension *extension);
typedef void (*activate_func_t)(void);
typedef void (*deactivate_func_t)(void);
 
typedef void (*message_handler_func_t)(int message, void *arg);
 
typedef void (*op_array_handler_func_t)(zend_op_array *op_array);
 
typedef void (*statement_handler_func_t)(zend_op_array *op_array);
typedef void (*fcall_begin_handler_func_t)(zend_op_array *op_array);
typedef void (*fcall_end_handler_func_t)(zend_op_array *op_array);
 
typedef void (*op_array_ctor_func_t)(zend_op_array *op_array);
typedef void (*op_array_dtor_func_t)(zend_op_array *op_array);
 
typedef struct _zend_extension {
    char *name;
    char *version;
    char *author;
    char *URL;
    char *copyright;
 
    startup_func_t startup;
    shutdown_func_t shutdown;
    activate_func_t activate;
    deactivate_func_t deactivate;
 
    message_handler_func_t message_handler;
 
    op_array_handler_func_t op_array_handler;
 
    statement_handler_func_t statement_handler;
    fcall_begin_handler_func_t fcall_begin_handler;
    fcall_end_handler_func_t fcall_end_handler;
 
    op_array_ctor_func_t op_array_ctor;
    op_array_dtor_func_t op_array_dtor;
 
    int (*api_no_check)(int api_no);
    int (*build_id_check)(const char* build_id);
    void *reserved3;
    void *reserved4;
    void *reserved5;
    void *reserved6;
    void *reserved7;
    void *reserved8;
 
    DL_HANDLE handle;
    int resource_number;
} zend_extension;
 
typedef struct _zend_extension_version_info {
    int zend_extension_api_no;
    char *build_id;
} zend_extension_version_info;
以xdebug 2.4.0为例，来说明为何其既是PHP扩展，也是Zend扩展。

首先看xdebug.c的2710行处，定义了一个Zend扩展的结构体：

ZEND_DLEXPORT zend_extension zend_extension_entry = {
    XDEBUG_NAME,
    XDEBUG_VERSION,
    XDEBUG_AUTHOR,
    XDEBUG_URL_FAQ,
    XDEBUG_COPYRIGHT_SHORT,
    xdebug_zend_startup,
    xdebug_zend_shutdown,
    NULL,           /* activate_func_t */
    NULL,           /* deactivate_func_t */
    NULL,           /* message_handler_func_t */
    NULL,           /* op_array_handler_func_t */
    xdebug_statement_call, /* statement_handler_func_t */
    NULL,           /* fcall_begin_handler_func_t */
    NULL,           /* fcall_end_handler_func_t */
    xdebug_init_oparray,   /* op_array_ctor_func_t */
    NULL,           /* op_array_dtor_func_t */
    STANDARD_ZEND_EXTENSION_PROPERTIES
};
因为xdebug提供了单步调试、性能优化等高级功能，这些是需要hook到Zend引擎才能做到的，而zend_extension这个结构体就提供了hook到Zend引擎的钩子，例如xdebug_statement_call会在每一条PHP语句执行之后调用。

再看xdebug.c的160行处，定义了一个PHP扩展的结构体：

zend_module_entry xdebug_module_entry = {
    STANDARD_MODULE_HEADER,
    "xdebug",
    xdebug_functions,
    PHP_MINIT(xdebug),
    PHP_MSHUTDOWN(xdebug),
    PHP_RINIT(xdebug),
    PHP_RSHUTDOWN(xdebug),
    PHP_MINFO(xdebug),
    XDEBUG_VERSION,
    NO_MODULE_GLOBALS,
    ZEND_MODULE_POST_ZEND_DEACTIVATE_N(xdebug),
    STANDARD_MODULE_PROPERTIES_EX
}
xdebug提供了许多函数xdebug_enable()、xdebug_disable()、xdebug_call_class()……，这些函数都是定义在xdebug_module_entry中的xdebug_functions。

所以，以我的理解，向用户层面提供一些C实现的PHP函数，需要用到zend_module_entry（即作为PHP扩展），而需要hook到Zend引擎的话，就得用到zend_extension（即作为Zend扩展），xdebug在这里两种都需要。

但是xdebug是通过Zend扩展加载的（zend_extension=*），其PHP扩展部分是如何被加载的呢？其实在Zend扩展加载的时候会调用zend_extension中的startup，即xdebug的xdebug_zend_startup函数，在该函数中：

ZEND_DLEXPORT int xdebug_zend_startup(zend_extension *extension)
{
    /* Hook output handlers (header and output writer) */
    xdebug_hook_output_handlers();

    zend_xdebug_initialised = 1;
                                                                                                                                                                                                               
    return zend_startup_module(&xdebug_module_entry);//对xdebug PHP扩展的加载
}
加载顺序区别
分清PHP扩展和Zend扩展的差异后，接着看看扩展是如何加载的，其加载顺序和依赖是怎么处理的：

相关的主要代码片段在main/main.c中：

int php_module_startup(sapi_module_struct *sf, zend_module_entry *additional_modules, uint num_additional_modules)
{
    . . . . . .

    /* this will read in php.ini, set up the configuration parameters,
       load zend extensions and register php function extensions
       to be loaded later */
    if (php_init_config(TSRMLS_C) == FAILURE) {
        return FAILURE;
    }

    . . . . . .

    /* startup extensions statically compiled in */
    if (php_register_internal_extensions_func(TSRMLS_C) == FAILURE) {
        php_printf("Unable to start builtin modules\n");
        return FAILURE;
    }

    /* start additional PHP extensions */
    php_register_extensions_bc(additional_modules, num_additional_modules TSRMLS_CC);

    /* load and startup extensions compiled as shared objects (aka DLLs)
       as requested by php.ini entries
       theese are loaded after initialization of internal extensions
       as extensions *might* rely on things from ext/standard
       which is always an internal extension and to be initialized
       ahead of all other internals
     */
    php_ini_register_extensions(TSRMLS_C);
    zend_startup_modules(TSRMLS_C);

    /* start Zend extensions */
    zend_startup_extensions();

    . . . . . .
}
1、php_init_config解析php.ini文件，获取需要加载的PHP扩展和Zend扩展

2、php_register_internal_extensions_func加载静态编译的扩展，静态编译进PHP中的扩展，例如date、ereg、pcre等，这些扩展包含在main/internal_functions.c中的zend_module_entry *php_builtin_extensions[]中

3、php_register_extensions_bc注册SAPI的扩展模块，即additional_modules中的扩展，例如Apache SAPI会注册一些与Apache功能相关的扩展，CLI模式下additional_modules为NULL

4、php_ini_register_extensions中会先加载Zend扩展，之后再加载PHP扩展

void php_ini_register_extensions(TSRMLS_D)
{
    zend_llist_apply(&extension_lists.engine, php_load_zend_extension_cb TSRMLS_CC);//加载Zend扩展
    zend_llist_apply(&extension_lists.functions, php_load_php_extension_cb TSRMLS_CC);//加载PHP扩展

    zend_llist_destroy(&extension_lists.engine);
    zend_llist_destroy(&extension_lists.functions);
}
Zend扩展的加载（php_load_zend_extension_cb->zend_load_extension）：

在zend_load_extension中会判断扩展是否合法（能否加载？版本信息是否符合要求？……），如果合法的话将调用zend_register_extension，向全局变量zend_extensions添加当前Zend扩展，并向其他已加载的Zend扩展广播一条消息（实现了message_handler的Zend扩展将接收到）说明自己被加载了，该功能的用处在于，如果某些扩展存在冲突，则扩展能够发出警告信息，并停止加载。

PHP扩展的加载（php_load_php_extension_cb->php_load_extension）：

在php_load_extension中会判断扩展是否合法（能否加载？版本信息是否符合要求？……），如果合法的话将调用zend_register_module_ex，先检查扩展的依赖（_zend_module_entry中的deps），这里只检查当前扩展是否与已加载的扩展冲突，如果没有冲突，则向全局变量module_registry添加当前PHP扩展，接着注册当前PHP扩展的函数（_zend_module_entry中的functions）

这里需要注意的点是deps依赖，该依赖声明了：

当前扩展会与哪些扩展冲突
哪些扩展需要先于当前扩展加载
但是zend_register_module_ex中只对冲突做检查，如果A声明自己与B冲突，但是B却在A之后加载，那么就无法检查到该冲突了。为了解决这个问题，可以在php.ini中排好PHP扩展的顺序（即B放到A之前），加载的时候会按照该顺序加载。

5、扩展初始化阶段：

先激活PHP扩展，在zend_startup_modules中，会先对PHP扩展进行排序（根据每个PHP扩展中deps声明的依赖），然后执行zend_startup_module_ex，调用PHP扩展的MINIT

再激活Zend扩展，在zend_startup_extensions中，对每个Zend扩展调用其startup()

6、请求初始化阶段：

先调用Zend扩展的activate()：php_request_startup->zend_activate->init_executor->zend_extension_activator->activate

再调用PHP扩展的RINIT：php_request_startup->zend_activate_modules->request_startup_func

7、请求结束阶段：

先调用PHP扩展的RSHUTDOWN：php_request_shutdown->zend_deactivate_modules->request_shutdown_func

再调用Zend扩展的deactivate：php_request_shutdown->zend_deactivate->shutdown_executor->zend_extension_deactivator->deactivate

8、扩展关闭阶段：

先调用PHP扩展的MSHUTDOWN：zend_shutdown->zend_destroy_modules->zend_hash_graceful_reverse_destroy->module_destructor->module_shutdown_func

再调用Zend扩展的shutdown()：zend_shutdown->zend_shutdown_extensions->zend_extension_shutdown->shutdown

一张图了解扩展的生命周期
	<img src="{{site.url}}{{site.baseurl}}/img/zend_extension.png"/>
	http://yangxikun.github.io/php/2016/07/10/php-zend-extension.html
	
	extension意为基于php引擎的扩展

zend_extension意为基于zend引擎的扩展

注：php是基于zend引擎的。晕了吧。


 

不同的扩展安装后，在php.ini里是用extension还是zend_extension，是取决于该扩展，有的扩展可能只能用zend_extension，如xdebug，也有的扩展可以用extension或zend_extension，如mmcache。


注：上面的结论不保证准确。zend_extension加载php扩展时需用全路径，而extension加载时可以用相对extension_dir的路径。


确定可以用zend_extension之后，还有下面几种可能：

根据 PHP 版本，zend_extension 指令可以是以下之一：

 

zend_extension (non ZTS, non debug build)zend_extension_ts ( ZTS, non debug build)zend_extension_debug (non ZTS, debug build)zend_extension_debug_ts ( ZTS, debug build) ZTS：ZEND Thread Safety 可通过phpinfo()查看ZTS是否启用，从而决定用zend_extension还是zend_extension_ts，当然试一下怎么生效也可以。

https://www.cnblogs.com/breg/p/3544799.html

xdebug下载地址:http://www.xdebug.org/download.php
版本选择: xdebug有Non-thread-safe(非线程安全)、thread-safe(线程安全)
写一个test.php，内容为<?php phpinfo(); ?>，搜索"Thread Safety" enable为线程安全版、disable为非线程安全版
选择VC6还是VC9?

标明 MSVC9 (Visual C++ 2008) 的是VC9
如果你在apache1或者apache2下使用PHP，你应该选择VC6的版本
如果你在IIS下使用PHP应该选择VC9的版本
VC6的版本使用visual studio6编译
VC9使用Visual Studio 2008编译，并且改进了性能和稳定性。VC9版本的PHP需要你安装Microsoft 2008 C++ Runtime
不要在apache下使用VC9的版本

Xdebug安装:
将下载的php_xdebug-2.1.0-5.2-vc6.dll放到C:\php5\ext目录，重命名为php_xdebug.dll；

编辑php.ini，加入下面几行：
[Xdebug]
extension=php_xdebug.dll
xdebug.profiler_enable=on
xdebug.trace_output_dir="X:\Projects\xdebug"
xdebug.profiler_output_dir="X:\Projects\xdebug"
后面的目录“I:\Projects\xdebug”为你想要放置Xdebug输出的数据文件的目录，可自由设置。
4． 重启Apache；
5． 写一个test.php，内容为<?php phpinfo(); ?>，如果输出的内容中有看到xdebug，说明安装配置成功。


Xdebug使用
Xdebug使用之开始调试：
我们先写一个可以导致执行出错的程序，例如尝试包含一个不存在的文件。
testXdebug.php
<?php
require_once(‘abc.php’);
?>
然后通过浏览器访问，我们惊奇地发现，出错信息变成了彩色的了：

不过除了样式改变，和我们平时打印的出错信息内容没什么不同，意义不大。好，我们继续改写程序：
testXdebug2.php
<?php
testXdebug();
function testXdebug() {
       require_once('abc.php');
}
?>
输出信息：

发现了什么？　Xdebug跟踪代码的执行，找到了出错的函数testXdebug()。
我们把代码再写得复杂一些：　
testXdebug3.php
<?php
testXdebug();
function testXdebug() {
       requireFile();     
}
function requireFile() {
       require_once('abc.php');
}
?>
输出信息：
'
呵呵，也就是说Xdebug具有类似于Java的Exception的“跟踪回溯”的功能，可以根据程序的执行一步步跟踪到出错的具体位置，哪怕程序中的调用很复杂，我们也可以通过这个功能来理清代码关系，迅速定位，快速排错。


Xdebug配置

第一部分：基本特征:
相关参数设置
xdebug.default_enable
类型：布尔型 默认值：On
如果这项设置为On，堆栈跟踪将被默认的显示在错误事件中。你可以通过在代码中使用xdebug_disable()来禁止堆叠跟踪的显示。因为这是xdebug基本功能之一，将这项参数设置为On是比较明智的。
xdebug.max_nesting_level
类型：整型 默认值：100
The value of this setting is the maximum level of nested functions that are allowed before the  script  will be aborted.
限制无限递归的访问深度。这项参数设置的值是脚本失败前所允许的嵌套程序的最大访问深度。

第二部分：堆栈跟踪:
相关参数设置
xdebug.dump_globals
类型：布尔型 默认值：1
限制是否显示被xdebug.dump.*设置定义的超全局变量的值
例 如，xdebug.dump.SERVER = REQUEST_METHOD,REQUEST_URI,HTTP_USER_AGENT 将打印 PHP 超全局变量 $_SERVER['REQUEST_METHOD']、$_SERVER['REQUEST_URI'] 和 $_SERVER['HTTP_USER_AGENT']。
xdebug.dump_once
类型：布尔型 默认值：1
限制是否超全局变量的值应该转储在所有出错环境(设置为Off时)或仅仅在开始的地方(设置为On时)
xdebug.dump_undefined
类型：布尔型 默认值：0
如果你想从超全局变量中转储未定义的值，你应该把这个参数设置成On，否则就设置成Off
xdebug.show_exception_trace
类型：整型 默认值：0
当这个参数被设置为1时，即使捕捉到异常，xdebug仍将强制执行异常跟踪当一个异常出现时。
xdebug.show_local_vars
类型：整型 默认值：0
当这个参数被设置为不等于0时，xdebug在错环境中所产生的堆栈转储还将显示所有局部变量，包括尚未初始化的变量在最上面。要注意的是这将产生大量的信息，也因此默认情况下是关闭的。

第三部分：分析PHP脚本
相关参数设置
xdebug.profiler_append
类型：整型 默认值：0
当这个参数被设置为1时，文件将不会被追加当一个新的需求到一个相同的文件时(依靠xdebug.profiler_output_name的设置)。相反的设置的话，文件将被附加成一个新文件。
xdebug.profiler_enable
类型：整型 默认值：0
开放xdebug文件的权限，就是在文件输出目录中创建文件。那些文件可以通过KCacheGrind来阅读来展现你的数据。这个设置不能通过在你的脚本中调用ini_set()来设置。
xdebug.profiler_output_dir
类型：字符串 默认值：/tmp
这个文件是profiler文件输出写入的，确信PHP用户对这个目录有写入的权限。这个设置不能通过在你的脚本中调用ini_set()来设置。
xdebug.profiler_output_name
类型：字符串 默认值：cachegrind.out%p
这个设置决定了转储跟踪写入的文件的名称。

第四部分：远程Debug
相关参数设置
xdebug.remote_autostart
类型：布尔型 默认值：0
一般来说，你需要使用明确的HTTP GET/POST变量来开启远程debug。而当这个参数设置为On，xdebug将经常试图去开启一个远程debug session并试图去连接客户端，即使GET/POST/COOKIE变量不是当前的。
xdebug.remote_enable
类型：布尔型 默认值：0
这个开关控制xdebug是否应该试着去连接一个按照xdebug.remote_host和xdebug.remote_port来设置监听主机和端口的debug客户端。
xdebug.remote_host
类型：字符串 默认值：localhost
选择debug客户端正在运行的主机，你不仅可以使用主机名还可以使用IP地址
xdebug.remote_port
类型：整型 默认值：9000
这个端口是xdebug试着去连接远程主机的。9000是一般客户端和被绑定的debug客户端默认的端口。许多客户端都使用这个端口数字，最好不要去修改这个设置。
注意：所有以上参数修改后，要重启Apache才能生效！

Xdebug调试
其实PHP函数debug_backtrace()也有类似的功能，但是要注意debug_backtrace()函数只在PHP4.3.0之后版本及
PHP5中才生效。这个函数是PHP开发团队在PHP5中新增的函数，然后又反向移植到PHP4.3中。
Xdebug使调试信息更加美观
Xdebug 扩展加载后，Xdebug会对原有的某些PHP函数进行覆写，以便好更好地进行Debug。比如var_dump()函数，我们知道通常我们需要在函数前 后加上”<pre>…</pre>”才能够让输出的变量信息比较美观、可读性好。但是加载了Xdebug后，我们不再需要这样做 了，Xdebug不但自动给我们加上了<pre>标签，还给变量加上颜色。
例：
<?php
$arrTest=array(
       "test"=>"abc",
       "test2"=>"abc2"
);
var_dump($arrTest);
?>
输出：

看到了吗？　数组元素的值自动显示颜色。
如果你还是希望使用PHP的var_dump函数 只要在php.ini关于xdebug的配置中加上 xdebug.overload_var_dump = Off 即可
Xdebug测试脚本执行时间
测试某段脚本的执行时间，通常我们都需要用到microtime()函数来确定当前时间。例如PHP手册上的例子：
<?php
/**
* Simple function to replicate PHP 5 behaviour
*/
function microtime_float()
{
    list($usec, $sec) = explode(" ", microtime());
return ((float)$usec + (float)$sec);
}
$time_start = microtime_float();
// Sleep for a while
usleep(100);
$time_end = microtime_float();
$time = $time_end - $time_start;
echo "Did nothing in $time seconds\n";
?>
但是microtime()返回的值是微秒数及绝对时间戳（例如“0.03520000 1153122275”），没有可读性。所以如上程序，我们需要另外写一个函数microtime_float()，来将两者相加。
Xdebug自带了一个函数xdebug_time_index()来显示时间。

PHP脚本占用的内存
有时候我们想知道程序执行到某个特定阶段时到底占用了多大内存，为此PHP提供了函数memory_get_usage()。这个函数只有当PHP编译时使用了--enable-memory-limit参数时才有效。　
Xdebug同样提供了一个函数xdebug_memory_usage()来实现这样的功能，另外xdebug还提供了一个xdebug_peak_memory_usage()函数来查看内存占用的峰值。


WinCacheGrind
有 时候代码没有明显的编写错误，没有显示任何错误信息（如error、warning、notice等），但是这不表明代码就是正确无误的。有时候可能某段 代码执行时间过长，占用内存过多以致于影响整个系统的效率，我们没有办法直接看出来是哪部份代码出了问题。这时候我们希望把代码的每个阶段的运行情况都监 控起来，写到日志文件中去，运行一段时间后再进行分析，找到问题所在。
回忆一下，之前我们编辑php.ini文件
加入
[Xdebug]
xdebug.profiler_enable=on
xdebug.trace_output_dir="I:\Projects\xdebug"
xdebug.profiler_output_dir="I:\Projects\xdebug"
这 几行，目的就在于把执行情况的分析文件写入到”I:\Projects\xdebug”目录中去（你可以替换成任何你想设定的目录）。如果你执行某段程序 后，再打开相应的目录，可以发现生成了一堆文件，例如cachegrind.out.1169585776这种格式命名的文件。这些就是Xdebug生成 的分析文件。用编辑器打开你可以看到很多程序运行的相关细节信息，不过很显然这样看太累了，我们需要用图形化的软件来查看。
WinCacheGrind下载
在Windows平台下，可以用WinCacheGrind(wincachegrind.souceforge.net)这个软件来打开这些文件。可以直观漂亮地显示其中内容：

WinCacheGrind小结：
Xdebug提供了各种自带的函数，并对已有的某些PHP函数进行覆写，可以方便地用于调试排错；Xdebug还可以跟踪程序的运行，通过对日志文件的分析，我们可以迅速找到程序运行的瓶颈所在，提高程序效率，从而提高整个系统的性能。


如果在linux下 你可以使用 kcachegrind 查看xdebug的日志

 

 

 

php扩展之关于extension,zend_extension和zend_extension_ts

extension意为基于php引擎的扩展

zend_extension意为基于zend引擎的扩展

注：php是基于zend引擎的。晕了吧。

 

不同的扩展安装后，在php.ini里是用extension还是zend_extension，是取决于该扩展，有的扩展可能只能用zend_extension，如xdebug，也有的扩展可以用extension或zend_extension，如mmcache。

 

注：上面的结论不保证准确。zend_extension加载php扩展时需用全路径，而extension加载时可以用相对extension_dir的路径。

 

确定可以用zend_extension之后，还有下面几种可能：

根据 PHP 版本，zend_extension 指令可以是以下之一：

zend_extension (non ZTS, non debug build)
zend_extension_ts ( ZTS, non debug build)
zend_extension_debug (non ZTS, debug build)
zend_extension_debug_ts ( ZTS, debug build)

ZTS：ZEND Thread Safety

可通过phpinfo()查看ZTS是否启用，从而决定用
zend_extension还是zend_extension_ts，当然试一下怎么生效也可以。

https://blog.51cto.com/higgs/996794
{% endraw %}
https://www.cnblogs.com/ghjbk/p/6953576.html
https://www.php.cn/php-weizijiaocheng-329268.html
