---
title: php扩展实现多线程
layout: post
category: php
author: 夏泽民
---
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
<!-- more -->

{% raw %}
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

{% endraw %}

https://blog.csdn.net/weixin_34419321/article/details/89225952

http://www.hongweipeng.com/index.php/archives/1026/

因为PHP的 call_user_method()和 call_user_method_array()被标记为不赞成我想知道什么替代方案被推荐？
一种方法是使用call_user_func()，因为通过给一个具有对象的数组和方法名作为第一个参数，就像已弃用的函数一样。由于这个功能没有被标记为不赞成，我认为这个原因不是非OOP时尚的使用方式吗？

另外我可以想到的是使用Reflection API，这可能是最舒适和面向未来的替代方案。然而，它是更多的代码，我可以形象，它比使用上面提到的功能慢。

我感兴趣的是

>有没有一种全新的技术来通过名称来调用对象的方法？
>哪个是最快/最好/官方的更换？
>弃用的原因是什么？

如您所说 call_user_func可以轻松地复制此功能的行为。有什么问题？
call_user_method页甚至列出了替代方案：

<?php
call_user_func(array($obj, $method_name), $parameter /* , ... */);
call_user_func(array(&$obj, $method_name), $parameter /* , ... */); // PHP 4
?>
至于为什么这是不赞成的，this posting explains it:

This is
because the call_user_method() and call_user_method_array() functions
can easily be duplicated by:

old way:
call_user_method($func, $obj, "method", "args", "go", "here");

new way:
call_user_func(array(&$obj, "method"), "method", "args", "go", "here");

就个人而言，我可能会用乍得发布的变量变量建议。
