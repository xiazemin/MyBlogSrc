---
title: call_user_function_ex 扩展调用php脚本函数
layout: post
category: php
author: 夏泽民
---
其中可能我认为最关键的应该是在PHP扩展里怎么调用用户空间里的函数了！对于一个framework来说，最基本的功能就是路由到请求对应的Action了。

在PHP扩展里是通过 call_user_function_ex 函数来调用用户空间的函数的。下面我们来分析下这个函数的使用方式吧。

下面这个是call_user_function_ex 函数的定义：

ZEND_API int call_user_function_ex(
	HashTable *function_table, 
	zval **object_pp, 
	zval *function_name, 
	zval **retval_ptr_ptr, 
	zend_uint param_count, 
	zval **params[], 
	int no_seperation, 
	HashTable *symbol_table TSRMLS_DC);
function_table is the hash table where the function you wish to call is located. If you\'re using object_pp, set this to NULL. If the function is global, most likely it\'s located in the hash table returned by the macro CG() with the parameter `function_table\', i.e.

CG(function_table)
object_pp is a pointer to a zval pointer where an initialized object is located. If you use this, set function_table to NULL, as previously noted.
function_name is a pointer to a zval which contains the name of the function in string form.
retval_ptr_ptr is a pointer to a zval pointer which will contain the return value of the function. The zval passed doesn\'t need to be initialized, and it may cause problems if you initialize it when it\'s not neccesary. You  must  always pass a real pointer to a zval pointer, you may not use NULL for this as it will cause a segmentation fault.
param_count is the number of parameters you wish to pass to the function being called.
params is an array of pointers to zval pointers. Note: this is _not_ a PHP/zval array, it is a C array. Example:
zval *foo; zval *bar; zval **params[2]; params[0] = &foo; params[1] = bar;
no_seperation is either 1 or 0, 0 being no zval seperation, 1 enabling zval seperation. 
symbol_table is the hash table for symbols. I currently don\'t know what this is, so when I find out, I\'ll edit the post and put it here
After the symbol_table parameter, you should put TSRMLS_CC to make it threadsafe.
<!-- more -->
https://my.oschina.net/jackin/blog/172926
https://forums.phpfreaks.com/topic/1303-call_user_function_ex-documentation/


<?php
class demo {
    public function get_site_name ($prefix) {
        return $prefix."信海龙的博客\n";
    }
}
function get_site_url ($prefix) {
    return $prefix."www.bo56.com\n";
}

function call_function ($obj, $fun, $param) {
    if ($obj == null) {
        $result = $fun($param);
    } else {
        $result = $obj->$fun($param);
    }
    return $result;
}
$demo = new demo();
echo call_function($demo, "get_site_name", "site name:");
echo call_function(null, "get_site_url", "site url:");
?>

我们将要使用扩展实现call_function方法的功能。

代码
基础代码
这个扩展，我们将在say扩展上增加call_function()。say扩展相关代码大家请看这篇博文。PHP7扩展开发之hello word 文中已经详细介绍了如何创建一个扩展和提供了源码下载。

代码实现
call_function的源码如下：

PHP_FUNCTION(call_function)
{
    zval            *obj = NULL;
    zval            *fun = NULL;
    zval            *param = NULL;
    zval            retval;
    zval            args[1];

#ifndef FAST_ZPP
    /* Get function parameters and do error-checking. */
    if (zend_parse_parameters(ZEND_NUM_ARGS(), "zzz", &obj, &fun, &param) == FAILURE) {
        return;
    }
#else
    ZEND_PARSE_PARAMETERS_START(3, 3)
        Z_PARAM_ZVAL(obj)
        Z_PARAM_ZVAL(fun)
        Z_PARAM_ZVAL(param)
    ZEND_PARSE_PARAMETERS_END();
#endif

    args[0] = *param;
    if (obj == NULL || Z_TYPE_P(obj) == IS_NULL) {
        call_user_function_ex(EG(function_table), NULL, fun, &retval, 1, args, 0, NULL);
    } else {
        call_user_function_ex(EG(function_table), obj, fun, &retval, 1, args, 0, NULL);
    }
    RETURN_ZVAL(&retval, 0, 1);
}

代码解读
参数的接受之前有过文章详细说明过，这里就不再说了。这次我们主要说下call_user_function_ex方法的使用。

call_user_function_ex方法用于调用函数和方法。参数说明如下：
* 第一个参数：方法表。通常情况下，写 EG(function_table) 更多信息查看
* 第二个参数：对象。如果不是调用对象的方法，而是调用函数，填写NULL
* 第三个参数：方法名。
* 第四个参数：返回值。
* 第五个参数：参数个数。
* 第六个参数：参数值。是一个zval数组。
* 第七个参数：参数是否进行分离操作。详细的，你可以搜索下 PHP 参数分离。查看相关文章
* 第八个参数：符号表。一般情况写设置为NULL即可。

https://blog.csdn.net/u011957758/article/details/72513935

https://www.cnblogs.com/wuhen781/p/6132878.html
https://forums.phpfreaks.com/topic/1303-call_user_function_ex-documentation/

调用类的内部方法和 call_user_func 函数的调用方式一样，都是使用了数组的形式来调用。

call_user_func ( callback $function [, mixed $parameter [, mixed $... ]] )

调用第一个参数所提供的用户自定义的函数。
返回值：返回调用函数的结果，或FALSE。

example

：

<?php
function eat($fruit) //参数可以为多个
{
echo "You want to eat $fruit, no problem";
}
call_user_func('eat', "apple"); //print: You want to eat apple, no problem;
call_user_func('eat', "orange"); //print: You want to eat orange,no problem;
?>

 


调用类的内部方法：

<?php
class myclass {
function say_hello($name)
{
echo "Hello!$name";
}
}
$classname = "myclass";
//调用类内部的函数需要使用数组方式 array(类名，方法名)
call_user_func(array($classname, 'say_hello'), 'dain_sun');
//print Hello! dain_sun
?>

 


call_user_func_array 函数和 call_user_func 很相似，只是
使
用了数组
的传递参数形式，让参数的结构更清晰:

call_user_func_array

( callback

$function

, array
$param_arr

)

类中的回调函数最好是static的

数组中的第一个参数（类名）也可以用实例化的对象来代替。

调用用户定义的函数，参数为数组形式。
返回值：返回调用函数的结果，或FALSE。

<?php
function debug($var, $val)
{
echo "variable: $var <br> value: $val <br>";
echo "<hr>";
}
$host = $_SERVER["SERVER_NAME"];
$file = $_SERVER["PHP_SELF"];
call_user_func_array('debug', array("host", $host));
call_user_func_array('debug', array("file", $file));
?>

调用类的内部方法和 call_user_func 函数的调用方式一样，都是使用了数组的形式来调用。

exmaple:


<?php
class test
{
function debug($var, $val)
{
echo "variable: $var <br> value: $val <br>";
echo "<hr>";
}
}
$host = $_SERVER["SERVER_NAME"];
$file = $_SERVER["PHP_SELF"];
call_user_func_array(array('test', 'debug'), array("host", $host));
call_user_func_array(array('test', 'debug'), array("file", $file));
?>



注：call_user_func
函数和call_user_func_array函数都支持引用。

<?php
function increment(&$var)
{
$var++;
}
$a = 0;
call_user_func('increment', $a);
echo $a; // 0
call_user_func_array('increment', array(&$a)); // You can use this instead
echo $a; // 1
?>


很多时候，需要把控制权限交给用户，或者在扩展里完成某件事后去回调用户的方法。

在PHP扩展里是通过 call_user_function_ex 函数来调用用户空间的函数的。

定义
它的定义在 Zend/zend_API.h :

#define call_user_function_ex(function_table, object, function_name, retval_ptr, param_count, params, no_separation, symbol_table)

    _call_user_function_ex(object, function_name, retval_ptr, param_count, params, no_separation)
通过宏定义替换为_call_user_function_ex,其中参数 function_table 被移除了，它之所以在API才存在大概是为了兼容以前的写法。函数的真正定义是：

ZEND_API int _call_user_function_ex(
    zval *object, 
    zval *function_name, 
    zval *retval_ptr, 
    uint32_t param_count, 
    zval params[], 
    int no_separation);
参数分析：

zval *object:这个是用来我们调用类里的某个方法的对象。

zval *function_name:要调用的函数的名字。

zval *retval_ptr：收集回调函数的返回值。

uint32_t param_count：回调函数需要传递参数的个数。

zval params[]: 参数列表。

int no_separation：是否对zval进行分离，如果设为1则直接会出错，分离的作用是为了优化空间。

回调功能的实现
PHP_FUNCTION(hello_callback)
{
    zval *function_name;
    zval retval;
    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "z", &function_name) == FAILURE) {
        return;
    }
    if (Z_TYPE_P(function_name) != IS_STRING) {
        php_printf("Function require string argumnets!");
        return;
    }
    //TSRMLS_FETCH();
    if (call_user_function_ex(EG(function_table), NULL, function_name, &retval, 0, NULL, 0, NULL TSRMLS_CC) != SUCCESS) {
        php_printf("Function call failed!");
        return;
    }

    *return_value = retval;
    zval_copy_ctor(return_value);
    zval_ptr_dtor(&retval);

}
zval_copy_ctor()原始（zval）的内容拷贝给它。zval_ptr_dtor()释放空间。return_value不是一个函数外的变量，它的由函数声明里的变量。PHP_FUNCTION(hello_callback)这个声明是简写，最终会被预处理宏替换为：

void zif_hello_callback(zend_execute_data *execute_data, zval *return_value)
return_value变量其实也就是最终返回给调用脚本的，RETURN_STR(s) 等返回函数最终也都是宏替换为对该变量的操作。

测试脚本：

<?php
function fun1() {
    for ($i = 0; $i < 5; $i++) {
        echo 'fun1:'.$i."\n";
    }
    return 'call end';
}

echo hello_callback('fun1');
一个并行扩展
早期的php不支持多进程多线程的，现在随着发展有很多扩展不断完善它，诸如pthread,swoole等，不仅能多线程，而且能实现异步。

利用c语言多线程pthread库来实现一个简单的并行扩展。

先声明我们一会用到的结构：

struct myarg
{
    zval *fun;
    zval ret;
};
线程函数：

static void my_thread(struct myarg *arg) {
    zval *fun = arg->fun;
    zval ret = arg->ret;
    if (call_user_function_ex(EG(function_table), NULL, fun, &ret, 0, NULL, 0, NULL TSRMLS_CC) != SUCCESS) {
        return;
    }
}
函数的实现：

PHP_FUNCTION(hello_thread)
{
    pthread_t tid;
    zval *fun1, *fun2;
    zval ret1, ret2;
    struct myarg arg;
    int ret;
    if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "zz", &fun1, &fun2) == FAILURE) {
        return;
    }
    arg.fun = fun1;
    arg.ret = ret1;
    ret = pthread_create(&tid, NULL, (void*)my_thread, (void*)&arg);
    if(ret != 0) {
        php_printf("Thread Create Error\n");
        exit(0);
    }
    if (call_user_function_ex(EG(function_table), NULL, fun2, &ret2, 0, NULL, 0, NULL TSRMLS_CC) != SUCCESS) {
        return;
    }
    pthread_join(tid, NULL);
    RETURN_NULL();

}
测试脚本：

<?php
function fun1() {
    for ($i = 0; $i < 5; $i++) {
        echo 'fun1:'.$i.'\n';
    }
}

function fun2() {
    for ($i = 0; $i < 5; $i++) {
        echo 'fun2:'.$i.'\n';
    }
}

hello_thread('fun1', 'fun2');
echo 'after 多并发';

https://segmentfault.com/a/1190000007648157
http://www.hongweipeng.com/index.php/archives/1026/?utm_source=tuicool&utm_medium=referral
https://blog.csdn.net/qq_32783703/article/details/80641355

