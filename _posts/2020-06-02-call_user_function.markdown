---
title: call_user_function
layout: post
category: php
author: 夏泽民
---
扩展可能需要调用用户自定义的函数或者其他扩展定义的内部函数,PHP提供的函数调用API的使用:

ZEND_API int call_user_function(HashTable *function_table, zval *object,zval *function_name, zval *retval_ptr, uint32_t param_count, zval params[
]);
各参数的含义：

__function_table:__ 函数符号表，普通函数是EG(function_table)，如果是成员方法则是zend_class_entry.function_table
object: 调用成员方法时的对象
__function_name:__ 调用的函数名称
__retval_ptr:__ 函数返回值地址
__param_count:__ 参数数量
params: 参数数组
从接口的定义看其使用还是很简单的，不需要我们关心执行过程中各阶段复杂的操作。
<!-- more -->
{% raw %}
下面从一个具体的例子看下其使用：

（1）在PHP中定义了一个普通的函数，将参数$i加上100后返回：

function mySum($i){
    return $i+100;
}

（2）接下来在扩展中调用这个函数：

PHP_FUNCTION(my_func_2)
{
    zend_long i;
    zval call_func_name, call_func_ret, call_func_params[1];
    uint32_t call_func_param_cnt = 1;
    zend_string *call_func_str;
    char *func_name = "mySum";
    if(zend_parse_parameters(ZEND_NUM_ARGS(), "l", &i) == FAILURE){
        RETURN_FALSE;
    }
    //分配zend_string:调用完需要释放
    call_func_str = zend_string_init(func_name, strlen(func_name), 0);
    //设置到zval
    ZVAL_STR(&call_func_name, call_func_str);
    
    //设置参数
    ZVAL_LONG(&call_func_params[0], i);
    
    //call
    if(SUCCESS != call_user_function(EG(function_table), NULL, &call_func_name, &call_func_ret, call_func_param_cnt, call_func_params)){
        zend_string_release(call_func_str);
        RETURN_FALSE;
    }
    
    zend_string_release(call_func_str);
    RETURN_LONG(Z_LVAL(call_func_ret));
}
（3）最后调用这个内部函数：

function mySum($i){
    return $i+100;
}
echo my_func_2(60);
===========[output]===========
160
call_user_function() 并不是只能调用PHP脚本中定义的函数，内核或其它扩展注册
的函数同样可以通过此函数调用，比如：array_merge()。

PHP_FUNCTION(my_func_1)
{
    zend_array *arr1, *arr2;
    zval call_func_name, call_func_ret, call_func_params[2];
    uint32_t call_func_param_cnt = 2;
    zend_string *call_func_str;
    char *func_name = "array_merge";
    if(zend_parse_parameters(ZEND_NUM_ARGS(), "hh", &arr1, &arr2) == FAILURE){
        RETURN_FALSE;
    }
    
    //分配zend_string
    call_func_str = zend_string_init(func_name, strlen(func_name), 0);
    //设置到zval
    ZVAL_STR(&call_func_name, call_func_str);
    
    ZVAL_ARR(&call_func_params[0], arr1);
    ZVAL_ARR(&call_func_params[1], arr2);
    
    if(SUCCESS != call_user_function(EG(function_table), NULL, &call_func_name, &call_func_ret, call_func_param_cnt, call_func_params)){
    zend_string_release(call_func_str);
    RETURN_FALSE;
    }
    zend_string_release(call_func_str);
    RETURN_ARR(Z_ARRVAL(call_func_ret));
}
$arr1 = array(1,2);
$arr2 = array(3,4);
$arr = my_func_1($arr1, $arr2);
var_dump($arr);
{% endraw %}
https://blog.csdn.net/rorntuck7/article/details/86240202
