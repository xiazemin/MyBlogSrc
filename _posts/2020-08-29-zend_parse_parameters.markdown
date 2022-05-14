---
title: zend_parse_parameters php扩展程序的参数传递 扩展类传参数
layout: post
category: php
author: 夏泽民
---
在php扩展程序的开发中，涉及参数接受处理时，第一步便是要对传入参数进行判断，如生成的扩展示例代码：
if (zend_parse_parameters(ZEND_NUM_ARGS(), "s", &arg, &arg_len) == FAILURE) {
        return;
    }如上述示例代码，其判断有
1：判断是否有入参，如果没有入参就会报缺少参数错误。
2：判断入参是不是字符串，如果不是字符串就会把参数类型错误。
    先说一下参数类型吧，上面的例子中只有字符串，没有其它类型。实际PHP扩展程序中的类型不少，有整型，浮点型，还有zval类型。zval是Zend引擎的值容器，无论这个变量是个简单的布尔值，字符串或者其他任何类型值，其信息总是一个完整的zval结构。可以认为是一个简单数据的底层复杂描述的结构。
    
    
    PHP_FUNCTION(kermitcal)
{   
    char *username;
    size_t username_len;
    char *age;
    size_t age_len;
    char *email = "admin@04007.cn";
    size_t email_len = sizeof("admin@04007.cn") -1;
    zend_string *strg;
    #使用sl|s表示|后的这个email参数可以不传递，使用默认值。
    if (zend_parse_parameters(ZEND_NUM_ARGS(), "sl|s",&username, &username_len,&age, &age_len, &email, &email_len) == FAILURE){
        php_printf("need params username(string) and age(int).!");
        RETURN_NULL();
    }
    strg = strpprintf(0, "大家好，我叫%s, 今年%d岁, 我的邮箱是:%s \n", username, age, email);
    RETURN_STR(strg);
}

<!-- more -->
https://www.cnblogs.com/dearmrli/p/6553542.html

1、zend_parse_parameters
如果你实际操作中看了PHP_FUNCTION(confirm_helloworld_compiled)中的内容，就会发现这个函数，它的作用就是接收参数，可以类比PHP的sscanf()

zend_parse_parameters(int num_args TSRMLS_DC, char *type_spec, &参数1,&参数2…);
第一个参数是传递给函数的参数个数
第二个参数是为了线程安全，总是传递TSRMLS_CC宏。关于TSRMLS_CC，感兴趣的小伙伴可以阅读下鸟哥的博客揭秘TSRM
第三个参数是一个字符串，指定了函数期望的参数类型，后面紧跟着需要随参数值更新的变量列表。由于PHP是弱类型语言，采用松散的变量定义和动态的类型判断，而c语言是强类型的，而zend_parse_parameters()就是为了把不同类型的参数转化为期望的类型。(当实际的值无法强制转成期望的类型的时候，会发出一个警告)
使用示例
zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "s",&name, &name_len)
参数表

额外用法
第四第五直到第n个参数，都是要传进来的值的数值。
2、zend_get_arguments
ZEND_FUNCTION(my_function) {
    zval *email;
    if (zend_get_parameters(ZEND_NUM_ARGS(), 1, &email)== FAILURE) {
        php_error_docref(NULL TSRMLS_CC, E_WARNING,"至少需要一个参数");
        RETURN_NULL();
    }
    // ... 
}
1.能够兼容老版本的PHP，并且只以zval为载体来接收参数。 
2.直接获取，而不做解析，不会进行类型转换，所有的参数在扩展实现中的载体都需要是zval类型的。 
3.接受失败的时候，不会自己抛出错误，也不能方便的处理有默认值的参数。 
4.会自动的把所有符合copy-on-write的zval进行强制分离，生成一个崭新的copy送到函数内部。

综合评价：还是用zend_parse_parameters吧，这个函数了解下即可，不给力。
3、zend_get_parameters_ex()
ZEND_FUNCTION(my_function) {
    zval **msg;
    if (zend_get_parameters_ex(1, &msg) == FAILURE) {
        php_error_docref(NULL TSRMLS_CC, E_WARNING,"至少需要一个参数");
        RETURN_NULL();
    }
    // ...
}
1.不需要ZEND_NUM_ARGS()作为参数，因为它是在是在后期加入的，那个参数已经不再需要了。
2.此函数基本同zend_get_parameters()。 
3.唯一不同的是它不会自动的把所有符合copy-on-write的zval进行强制分离，会用到老的zval的特性
综合评价：极端情况下可能会用到，这个函数了解下即可。
4、zend_get_parameters_array_ex()和zend_get_parameters_array()

//程序首先获取参数数量，然后通过safe_emalloc函数申请了相应大小的内存来存放这些zval**类型的参数。这里使用了zend_get_parameters_array_ex()函数来把传递给函数的参数填充到args中。 
//这个参数专门用于解决像php里面的var_dump的一样，可以无限传参数进去的函数的实现

ZEND_FUNCTION(var_dump) {
    int i, argc = ZEND_NUM_ARGS();
    zval ***args;

    args = (zval ***)safe_emalloc(argc, sizeof(zval **), 0);
    if (ZEND_NUM_ARGS() == 0 || zend_get_parameters_array_ex(argc, args) == FAILURE) {
        efree(args);
        WRONG_PARAM_COUNT;
    }
    for (i=0; i<argc; i++) {
        php_var_dump(args[i], 1 TSRMLS_CC);
    }
    efree(args);
}
1.用于应对无限参数的扩展函数的实现。 
2.zend_get_parameters_array与zend_get_parameters_array_ex唯一不同的是它将zval*类型的参数填充到args中，并且需要ZEND_NUM_ARGS()作为参数。
综合评价:当遇到确实需要处理无限参数的时候，真的要用这个函数了。



扩展类传参数

https://blog.csdn.net/weixin_30265103/article/details/97891064



https://blog.csdn.net/rorntuck7/article/details/86307015

https://segmentfault.com/a/1190000007575322
