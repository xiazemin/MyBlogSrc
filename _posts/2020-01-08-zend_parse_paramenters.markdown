---
title: zend_parse_paramenters
layout: post
category: lang
author: 夏泽民
---
基本参数
最简单的获取函数调用者传递过来的参数便是使用 zend_parse_parameters() 函数。 zend_parse_parameters() 函数的前几个参数我们直接用内核里的宏来生成便可以了，形式为：ZEND_NUM_ARGS() TSRMLS_CC，注意两者之间有个空格，但是没有逗号。从名字可以看出，ZEND_NUM_ARGS() 代表着参数的个数。紧接着需要传递给 zend_parse_parameters() 函数的参数是一个用于格式化的字符串，就像 printf 的第一个参数一样。下面列出了最常用的几个符号：

参数	代表着的类型
b	Boolean
l	Integer
d	Float
s	String
r	Resource
a	Array
o	Object
O	特定类型的Object
z	任意类型
Z	zval**类型
f	表示函数、方法名称
这个函数就像 printf() 函数一样，后面的参数是与格式化字符串里的格式一一对应的。一些基础类型的数据会直接映射成 C 语言里的类型。

1
ZEND_FUNCTION(sample_getlong) {
2
    long foo;
3
    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "l", &foo) == FAILURE)
4
    {
5
        RETURN_NULL();
6
    }
7
    php_printf("The integer value of the parameter is: %ld\n", foo);
8
    RETURN_TRUE;
9
}
一般来说，int 和 long 这两种数据类型的数据往往是相同的，但也有例外情况。所以我们不应该把long 的数组放在一个 int 里，尤其是在 64 位平台里，那将引发一些不容易排查的 Bug。所以通过zend_parse_parameter() 函数接收参数时，我们应该使用内核约定好的类型变量作为载体：

参数	对应C里的数据类型
b	zend_bool
l	long
d	double
s	char*, int 前者接收指针，后者接收长度
r	zval*
a	zval*
o	zval*
O	zval*, zend_class_entry*
z	zval*
Z	zval**
注意，所有的 PHP 语言中的复合类型参数都需要 zval* 类型来作为载体，因为它们都是内核自定义的一些数据结构。我们一定要确认参数和载体的类型一致，如果需要，它可以进行类型转换，比如把 array 转换成 stdClass 对象。 s 和 O (字母大写)类型需要特殊一些，因为它们都需要两个载体。我们将在接下来的章节里了解 PHP 中对象的具体实现。这样我们改写一下我们之前定义的一个函数：

1
<?php
2
function sample_hello_world($name) {
3
    echo "Hello $name!\n";
4
}
在编写扩展时，我们需要用 zend_parse_parameters() 来接收这个字符串:

1
ZEND_FUNCTION(sample_hello_world) 
2
{
3
    char *name;
4
    int name_len;
5
​
6
    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "s", &name, &name_len) == FAILURE)
7
    {
8
            RETURN_NULL();
9
    }
10
       php_printf("Hello ");
11
       PHPWRITE(name, name_len);
12
       php_printf("!\n");
13
}
如果传递给函数的参数数量小于 zend_parse_parameters() 要接收的参数数量，它便会执行失败，并返回 FAILURE。

如果我们需要接收多个参数，可以直接在 zend_parse_paramenters() 的参数里罗列接收载体便可以了，如：

1
<?php
2
function sample_hello_world($name, $greeting) {
3
    echo "Hello $greeting $name!\n";
4
}
5
sample_hello_world('John Smith', 'Mr.');
在 PHP 扩展里应该这样来实现：

1
ZEND_FUNCTION(sample_hello_world) {
2
    char *name;
3
    int name_len;
4
    char *greeting;
5
    int greeting_len;
6
​
7
    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "ss", &name, &name_len, &greeting, &greeting_len) == FAILURE) {
8
        RETURN_NULL();
9
    }
10
​
11
    php_printf("Hello ");
12
    PHPWRITE(greeting, greeting_len);
13
    php_printf(" ");
14
    PHPWRITE(name, name_len);
15
    php_printf("!\n");
16
}
除了上面定义的参数，还有其它的三个参数来增强我们接收参数的能力,如下：

|：它之前的参数都是必须的，之后的都是非必须的，也就是有默认值的。
!：如果接收了一个 PHP 语言里的 NULL 变量，则直接把其转成 C 语言里的 NULL，而不是封装成IS_NULL 类型的 zval。
/：如果传递过来的变量与别的变量共用一个 zval，而且不是引用，则进行强制分离，新的 zval 的is_ref__gc 等于 0，并且 refcount__gc 等于 1。
默认参数值
现在让我们继续改写 sample_hello_world(), 接下来我们使用一些参数的默认值，在 PHP 语言里就像下面这样：

1
<?php
2
function sample_hello_world($name, $greeting='Mr./Ms.') {
3
    echo "Hello $greeting $name!\n";
4
}
5
sample_hello_world('Ginger Rogers','Ms.');
6
sample_hello_world('Fred Astaire');
此时既可以只向 sample_hello_world 中传递一个参数，也可以传递完整的两个参数。那同样的功能我们怎样在扩展函数里实现呢？我们需要借助 zend_parse_parameters 中的 | 参数，这个参数之前的参数被认为是必须的，之后的便认为是非必须的了，如果没有传递，则不会去修改载体。

1
ZEND_FUNCTION(sample_hello_world) {
2
    char *name;
3
    int name_len;
4
    char *greeting = "Mr./Mrs.";
5
    int greeting_len = sizeof("Mr./Mrs.") - 1;
6
​
7
    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "s|s",
8
      &name, &name_len, &greeting, &greeting_len) == FAILURE) {
9
        RETURN_NULL();
10
    }
11
​
12
    php_printf("Hello ");
13
    PHPWRITE(greeting, greeting_len);
14
    php_printf(" ");
15
    PHPWRITE(name, name_len);
16
    php_printf("!\n");
17
}
如果你不传递第二个参数，则扩展函数会被认为默认而不去修改载体。所以，我们需要自己来预先设置有载体的值，它往往是 NULL，或者一个与函数逻辑有关的值。

NULL参数值
每个 zval，包括 IS_NULL 型的 zval，都需要占用一定的内存空间，并且需要 CPU 的计算资源来为它申请内存、初始化，并在它们完成工作后释放掉。但是很多代码都没有意识到这一点。有很多代码都会把一个 NULL 类型的值包裹成 zval 的 IS_NULL 类型，在扩展开发里这种操作是可以优化的，我们可以把参数接收成 C 语言里的 NULL。我们就这一个问题看以下代码：

1
ZEND_FUNCTION(sample_arg_fullnull) {
2
    zval *val;
3
    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "z", &val) == FAILURE) {
4
        RETURN_NULL();
5
    }
6
    if (Z_TYPE_P(val) == IS_NULL) {
7
        val = php_sample_make_defaultval(TSRMLS_C);
8
    }
9
    ...
10
}
11
​
12
ZEND_FUNCTION(sample_arg_nullok) {
13
    zval *val;
14
    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "z!", &val) == FAILURE) {
15
        RETURN_NULL();
16
    }
17
    if (!val) {
18
        val = php_sample_make_defaultval(TSRMLS_C);
19
    }
20
}
这两段代码乍看起来并没有什么很大的不同，但是第一段代码确实需要更多的 CPU 和内存资源。可能这个技巧在平时并没多大用，不过技多不压身，知道总比不知道好。

非引用参数强制分离
当一个变量被传递给函数时候，无论它是否被引用，它的 refcoung__gc 属性都会加 1，至少成为 2：一份是它自己，另一份是传递给函数的拷贝。在改变这个 zval 之前，有时会需要提前把它分离成实际意义上的两份拷贝。这就是 / 修饰符的作用。它把写时复制的 zval 提前分离成两个完整独立的拷贝，从而使我们可以在后面的代码中随意对其进行操作，否则我们可能需要不停的提醒自己对接收的参数进行分离等操作。和 ! 一样，该修饰符位于其所影响的类型之后。

zend_get_arguments()
如果你想让你的扩展能够兼容老版本的 PHP，或者你只想以 zval 为载体来接收参数，可以考虑使用zend_get_parameters() 函数来接收参数。zend_get_parameters() 与zend_parse_parameters() 不同，从名字上我们便可以看出，它直接获取，而不做解析。

首先，它不会自动进行类型转换，所有的参数在扩展实现中的载体都需要是 zval 类型的，下面让我们来看一个最简单的例子：

1
ZEND_FUNCTION(sample_onearg) {
2
    zval *firstarg;
3
​
4
    if (zend_get_parameters(ZEND_NUM_ARGS(), 1, &firstarg)== FAILURE) {
5
        php_error_docref(NULL TSRMLS_CC, E_WARNING, "Expected at least 1 parameter.");
6
        RETURN_NULL();
7
    }
8
    /* Do something with firstarg... */
9
}
其次，zend_get_parameters() 在接收失败的时候，并不会自己抛出错误，它也不能方便地处理具有默认值的参数。最后一点与 zend_parse_parameters 不同的是，它会自动的把所有符合写时复制的 zval 进行强制分离，生成一个崭新的拷贝送到函数内部。如果你希望用它其它的特性，而唯独不需要这个功能，可以去尝试一下用 zend_get_parameters_ex() 函数来接收参数。为了不对写时复制的变量进行分离操作，zend_get_parameters_ex() 的参数是 zval** 类型的，而不是zval*。这个函数不太经常用，可能只会在你碰到一些极端问题时候才会想到它，而它用起来却很简单：

1
ZEND_FUNCTION(sample_onearg) {
2
    zval **firstarg;
3
    if (zend_get_parameters_ex(1, &firstarg) == FAILURE) {
4
        WRONG_PARAM_COUNT;
5
    }
6
    /* Do something with firstarg... */
7
}
注：zend_get_parameters_ex 不需要 ZEND_NUM_ARGS() 作为参数，因为它是在是在后期加入的，那个参数已经不再需要了。

上面例子中还使用了 WRONG_PARAM_COUNT 宏，它的功能是抛出一个 E_WARNING 级别的错误信息，并自动 return。

可变参数实现
还有一种 zend_get_parameter_** 函数，专门用来解决参数很多或者无法提前知道参数数目的问题。试想一下 PHP 语言中 var_dump() 函数的用法，我们可以向其传递任意数量的参数，它在内核中的实现其实是这样的：

1
ZEND_FUNCTION(var_dump) {
2
    int i, argc = ZEND_NUM_ARGS();
3
    zval ***args;
4
​
5
    args = (zval ***)safe_emalloc(argc, sizeof(zval **), 0);
6
    if (ZEND_NUM_ARGS() == 0 || zend_get_parameters_array_ex(argc, args) == FAILURE) {
7
        efree(args);
8
        WRONG_PARAM_COUNT;
9
    }
10
​
11
    for (i = 0; i < argc; i++) {
12
            php_var_dump(args[i], 1 TSRMLS_CC);
13
    }
14
​
15
    efree(args);
16
}
https://xueyuanjun.com/link/7233#bkmrk-%E7%A8%8B%E5%BA%8F%E9%A6%96%E5%85%88%E8%8E%B7%E5%8F%96%E5%8F%82%E6%95%B0%E6%95%B0%E9%87%8F%EF%BC%8C%E7%84%B6%E5%90%8E%E9%80%9A%E8%BF%87-safe
 
程序首先获取参数数量，然后通过 safe_emalloc 函数申请了相应大小的内存来存放这些 zval** 类型的参数。这里使用了 zend_get_parameters_array_ex() 函数来把传递给函数的参数填充到 args 中。你可能已经立即想到，还存在一个名为 zend_get_parameters_array() 的函数，唯一不同的是它将 zval* 类型的参数填充到 args 中，并且需要 ZEND_NUM_ARGS() 作为参数。
<!-- more -->
在经过词语分析，语法分析后，我们知道对于函数的参数检查是通过 zend_do_receive_arg 函数来实现的。在此函数中对于参数的关键代码如下：

CG(active_op_array)->arg_info = erealloc(CG(active_op_array)->arg_info,
        [sizeof](http://www.php.net/sizeof)(zend_arg_info)*(CG(active_op_array)->num_args));
cur_arg_info = &CG(active_op_array)->arg_info[CG(active_op_array)->num_args-1];
cur_arg_info->name = estrndup(varname->u.[constant](http://www.php.net/constant).value.str.val,
        varname->u.[constant](http://www.php.net/constant).value.str.len);
cur_arg_info->name_len = varname->u.[constant](http://www.php.net/constant).value.str.len;
cur_arg_info->array_type_hint = 0;
cur_arg_info->allow_null = 1;
cur_arg_info->pass_by_reference = pass_by_reference;
cur_arg_info->class_name = NULL;
cur_arg_info->class_name_len = 0;
整个参数的传递是通过给中间代码的arg_info字段执行赋值操作完成。关键点是在arg_info字段。arg_info字段的结构如下：

typedef struct _zend_arg_info {
    const char *name;   /* 参数的名称*/
    zend_uint name_len;     /* 参数名称的长度*/
    const char *class_name; /* 类名 */
    zend_uint class_name_len;   /* 类名长度*/
    zend_bool array_type_hint;  /* 数组类型提示 */
    zend_bool allow_null;   /* 是否允许为NULL　*/
    zend_bool pass_by_reference;    /*　是否引用传递 */
    zend_bool return_reference; 
    int required_num_args;  
} zend_arg_info;
参数的值传递和参数传递的区分是通过 pass_by_reference参数在生成中间代码时实现的。

对于参数的个数，中间代码中包含的arg_nums字段在每次执行 **zend_do_receive_arg×× 时都会加1.如下代码：

CG(active_op_array)->num_args++;
并且当前参数的索引为CG(active_op_array)->num_args-1 .如下代码：

cur_arg_info = &CG(active_op_array)->arg_info[CG(active_op_array)->num_args-1];
以上的分析是针对函数定义时的参数设置，这些参数是固定的。而在实际编写程序时可能我们会用到可变参数。此时我们会使用到函数 func_num_args 和 func_get_args。它们是以内部函数存在。在 Zend\zend_builtin_functions.c 文件中找到这两个函数的实现。首先我们来看func_num_args函数的实现。其代码如下：

/* \{\{\{ proto int func_num_args(void) Get the number of arguments that were passed to the function */
ZEND_FUNCTION(func_num_args)
{
    zend_execute_data *ex = EG(current_execute_data)->prev_execute_data;
 
    if (ex && ex->function_state.arguments) {
        RETURN_LONG((long)(zend_uintptr_t)*(ex->function_state.arguments));
    } else {
        zend_error(E_WARNING,
"func_num_args(): Called from the global scope - no function context");
        RETURN_LONG(-1);
    }
}
/* }}} */
在存在 ex->function_state.arguments的情况下，即函数调用时，返回ex->function_state.arguments转化后的值 ，否则显示错误并返回-1。这里最关键的一点是EG(current_execute_data)。这个变量存放的是当前执行程序或函数的数据。此时我们需要取前一个执行程序的数据，为什么呢？因为这个函数的调用是在进入函数后执行的。函数的相关数据等都在之前执行过程中。于是调用的是：

zend_execute_data *ex = EG(current_execute_data)->prev_execute_data;
function_state等结构请参照本章第一小节。

在了解func_num_args函数的实现后，func_get_args函数的实现过程就简单了，它们的数据源是一样的，只是前面返回的是长度，而这里返回了一个创建的数组。数组中存放的是从ex->function_state.arguments转化后的数据。

内部函数的参数
以上我们所说的都是用户自定义函数中对于参数的相关内容。下面我们开始讲解内部函数是如何传递参数的。以常见的count函数为例。其参数处理部分的代码如下：

/* \{\{\{ proto int count(mixed var [, int mode]) Count the number of elements in a variable (usually an array) */
PHP_FUNCTION(count)
{
    zval *array;
    long mode = COUNT_NORMAL;
 
    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "z|l",
         &array, &mode) == FAILURE) {
        return;
    }
    ... //省略
}
这包括了两个操作：一个是取参数的个数，一个是解析参数列表。

取参数的个数

取参数的个数是通过ZEND_NUM_ARGS()宏来实现的。其定义如下：

#define ZEND_NUM_ARGS() (ht)
PHP3 中使用的是宏 ARG_COUNT

ht是在 Zend/zend.h文件中定义的宏 INTERNAL_FUNCTION_PARAMETERS 中的ht，如下：

#define INTERNAL_FUNCTION_PARAMETERS int ht, zval *return_value,
zval **return_value_ptr, zval *this_ptr, int return_value_used TSRMLS_DC
解析参数列表

PHP内部函数在解析参数时使用的是 zend_parse_parameters。它可以大大简化参数的接收处理工作，虽然它在处理可变参数时还有点弱。

其声明如下：

ZEND_API int zend_parse_parameters(int num_args TSRMLS_DC, char *type_spec, ...)
第一个参数num_args表明表示想要接收的参数个数，我们经常使用ZEND_NUM_ARGS() 来表示对传入的参数“有多少要多少”。
第二参数应该总是宏 TSRMLS_CC 。
第三个参数 type_spec 是一个字符串，用来指定我们所期待接收的各个参数的类型，有点类似于 printf 中指定输出格式的那个格式化字符串。
剩下的参数就是我们用来接收PHP参数值的变量的指针。
zend_parse_parameters() 在解析参数的同时会尽可能地转换参数类型，这样就可以确保我们总是能得到所期望的类型的变量。任何一种标量类型都可以转换为另外一种标量类型，但是不能在标量类型与复杂类型（比如数组、对象和资源等）之间进行转换。如果成功地解析和接收到了参数并且在转换期间也没出现错误，那么这个函数就会返回 SUCCESS，否则返回 FAILURE。如果这个函数不能接收到所预期的参数个数或者不能成功转换参数类型时就会抛出一些错误信息。

第三个参数指定的各个参数类型列表如下所示：

l - 长整形
d - 双精度浮点类型
s - 字符串 (也可能是空字节)和其长度
b - 布尔型
r - 资源，保存在 zval*
a - 数组，保存在 zval*
o - （任何类的）对象，保存在 zval *
O - （由class entry 指定的类的）对象，保存在 zval *
z - 实际的 zval*
除了各个参数类型，第三个参数还可以包含下面一些字符，它们的含义如下：

| - 表明剩下的参数都是可选参数。如果用户没有传进来这些参数值，那么这些值就会被初始化成默认值。
/ - 表明参数解析函数将会对剩下的参数以 SEPARATE_ZVAL_IF_NOT_REF() 的方式来提供这个参数的一份拷贝，除非这些参数是一个引用。
! - 表明剩下的参数允许被设定为 NULL（仅用在 a、o、O、r和z身上）。如果用户传进来了一个 NULL 值，则存储该参数的变量将会设置为 NULL。
参数的传递
在PHP的运行过程中，如果函数有参数，当执行参数传递时，所传递参数的引用计数会发生变化。如和Xdebug的作者Derick Rethans在其文章php variables中的示例的类似代码：

function do_something($s) {
       xdebug_debug_zval('s');
        $s = 100;
        return $s;
}
 
$a = 1111;
$b = do_something($a);
[echo](http://www.php.net/echo) $b;
如果你安装了xdebug，此时会输出s变量的refcount为3，如果使用debug_zval_dump，会输出4。因为此内部函数调用也对refcount执行了加1操作。这里的三个引用计数分别是：

function stack中的引用
function symbol table中引用
原变量$a的引用。
这个函数符号表只有用户定义的函数才需要，内置和扩展里的函数不需要此符号表。debug_zval_dump()是内置函数，并不需要符号表，所以只增加了1。 xdebug_debug_zval()传递的是变量名字符串，所以没有增加refcount。

每个PHP脚本都有自己专属的全局符号表，而每个用户自定义的函数也有自己的符号表，这个符号表用来存储在这个函数作用域下的属于它自己的变量。当调用每个用户自定义的函数时，都会为这个函数创建一个符号表，当这个函数返回时都会释放这个符号表。

当执行一个拥有参数的用户自定义的函数时，其实它相当于赋值一个操作，即$s = $a;只是这个赋值操作的引用计数会执行两次，除了给函数自定义的符号表，还有一个是给函数栈。

参数的传递的第一步是SEND_VAR操作，这一步操作是在函数调用这一层级，如示例的PHP代码通过VLD生成的中间代码:

compiled vars:  !0 = $a, !1 = $b
line     # *  op                           fetch          ext  return  operands
--------------------------------------------------------------------------------
-
   2     0  >   EXT_STMT
         1      NOP
   7     2      EXT_STMT
         3      ASSIGN                                                   !0, 1111
   8     4      EXT_STMT
         5      EXT_FCALL_BEGIN
         6      SEND_VAR                                                 !0
         7      DO_FCALL                                      1          'demo'
         8      EXT_FCALL_END
         9      ASSIGN                                                   !1, $1
   9    10      EXT_STMT
        11      ECHO                                                     !1
        12    > RETURN                                                   1
 
branch: #  0; line:     2-    9; sop:     0; eop:    12
path #1: 0,
Function demo:
函数调用是DO_FCALL，在此中间代码之前有一个SEND_VAR操作，此操作的作用是将实参传递给函数，并且将它添加到函数栈中。最终调用的具体代码参见zend_send_by_var_helper_SPEC_CV函数，在此函数中执行了引用计数加1（Z_ADDREF_P）操作和函数栈入栈操作（zend_vm_stack_push）。

与第一步的SEND操作对应，第二步是RECV操作。RECV操作和SEND_VAR操作不同，它是归属于当前函数的操作，仅为此函数服务。它的作用是接收SEND过来的变量，并将它们添加到当前函数的符号表。示例函数生成的中间代码如下：

compiled vars:  !0 = $s
line     # *  op                           fetch          ext  return  operands
--------------------------------------------------------------------------------
-
   2     0  >   EXT_NOP
         1      RECV                                                     1
   3     2      EXT_STMT
         3      ASSIGN                                                   !0, 10
   4     4      EXT_STMT
         5    > RETURN                                                   !0
   5     6*     EXT_STMT
         7*   > RETURN                                                   null
 
branch: #  0; line:     2-    5; sop:     0; eop:     7
参数和普通局部变量一样 ，都需要进行操作，都需要保存在符号表（或CVs里，不过查找一般都是直接从变量变量数组里查找的）。如果函数只是需要读这个变量，如果我们将这个变量复制一份给当前函数使用的话，在内存使用和性能方面都会有问题，而现在的方案却避免了这个问题，如我们的示例：使用类似于赋值的操作，将原变量的引用计数加一，将有变化时才将原变量引用计数减一，并新建变量。其最终调用是ZEND_RECV_SPEC_HANDLER。

参数的压栈操作用户自定义的函数和内置函数都需要，而RECV操作仅用户自定义函数需要。

函数参数解析
之前我们定义的函数没有接收任何参数，那么扩展定义的内部函数如何读取参数呢？用户自定义函数在编译时会为每个参数创建一个zend_arg_info结构，这个结构用来记录参数的名称、是否引用传参、是否为可变参数等，在存储上函数参数与局部变量相同，都分配在zend_execute_data上，且最先分配的就是函数参数，调用函数时首先会进行参数传递，按参数次序依次将参数的value从调用空间传递到被调函数的zend_execute_data，函数内部像访问普通局部变量一样通过存储位置访问参数，这是用户自定义函数的参数实现。

/* arg_info for user functions */
typedef struct _zend_arg_info {
	zend_string *name;//参数名
	zend_string *class_name;
	zend_uchar type_hint;//显式声明的参数类型，比如(array $param_1)
	zend_uchar pass_by_reference;//是否引用传参，参数前加&的这个值就是1
	zend_bool allow_null;//是否允许为NULL,注意：这个值并不是用来表示参数是否为必传的
	zend_bool is_variadic;//是否为可变参数，即...用法，与golang的用法相同，5.6以上新增的一个用法：function my_func($a, ...$b){...}
} zend_arg_info;
PHP中通过 zend_parse_parameters() 这个函数解析zend_execute_data上保存的参数：

zend_parse_parameters(int num_args, const char *type_spec, ...);
num_args为实际传参数，通过 ZEND_NUM_ARGS()获取；
type_spec是一个字符串，用来标识解析参数的类型，比如:"la"表示第一个参数为整形，第二个为数组，将按照这个解析到指定变量；
后面是一个可变参数，用来指定解析到的变量，这个值与type_spec配合使用，即type_spec用来指定解析的变量类型，可变参数用来指定要解析到的变量，这个值必须是指针。
解析的过程也比较容易理解，调用函数时首先会把参数拷贝到调用函数的zend_execute_data上，所以解析的过程就是按照type_spec指定的各个类型，依次从zend_execute_data上获取参数，然后将参数地址赋给目标变量。

参数类型

整形：l、L

整形通过"l"或"L"标识，表示解析的参数为整形，解析到的变量类型必须是 zend_long ，不能解析其它类型，如果输入的参数不是整形将按照类型转换规则将其转为整形：

zend_long lval;
if(zend_parse_parameters(ZEND_NUM_ARGS(), "l", &lval){
...
}
printf("lval:%d\n", lval);
如果在标识符后加"!"，即："l!"、"L!"，则必须再提供一个zend_bool变量的地址，通过这个值可以判断传入的参数是否为NULL，如果为NULL则将要解析到的zend_long值设置为0，同时zend_bool设置为1：

zend_long lval; //如果参数为NULL则此值被设为0
zend_bool is_null;//如果参数为NULL则此值为1，否则为0
if(zend_parse_parameters(ZEND_NUM_ARGS(), "l!", &lval, &is_null){
..
}
布尔型：b

通过"b"标识符表示将传入的参数解析为布尔型，解析到的变量必须是zend_bool：

zend_bool ok;
if(zend_parse_parameters(ZEND_NUM_ARGS(), "b", &ok, &is_null) == FAILURE){
..
}
"b!"的用法与整形的完全相同，也必须再提供一个zend_bool的地址用于获取传参是否为NULL，如果为NULL，则zend_bool为0，用于获取是否NULL的zend_bool为1。


浮点型：d
通过"d"标识符表示将参数解析为浮点型，解析的变量类型必须为double：

double dval;
 
if(zend_parse_parameters(ZEND_NUM_ARGS(), "d", &dval) == FAILURE){
..
}
具体解析过程不再展开，"d!"与整形、布尔型用法完全相同。


字符串：s、S、p、P
字符串解析有两种形式：char、zend_string，其中"s"将参数解析到`char`，且需要额外提供一个size_t类型的变量用于获取字符串长度，"S"将解析到zend_string：

char *str;
size_t str_len;
if(zend_parse_parameters(ZEND_NUM_ARGS(), "s", &str, &str_len) == FAILURE){
...
}
zend_string *str;
if(zend_parse_parameters(ZEND_NUM_ARGS(), "S", &str) == FAILURE){
...
}
"s!"、"S!"与整形、布尔型用法不同，字符串时不需要额外提供zend_bool的地址，如果参数为NULL，则char*、zend_string将设置为NULL。除了"s"、"S"之外还有两个类似的："p"、"P"，从解析规则来看主要用于解析路径，实际与普通字符串没什么区别，尚不清楚这俩有什么特殊用法。

数组：a、A、h、H
数组的解析也有两类，一类是解析到zval层面，另一类是解析到HashTable，其中"a"、"A"解析到的变量必须是zval，"h"、"H"解析到HashTable，这两类是等价的：

zval *arr; //必须是zval指针，不能是zval arr，因为参数保存在zend_execute_data上，arr为此空间上参数的地址
HashTable *ht;
if(zend_parse_parameters(ZEND_NUM_ARGS(), "ah", &arr, &ht) == FAILURE){
...
}
"a!"、"A!"、"h!"、"H!"的用法与字符串一致，也不需要额外提供别的地址，如果传参为NULL，则对应解析到的zval、HashTable也为NULL.

对象：o、O
如果参数是一个对象则可以通过"o"、"O"将其解析到目标变量，注意：只能解析为zval，无法解析为zend_object。

zval *obj;
if(zend_parse_parameters(ZEND_NUM_ARGS(), "o", &obj) == FAILURE){
...
}
O"是要求解析指定类或其子类的对象，类似传参时显式的声明了参数类型的用法function my_func(MyClass $obj){...} ，如果参数不是指定类的实例化对象则无法解析。
"o!"、"O!"与字符串用法相同

资源：r
如果参数为资源则可以通过"r"获取其zval的地址，但是无法直接解析到zend_resource的地址，与对象相同。

zval *res;
if(zend_parse_parameters(ZEND_NUM_ARGS(), "r", &res) == FAILURE){
...
}
"r!"与字符串用法相。

类：C
如果参数是一个类则可以通过"C"解析出zend_class_entry地址：function my_func(stdClass){...} ，这里有个地方比较特殊，解析到的变量可以设定为一个类，这种情况下解析时将会找到的类与指定的类之间的父子关系，只有存在父子关系才能解析，如果只是想根据参数获取类型的zend_class_entry地址，记得将解析到的地址初始化为NULL，否则将会不可预料的错误。

zend_class_entry *ce = NULL; //初始为NULL
if(zend_parse_parameters(ZEND_NUM_ARGS(), "C", &ce) == FAILURE){
RETURN_FALSE;
}
callable：f
callable指函数或成员方法，如果参数是函数名称字符串、array(对象/类,成员方法)，则可以通过"f"标识符解析出 zend_fcall_info 结构，这个结构是调用函数、成员方法时的唯一输入。

zend_fcall_info callable; //注意，这两个结构不能是指针
zend_fcall_info_cache call_cache;
if(zend_parse_parameters(ZEND_NUM_ARGS(), "f", &callable, &call_cache) ==FAILURE){
RETURN_FALSE;
}
函数调用：

my_func_1("func_name");//或
my_func_1(array('class_name', 'static_method'));//或
my_func_1(array($object, 'method'));
解析出 zend_fcall_info 后就可以通过 zend_call_function() 调用函数、成员方法了，提供"f"解析到 zend_fcall_info 的用意是简化函数调用的操作，否则需要我们自己去查找函数、检查是否可被调用等工作，关于这个结构稍后介绍函数调用时再作详细说明。

任意类型：z
"z"表示按参数实际类型解析，比如参数为字符串就解析为字符串，参数为数组就解析为数组，这种实际就是将zend_execute_data上的参数地址拷贝到目的变量了，没有做任何转化。
"z!"与字符串用法相同。
 

引用传参
函数中解析参数还有一种就是引用传参。

如果函数需要使用引用类型的参数或返回引用就需要创建函数的参数数组，这个数组通过：

ZEND_BEGIN_ARG_INFO()或ZEND_BEGIN_ARG_INFO_EX() 、
ZEND_END_ARG_INFO() 宏定义：

#define ZEND_BEGIN_ARG_INFO_EX(name, _unused, return_reference, required_num_args)
#define ZEND_BEGIN_ARG_INFO(name, _unused)
name: 参数数组名，注册函数 PHP_FE(function, arg_info) 会用到
_unused: 保留值，暂时无用
__return_reference:__ 返回值是否为引用，一般很少会用到
required_num_args: required参数数
这两个宏需要与 ZEND_END_ARG_INFO() 配合使用：

ZEND_BEGIN_ARG_INFO_EX(arginfo_my_func_1, 0, 0, 2)
...
ZEND_END_ARG_INFO
 
接着就是在上面两个宏中间定义每一个参数的zend_internal_arg_info，PHP提供的宏有：

//pass_by_ref表示是否引用传参，name为参数名称
#define ZEND_ARG_INFO(pass_by_ref, name) { #name, NULL, 0, pass_by_ref, 0, 0 },
 
//只声明此参数为引用传参
#define ZEND_ARG_PASS_INFO(pass_by_ref) { NULL, NULL, 0, pass_by_ref, 0, 0 },
 
//显式声明此参数的类型为指定类的对象，等价于PHP中这样声明：MyClass $obj
#define ZEND_ARG_OBJ_INFO(pass_by_ref, name, classname, allow_null) { #name, #classname, IS_OBJECT, pass_by_ref, allow_null, 0 },
 
//显式声明此参数类型为数组，等价于：array $arr
#define ZEND_ARG_ARRAY_INFO(pass_by_ref, name, allow_null) { #name, NULL, IS_ARRAY, pass_by_ref, allow_null, 0 },
 
//显式声明为callable，将检查函数、成员方法是否可调
#define ZEND_ARG_CALLABLE_INFO(pass_by_ref, name, allow_null) { #name, NULL, IS_CALLABLE, pass_by_ref, allow_null, 0 },
 
//通用宏，自定义各个字段
#define ZEND_ARG_TYPE_INFO(pass_by_ref, name, type_hint, allow_null) { #name, NULL, type_hint, pass_by_ref, allow_null, 0 },
 
//声明为可变参数
#define ZEND_ARG_VARIADIC_INFO(pass_by_ref, name) { #name, NULL, 0, pass_by_ref, 0, 1 },
看个例子：

function my_func_1(&$a, Exception $c){..}
用内核实现则可以这么定义：

ZEND_BEGIN_ARG_INFO_EX(arginfo_my_func_1, 0, 0, 1)
	ZEND_ARG_INFO(1, a) //引用
	ZEND_ARG_OBJ_INFO(0, b, Exception, 0) //注意：这里不要把字符串加""
ZEND_END_ARG_INFO()
展开后：

static const zend_internal_arg_info name[] = {
	//多出来的这个是给返回值用的
	{ (const char*)(zend_uintptr_t)(2), NULL, 0, 0, 0, 0 },
	{ "a", NULL, 0, 0, 0, 0 },
	{ "b", "Exception", 8, 1, 0, 0 },
}
第一个数组元素用于记录必传参数的数量以及返回值是否为引用。定义完这个数组接下来就需要把这个数组告诉函数：

const zend_function_entry mytest_functions[] = {
	PHP_FE(my_func_1, arginfo_my_func_1)
	PHP_FE(my_func_2, NULL)
	PHP_FE_END //末尾必须加这个
};
引用参数通过 zend_parse_parameters() 解析时只能使用"z"解析，不能再直接解析为zend_value了，否则引用将失效：

PHP_FUNCTION(my_func_1)
{
	zval *lval; //必须为zval，定义为zend_long也能解析出，但不是引用
	zval *obj;
	if(zend_parse_parameters(ZEND_NUM_ARGS(), "zo", &lval, &obj) == FAILURE){
	    RETURN_FALSE;
	}
	//lval的类型为IS_REFERENCE
	zval *real_val = Z_REFVAL_P(lval); //获取实际引用的zval地址：&(lval.value->ref.val)
	Z_LVAL_P(real_val) = 100; //设置实际引用的类型
}
$a = 90;
$b = new Exception;
my_func_1($a, $b);
echo $a; //100
函数返回值
调用内部函数时其返回值指针作为参数传入，这个参数为 zval *return_value ，如果函数有返回值直接设置此指针即可，需要特别注意的是设置返回值时需要增加其引用计数，举个例子来看:

PHP_FUNCTION(my_func_1)
{
    zval *arr;
    if(zend_parse_parameters(ZEND_NUM_ARGS(), "a", &arr) == FAILURE){
        RETURN_FALSE;
    }
    //增加引用计数
    Z_ADDREF_P(arr);
    //设置返回值为数组：
    ZVAL_ARR(return_value, Z_ARR_P(arr));
}
此函数接收一个数组，然后直接返回该数组，相当于：

function my_func_1($arr){
    return $arr;
}
虽然可以直接设置return_value，但实际使用时并不建议这么做，因为PHP提供了很多专门用于设置返回值的宏，这些宏定义在 zend_API.h 中：

//返回布尔型，b：IS_FALSE、IS_TRUE
#define RETURN_BOOL(b) { RETVAL_BOOL(b); return; }
 
//返回NULL
#define RETURN_NULL() { RETVAL_NULL(); return;}
 
//返回整形，l类型：zend_long
#define RETURN_LONG(l) { RETVAL_LONG(l); return; }
 
//返回浮点值，d类型：double
#define RETURN_DOUBLE(d) { RETVAL_DOUBLE(d); return; }
 
//返回字符串，可返回内部字符串，s类型为：zend_string *
#define RETURN_STR(s) { RETVAL_STR(s); return; }
 
//返回内部字符串，这种变量将不会被回收，s类型为：zend_string *
#define RETURN_INTERNED_STR(s) { RETVAL_INTERNED_STR(s); return;
}
 
//返回普通字符串，非内部字符串，s类型为：zend_string *
#define RETURN_NEW_STR(s) { RETVAL_NEW_STR(s); return; }
 
//拷贝字符串用于返回，这个会自己加引用计数，s类型为：zend_string *
#define RETURN_STR_COPY(s) { RETVAL_STR_COPY(s); return; }
 
//返回char *类型的字符串，s类型为char *
#define RETURN_STRING(s) { RETVAL_STRING(s); return; }
 
//返回char *类型的字符串，s类型为char *，l为字符串长度，类型为size_t
#define RETURN_STRINGL(s, l) { RETVAL_STRINGL(s, l); return; }
 
//返回空字符串
#define RETURN_EMPTY_STRING() { RETVAL_EMPTY_STRING(); return;
}
 
//返回资源，r类型：zend_resource *
#define RETURN_RES(r) { RETVAL_RES(r); return; }
 
//返回数组，r类型：zend_array *
#define RETURN_ARR(r) { RETVAL_ARR(r); return; }
 
//返回对象，r类型：zend_object *
#define RETURN_OBJ(r) { RETVAL_OBJ(r); return; }
 
//返回zval
#define RETURN_ZVAL(zv, copy, dtor) { RETVAL_ZVAL(zv, copy, dtor); re
turn; }
 
//返回false
#define RETURN_FALSE { RETVAL_FALSE; return; }
 
//返回true
#define RETURN_TRUE { RETVAL_TRUE; return;}

https://www.fzb.me/2015-2-27-zend-api-manual.html
