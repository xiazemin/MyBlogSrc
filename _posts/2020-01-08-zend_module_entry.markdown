---
title: php 扩展加载顺序
layout: post
category: lang
author: 夏泽民
---
在PHP开发中使用扩展会大大提高编码和运行效率，但在使用过程中（尤其是一些开源扩展）难免遇到各种各样的问题，其中一个就是模块间冲突或依赖

一般来讲，比较好的扩展会提示你跟它冲突或有依赖的模块，运行php -m会有类似下面的这个提示
说明扩展testso所依赖的apm.so并没有被安装

如果我们继续深入了解这种功能如何实现的，那么就要研究一下该扩展的源码。其实做到这点并不复杂，也不需要太多的扩展开发知识，只需要了解几个结构体的相关功能即可

首先，每个扩展都会声明一个zend_module_entry结构的扩展说明模块，告诉Zend所有该模块相关信息。自然也包含所依赖和冲突的模块。当依赖关系出错的时候Zend就会报错，给出详细的出错信息，并停止运行

struct _zend_module_entry {
	unsigned short size;
	unsigned int zend_api;
	unsigned char zend_debug;
	unsigned char zts;
	const struct _zend_ini_entry *ini_entry;
	const struct _zend_module_dep *deps;         //我们要看的就是这个属性
	const char *name;
	const struct _zend_function_entry *functions;
	int (*module_startup_func)(INIT_FUNC_ARGS);
	int (*module_shutdown_func)(SHUTDOWN_FUNC_ARGS);
	int (*request_startup_func)(INIT_FUNC_ARGS);
	int (*request_shutdown_func)(SHUTDOWN_FUNC_ARGS);
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
deps就是用来注册依赖、冲突模块的，它的结构zend_module_dep如下：

struct _zend_module_dep {
	const char *name;		/* module name */
	const char *rel;		        /* version relationship: NULL (exists), lt|le|eq|ge|gt (to given version) */
	const char *version;	        /* version */
	unsigned char type;		/* dependency type */
};
我们可以通过研究代码及注释来了解这个结构体的使用方法，不过还有更好的方法，原来PHP贴心的提供了几个宏来完成zend_module_dep的填写：

ZEND_MOD_REQUIRED
ZEND_MOD_CONFLICTS
ZEND_MOD_OPTIONAL
比如在刚开始说到的testso中，我就是这么写的：

static const zend_module_dep test_deps[] = {
        ZEND_MOD_REQUIRED("apm")
        ZEND_MOD_END
};

zend_module_entry testso_module_entry = {
        STANDARD_MODULE_HEADER_EX,
        NULL,
        test_deps,                            //zend_module_dep
        "testso",
        testso_functions,
        PHP_MINIT(testso),
        PHP_MSHUTDOWN(testso),
        PHP_RINIT(testso),              
        PHP_RSHUTDOWN(testso), 
        PHP_MINFO(testso),
        PHP_TESTSO_VERSION,
        STANDARD_MODULE_PROPERTIES
};
至于第三个宏ZEND_MOD_OPTIONAL怎么用，大家有兴趣可以自行研究源码。在PHP这种开源项目中，源码本身就向你说明了一切

最后再讲一下扩展加载顺序的问题，这个规则更简单，只有一条：写在php.ini中的扩展按从上到下的顺序加载

比如在上面的情况中，testso的顺序就在apm之前。但是尽管放心，扩展依赖并不取决于加载顺序，所以把apm放在后面，PHP仍然可以正常执行
<!-- more -->
PHP扩展与Zend扩展区别
通常在php.ini中，通过extension=*加载的扩展我们称为PHP扩展，通过zend_extension=*加载的扩展我们称为Zend扩展，但从源码的角度来讲，PHP扩展应该称为“模块”（源码中以module命名），而Zend扩展称为“扩展”（源码中以extension命名）。

两者最大的区别在于向引擎注册的钩子。少数的扩展，例如xdebug、opcache，既是PHP扩展，也是Zend扩展，但它们在php.ini中的加载方式得用zend_extension=*，具体原因下文会说明。

先来看看两种类型扩展主要的数据结构：

PHP扩展的结构体：

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
<img src="{{site.url}}{{site.baseurl}}/img/phplifecircle.png"/>

我们从未手动开启过PHP的相关进程，它是随着Apache的启动而运行的；
PHP通过mod_php5.so模块和Apache相连（具体说来是SAPI，即服务器应用程序编程接口）；
PHP总共有三个模块：内核、Zend引擎、以及扩展层；
PHP内核用来处理请求、文件流、错误处理等相关操作；
Zend引擎（ZE）用以将源文件转换成机器语言，然后在虚拟机上运行它；
扩展层是一组函数、类库和流，PHP使用它们来执行一些特定的操作。比如，我们需要mysql扩展来连接MySQL数据库；
当ZE执行程序时可能会需要连接若干扩展，这时ZE将控制权交给扩展，等处理完特定任务后再返还；
最后，ZE将程序运行结果返回给PHP内核，它再将结果传送给SAPI层，最终输出到浏览器上。
深入探讨
　　等等，没有这么简单。以上过程只是个简略版，让我们再深入挖掘一下，看看幕后还发生了些什么。

Apache启动后，PHP解释程序也随之启动；
PHP的启动过程有两步；
第一步是初始化一些环境变量，这将在整个SAPI生命周期中发生作用；
第二步是生成只针对当前请求的一些变量设置。
PHP启动第一步
　　不清楚什么第一第二步是什么？别担心，我们接下来详细讨论一下。让我们先看看第一步，也是最主要的一步。要记住的是，第一步的操作在任何请求到达之前就发生了。

启动Apache后，PHP解释程序也随之启动；
PHP调用各个扩展的MINIT方法，从而使这些扩展切换到可用状态。看看php.ini文件里打开了哪些扩展吧；
MINIT的意思是“模块初始化”。各个模块都定义了一组函数、类库等用以处理其他请求。
　　一个典型的MINIT方法如下：
PHP_MINIT_FUNCTION(extension_name){
/* Initialize functions, classes etc */
}
PHP启动第二步

当一个页面请求发生时，SAPI层将控制权交给PHP层。于是PHP设置了用于回复本次请求所需的环境变量。同时，它还建立一个变量表，用来存放执行过程中产生的变量名和值。
PHP调用各个模块的RINIT方法，即“请求初始化”。一个经典的例子是Session模块的RINIT，如果在php.ini中启用了Session模块，那在调用该模块的RINIT时就会初始化$_SESSION变量，并将相关内容读入；
RINIT方法可以看作是一个准备过程，在程序执行之间就会自动启动。
一个典型的RINIT方法如下：
PHP_RINIT_FUNCTION(extension_name) {
/* Initialize session variables, pre-populate variables, redefine global variables etc */
}
PHP关闭第一步
　　如同PHP启动一样，PHP的关闭也分两步：

一旦页面执行完毕（无论是执行到了文件末尾还是用exit或die函数中止），PHP就会启动清理程序。它会按顺序调用各个模块的RSHUTDOWN方法。
RSHUTDOWN用以清除程序运行时产生的符号表，也就是对每个变量调用unset函数。
　　一个典型的RSHUTDOWN方法如下：
PHP_RSHUTDOWN_FUNCTION(extension_name) {
/* Do memory management, unset all variables used in the last PHP call etc */
}
PHP关闭第二步
　　最后，所有的请求都已处理完毕，SAPI也准备关闭了，PHP开始执行第二步：

PHP调用每个扩展的MSHUTDOWN方法，这是各个模块最后一次释放内存的机会。
　　一个典型的RSHUTDOWN方法如下：
PHP_MSHUTDOWN_FUNCTION(extension_name) {
/* Free handlers and persistent memory etc */
}
　　这样，整个PHP生命周期就结束了。要注意的是，只有在服务器没有请求的情况下才会执行“启动第一步”和“关闭第二步”。
　　
　　php从下到上是一个4层体系

①Zend引擎

Zend整体用纯c实现，是php的内核部分，它将php代码翻译（词法、语法解析等一系列编译过程）为可执行opcode的处理并实现相应的处理方法、实现了基本的数据结构（如hashtable、oo）、内存分配及管理、提供了相应的api方法供外部调用，是一切的核心，所有的外围功能均围绕zend实现。

②Extensions

围绕着zend引擎，extensions通过组件式的方式提供各种基础服务，我们常见的各种内置函数（如array系列）、标准库等都是通过extension来实现，用户也可以根据需要实现自己的extension以达到功能扩展、性能优化等目的（如贴吧正在使用的php中间层、富文本解析就是extension的典型应用）。

③Sapi

Sapi全称是Server Application Programming Interface，也就是服务端应用编程接口，sapi通过一系列钩子函数，使得php可以和外围交互数据，这是php非常优雅和成功的一个设计，通过sapi成功的将php本身和上层应用解耦隔离，php可以不再考虑如何针对不同应用进行兼容，而应用本身也可以针对自己的特点实现不同的处理方式。后面将在sapi章节中介绍

④上层应用

这就是我们平时编写的php程序，通过不同的sapi方式得到各种各样的应用模式，如通过webserver实现web应用、在命令行下以脚本方式运行等等。

构架思想：

引擎(Zend)+组件(ext)的模式降低内部耦合

中间层(sapi)隔绝web server和php

**************************************************************************

如果php是一辆车，那么

车的框架就是php本身

Zend是车的引擎（发动机）

Ext下面的各种组件就是车的轮子

Sapi可以看做是公路，车可以跑在不同类型的公路上

而一次php程序的执行就是汽车跑在公路上。

因此，我们需要：性能优异的引擎+合适的车轮+正确的跑道

Apache和php的关系
Apache对于php的解析，就是通过众多Module中的php Module来完成的。

把php最终集成到Apache系统中，还需要对Apache进行一些必要的设置。这里，我们就以php的mod_php5 SAPI运行模式为例进行讲解，至于SAPI这个概念后面我们还会详细讲解。

假定我们安装的版本是Apache2 和 Php5，那么需要编辑Apache的主配置文件http.conf，在其中加入下面的几行内容：

Unix/Linux环境下：

LoadModule php5_module modules/mod_php5.so

AddType application/x-httpd-php .php

注：其中modules/mod_php5.so 是X系统环境下mod_php5.so文件的安装位置。

Windows环境下：

LoadModule php5_module d:/php/php5apache2.dll

AddType application/x-httpd-php .php

注：其中d:/php/php5apache2.dll 是在Windows环境下php5apache2.dll文件的安装位置。

这两项配置就是告诉Apache Server，以后收到的Url用户请求，凡是以php作为后缀，就需要调用php5_module模块（mod_php5.so/ php5apache2.dll）进行处理。

Apache的生命周期
wps_clip_image-8490

Apach的请求处理流程
wps_clip_image-17917

Apache请求处理循环详解 
    Apache请求处理循环的11个阶段都做了哪些事情呢？

1、Post-Read-Request阶段

    在正常请求处理流程中，这是模块可以插入钩子的第一个阶段。对于那些想很早进入处理请求的模块来说，这个阶段可以被利用。

    2、URI Translation阶段 
    Apache在本阶段的主要工作：将请求的URL映射到本地文件系统。模块可以在这阶段插入钩子，执行自己的映射逻辑。mod_alias就是利用这个阶段工作的。

    3、Header Parsing阶段 
    Apache在本阶段的主要工作：检查请求的头部。由于模块可以在请求处理流程的任何一个点上执行检查请求头部的任务，因此这个钩子很少被使用。mod_setenvif就是利用这个阶段工作的。

    4、Access Control阶段 
    Apache在本阶段的主要工作：根据配置文件检查是否允许访问请求的资源。Apache的标准逻辑实现了允许和拒绝指令。mod_authz_host就是利用这个阶段工作的。

    5、Authentication阶段 
     Apache在本阶段的主要工作：按照配置文件设定的策略对用户进行认证，并设定用户名区域。模块可以在这阶段插入钩子，实现一个认证方法。

    6、Authorization阶段 
    Apache在本阶段的主要工作：根据配置文件检查是否允许认证过的用户执行请求的操作。模块可以在这阶段插入钩子，实现一个用户权限管理的方法。

    7、MIME Type Checking阶段 
    Apache在本阶段的主要工作：根据请求资源的MIME类型的相关规则，判定将要使用的内容处理函数。标准模块mod_negotiation和mod_mime实现了这个钩子。

    8、FixUp阶段 
    这是一个通用的阶段，允许模块在内容生成器之前，运行任何必要的处理流程。和Post_Read_Request类似，这是一个能够捕获任何信息的钩子，也是最常使用的钩子。

    9、Response阶段 
    Apache在本阶段的主要工作：生成返回客户端的内容，负责给客户端发送一个恰当的回复。这个阶段是整个处理流程的核心部分。

    10、Logging阶段 
    Apache在本阶段的主要工作：在回复已经发送给客户端之后记录事务。模块可能修改或者替换Apache的标准日志记录。

11、CleanUp阶段 
    Apache在本阶段的主要工作：清理本次请求事务处理完成之后遗留的环境，比如文件、目录的处理或者Socket的关闭等等，这是Apache一次请求处理的最后一个阶段。

LAMP架构：

wps_clip_image-24435

从下往上四层：

①liunx 属于操作系统的底层

②apache服务器，属于次服务器，沟通linux和PHP

③php:属于服务端编程语言，通过php_module 模块 和apache关联

    ④mysql和其他web服务：属于应用服务，通过PHP的Extensions外 挂模块和mysql关联

Android系统架构图
lamp和安卓的架构图比较一下，貌似和lamp架构有点相似，本人不懂安卓，只是感觉上有点相似，高手可以指出区别，小弟在此不胜感谢

wps_clip_image-27187

从上往下：

安卓架构--------------说明--------LAMP架构

1.应用程序 --------具体应用--------web应用

2.应用程序框架 ----java-------------PHP语言和库

3.系统运行库 ：----虚拟机---------WEB服务器

⒋Linux 内核 ：---操作系统-------lamp架构中的L

1、PHP的运行模式：
    PHP两种运行模式是WEB模式、CLI模式。无论哪种模式，PHP工作原理都是一样的，作为一种SAPI运行。

1、当我们在终端敲入php这个命令的时候，它使用的是CLI。

它就像一个web服务器一样来支持php完成这个请求，请求完成后再重新把控制权交给终端。

2、当使用Apache或者别web服务器作为宿主时，当一个请求到来时，PHP会来支持完成这个请求。一般有：

 

    多进程(通常编译为apache的模块来处理PHP请求)

    多线程模式

2、一切的开始: SAPI接口
    通常我们编写php Web程序都是通过Apache或者Nginx这类Web服务器来测试脚本. 或者在命令行下通过php程序来执行PHP脚本. 执行完成脚本后，服务器应答，浏览器显示应答信息,或者在命令结束后在标准输出显示内容. 我们很少关心PHP解释器在哪里. 虽然通过Web服务器和命令行程序执行脚本看起来很不一样. 实际上她们的工作是一样的. 命令行程序和Web程序类似, 命令行参数传递给要执行的脚本,相当于通过url 请求一个PHP页面. 脚本戳里完成后返回响应结果,只不过命令行响应的结果是显示在终端上. 脚本执行的开始都是通过SAPI接口进行的. 

        1)、启动apache：当给定的SAPI启动时，例如在对/usr/local/apache/bin/apachectl start的响应中，PHP由初始化其内核子系统开始。在接近启动例程的末尾，它加载每个扩展的代码并调用其模块初始化例程（MINIT）。这使得每个扩展可以初始化内部变量、分配资源、注册资源处理器，以及向ZE注册自己的函数，以便于脚本调用这其中的函数时候ZE知道执行哪些代码。

       2)、请求处理初始化：接下来，PHP等待SAPI层请求要处理的页面。对于CGI或CLI等SAPI，这将立刻发生且只发生一次。对于Apache、IIS或其他成熟的web服务器SAPI，每次远程用户请求页面时都将发生，因此重复很多次，也可能并发。不管请求如何产生，PHP开始于要求ZE建立脚本的运行环境，然后调用每个扩展的请求初始化 （RINIT）函数。RINIT使得扩展有机会设定特定的环境变量，根据请求分配资源，或者执行其他任务，如审核。 session扩展中有个RINIT作用的典型示例，如果启用了session.auto_start选项，RINIT将自动触发用户空间的session_start()函数以及预组装$_SESSION变量。

      3)、执行php代码： 一旦请求被初始化了，ZE开始接管控制权，将PHP脚本翻译成符号，最终形成操作码并逐步运行之。如任一操作码需要调用扩展的函数，ZE将会把参数绑定到该函数，并且临时交出控制权直到函数运行结束。

       4)、脚本结束：脚本运行结束后，PHP调用每个扩展的请求关闭（RSHUTDOWN）函数以执行最后的清理工作（如将session变量存入磁盘）。接下来，ZE执行清理过程（垃圾收集）－有效地对之前的请求期间用到的每个变量执行unset()。

       5)、sapi关闭：一旦完成，PHP继续等待SAPI的其他文档请求或者是关闭信号。对于CGI和CLI等SAPI，没有“下一个请求”，所以SAPI立刻开始关闭。关闭期间，PHP再次遍历每个扩展，调用其模块关闭（MSHUTDOWN）函数，并最终关闭自己的内核子系统。

       简要的过程如下：

      1. PHP是随着Apache的启动而运行的；
      2. PHP通过mod_php5.so模块和Apache相连（具体说来是SAPI，即服务器应用程序编程接口）；
      3. PHP总共有三个模块：内核、Zend引擎、以及扩展层；
      4. PHP内核用来处理请求、文件流、错误处理等相关操作；
      5. Zend引擎（ZE）用以将源文件转换成机器语言，然后在虚拟机上运行它；
      6. 扩展层是一组函数、类库和流，PHP使用它们来执行一些特定的操作。比如，我们需要mysql扩展来连接MySQL数据库；
      7. 当ZE执行程序时可能会需要连接若干扩展，这时ZE将控制权交给扩展，等处理完特定任务后再返还；
      8. 最后，ZE将程序运行结果返回给PHP内核，它再将结果传送给SAPI层，最终输出到浏览器上。

3、PHP的开始和结束阶段
开始阶段有两个过程：
     第一个过程：apache启动的过程，即在任何请求到达之前就发生。是在整个SAPI生命周期内(例如Apache启动以后的整个生命周期内或者命令行程序整个执行过程中)的开始阶段(MINIT),该阶段只进行一次.。启动Apache后，PHP解释程序也随之启动； PHP调用各个扩展（模块）的MINIT方法，从而使这些扩展切换到可用状态。看看php.ini文件里打开了哪些扩展吧； MINIT的意思是“模块初始化”。各个模块都定义了一组函数、类库等用以处理其他请求。 模块在这个阶段可以进行一些初始化工作,例如注册常量, 定义模块使用的类等等.典型的的模块回调函数MINIT方法如下：

 

 

[cpp] view plaincopyprint?
PHP_MINIT_FUNCTION(myphpextension) { /* Initialize functions, classes etc */ }  
{  
    // 注册常量或者类等初始化操作   
    return SUCCESS;   
}  
     第二个过程发生在请求阶段,当一个页面请求发生时.则在每次请求之前都会进行初始化过程(RINIT请求开始).

请求到达之后，SAPI层将控制权交给PHP层，PHP初始化本次请求执行脚本所需的环境变量,例如创建一个执行环境,包括保存php运行过程中变量名称和变量值内容的符号表. 以及当前所有的函数以及类等信息的符号表.例如是Session模块的RINIT，如果在php.ini中启用了Session 模块，那在调用该模块的RINIT时就会初始化$_SESSION变量，并将相关内容读入；  然后PHP会调用所有模块RINIT函数,即“请求初始化”。 在这个阶段各个模块也可以执行一些相关的操作, 模块的RINIT函数和MINIT函数类似 ，RINIT方法可以看作是一个准备过程，在程序执行之间就会自动启动。

 

[cpp] view plaincopyprint?
PHP_RINIT_FUNCTION(myphpextension)  
{  
    // 例如记录请求开始时间   
    // 随后在请求结束的时候记录结束时间.这样我们就能够记录下处理请求所花费的时间了  
    return SUCCESS;   
}  
结束阶段分为两个环节：
请求处理完后就进入了结束阶段, 一般脚本执行到末尾或者通过调用exit()或者die()函数,PHP都将进入结束阶段. 和开始阶段对应,结束阶段也分为两个环节,一个在请求结束后(RSHUWDOWN),一个在SAPI生命周期结束时(MSHUTDOWN).
 
第一个环节：请求处理完后结束阶段：请求处理完后就进入了结束阶段，PHP就会启动清理程序。它会按顺序调用各个模块的RSHUTDOWN方法。 RSHUTDOWN用以清除程序运行时产生的符号表，也就是对每个变量调用unset函数。典型的RSHUTDOWN方法如下：
[cpp] view plaincopyprint?
PHP_RSHUTDOWN_FUNCTION(myphpextension)  
{  
    // 例如记录请求结束时间, 并把相应的信息写入到日至文件中.  
    return SUCCESS;   
}  
第二个环节：最后，所有的请求都已处理完毕，SAPI也准备关闭了， PHP调用每个扩展的MSHUTDOWN方法，这是各个模块最后一次释放内存的机会。（这个是对于CGI和CLI等SAPI，没有“下一个请求”，所以SAPI立刻开始关闭。）

典型的RSHUTDOWN方法如下：

 

[plain] view plaincopyprint?
PHP_MSHUTDOWN_FUNCTION(extension_name) {   
    /* Free handlers and persistent memory etc */   
    return SUCCESS;   
}  
 

这样，整个PHP生命周期就结束了。要注意的是，只有在服务器没有请求的情况下才会执行“启动第一步”和“关闭第二步”。

SAPI运行PHP都经过下面几个阶段:
       1、模块初始化阶段(Module init)     ：
           即调用每个拓展源码中的的PHP_MINIT_FUNCTION中的方法初始化模块,进行一些模块所需变量的申请,内存分配等。
        2、请求初始化阶段(Request init)  ：
           即接受到客户端的请求后调用每个拓展的PHP_RINIT_FUNCTION中的方法,初始化PHP脚本的执行环境。
        3、执行PHP脚本
        4、请求结束(Request Shutdown) ：
          这时候调用每个拓展的PHP_RSHUTDOWN_FUNCTION方法清理请求现场,并且ZE开始回收变量和内存。
        5、关闭模块(Module shutdown)     ：
           Web服务器退出或者命令行脚本执行完毕退出会调用拓展源码中的PHP_MSHUTDOWN_FUNCTION 方法

 

4、单进程SAPI生命周期
 

CLI/CGI模式的PHP属于单进程的SAPI模式。这类的请求在处理一次请求后就关闭。也就是只会经过如下几个环节： 开始 - 请求开始 - 请求关闭 - 结束 SAPI接口实现就完成了其生命周期。如图所示：

                                      

 

 

5、多进程SAPI生命周期
 

通常PHP是编译为apache的一个模块来处理PHP请求。Apache一般会采用多进程模式， Apache启动后会

fork出多个子进程，每个进程的内存空间独立，每个子进程都会经过开始和结束环节， 不过每个进程的开始阶

段只在进程fork出来以来后进行，在整个进程的生命周期内可能会处理多个请求。 只有在Apache关闭或者进程

被结束之后才会进行关闭阶段，在这两个阶段之间会随着每个请求重复请求开始-请求关闭的环节。 

如图所示：

                                     

 

6、多线程的SAPI生命周期
多线程模式和多进程中的某个进程类似，不同的是在整个进程的生命周期内会并行的重复着 请求开始-请求关闭的环节.

 

 

在这种模式下，只有一个服务器进程在运行着，但会同时运行很多线程，这样可以减少一些资源开销，向Module init和Module shutdown就只需要运行一遍就行了，一些全局变量也只需要初始化一次，因为线程独具的特质，使得各个请求之间方便的共享一些数据成为可能。

 多线程工作方式如下图

                              

 

7、Apache一般使用多进程模式prefork
        在linux下使用#http –l 命令可以查看当前使用的工作模式。也可以使用#apachectl -l命令。
        看到的prefork.c，说明使用的prefork工作模式。

        prefork 进程池模型，用在 UNIX 和类似的系统上比较多，主要是由于写起来方便，也容易移植，还不容易出问题。要知道，如果采用线程模型的话，用户线程、内核线程和混合型线程有不同的特性，移植起来就麻烦。prefork 模型，即预先 fork() 出来一些子进程缓冲一下，用一个锁来控制同步，连接到来了就放行一个子进程，让它去处理。

 

    prefork MPM 使用多个子进程，每个子进程只有一个线程。每个进程在某个确定的时间只能维持一个连接。在大多数平台上，Prefork MPM在效率上要比Worker MPM要高，但是内存使用大得多。prefork的无线程设计在某些情况下将比worker更有优势：他能够使用那些没有处理好线程安全的第三方模块，并 且对于那些线程调试困难的平台而言，他也更容易调试一些。
 PHP生命周期
一切的开始：SAPI（Server Application Programming Interface）
    SAPI指的是PHP具体应用的编程接口，就像PC一样，无论安装哪些操作系统，只要满足了PC的接口规范都可以在PC上正常运行。PHP脚本要执行有很多种方式，通过Web服务器，或者直接在命令行下，也可以嵌入在其他程序中。

    通常，我们使用Apache或者Nginx这类Web服务器来测试PHP脚本，或者在命令行下通过PHP解释器程序来执行。 脚本执行完后，Web服务器应答，浏览器显示应答信息，或者在命令行标准输出上显示内容。

    我们很少关心PHP解释器在哪里。虽然通过Web服务器和命令行程序执行脚本看起来很不一样， 实际上它们的工作流程是一样的。命令行参数传递给PHP解释器要执行的脚本， 相当于通过url请求一个PHP页面。脚本执行完成后返回响应结果，只不过命令行的响应结果是显示在终端上。   

    脚本执行的开始都是以SAPI接口实现开始的。只是不同的SAPI接口实现会完成他们特定的工作， 例如Apache的mod_php SAPI实现需要初始化从Apache获取的一些信息，在输出内容是将内容返回给Apache， 其他的SAPI实现也类似。

开始和结束
PHP开始执行以后会经过两个主要的阶段：处理请求之前的开始阶段和请求之后的结束阶段。 开始阶段有两个过程：第一个过程是模块初始化阶段（MINIT）， 在整个SAPI生命周期内(例如Apache启动以后的整个生命周期内或者命令行程序整个执行过程中)， 该过程只进行一次。第二个过程是模块激活阶段（RINIT），该过程发生在请求阶段， 例如通过url请求某个页面，则在每次请求之前都会进行模块激活（RINIT请求开始）。 例如PHP注册了一些扩展模块，则在MINIT阶段会回调所有模块的MINIT函数。 模块在这个阶段可以进行一些初始化工作，例如注册常量，定义模块使用的类等等。 模块在实现时可以通过如下宏来实现这些回调函数：

    

PHP_MINIT_FUNCTION(myphpextension)
{
    // 注册常量或者类等初始化操作
    return SUCCESS; 
}
请求到达之后PHP初始化执行脚本的基本环境，例如创建一个执行环境，包括保存PHP运行过程中变量名称和值内容的符号表， 以及当前所有的函数以及类等信息的符号表。然后PHP会调用所有模块的RINIT函数， 在这个阶段各个模块也可以执行一些相关的操作，模块的RINIT函数和MINIT回调函数类似：

复制代码
PHP_RINIT_FUNCTION(myphpextension)
{
    // 例如记录请求开始时间
    // 随后在请求结束的时候记录结束时间。这样我们就能够记录下处理请求所花费的时间了
    return SUCCESS; 
}
复制代码
 

请求处理完后就进入了结束阶段，一般脚本执行到末尾或者通过调用exit()或die()函数， PHP都将进入结束阶段。和开始阶段对应，结束阶段也分为两个环节，一个在请求结束后停用模块(RSHUTDOWN，对应RINIT)， 一个在SAPI生命周期结束（Web服务器退出或者命令行脚本执行完毕退出）时关闭模块(MSHUTDOWN，对应MINIT)。

PHP_RSHUTDOWN_FUNCTION(myphpextension)
{
    // 例如记录请求结束时间，并把相应的信息写入到日至文件中。
    return SUCCESS; 
}
 

想要了解扩展开发的相关内容，请参考第十三章 扩展开发

单进程SAPI生命周期
CLI/CGI模式的PHP属于单进程的SAPI模式。这类的请求在处理一次请求后就关闭。也就是只会经过如下几个环节： 开始 - 请求开始 - 请求关闭 - 结束 SAPI接口实现就完成了其生命周期。如图2.1所示：

 

图2.1 单进程SAPI生命周期
图2.1 单进程SAPI生命周期
 

如上的图是非常简单，也很好理解。只是在各个阶段之间PHP还做了许许多多的工作。这里做一些补充：

启动

在调用每个模块的模块初始化前，会有一个初始化的过程，它包括：

初始化若干全局变量
这里的初始化全局变量大多数情况下是将其设置为NULL，有一些除外，比如设置zuf（zend_utility_functions）， 以zuf.printf_function = php_printf为例，这里的php_printf在zend_startup函数中会被赋值给zend_printf作为全局函数指针使用， 而zend_printf函数通常会作为常规字符串输出使用，比如显示程序调用栈的debug_print_backtrace就是使用它打印相关信息。

初始化若干常量
这里的常量是PHP自己的一些常量，这些常量要么是硬编码在程序中,比如PHP_VERSION，要么是写在配置头文件中， 比如PEAR_EXTENSION_DIR，这些是写在config.w32.h文件中。

初始化Zend引擎和核心组件
前面提到的zend_startup函数的作用就是初始化Zend引擎，这里的初始化操作包括内存管理初始化、 全局使用的函数指针初始化（如前面所说的zend_printf等），对PHP源文件进行词法分析、语法分析、 中间代码执行的函数指针的赋值，初始化若干HashTable（比如函数表，常量表等等），为ini文件解析做准备， 为PHP源文件解析做准备，注册内置函数（如strlen、define等），注册标准常量（如E_ALL、TRUE、NULL等）、注册GLOBALS全局变量等。

解析php.ini
php_init_config函数的作用是读取php.ini文件，设置配置参数，加载zend扩展并注册PHP扩展函数。此函数分为如下几步： 初始化参数配置表，调用当前模式下的ini初始化配置，比如CLI模式下，会做如下初始化：

INI_DEFAULT("report_zend_debug", "0");
INI_DEFAULT("display_errors", "1");
 

不过在其它模式下却没有这样的初始化操作。接下来会的各种操作都是查找ini文件：

判断是否有php_ini_path_override，在CLI模式下可以通过-c参数指定此路径（在php的命令参数中-c表示在指定的路径中查找ini文件）。
如果没有php_ini_path_override，判断php_ini_ignore是否为非空（忽略php.ini配置，这里也就CLI模式下有用，使用-n参数）。
如果不忽略ini配置，则开始处理php_ini_search_path（查找ini文件的路径），这些路径包括CWD(当前路径，不过这种不适用CLI模式)、 执行脚本所在目录、环境变量PATH和PHPRC和配置文件中的PHP_CONFIG_FILE_PATH的值。
在准备完查找路径后，PHP会判断现在的ini路径（php_ini_file_name）是否为文件和是否可打开。 如果这里ini路径是文件并且可打开，则会使用此文件， 也就是CLI模式下通过-c参数指定的ini文件的优先级是最高的， 其次是PHPRC指定的文件，第三是在搜索路径中查找php-%sapi-module-name%.ini文件（如CLI模式下应该是查找php-cli.ini文件）， 最后才是搜索路径中查找php.ini文件。
        php.ini 的搜索路径如下（按顺序）：

SAPI 模块所指定的位置（Apache 2 中的 PHPIniDir 指令，CGI 和 CLI 中的 -c 命令行选项，NSAPI 中的 php_ini 参数，THTTPD 中的PHP_INI_PATH 环境变量）。
PHPRC 环境变量。在 PHP 5.2.0 之前，其顺序在以下提及的注册表键值之后。
自 PHP 5.2.0 起，可以为不同版本的 PHP 指定不同的 php.ini 文件位置。将以下面的顺序检查注册表目录：[HKEY_LOCAL_MACHINE\SOFTWARE\PHP\x.y.z]，[HKEY_LOCAL_MACHINE\SOFTWARE\PHP\x.y] 和[HKEY_LOCAL_MACHINE\SOFTWARE\PHP\x]，其中的 x，y 和 z 指的是 PHP 主版本号，次版本号和发行批次。如果在其中任何目录下的IniFilePath 有键值，则第一个值将被用作 php.ini 的位置（仅适用于 windows）。
[HKEY_LOCAL_MACHINE\SOFTWARE\PHP] 内 IniFilePath 的值（Windows 注册表位置）。
当前工作目录（对于 CLI）。
web 服务器目录（对于 SAPI 模块）或 PHP 所在目录（Windows 下其它情况）。
Windows 目录（C:\windows 或 C:\winnt），或 --with-config-file-path 编译时选项指定的位置。
全局操作函数的初始化
php_startup_auto_globals函数会初始化在用户空间所使用频率很高的一些全局变量，如：$_GET、$_POST、$_FILES等。 这里只是初始化，所调用的zend_register_auto_global函数也只是将这些变量名添加到CG(auto_globals)这个变量表。

php_startup_sapi_content_types函数用来初始化SAPI对于不同类型内容的处理函数， 这里的处理函数包括POST数据默认处理函数、默认数据处理函数等。

初始化静态构建的模块和共享模块(MINIT)
php_register_internal_extensions_func函数用来注册静态构建的模块，也就是默认加载的模块， 我们可以将其认为内置模块。在PHP5.3.0版本中内置的模块包括PHP标准扩展模块（/ext/standard/目录， 这里是我们用的最频繁的函数，比如字符串函数，数学函数，数组操作函数等等），日历扩展模块、FTP扩展模块、 session扩展模块等。这些内置模块并不是一成不变的，在不同的PHP模板中，由于不同时间的需求或其它影响因素会导致这些默认加载的模块会变化， 比如从代码中我们就可以看到mysql、xml等扩展模块曾经或将来会作为内置模块出现。

模块初始化会执行两个操作： 1. 将这些模块注册到已注册模块列表（module_registry），如果注册的模块已经注册过了，PHP会报Module XXX already loaded的错误。 2. 将每个模块中包含的函数注册到函数表（ CG(function_table) ），如果函数无法添加，则会报 Unable to register functions, unable to load。

在注册了静态构建的模块后，PHP会注册附加的模块，不同的模式下可以加载不同的模块集，比如在CLI模式下是没有这些附加的模块的。

在内置模块和附加模块后，接下来是注册通过共享对象（比如DLL）和php.ini文件灵活配置的扩展。

在所有的模块都注册后，PHP会马上执行模块初始化操作（zend_startup_modules）。 它的整个过程就是依次遍历每个模块，调用每个模块的模块初始化函数， 也就是在本小节前面所说的用宏PHP_MINIT_FUNCTION包含的内容。

禁用函数和类
php_disable_functions函数用来禁用PHP的一些函数。这些被禁用的函数来自PHP的配置文件的disable_functions变量。 其禁用的过程是调用zend_disable_function函数将指定的函数名从CG(function_table)函数表中删除。

php_disable_classes函数用来禁用PHP的一些类。这些被禁用的类来自PHP的配置文件的disable_classes变量。 其禁用的过程是调用zend_disable_class函数将指定的类名从CG(class_table)类表中删除。

ACTIVATION

在处理了文件相关的内容，PHP会调用php_request_startup做请求初始化操作。 请求初始化操作，除了图中显示的调用每个模块的请求初始化函数外，还做了较多的其它工作，其主要内容如下：

激活Zend引擎
gc_reset函数用来重置垃圾收集机制，当然这是在PHP5.3之后才有的。

init_compiler函数用来初始化编译器，比如将编译过程中放在opcode里的数组清空，准备编译时需要用的数据结构等等。

init_executor函数用来初始化中间代码执行过程。 在编译过程中，函数列表、类列表等都存放在编译时的全局变量中， 在准备执行过程时，会将这些列表赋值给执行的全局变量中，如：EG(function_table) = CG(function_table); 中间代码执行是在PHP的执行虚拟栈中，初始化时这些栈等都会一起被初始化。 除了栈，还有存放变量的符号表(EG(symbol_table))会被初始化为50个元素的hashtable，存放对象的EG(objects_store)被初始化了1024个元素。 PHP的执行环境除了上面的一些变量外，还有错误处理，异常处理等等，这些都是在这里被初始化的。 通过php.ini配置的zend_extensions也是在这里被遍历调用activate函数。

激活SAPI
sapi_activate函数用来初始化SG(sapi_headers)和SG(request_info)，并且针对HTTP请求的方法设置一些内容， 比如当请求方法为HEAD时，设置SG(request_info).headers_only=1； 此函数最重要的一个操作是处理请求的数据，其最终都会调用sapi_module.default_post_reader。 而sapi_module.default_post_reader在前面的模块初始化是通过php_startup_sapi_content_types函数注册了 默认处理函数为main/php_content_types.c文件中php_default_post_reader函数。 此函数会将POST的原始数据写入$HTTP_RAW_POST_DATA变量。

在处理了post数据后，PHP会通过sapi_module.read_cookies读取cookie的值， 在CLI模式下，此函数的实现为sapi_cli_read_cookies，而在函数体中却只有一个return NULL;

如果当前模式下有设置activate函数，则运行此函数，激活SAPI，在CLI模式下此函数指针被设置为NULL。

环境初始化
这里的环境初始化是指在用户空间中需要用到的一些环境变量初始化，这里的环境包括服务器环境、请求数据环境等。 实际到我们用到的变量，就是$_POST、$_GET、$_COOKIE、$_SERVER、$_ENV、$_FILES。 和sapi_module.default_post_reader一样，sapi_module.treat_data的值也是在模块初始化时， 通过php_startup_sapi_content_types函数注册了默认数据处理函数为main/php_variables.c文件中php_default_treat_data函数。

以$_COOKIE为例，php_default_treat_data函数会对依据分隔符，将所有的cookie拆分并赋值给对应的变量。

模块请求初始化
PHP通过zend_activate_modules函数实现模块的请求初始化，也就是我们在图中看到Call each extension's RINIT。 此函数通过遍历注册在module_registry变量中的所有模块，调用其RINIT方法实现模块的请求初始化操作。

运行

php_execute_script函数包含了运行PHP脚本的全部过程。

当一个PHP文件需要解析执行时，它可能会需要执行三个文件，其中包括一个前置执行文件、当前需要执行的主文件和一个后置执行文件。 非当前的两个文件可以在php.ini文件通过auto_prepend_file参数和auto_append_file参数设置。 如果将这两个参数设置为空，则禁用对应的执行文件。

对于需要解析执行的文件，通过zend_compile_file（compile_file函数）做词法分析、语法分析和中间代码生成操作，返回此文件的所有中间代码。 如果解析的文件有生成有效的中间代码，则调用zend_execute（execute函数）执行中间代码。 如果在执行过程中出现异常并且用户有定义对这些异常的处理，则调用这些异常处理函数。 在所有的操作都处理完后，PHP通过EG(return_value_ptr_ptr)返回结果。

DEACTIVATION

PHP关闭请求的过程是一个若干个关闭操作的集合，这个集合存在于php_request_shutdown函数中。 这个集合包括如下内容：

调用所有通过register_shutdown_function()注册的函数。这些在关闭时调用的函数是在用户空间添加进来的。 一个简单的例子，我们可以在脚本出错时调用一个统一的函数，给用户一个友好一些的页面，这个有点类似于网页中的404页面。
执行所有可用的__destruct函数。 这里的析构函数包括在对象池（EG(objects_store）中的所有对象的析构函数以及EG(symbol_table)中各个元素的析构方法。
将所有的输出刷出去。
发送HTTP应答头。这也是一个输出字符串的过程，只是这个字符串可能符合某些规范。
遍历每个模块的关闭请求方法，执行模块的请求关闭操作，这就是我们在图中看到的Call each extension's RSHUTDOWN。
销毁全局变量表（PG(http_globals)）的变量。
通过zend_deactivate函数，关闭词法分析器、语法分析器和中间代码执行器。
调用每个扩展的post-RSHUTDOWN函数。只是基本每个扩展的post_deactivate_func函数指针都是NULL。
关闭SAPI，通过sapi_deactivate销毁SG(sapi_headers)、SG(request_info)等的内容。
关闭流的包装器、关闭流的过滤器。
关闭内存管理。
重新设置最大执行时间
结束

最终到了要收尾的地方了。

flush
sapi_flush将最后的内容刷新出去。其调用的是sapi_module.flush，在CLI模式下等价于fflush函数。

关闭Zend引擎
zend_shutdown将关闭Zend引擎。

此时对应图中的流程，我们应该是执行每个模块的关闭模块操作。 在这里只有一个zend_hash_graceful_reverse_destroy函数将module_registry销毁了。 当然，它最终也是调用了关闭模块的方法的，其根源在于在初始化module_registry时就设置了这个hash表析构时调用ZEND_MODULE_DTOR宏。 而ZEND_MODULE_DTOR宏对应的是module_destructor函数。 在此函数中会调用模块的module_shutdown_func方法，即PHP_RSHUTDOWN_FUNCTION宏产生的那个函数。

在关闭所有的模块后，PHP继续销毁全局函数表，销毁全局类表、销售全局变量表等。 通过zend_shutdown_extensions遍历zend_extensions所有元素，调用每个扩展的shutdown函数。

多进程SAPI生命周期
通常PHP是编译为apache的一个模块来处理PHP请求。Apache一般会采用多进程模式， Apache启动后会fork出多个子进程，每个进程的内存空间独立，每个子进程都会经过开始和结束环节， 不过每个进程的开始阶段只在进程fork出来以来后进行，在整个进程的生命周期内可能会处理多个请求。 只有在Apache关闭或者进程被结束之后才会进行关闭阶段，在这两个阶段之间会随着每个请求重复请求开始-请求关闭的环节。 如图2.2所示：

 

图2.2 多进程SAPI生命周期
图2.2 多进程SAPI生命周期
 

多线程的SAPI生命周期
多线程模式和多进程中的某个进程类似，不同的是在整个进程的生命周期内会并行的重复着 请求开始-请求关闭的环节

 

图2.3 多线程SAPI生命周期
图2.3 多线程SAPI生命周期
 

Zend引擎
Zend引擎是PHP实现的核心，提供了语言实现上的基础设施。例如：PHP的语法实现，脚本的编译运行环境， 扩展机制以及内存管理等，当然这里的PHP指的是官方的PHP实现(除了官方的实现， 目前比较知名的有facebook的hiphop实现，不过到目前为止，PHP还没有一个标准的语言规范)， 而PHP则提供了请求处理和其他Web服务器的接口(SAPI)。

目前PHP的实现和Zend引擎之间的关系非常紧密，甚至有些过于紧密了，例如很多PHP扩展都是使用的Zend API， 而Zend正是PHP语言本身的实现，PHP只是使用Zend这个内核来构建PHP语言的，而PHP扩展大都使用Zend API， 这就导致PHP的很多扩展和Zend引擎耦合在一起了，在笔者编写这本书的时候PHP核心开发者就提出将这种耦合解开，

目前PHP的受欢迎程度是毋庸置疑的，但凡流行的语言通常都会出现这个语言的其他实现版本， 这在Java社区里就非常明显，目前已经有非常多基于JVM的语言了，例如IBM的Project Zero就实现了一个基于JVM的PHP实现， .NET也有类似的实现，通常他们这样做的原因无非是因为：他们喜欢这个语言，但又不想放弃原有的平台， 或者对现有的语言实现不满意，处于性能或者语言特性等（HipHop就是这样诞生的）。

很多脚本语言中都会有语言扩展机制，PHP中的扩展通常是通过Pear库或者原生扩展，在Ruby中则这两者的界限不是很明显， 他们甚至会提供两套实现，一个主要用于在无法编译的环境下使用，而在合适的环境则使用C实现的原生扩展， 这样在效率和可移植性上都可以保证。目前这些为PHP编写的扩展通常都无法在其他的PHP实现中实现重用， HipHop的做法是对最为流行的扩展进行重写。如果PHP扩展能和ZendAPI解耦，则在其他语言中重用这些扩展也将更加容易了。

为什么必须装入Xdebug扩展作为一个Zend扩展？

什么是Zend扩展和定期PHP扩展和Zend扩展differents是什么？

让我们开始从扩展加载过程。

PHP是可以被扩展的，PHP的核心引擎Zend引擎也是可以被扩展的，如果你也对Apache模块的编写也有所了解的话，那么，你就会对如下的结构很熟悉了：

结构 _zend_extension  {
       字符 *名称;
       字符 *版本;
       字符 *作者;
       字符 *网址;
       字符 *版权;
       startup_func_t 启动;
       shutdown_func_t 关机;
       activate_func_t 激活;
       deactivate_func_t 停用;
       message_handler_func_t  MESSAGE_HANDLER ;
       op_array_handler_func_t  op_array_handler ;
       statement_handler_func_t  statement_handler ;
       fcall_begin_handler_func_t  fcall_begin_handler ;
       fcall_end_handler_func_t  fcall_end_handler ;
       op_array_ctor_func_t  op_array_ctor ;
       op_array_dtor_func_t  op_array_dtor ;
        （* api_no_check ）（INT api_no ）;
       无效 * reserved2 ;
       无效 * reserved3 ;
       无效 * reserved4 ;
       无效 * reserved5 ;
       无效 * reserved6 ;
       无效 * reserved7 ;
       无效 * reserved8 ;
       DL_HANDLE 处理;
       INT resource_number ;
} ;
然后，让我们对比下的PHP扩展的模块的入口：

    结构 _zend_module_entry  {
         无符号 短的大小;
         无符号的 诠释 zend_api ;
         无符号 字符 zend_debug ;
         无符号的 字符特稿;
          _zend_ini_entry  * ini_entry ;
          _zend_module_dep  * DEPS ;
         字符 *名称;
          _zend_function_entry  *职能;
         （  * module_startup_func ）（ INIT_FUNC_ARGS ）;
         （  * module_shutdown_func ）（ SHUTDOWN_FUNC_ARGS ）;
         （  * request_startup_func ）（ INIT_FUNC_ARGS ）;
         （  * request_shutdown_func ）（ SHUTDOWN_FUNC_ARGS ）;
          （* info_func ）（ ZEND_MODULE_INFO_FUNC_ARGS ）无效;
         字符 *版本;
         为size_t  globals_size ;
    ＃IFDEF特稿
         ts_rsrc_id * globals_id_ptr ;
    ＃ELSE
         * globals_ptr 无效;
    ＃ENDIF
         无效 （* globals_ctor ）（无效 * 全球  TSRMLS_DC ）;
         无效 （* globals_dtor ）（无效 * 全球  TSRMLS_DC ）;
         INT  （* post_deactivate_func ）（无效）;
         INT module_started ;
         无符号的 字符类型;
         无效 *处理;
         INT module_number ;
    } ;
上面的结构，可以结合我之前的处理用的C / C + +扩展你的PHP来帮助理解。

恩，回到主题：既然Xdebug的要以Zend扩展方式加载，那么栭必然有基于Zend扩展的需求，会是什么呢？

恩，我们知道的Xdebug有个人资料的PHP的功能，对，就是statement_handler：
语句处理程序回调插入一个额外的操作码，在每一个回调调用的脚本语句的结束。这种回调的主要用途之一是要落实每行的分析，加强调试器，代码覆盖工具 。

并且，因为Xdebug的也提供了给使用者脚本使用的函数，所以，它也会有部分PHP扩展的实现，

最后，将PHP扩展的载入过程罗列如下（我会慢慢加上注释），当然，如果你等不及想知道，也欢迎你直接在我的博客风雪之隅留言探讨。

以apache/mod_php5.c为例

1。在mod_php5.c中，定义了Apache的模块结构：

    模块 MODULE_VAR_EXPORT  php5_module =
    {
         STANDARD_MODULE_STUFF ，
         php_init_handler ，          / *初始化程序* /
         php_create_dir ，            / *每个目录配置的创造者* /
         php_merge_dir ，             / *目录合并* /
         NULL ，                      / *每个服务器配置的创造者* /
         NULL ，                      / *合并服务器配置* /
         php_commands ，              / *命令表* /
         php_handlers ，              / *处理程序* /
         NULL ，                      / *文件名 ​​翻译* /
         NULL ，                      / * check_user_id * /
         NULL ，                      / *检查AUTH * /
         NULL ，                      / *检查访问* /
         NULL ，                      / * type_checker * /
         NULL ，                      / *链接地址* /
         NULL，                        / *记录* /
    ＃如果 MODULE_MAGIC_NUMBER > = 19970103
         NULL                      / *头分析器* /
    ＃ENDIF
    ＃如果 MODULE_MAGIC_NUMBER > = 19970719
         NULL                      / * child_init * /
    ＃ENDIF
    ＃如果 MODULE_MAGIC_NUMBER > = 19970728
         php_child_exit_handler        / * child_exit * /
    ＃ENDIF
    ＃如果 MODULE_MAGIC_NUMBER > = 19970902
         ， NULL                      / *后读请求* /
    ＃ENDIF
    } ;
/ *}}} * /
可见，最开始被调用的将会是php_init_handler，

静态 无效 php_init_handler （server_rec  * 游泳池* P ）  
{
   register_cleanup(p, NULL, (void (*)(void *))apache_php_module_shutdown_wrapper, (void (*)(void *))php_module_shutdown_for_exec);
    if (!apache_php_initialized) {
       apache_php_initialized = 1;
#ifdef ZTS
       tsrm_startup(1, 1, 0, NULL);
#endif
       sapi_startup(&apache_sapi_module);
       php_apache_startup(&apache_sapi_module);
    }
#if MODULE_MAGIC_NUMBER >= 19980527
    {
       TSRMLS_FETCH();
       if (PG(expose_php)) {
           ap_add_version_component("PHP/" PHP_VERSION);
       }
    }
#endif
}
这里, 调用了sapi_startup, 这部分是初始化php的apache sapi,
然后是调用,php_apache_startup:

static int php_apache_startup(sapi_module_struct *sapi_module)
{
    if (php_module_startup(sapi_module, &apache_module_entry, 1) == FAILURE) {
       return FAILURE;
    } else {
       return SUCCESS;
    }
}
这个时候,调用了php_module_startup, 其中有:


    if (php_init_config(TSRMLS_C) == FAILURE) {
       return FAILURE;
    }
调用了php_init_config, 这部分读取所有的php.ini和关联的ini文件, 然后对于每一条配置指令调用:

    ....
 if (sapi_module.ini_entries) {
       zend_parse_ini_string(sapi_module.ini_entries, 1, php_config_ini_parser_cb, &extension_lists);
    }
然后在php_config_ini_parser_cb中:
              if (!strcasecmp(Z_STRVAL_P(arg1), "extension")) {
                   zval copy;
 
                   copy = *arg2;
                   zval_copy_ctor(&copy);
                   copy.refcount = 0;
                   zend_llist_add_element(&extension_lists.functions, &copy);
               } else if (!strcasecmp(Z_STRVAL_P(arg1), ZEND_EXTENSION_TOKEN)) {
                   char *extension_name = estrndup(Z_STRVAL_P(arg2), Z_STRLEN_P(arg2));
 
                   zend_llist_add_element(&extension_lists.engine, &extension_name);
               } else {
                   zend_hash_update(&configuration_hash, Z_STRVAL_P(arg1), Z_STRLEN_P(arg1) + 1, arg2, sizeof(zval), (void **) &entry);
                   Z_STRVAL_P(entry) = zend_strndup(Z_STRVAL_P(entry), Z_STRLEN_P(entry));
               }
这里记录下来所有要载入的php extension和zend extension,
然后, 让我们回到php_module_startup, 后面有调用到了
php_ini_register_extensions(TSRMLS_C);

void php_ini_register_extensions(TSRMLS_D)
{
   zend_llist_apply(&extension_lists.engine, php_load_zend_extension_cb TSRMLS_CC);
   zend_llist_apply(&extension_lists.functions, php_load_function_extension_cb TSRMLS_CC);
 
   zend_llist_destroy(&extension_lists.engine);
   zend_llist_destroy(&extension_lists.functions);
}
我们可以看到, 对于每一个扩展记录, 都调用了一个回叫函数, 我们这里只看php_load_function_extension_cb:

static void php_load_function_extension_cb(void *arg TSRMLS_DC)
{
    zval *extension = (zval *) arg;
    zval zval;
 
   php_dl(extension, MODULE_PERSISTENT, &zval, 0 TSRMLS_CC);
}
最后, 就是核心的载入逻辑了:

void php_dl(zval *file, int type, zval *return_value, int start_now TSRMLS_DC)
{
       void *handle;
       char *libpath;
       zend_module_entry *module_entry;
       zend_module_entry *(*get_module)(void);
       int error_type;
       char *extension_dir;
 
       if (type == MODULE_PERSISTENT) {
               extension_dir = INI_STR("extension_dir");
       } else {
               extension_dir = PG(extension_dir);
       }
 
       if (type == MODULE_TEMPORARY) {
               error_type = E_WARNING;
       } else {
               error_type = E_CORE_WARNING;
       }
 
       if (extension_dir && extension_dir[0]){
               int extension_dir_len = strlen(extension_dir);
 
               if (type == MODULE_TEMPORARY) {
                       if (strchr(Z_STRVAL_P(file), '/') != NULL || strchr(Z_STRVAL_P(file), DEFAULT_SLASH) != NULL) {
                               php_error_docref(NULL TSRMLS_CC, E_WARNING, "Temporary module name should contain only filename");
                               RETURN_FALSE;
                       }
               }
 
               if (IS_SLASH(extension_dir[extension_dir_len-1])) {
                       spprintf(&libpath, 0, "%s%s", extension_dir, Z_STRVAL_P(file));
               } else {
                       spprintf(&libpath, 0, "%s%c%s", extension_dir, DEFAULT_SLASH, Z_STRVAL_P(file));
               }
       } else {
               libpath = estrndup(Z_STRVAL_P(file), Z_STRLEN_P(file));
       }
 
       
       handle = DL_LOAD(libpath);
       if (!handle) {
               php_error_docref(NULL TSRMLS_CC, error_type, "Unable to load dynamic library '%s' - %s", libpath, GET_DL_ERROR());
               GET_DL_ERROR();
               efree(libpath);
               RETURN_FALSE;
       }
 
       efree(libpath);
 
       get_module = (zend_module_entry *(*)(void)) DL_FETCH_SYMBOL(handle, "get_module");
 
       
 
       if (!get_module)
               get_module = (zend_module_entry *(*)(void)) DL_FETCH_SYMBOL(handle, "_get_module");
 
       if (!get_module) {
               DL_UNLOAD(handle);
               php_error_docref(NULL TSRMLS_CC, error_type, "Invalid library (maybe not a PHP library) '%s' ", Z_STRVAL_P(file));
               RETURN_FALSE;
       }
       module_entry = get_module();
       if ((module_entry->zend_debug != ZEND_DEBUG) || (module_entry->zts != USING_ZTS)
               || (module_entry->zend_api != ZEND_MODULE_API_NO)) {
               
                       struct pre_4_1_0_module_entry {
                                 char *name;
                                 zend_function_entry *functions;
                                 int (*module_startup_func)(INIT_FUNC_ARGS);
                                 int (*module_shutdown_func)(SHUTDOWN_FUNC_ARGS);
                                 int (*request_startup_func)(INIT_FUNC_ARGS);
                                 int (*request_shutdown_func)(SHUTDOWN_FUNC_ARGS);
                                 void (*info_func)(ZEND_MODULE_INFO_FUNC_ARGS);
                                 int (*global_startup_func)(void);
                                 int (*global_shutdown_func)(void);
                                 int globals_id;
                                 int module_started;
                                 unsigned char type;
                                 void *handle;
                                 int module_number;
                                 unsigned char zend_debug;
                                 unsigned char zts;
                                 unsigned int zend_api;
                       };
 
                       char *name;
                       int zend_api;
                       unsigned char zend_debug, zts;
 
                       if ((((struct pre_4_1_0_module_entry *)module_entry)->zend_api > 20000000) &&
                               (((struct pre_4_1_0_module_entry *)module_entry)->zend_api < 20010901)
                       ) {
                               name      = ((struct pre_4_1_0_module_entry *)module_entry)->name;
                               zend_api   = ((struct pre_4_1_0_module_entry *)module_entry)->zend_api;
                               zend_debug = ((struct pre_4_1_0_module_entry *)module_entry)->zend_debug;
                               zts       = ((struct pre_4_1_0_module_entry *)module_entry)->zts;
                       } else {
                               name      = module_entry->name;
                               zend_api   = module_entry->zend_api;
                               zend_debug = module_entry->zend_debug;
                               zts       = module_entry->zts;
                       }
 
                       php_error_docref(NULL TSRMLS_CC, error_type,
                                         "%s: Unable to initialize module\n"
                                         "Module compiled with module API=%d, debug=%d, thread-safety=%d\n"
                                         "PHP    compiled with module API=%d, debug=%d, thread-safety=%d\n"
                                         "These options need to match\n",
                                         name, zend_api, zend_debug, zts,
                                         ZEND_MODULE_API_NO, ZEND_DEBUG, USING_ZTS);
                       DL_UNLOAD(handle);
                       RETURN_FALSE;
       }
       module_entry->type = type;
       module_entry->module_number = zend_next_free_module();
       module_entry->handle = handle;
 
       if ((module_entry = zend_register_module_ex(module_entry TSRMLS_CC)) == NULL) {
               DL_UNLOAD(handle);
               RETURN_FALSE;
       }
 
       if ((type == MODULE_TEMPORARY || start_now) && zend_startup_module_ex(module_entry TSRMLS_CC) == FAILURE) {
               DL_UNLOAD(handle);
               RETURN_FALSE;
       }
 
       if ((type == MODULE_TEMPORARY || start_now) && module_entry->request_startup_func) {
               if (module_entry->request_startup_func(type, module_entry->module_number TSRMLS_CC) == FAILURE) {
                       php_error_docref(NULL TSRMLS_CC, error_type, "Unable to initialize module '%s'", module_entry->name);
                       DL_UNLOAD(handle);
                       RETURN_FALSE;
               }
       }
       RETURN_TRUE;
}

1. 扩展的加载顺序是和它出现在配置文件中的先后顺序相关的, 也就是说, 如果在配置文件中的顺序如下,

extension=mysql.so
extension=pdo.so
   
那么, mysql扩展就会比pdo扩展先载入.

同理,对于单独的配置文件,则和这个文件的载入顺序相关. 一般来说,这个时候和这个文件的命名相关.

2. 那么如果顺序出错, 我们又要怎么保证正确的加载, 或者告诉Zend此时出错了呢?

回忆, 我们在写扩展的时候, 都会申明一个zend_module_entry结构的扩展说明模块, 这个结构用来告诉Zend所有和当前模块相关的, Zend关心的信息, 也就是在这里, 我们可以申明我们的模块所依赖的模块, 这样当依赖关系出错的时候Zend就会报错, 给出详细的出错信息, 并停止运行.

struct _zend_module_entry {
    unsigned short size;
    unsigned int zend_api;
    unsigned char zend_debug;
    unsigned char zts;
    struct _zend_ini_entry *ini_entry;
    struct _zend_module_dep *deps;  //关键属性
    char *name;
    struct _zend_function_entry *functions;
    int (*module_startup_func)(INIT_FUNC_ARGS);
    int (*module_shutdown_func)(SHUTDOWN_FUNC_ARGS);
    int (*request_startup_func)(INIT_FUNC_ARGS);
    int (*request_shutdown_func)(SHUTDOWN_FUNC_ARGS);
    void (*info_func)(ZEND_MODULE_INFO_FUNC_ARGS);
    char *version;
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
};
  
第5个属性, zend_module_dep是一个接受zend_module_dep结构数组的指针, 用来指出当前模块所以来的所有模块, 而zend_module_dep的结构如下:

struct _zend_module_dep {
    char *name;         /* module name */
    char *rel;          /* version relationship: NULL (exists), lt|le|eq|ge|gt (to given version) */
    char *version;      /* version */
    unsigned char type; /* dependency type */
};
 
但, 很PHP特色的, 我们不需要知道详细的这个结构细节, 也不需要去和这个结构直接接触, 我们可以使用PHP提供的宏:GET_MOD_REQUIRE来完成zend_module_dep的填写,　比如在pdo_mysql中:

static zend_module_dep pdo_mysql_deps[] = {
    ZEND_MOD_REQUIRED("pdo")
    {NULL, NULL, NULL}
};
 
申明了, pdo_mysql必须依赖于pdo扩展模块.

使用PHP扩展的原因：
①如果应用注重效率，使用非常复杂的算法，推荐使用PHP扩展。
②有些系统调用PHP不能直接访问（如Linux的fork()函数创建进程），需要编写成PHP扩展。
③应用不想暴露关键代码，可以创建扩展使用。
准备工作
一：了解PHP源码目录
网上下载下来PHP 5.4版本源代码，目录结构如下：

php-5.4.30
  |____build    --和编译有关的目录，里面包括wk，awk和sh脚本用于编译处理，其中m4文件是linux下编译程序自动生成的文件，可以使用buildconf命令操作具体的配置文件。
  |____ext      --扩展库代码，例如Mysql，gd，zlib，xml，iconv 等我们熟悉的扩展库，ext_skel是linux下扩展生成脚本，windows下使用ext_skel_win32.php。
  |____main     --主目录，包含PHP的主要宏定义文件，php.h包含绝大部分PHP宏及PHP API定义。
  |____netware  --网络目录，只有sendmail_nw.h和start.c，分别定义SOCK通信所需要的头文件和具体实现。
  |____pear     --扩展包目录，PHP Extension and Application Repository。
  |____sapi     --各种服务器的接口调用，如Apache，IIS等。
  |____scripts  --linux下的脚本目录。
  |____tests    --测试脚本目录，主要是phpt脚本，由--TEST--，--POST--，--FILE--，--EXPECT--组成，需要初始化可添加--INI--部分。
  |____TSRM     --线程安全资源管理器，Thread Safe Resource Manager保证在单线程和多线程模型下的线程安全和代码一致性。
  |____win32    --Windows下编译PHP 有关的脚本。
  |____Zend     --包含Zend引擎的所有文件，包括PHP的生命周期，内存管理，变量定义和赋值以及函数宏定义等等。
二：自动构建工具
本篇针对Linux环境下创建PHP扩展，使用扩展自动构建工具为ext_skel，Windows下使用ext_skel_win32.php，构建方式略有不同，其余开发无差别。
构建PHP扩展的步骤如下（不唯一）：

①cd php_src/ext
②./ext_skel --extname=XXX
    此时当前目录下会生成一个名为XXX的文件夹
③cd XXX/
④vim config.m4
    会有这段文字：
    dnl If your extension references something external, use with:
    dnl PHP_ARG_WITH(say, for say support,
    dnl Make sure that the comment is aligned:
    dnl [  --with-say             Include say support])
    dnl Otherwise use enable:
    dnl PHP_ARG_ENABLE(say, whether to enable say support,
    dnl Make sure that the comment is aligned:
    dnl [  --enable-say           Enable say support])
其中，dnl 是注释符号。上面的代码说，如果你所编写的扩展如果依赖其它的扩展或者lib库，需要去掉PHP_ARG_WITH相关代码的注释。否则，去掉 PHP_ARG_ENABLE 相关代码段的注释。本篇的扩展不依赖其他扩展，故修改为：
    dnl If your extension references something external, use with:
    dnl PHP_ARG_WITH(say, for say support,
    dnl Make sure that the comment is aligned:
    dnl [  --with-say             Include say support])
    dnl Otherwise use enable:
    PHP_ARG_ENABLE(say, whether to enable say support,
    Make sure that the comment is aligned:
    [  --enable-XXX           Enable say support])
⑤在XXX.c中具体实现
⑥编译安装
    phpize
    ./configure --with-php-config=php_path/bin/php-config
    make && make install
⑦修改php.ini文件
    增加
    [XXX]
    extension = XXX.so
三：了解PHP生命周期
任何一个PHP实例都会经过Module init、Request init、Request shutdown和Module shutdown四个过程。
1.Module init
在所有请求到达前发生，例如启动Apache服务器，PHP解释器随之启动，相关的各个模块（Redis、Mysql等）的MINIT方法被调用。仅被调用一次。创建XXX扩展后，相应的XXX.c文件中将自动生成该方法：

PHP_MINIT_FUNCTION(XXX) {  
    return SUCCESS;   
}
2.Request init
每个请求达到时都被触发。SAPI层将控制权交由PHP层，PHP初始化本次请求执行脚本所需的环境变量，函数列表等，调用所有模块的RINIT函数。XXX.c中对应函数如下:

PHP_RINIT_FUNCTION(XXX){
    return SUCCESS;
}
3.Request shutdown
每个请求结束，PHP就会自动清理程序，顺序调用各个模块的RSHUTDOWN方法，清除程序运行期间的符号表。典型的RSHUTDOWN方法如：

PHP_RSHUTDOWN_FUNCTION(XXX){
    return SUCCESS;
}
4.Module shutdown
所有请求处理完毕后，SAPI也关闭了（即服务器关闭），PHP调用各个模块的MSHUTDOWN方法释放内存。

PHP_MSHUTDOWN_FUNCTION(XXX){
    return SUCCESS;
}
PHP的生命周期常见如下几种

①单进程SAPI生命周期
②多进程SAPI生命周期
③多线程SAPI声明周期
这与PHP的运行模式有很大关系，常见的运行模式有CLI、CGI、FastCGI和mod_php。

①CLI模式——单进程SAPI生命周期
所谓CLI模式，即在终端界面通过php+文件名的方式执行PHP文件


单进程SAPI生命周期

输入命令后，依次调用MINIT，RINIT，RSHUTDOWN，MSHUTDOWN即完成生命周期，一次只处理一个请求。

②CGI模式——单进程SAPI生命周期
和CLI模式一样，请求到达时，为每个请求fork一个进程，一个进程只对一个请求做出响应，请求结束后，进程也就结束了。
与CLI模式不同的是，CGI可以看作是规定了Web Server与PHP的交流规则，相当于执行response = exec("php -f xxx.php -url=xxx -cookie=xxx -xxx=xxx")。

③FastCGI模式——多进程SAPI生命周期
CGI模式存在明显缺点，每个进程处理一个请求及结束，新请求过来需要重新加载php.ini，调用MINIT等函数。FastCGI相当于可以执行多个请求的CGI，处理完一个请求后进程不结束，等待下一个请求到来。
服务启动时，FastCGI先启动多个子进程等待处理请求，避免了CGI模式请求到来时fork()进程（即fork-and-execute），提高效率。


多进程SAPI生命周期
④mod_php模式——多进程SAPI生命周期
该模式将PHP嵌入到Apache中，相当于给Apache增加了解析PHP的功能。PHP随服务器的启动而启动，两者之间存在从属关系。
证明：
CGI模式下，修改php.ini无需重启服务器，每个请求结束后，进程自动结束，新请求到来时会重新读取php.ini文件创建新进程。
mod_php下，进程是启动即创建，只有结束现有进程，重新启动服务器读取PHP配置创建新进程，修改才有效。

多线程SAPI模式
多线程模式和多进程模式的某个进程类似，在整个生命周期中会并行重复着请求开始，请求结束的环节。
只有一个服务器进程运行，但同时运行多个线程，优点是节省资源开销，MINIT和MSHUTDOWN只需在Web Server启动和结束时执行一次。由于线程的特质，使得各个请求之间共享数据成为可能。


多线程SAPI模式
四：PHP内核中的变量
PHP变量的弱类型实现在之前的文章中讲到，可以参读：https://www.jianshu.com/p/ef0c91be06a0 PHP的实现方式即PHP变量在内核中的存储。
PHP提供了一系列内核变量的访问宏。推荐使用它们设置和访问PHP的变量类型和值。
变量类型：

Z_TYPE(zval)  可以获取和设置变量类型
Z_TYPE(zval)函数返回变量的类型，PHP变量类型有：
IS_NULL(空类型)，IS_LONG(整型)，IS_DOUBLE(浮点型)，IS_STRING(字符串)，
IS_ARRAY(数组类型)，IS_OBJECT(对象类型)，IS_BOOL(布尔类型)，IS_RESOURCE(资源类型)

可以通过
Z_TYPE(zval) = IS_STRING的方式直接设置变量类型
变量值对应的访问宏：

整数类型  Z_LVAL(zval)  对应zval的实体；Z_VAL_P(&zval)  对应结构体的指针；Z_VAL_PP(&&zval)  对应结构体的二级指针
浮点数类型 Z_DVAL(zval)    Z_DVAL_P(&zval)    Z_DVAL_PP(&&zval)
布尔类型 Z_BVAL(zval)    Z_BVAL_P(&zval)    Z_BVAL_PP(&&zval)
字符串类型 
    获取值：Z_STRVAL(zval)    Z_STRVAL_P(&zval)    Z_STRVAL_PP(&&zval)
    获取长度：Z_STRLEN(zval)    Z_STRLEN_P(&zval)    Z_STRLEN_PP(&&zval)
数组类型 Z_ARRVAL(zval)    Z_ARRVAL_P(&zval)    Z_ARRVAL_PP(&&zval)
资源类型 Z_RESVAL(zval)    Z_RESVAL_P(&zval)    Z_RESVAL_PP(&&zval)
五：了解Zend API
1.Zend引擎
Zend引擎就是脚本语言引擎（解释器+虚拟机），负责解析、翻译和执行PHP脚本。其工作流程大致如下：

①Zend Engine Compiler编译PHP脚本为Opcode
②Opcode由Zend Engine Executor解析执行，期间Zend Engine Executor负责调用使用到的PHP extension
2.Zend内存管理
使用C语言开发PHP扩展，需要注意内存管理。忘记释放内存将造成内存泄漏，释放多次则产生系统错误。Zend引擎提供了一些内存管理的接口，使用这些接口申请内存交由Zend管理。
常见接口：

emalloc(size_t size)    申请size大小的内存
efree(void *ptr)    释放ptr指向的内存块
estrdup(char *str)    申请str大小的内存，并将str内容复制进去
estrndup(char *str, int slen)    同上，但指定长度复制
ecalloc(size_t numOfElem, size_t sizeOfElem)    复制numOfElem个sizeOfElem大小的内存块
erealloc(void *ptr, size_t nsize)    ptr指向内存块的大小扩大到nsize
内存管理申请的所有内存，将在脚本执行完毕和处理请求终止时被释放。

3.PHP扩展架构
使用准备工作（二）中命令生成基本架构后，生成的对应目录中会有两个文件，php_XXX.h和XXX.c，其中php_XXX.h文件用于声明扩展的一些基本信息和实现的函数，注意，只是声明。具体的实现在XXX.c中。
php_XXX.h大致结构如下：

#ifndef PHP_XXX_H
#define PHP_XXX_H

extern zend_module_entry php_list_module_entry;
#define phpext_php_list_ptr &php_list_module_entry

#define PHP_XXX_VERSION "0.1.0"

#ifdef PHP_WIN32
#   define PHP_XXX_API __declspec(dllexport)
#elif defined(__GNUC__) && __GNUC__ >= 4
#   define PHP_XXX_API __attribute__ ((visibility("default")))
#else
#   define PHP_XXX_API
#endif
#ifdef ZTS
#include "TSRM.h"
#endif

PHP_MINIT_FUNCTION(XXX);
PHP_MSHUTDOWN_FUNCTION(XXX);
PHP_RINIT_FUNCTION(XXX);
PHP_RSHUTDOWN_FUNCTION(XXX);
PHP_MINFO_FUNCTION(XXX);

PHP_FUNCTION(confirm_XXX_compiled);

#ifdef ZTS
#define PHP_LIST_G(v) TSRMG(php_list_globals_id, zend_php_list_globals *, v)
#else
#define PHP_LIST_G(v) (php_list_globals.v)
#endif
#endif
大致信息有版本号，MINIT、RINIT、RSHUTDOWN、MSHUTDOWN函数等，如果声明自定义函数，可以在之后以PHP_FUNCTION(XXX);的方式声明函数，并在XXX.c中具体实现。
XXX.c内容如下:

---------------------------------------头文件--------------------------------------
#ifdef HAVE_CONFIG_H
#include "config.h"
#endif

#include "php.h"
#include "php_ini.h"
#include "ext/standard/info.h"
#include "php_XXX.h"

static int le_XXX;
---------------------------------------Zend函数快--------------------------------------
const zend_function_entry XXX_functions[] = {
    PHP_FE(confirm_php_list_compiled,   NULL)       /* For testing, remove later. */
    PHP_FE_END  /* Must be the last line in php_list_functions[] */
};
---------------------------------------Zend模块--------------------------------------
zend_module_entry XXX_module_entry = {
#if ZEND_MODULE_API_NO >= 20010901
    STANDARD_MODULE_HEADER,
#endif
    "XXX",
    XXX_functions,
    PHP_MINIT(XXX),
    PHP_MSHUTDOWN(XXX),
    PHP_RINIT(XXX),     /* Replace with NULL if there's nothing to do at request start */
    PHP_RSHUTDOWN(XXX), /* Replace with NULL if there's nothing to do at request end */
    PHP_MINFO(XXX),
#if ZEND_MODULE_API_NO >= 20010901
    PHP_XXX_VERSION,
#endif
    STANDARD_MODULE_PROPERTIES
};
---------------------------------------实现get_module函数--------------------------------------
#ifdef COMPILE_DL_XXX
ZEND_GET_MODULE(XXX)
#endif
---------------------------------------生命周期函数--------------------------------------
PHP_MINIT_FUNCTION(XXX)
{
    return SUCCESS;
}

PHP_MSHUTDOWN_FUNCTION(XXX)
{
    return SUCCESS;
}

PHP_RINIT_FUNCTION(XXX)
{
    return SUCCESS;
}

PHP_RSHUTDOWN_FUNCTION(XXX)
{
    return SUCCESS;
}
---------------------------------------导出函数--------------------------------------
PHP_FUNCTION(confirm_XXX_compiled)
{
    char *arg = NULL;
    int arg_len, len;
    char *strg;

    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "s", &arg, &arg_len) == FAILURE) {
        return;
    }

    len = spprintf(&strg, 0, "Congratulations! You have successfully modified ext/%.78s/config.m4. Module %.78s is now compiled into PHP.", "php_list", arg);
    RETURN_STRINGL(strg, len, 0);
}
---------------------------------------负责扩展info显示--------------------------------------
PHP_MINFO_FUNCTION(XXX)
{
    php_info_print_table_start();
    php_info_print_table_header(2, "XXX support", "enabled");
    php_info_print_table_end();
}
①扩展头文件
所有扩展务必包含的头文件有且只有一个——php.h，以使用PHP定义的各种宏和API。如果在php_XXX.h中声明了扩展相关的宏和函数，需要将其引入。
②导出函数
先说导出函数，方便理解Zend函数块。PHP能够调用扩展中的类和方法都是通过导出函数实现。导出函数就是按照PHP内核要求编写的函数，形式如下：

void zif_extfunction(){  //extfunction即为扩展中实现的函数
    int ht;    //函数参数的个数
    zval *return_value;    //保存函数的返回值
    zval *this_ptr;    //指向函数所在对象
    int return_value_used;    //函数返回值脚本是否使用，0——不使用；1——使用
    zend_executor_globals *executor_globals;    //指向Zend引擎的全局设置
}
有如上定义，PHP脚本中即可使用扩展函数

<?php
    extfunction();
?>
由于导出函数格式固定，Zend引擎通过PHP_FUNCTION()宏声明

PHP_FUNCTION(extfunction);
即可产生之前代码块中的导出函数结构体。
③Zend函数块
作用是将编写的函数引入Zend引擎，通过zend_function_entry结构体引入。zend_function_entry结构体声明如下：

typedef struct _zend_function_entry{
    char *fname;    //指定在PHP脚本里定义的函数名
    void (*handler)(INTERNAL_FUNCTION_PARAMETERS);  //指向导出函数的句柄
    unsigned char *func_arg_types;  //标示一些参数是否强制按引用方式传递，通常设为NULL
} zend_function_entry;
Zend引擎通过zend_function_entry数组将导出函数引入内部。方式：

zend_function_entry XXX_functions[] = {
    PHP_FE(confirm_php_list_compiled,   NULL)       /* For testing, remove later. */
    PHP_FE_END  /* Must be the last line in php_list_functions[] */
};
PHP_FE宏会把zend_function_entry结构补充完整。

PHP_FE(extfunction);  ===>  ("extfunction", zif_extfunction, NULL);
PHP_FE_END是告知Zend引擎Zend函数块到此为止，有的版本可以使用{NULL, NULL, NULL}的方式。但推荐使用本文方式以兼容。
④Zend模块声明
Zend模块包含所有需要向Zend引擎提供的扩展模块信息，底层由zend_module_entry结构体大体实现

typedef struct _zend_module_entry{
    unsigned short size;       ---|
    unsigned int zend_api;        |
    unsigned char zend_debug;     |--> 通常用STANDARD_MODULE_HEADER填充
    unsigned char zts;         ---|
    char *name;    //模块名
    zend_function_entry *functions;    //函数声明
    int (*module_start_func)(INIT_FUNC_ARGS);    //MINIT函数
    int (*module_shutdown_func)(SHUTDOWN_FUNC_ARGS);    //MSHUTDOWN函数
    int (*request_start_func)(INIT_FUNC_ARGS);    //RINIT函数
    int (*request_shutdown_func)(SHUTDOWN_FUNC_ARGS);    //RSHUTDOWN函数
    char *version;    //版本号
    ...    //其余信息不做讨论
} zend_module_entry;
对比生成代码中的模块声明和①②③中所讲，Zend通过模块声明将所有信息读取并加入到引擎中。
⑤get_module函数
当扩展被加载时，调用get_module函数，该函数返回一个指向扩展模块声明的zend_module_entry指针。是PHP内核与扩展通信的渠道。
get_module函数被条件宏包围，故有些情况下不会执行get_module方法，当扩展被编译为PHP内建模块时get_module方法不被实现。
⑥实现导出函数
在php_XXX.h中声明的函数在XXX.c中具体实现，实现方式如自动生成的confirm_XXX_compiled导出函数形式一致

PHP_FUNCTION(extfunction){
    ...具体实现...
}
在函数中获取参数和返回结果，后续讲解。
⑦模块信息函数
PHP通过phpinfo查看PHP及其扩展信息，PHP_MINFO_FUNCTION负责实现。生成代码函数体是最基本的模块信息，可自行设置显示内容。

4.导出函数实现须知
这部分主要讲解函数具体实现过程中对参数和变量的处理。

（1）获取参数个数
通过ZEND_NUM_ARGS宏获取参数个数，这个宏实际上是获取zif_extfunction的ht字段，定义在Zend/zned_API.h下

#define ZEND_NUM_ARGS()  (ht)
（2）取得参数实体
Zend引擎提供获取参数实体的API，声明如下：

int zend_parse_parameters(int num_args TSRMLS_CC, char *type_spec)
num_args：传入的参数个数
type_spec：参数的类型，每种类型对应一个字符，当num_args>1时，接收参数通过字符串依次指明类型接收。
该函数成功将返回SUCCESS，失败返回FAILURE。
可以接受的参数类型如下：

普通：
l : 长整型
d : 双精度浮点类型
s : 字符串类型及其长度（需要两个变量保存！！！）
b : 布尔类型
r : 资源类型，保存在zva *l中
a : 数组，保存在zval *中
o : 对象（任何类型），保存在zval *中
O : 对象（class entry指定类型），保存在zval *中
z : zval *
特殊：
| : 当前及之后的参数为可选参数，有传即获取，否则设为默认值
/ : 当前及之后的参数将以SEPARATE_IF_NOT_REF的方式进行拷贝，除非是引用
! : 当前及之后的参数允许为NULL，仅用在r,a,o,O,z类型时
例：获取一个字符串和一个布尔型参数

char *str;
int strlen;
zend_bool b;
if(zend_parse_paramsters(ZEND_NUM_ARGS() TSRMLS_CC, "sb", &str, &strlen, &b) == FAILURE){  //字符串需要接收内容和长度
    return ;
}
例：获取一个数组和一个可选的长整型参数

zval *arr;
long l;
if(zend_parse_parameter(ZEND_NUM_ARGS() TSRMLS_CC, "a|l", &arr, &l) == FAILURE){
    return ;
}
例：获取一个对象或NULL

zval *obj;
if(zend_parse_parameter(ZEND_NUM_ARGS() TSRMLS_CC, "o!", &obj) == FAILURE){
    return ;
}
（3）获取可变参数
开发过程中会遇到方法可接受可变参数的情况，不指定类型的情况下接收参数的方式：

int num_arg = ZEND_NUM_ARGS();
zval **parameters[num_args];
if(zend_get_parameters_array_ex(num_arg, parameters) == FAILURE){
    return ;
}
（4）参数类型转换
（3）中PHP可以接受任意类型的参数，可能会导致在具体实现过程中出问题，因此Zend提供了一系列参数类型转换的API。


参数类型转换API
（5）处理通过引用传递的参数
"z"代表的zval类型的传参即为引用传递。PHP规定修改非引用传递的参数值不会影响原来变量的值，但PHP内核采用引用传递方式传参。PHP内核使用"zval分离"的方式避免这一问题。"zval分离"即写时复制，修改非引用类型的参数时，先复制一份新值，然后将引用指向新值，修改参数时不会影响原值。
判断参数是否为引用，通过PZVAL_IS_REF(zval *)，其定义为：

#define PZVAL_IS_REF(z) ((z)->is_ref)
使用宏SEPARATE_ZVAL(zval **)实现zval分离。

（6）扩展中创建变量
PHP扩展中创建变量需要以下三步：

①创建一个zval容器
②对zval容器进行填充
③引入到Zend引擎内部符号表中

//创建zval容器
zval *new_var;
//初始化和填充
MAKE_STD_ZVAL(new_var);
//引入符号表
ZEND_SET_SYMBOL(EG(active_symbol_table), "new_var", new_var);
MAKE_STD_ZVAL()宏作用是通过ALLOC_ZVAL()申请一个zval空间，之后通过INIT_ZVAL()进行初始化。

#define MAKE_STD_ZVAL(zv) \
  ALLOC_ZVAL(zv); \
  INIT_ZVAL(zv);

INIT_ZVAL()宏定义如下：
#DEFINE INIT_ZVAL(z) \
  (z) -> refcount = 1; \
  (z) -> is_ref = 0;
MAKE_STD_ZVAL()只是为变量分配了内存，设置了refcount和is_ref两个属性。
ZEND_SET_SYMBOL()宏将变量引入到符号表中，引入时先检查变量是否已经存在于表中，如果已经存在，销毁原有的zval并替换。
如果创建的是全局变量，前两步不变，只对引入操作做调整。局部变量引入active_symbol_table中，全局变量引入symbol_table中，通过

ZEND_SET_SYMBOL(&EG(symbol_table), "new_var", new_var);
注意：active_symbol_table是个指针，symbol_table不是指针，需要增加&取地址。
···
扩展中：
PHP_FUNCTION(extfunction){
zval *new_var;
ZEND_STD_ZVAL(new_var);
ZVAL_LONG(new_var, 10);
ZEND_SET_SYMBOL(&EG(symbol_table), "new_var", new_var);
}

PHP脚本中：
<?php
extfunction();
var_dump($new_var);
?>
结果输出：10
···

（7）变量赋值
①长整型（整型）赋值
PHP中所有整型都是保存在zval的value字段中，整数保存在value联合体的lval字段中，type为IS_LONG，赋值通过宏操作进行：

ZVAL_LONG(zval, 10);
②双精度浮点数类型赋值
浮点数保存在value的dval中，type对应IS_DOUBLE，通过宏操作

ZVAL_DOUBLE(zval, 3.14);
③字符串类型
value联合体的str结构体保存字符串值，val保存字符串，len保存长度，type为IS_STRING。

char *str = "hello world";
ZVAL_STRING(zval, str, 1);  //结尾参数表示字符串是否需要被复制。
④布尔类型
值存放在value.lval中，TRUE——1；FALSE——0，type对应IS_BOOL。

赋值为真：ZVAL_BOOL(zval, 1);
赋值为假：ZVAL_BOOL(zval, 0);
⑤数组类型变量
PHP数组基于HashTable实现，变量赋值为数组类型时先要创建一个HashTable，保存在value的ht字段中。Zend提供array_init()实现赋值。

array_init(zval);
同时Zend提供了一套完整的关联数组、索引数组API用于添加元素，这里不一一列举。
⑥对象类型变量
对象和数组类似，PHP中对象可以转换成数组，但数组无法转换成对象，会丢失方法。Zend通过object_init()函数初始化一个对象。

if(object_init(zval) != SUCCESS){
    RETURN_NULL();
}
Zend也提供了对象设置属性所需的API，和数组设置元素类似，用到时候找即可。
⑦资源类型
严格而言，资源不是数据类型，而是一个可以维护任何数据类型的抽象，类似C语言的指针。所有资源都保存在一个Zend内部的资源列表中，每份资源都有一个指向实际数据的指针。
为了及时回收无用的资源，Zend引擎会自动回收引用数为0的资源的析构函数，析构函数需要在扩展中自己定义。
Zend使用统一的zend_register_list_destructors_ex()为资源注册析构函数，该函数返回一个句柄，将资源与析构函数相关联。定义如下：

ZEND_ZPI int zend_register_list_destructors_ex(rsrc_dtor_func_t ld, 
  rsrc_dtor_func_t  pld, char *type_name, int module_number);
参数描述：
ld ： 普通资源的析构函数
pld : 持久化资源的析构函数
type_name ： 为资源类型起的名字，如：fopen()创建的资源名称为stream
module_number ： PHP_MINIT_FUNCTION函数会定义，可忽略
两种析构函数至少提供一个，为空可用NULL指定。
资源的析构函数必须如下定义：(resource_destruction_handler)函数名随意。

void resource_destruction_handler(zend_rsrc_entry *rsrc TSRMLS_DC){
    -----------具体实现代码------------
}
其中，rsrc是指向zend_rsrc_entry的指针，结构体结构为：

typedef struct _zend_rsrc_entry{
    void *ptr;  //资源的实际地址，析构时释放
    int type; 
    int refcount;
} zend_rsrc_entry ; 
通过zend_register_list_destructors_ex()函数返回的资源句柄，通过一个全局变量保存，ext_skel生成的扩展架构中，自动生成了一个'le_'为前缀的int型变量，zend_register_list_destructors_ex()在MINIT函数中使用并完成注册。如实现链表的析构：

---------------------------------phplist扩展-----------------------------------
static le_phplist;  //架构自动生成，保存资源句柄
//定义链表节点
struct ListNode{
    struct ListNode *next;
    void *data;
}
//析构函数具体实现
void phplist_destruction_handler(zend_rsrc_entry *rsrc TSRMLS_DC){
    ListNode *pre, *next;
    pre = (ListNode *)rsrc->ptr;  
    while(pre){
        next = pre -> next;
        free(pre);
        pre = next;
    }
}
//MINIT中注册析构函数
PHP_MINIT_FUNCTION(phplist){
    //完成注册
    le_phplist = zend_register_list_destructors_ex(phplist_destruction_handler,
      NULL, "php_list", module_number);
    return SUCCESS;
}
注册完析构函数，需要把资源和句柄关联起来，Zend提供zend_register_resource()函数或者ZEND_REGISTER_RESOURCE()宏完成这一操作。

int zend_register_resource(zval *rsrc_result, void *rsrc_pointer, int rsrc_type);
参数解释：
rsrc_result : 存储zend_register_resource返回的结果
rsrc_pointer ： 指向保存的资源
rsrc_type : 资源类型
该函数返回int型结果，该结果为资源的id。函数定义源码：

int zend_register_resource(zval *rsrc_result, void *rsrc_pointer, int rsrc_type){
    int rsrc_id;
    rsrc_id = zend_list_insert(rsrc_pointer, rsrc_type);  //该函数将资源加入资源列表，并返回资源在列表中的位置（即id）
    if(rsrc_result){
        rsrc_result -> value.lval = rsrc_id;
        rsrc_result -> type = IS_RESOURCE;
    }
    return rsrc_id;
}
用户根据资源的id冲资源列表中获取资源，Zend定义了ZEND_FETCH_RESOURCE()宏获取指定的资源。

ZEND_FETCH_RESOURCE(rsrc, rsrc_type, rsrc_id, default_rsrc_id, resource_type_name, resource_type);
其中
rsrc ： 保存返回的资源
rsrc_type ： 表明想要的资源类型，如 ListNode *等
rsrc_id ： 用户通过PHP脚本传来的资源id
default_rsrc_id ： 没有获取到资源时的标识符，一般用-1指定
resource_type_name ： 请求资源类的名称，用于找不到时抛出错误信息使用
resource_type ： 注册析构函数时的句柄，即le_phplist
例如获取用户指定的list

zval *lrc;
ListNode *list;
//获取用户参数
if(zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "r", &lrc) == FAILURE){
    RETURN_FALSE;
}
//获取对应资源
ZEND_FETCH_RESOURCE(list, ListNode *, &lrc, -1, "php list", le_phplist);
此时list即为所要获取的资源。
资源用完需要析构，当引用数为0时，Zend对资源进行回收，很多扩展对资源有相应的析构函数，比如mysql_connect()的mysql_close(),fopen()对应fclose()。PHP的unset()也可以直接释放一个资源。
如果想显示的定义函数释放资源，在自定义函数中调用zend_list_delete()函数即可

ZEND_API int zend_list_delete(int id TSRMLS_DC);
该函数的作用是根据id将资源的引用数-1，然后判断引用数是否大于0，是则触发析构函数清除资源。

（8）错误输出API
Zend推荐使用zend_error()函数输出错误信息，该函数定义如下：

ZEND_API void zend_error(int type, char *format, ...)
参数：
type ： PHP的6钟错误信息类型
    ①E_ERROR：抛出一个错误，脚本将停止运行
    ②E_WARNING ： 抛出警告，脚本继续执行
    ③E_NOTICE ： 抛出通知，脚本继续执行，一般情况下php.ini设置不显示
    ④E_CORE_ERROR ： 抛出PHP内核错误
    ⑤E_COMPILE_ERROR ： 抛出编译器内部错误
    ⑥E_COMPILE_WARNING ： 抛出编译器警告
    注意：后三种错误不应由自定义扩展模块抛出！！！
format ： 错误输出格式
（9）运行时信息函数
执行PHP脚本出错时，经常会有相关的运行信息，指出哪个文件，哪个函数，具体哪行有执行错误，Zend引擎有相关的实现接口。

查看当前执行的函数名
get_active_function_name(TSRMLS_C);
查看当前执行的文件名
zend_get_executed_filename(TSRMLS_C);
查看所在行
zend_get_executed_lineno(TSRMLS_C);
三个函数都需要以TSRMLS_C为参数，作为访问执行器（Executor）全局变量。TSRM_C是TSRM存储器，与线程安全相关，之后专门写篇博客讲讲。

（10）扩展调用脚本中用户自定义函数
这种情况比较少，但Zend功能全面，支持这类操作。
在扩展中使用用户自定义函数，通过call_user_function_ex()函数实现，函数原型：

int call_user_function_ex(HashTable *function_table,   //要访问的函数表指针
    zval **object_pp,   //调用方法的对象，没有设为NULL
    zval *function_name,   //函数名
    zval **retval_ptr_ptr,  //保存返回值的指针
    zend_uint param_count,  //参数个数
    zval **params[],  //参数
    int no_separation,  //是否禁止zcal分离操作
    HashTable symbol_table  //符号表，一般设为NULL
    TSRMLS_DC  
);  
其中no_separation为1会禁止zval分离，节省内存，但任何参数分离将导致操作失败，通常设为0。
脚本中定义用户函数

function userfunc(){
    return "call user function success";
}

-------------------------调用扩展方法----------------------------
$ret = call_user_function_in_ext();
var_dump($ret);
扩展中需要实现call_user_function_in_ext()函数

PHP_FUNCTION(call_user_function_in_ext){
    zval **funcName;
    zval *retval;
  
    if(ZEND_NUM_ARGS() != 1 || 
        zend_get_parameters_ex(1, &function_name) == FAILURE){
        zend_error(E_ERROR, "function %s call in extension fail", (*function_name)->value->str->val);
    }

    if((*function_name)->type != IS_STRING){
        zend_error(E_ERROR, "function name must be string");
    }

    if(call_user_function_ex(CG(function_table), NULL, *function_name, &retval, 0, NULL, 0, NULL TSRMLS_DC) != SUCCESS){
        zend_error(E_ERROR, "function call fail");
    }

    zval *ret_val = *retval;
    zval_copy_ctor(ret_val);
    zval_ptr_dtor(&retval);
}
此外Zend还有提供显示phpinfo的函数，比较简单，不做讲解。

创建扩展
创建一个链表操作的扩展，扩展名为phplist，生成架构先

cd php_src/ext
./ext_skel --extname=phplist
此时ext目录下生成phplist/，本例不依赖其他扩展或lib库，按准备工作（二）修改config.m4文件。
之后实现扩展的函数。在php_phplist.h中声明，具体实现在phplist.c中。
php_phplist.h如下：

#ifndef PHP_PHPLIST_H
#define PHP_PHPLIST_H

extern zend_module_entry phplist_module_entry;
#define phpext_phplist_ptr &phplist_module_entry

#define PHP_PHPLIST_VERSION "0.1.0" /* Replace with version number for your extension */

#ifdef PHP_WIN32
#   define PHP_PHPLIST_API __declspec(dllexport)
#elif defined(__GNUC__) && __GNUC__ >= 4
#   define PHP_PHPLIST_API __attribute__ ((visibility("default")))
#else
#   define PHP_PHPLIST_API
#endif

#ifdef ZTS
#include "TSRM.h"
#endif

PHP_MINIT_FUNCTION(phplist);
PHP_MSHUTDOWN_FUNCTION(phplist);
PHP_RINIT_FUNCTION(phplist);
PHP_RSHUTDOWN_FUNCTION(phplist);
PHP_MINFO_FUNCTION(phplist);

PHP_FUNCTION(confirm_phplist_compiled); /* For testing, remove later. */

/*声明扩展函数*/
PHP_FUNCTION(list_create);  //创建链表
PHP_FUNCTION(list_add_head);    //添加到链表头
PHP_FUNCTION(list_add_tail);    //添加到链表尾
PHP_FUNCTION(list_get_index);   //获取节点
PHP_FUNCTION(list_get_length);  //获取链表长度
PHP_FUNCTION(list_remove_index);    //移除节点

#ifdef ZTS
#define PHPLIST_G(v) TSRMG(phplist_globals_id, zend_phplist_globals *, v)
#else
#define PHPLIST_G(v) (phplist_globals.v)
#endif

#endif  /* PHP_PHPLIST_H */
在phplist.c中具体实现

#ifdef HAVE_CONFIG_H
#include "config.h"
#endif

#include "php.h"
#include "php_ini.h"
#include "ext/standard/info.h"
#include "php_phplist.h"

static int le_phplist;
static int isFree = 0;

const zend_function_entry phplist_functions[] = {
    PHP_FE(confirm_phplist_compiled,    NULL)       /* For testing, remove later. */
    PHP_FE(list_create, NULL)
    PHP_FE(list_add_head, NULL) 
    PHP_FE(list_add_tail, NULL) 
    PHP_FE(list_get_index, NULL)    
    PHP_FE(list_get_length, NULL)   
    PHP_FE(list_remove_index, NULL)
    PHP_FE(list_destroy, NULL)
    PHP_FE(list_get_head, NULL)
    PHP_FE_END  /* Must be the last line in phplist_functions[] */
};

zend_module_entry phplist_module_entry = {
#if ZEND_MODULE_API_NO >= 20010901
    STANDARD_MODULE_HEADER,
#endif
    "phplist",
    phplist_functions,
    PHP_MINIT(phplist),
    PHP_MSHUTDOWN(phplist),
    PHP_RINIT(phplist),     /* Replace with NULL if there's nothing to do at request start */
    PHP_RSHUTDOWN(phplist), /* Replace with NULL if there's nothing to do at request end */
    PHP_MINFO(phplist),
#if ZEND_MODULE_API_NO >= 20010901
    PHP_PHPLIST_VERSION,
#endif
    STANDARD_MODULE_PROPERTIES
};
/* }}} */

#ifdef COMPILE_DL_PHPLIST
ZEND_GET_MODULE(phplist)
#endif

//定义链表节点和链表头
typedef struct _ListNode{
    struct _ListNode *prev;
    struct _ListNode *next;
    zval *value;
}ListNode;
typedef struct _ListHead{
    struct _ListNode *head;
    struct _ListNode *tail;
    int size;
}ListHead;

//创建链表具体实现
ListHead * list_create(){

    ListHead *head;
    head = (ListHead *)malloc(sizeof(ListHead));
    if (head){
        head->size = 0;
        head->head = NULL;
        head->tail = NULL;
    }
    return head;
}

//向头部添加
int list_add_head(ListHead *head, zval *value){

    ListNode *node;
    node = (ListNode *)malloc(sizeof(*node));
    if (!node){
        return 0;
    }
    node->value = value;
    node->prev = NULL;
    node->next = head->head;
    if (head->head){
        head->head->prev = node;
    }
    head->head = node;
    if(!head->tail){
        head->tail = head->head;
    }
    head->size++;
    return 1;
}

//链表尾添加
int list_add_tail(ListHead *list, zval *value){

    ListNode *node;
    node = (ListNode *)malloc(sizeof(*node));
    if(!node){
        return 0;
    }
    node->value = value;
    node->next = NULL;
    node->prev = list->tail;
    if (list->tail){
        list->tail->next = node;
    }
    list->tail = node;
    if (!list->head){
        list->head = list->tail;
    }
    list->size++;
    return 1;
}

//获取指定元素
int list_get_index(ListHead *list, int index, zval **retval){

    ListNode *node;
    if(!list){
        return 0;
    }
    if (index <= 0 || list->size == 0 ||  index > list->size){
        return 0;
    }
    if (index < list->size / 2){
        node = list->head;
        while(node && index > 0){
            node = node->next;
            --index;
        }
    }else{
        node = list->tail;
        while(node && index > 0){
            node = node->prev;
            --index;
        }
    }
    *retval = node->value;
    return 1;
}

//获取链表长度
int list_get_length(ListHead *list){

    if (list){
        return list->size;
    }else{
        return 0;
    }
}

//删除节点
int list_remove_index(ListHead *list, int index){

    ListNode *node;
    if(!list){
        return 0;
    }
    if (index <= 0 || list->size == 0 || index > list->size){
        return 0;
    }
    if (index < list->size / 2){
        node = list->head;
        while(node && index > 0){
            node = node->next;
            --index;
        }
    }else{
        node = list->tail;
        while(node && index > 0){
            node = node->prev;
            --index;
        }
    }
    if (!node){
        return 0;
    }
    if (node->prev){
        node->prev->next = node->next;
    }else{
        list->head = node->next;
    }
    if(node->next){
        node->next->prev = node->prev;
    }else{
        list->tail = node->prev;
    }
    list->size--;
    return 1;
}

//释放链表
void list_destroy(ListHead *list){
    
    ListNode *pre, *next;
    pre = list->head;
    while(pre){
        next = pre->next;
        free(pre);
        pre = next;
    }
    free(list);
}

int list_get_head(ListHead *list, zval **retval){

    if (!list || !list->head){
        return 0;
    }
    *retval = list->head->value;
    return 1;
}

//析构函数实现
void phplist_destructor_handler(zend_rsrc_list_entry *rsrc TSRMLS_DC){
    if (!isFree){
        ListHead *list;
        list = (ListHead *)rsrc->ptr;
        list_destroy(list);
        isFree = 1;
    }
}

PHP_MINIT_FUNCTION(phplist)
{
    le_phplist = zend_register_list_destructors_ex(phplist_destructor_handler, NULL, "phplist", module_number);
    return SUCCESS;
}

PHP_MSHUTDOWN_FUNCTION(phplist)
{
    return SUCCESS;
}

PHP_RINIT_FUNCTION(phplist)
{
    return SUCCESS;
}

PHP_RSHUTDOWN_FUNCTION(phplist)
{
    return SUCCESS;
}

PHP_MINFO_FUNCTION(phplist)
{
    php_info_print_table_start();
    php_info_print_table_header(2, "phplist support", "enabled");
    php_info_print_table_end();
}

PHP_FUNCTION(confirm_phplist_compiled)
{
    char *arg = NULL;
    int arg_len, len;
    char *strg;

    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "s", &arg, &arg_len) == FAILURE) {
        return;
    }

    len = spprintf(&strg, 0, "Congratulations! You have successfully modified ext/%.78s/config.m4. Module %.78s is now compiled into PHP.", "phplist", arg);
    RETURN_STRINGL(strg, len, 0);
}

PHP_FUNCTION(list_create){

    ListHead *list;
    list = list_create();
    if (!list){
        RETURN_NULL();
    }
    ZEND_REGISTER_RESOURCE(return_value, list, le_phplist);
}

PHP_FUNCTION(list_add_head){

    zval *value;
    zval *lrc;
    ListHead *list;

    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "rz", &lrc, &value) == FAILURE){
        RETURN_FALSE;
    }
    ZEND_FETCH_RESOURCE(list, ListHead *, &lrc, -1, "php list", le_phplist);

    int ret = list_add_head(list, value);
    if(ret){
        RETURN_TRUE;
    }
    RETURN_FALSE;
}

PHP_FUNCTION(list_add_tail){

    zval *value;
    zval *lrc;
    ListHead *list;

    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "rz", &lrc, &value) == FAILURE){
        RETURN_FALSE;
    }
    ZEND_FETCH_RESOURCE(list, ListHead *, &lrc, -1, "php list", le_phplist);

    int ret = list_add_tail(list, value);
    if(ret){
        RETURN_TRUE;
    }
    RETURN_FALSE;
}

PHP_FUNCTION(list_get_index){

    zval *lrc;
    ListHead *list;
    long index;
    zval *retval;
    if(zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "rl", &lrc, &index) == FAILURE){
        RETURN_FALSE;
    }
    ZEND_FETCH_RESOURCE(list, ListHead *, &lrc, -1, "php list", le_phplist);
    int ret = list_get_index(list, index, &retval);
    if (ret){
        RETURN_ZVAL(retval, 1, 0);
    }
    RETURN_NULL();
}

PHP_FUNCTION(list_get_length){

    zval *lrc;
    ListHead *list;
    if(zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "r", &lrc) == FAILURE){
        RETURN_FALSE;
    }
    ZEND_FETCH_RESOURCE(list, ListHead *, &lrc, -1, "php list", le_phplist);
    RETURN_LONG(list_get_length(list));
}

PHP_FUNCTION(list_remove_index){

    zval *lrc;
    ListHead *list;
    long index;
    if(zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "rl", &lrc, &index) == FAILURE){
        RETURN_FALSE;
    }
    ZEND_FETCH_RESOURCE(list, ListHead *, &lrc, -1, "php list", le_phplist);
    int ret = list_remove_index(list, index);
    if (ret){
        RETURN_TRUE;
    }
    RETURN_FALSE;
}

PHP_FUNCTION(list_destroy){

    zval *lrc;
    ListHead *list;
    if(zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "r", &lrc) == FAILURE){
        RETURN_FALSE;
    }
    ZEND_FETCH_RESOURCE(list, ListHead *, &lrc, -1, "php list", le_phplist);
    list_destroy(list);
}

PHP_FUNCTION(list_get_head){

    zval *lrc;
    zval *retval;
    ListHead *list;
    if(zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "r", &lrc) == FAILURE){
        RETURN_FALSE;
    }
    ZEND_FETCH_RESOURCE(list, ListHead *, &lrc, -1, "php list", le_phplist);
    int ret = list_get_head(list, &retval);
    if (ret){
        RETURN_ZVAL(retval, 1, 0);
    }
    RETURN_NULL();
}
之后按照步骤在php.ini中引入扩展即可使用。