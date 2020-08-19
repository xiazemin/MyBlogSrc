---
title: log_errors display_errors
layout: post
category: php
author: 夏泽民
---
display_errors,开启log_errors,配置error_log路径
1）display_error  

display_errors ，错误回显，一般常用与开发环境。如果在生产环境中开启选项，错误回显会暴露出非常多的敏感信息，为攻击者下一步攻击提供便利。推荐关闭此选项。 

display_errors = On 

开启状态下，若出现错误，则报错，出现错误提示 

dispaly_errors = Off 

关闭状态下，若出现错误，则提示：服务器错误。但是不会出现错误提示 

对于PHP开发人员来说，一旦某个产品投入使用，那么第一件事就是应该将display_errors选项关闭，以免因为这些错误所透露的路径、数据库连接、数据表等信息而遭到黑客攻击。

             既然生产环境中不能出现错误提示信息，而当网站出现问题，我们有需要查看具体的错误信息时有需要怎么做呢？没错，这就用到了下面的错误日志记录。

（2）log_error

          log_error，错误日志，一般用于生产环境中。开发人员可以分析错误日志内容，进而发现并解决问题。

          log_error=on 开启错误日志  

          log_error=off  关闭错误日志

          日志默认是记录到WEB服务器的日志文件里，比如Apache的error.log文件。 当然也可以记录错误日志到指定的文件中。

     # vim /etc/php.inidisplay_errors = Off 
     log_errors = On 
     error_log = /var/log/php-error.log 
          在生产环境中，一旦开启了错误日志记录功能，个人强烈建议设置错误日志目录。
<!-- more -->
https://www.jianshu.com/p/5af8c0ba13e5