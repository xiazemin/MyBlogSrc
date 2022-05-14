---
title: magic_quotes_gpc
layout: post
category: php
author: 夏泽民
---
php中的magic_quotes_gpc是配置在php.ini中的，他的作用类似addslashes()，就是对输入的字符创中的字符进行转义处理。他可以对$_POST、$__GET以及进行数据库操作的sql进行转义处理，防止sql注入。

对于PHP magic_quotes_gpc=on的情况，

我们可以不对输入和输出数据库的字符串数据作

addslashes()和stripslashes()的操作,数据也会正常显示。

如果此时你对输入的数据作了addslashes()处理，
那么在输出的时候就必须使用stripslashes()去掉多余的反斜杠。

对于PHP magic_quotes_gpc=off 的情况

必须使用addslashes()对输入数据进行处理，但并不需要使用stripslashes()格式化输出
因为addslashes()并未将反斜杠一起写入数据库，只是帮助mysql完成了sql语句的执行。

程序中可以通过get_magic_quotes_gpc来获取magic_quotes_gpc环境变量的值，从而判断是否使用addslashes()和stripslashes()

addslashes -- 使用反斜线引用字符串

描述
string addslashes ( string str)

返回字符串，该字符串为了数据库查询语句等的需要在某些字符前加上了反斜线。这些字符是单引号（'）、双引号（"）、反斜线（''）与 NUL（NULL 字符）

stripslashes -- 函数删除由 addslashes() 函数添加的反斜杠。

描述
string addslashes ( string str)

返回字符串，该函数可用于清理从数据库中或者从 HTML 表单中取回的数
<!-- more -->
https://www.php.net/manual/zh/security.magicquotes.disabling.php

https://baike.baidu.com/item/magic_quotes_gpc

https://my.oschina.net/u/2268393/blog/521554

https://www.cnblogs.com/52php/p/5687219.html

https://www.cnblogs.com/timelesszhuang/p/3726736.html
