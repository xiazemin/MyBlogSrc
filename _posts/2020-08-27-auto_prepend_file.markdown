---
title: auto_prepend_file与auto_append_file使用方法
layout: post
category: php
author: 夏泽民
---
auto_prepend_file与auto_append_file使用方法

如果需要将文件require到所有页面的顶部与底部。



第一种方法：在所有页面的顶部与底部都加入require语句。

例如：

require('header.php');
页面内容
require('footer.php');
但这种方法如果需要修改顶部或底部require的文件路径，则需要修改所有页面文件。而且需要每个页面都加入require语句，比较麻烦。


第二种方法：使用auto_prepend_file与auto_append_file在所有页面的顶部与底部require文件。

php.ini中有两项

auto_prepend_file 在页面顶部加载文件

auto_append_file  在页面底部加载文件

使用这种方法可以不需要改动任何页面，当需要修改顶部或底部require文件时，只需要修改auto_prepend_file与auto_append_file的值即可。



例如：修改php.ini，修改auto_prepend_file与auto_append_file的值。

auto_prepend_file = "/home/fdipzone/header.php"
auto_append_file = "/home/fdipzone/footer.php"
修改后重启服务器，这样所有页面的顶部与底部都会require /home/fdipzone/header.php 与 /home/fdipzone/footer.php


注意：auto_prepend_file 与 auto_append_file 只能require一个php文件，但这个php文件内可以require多个其他的php文件。



如果不需要所有页面都在顶部或底部require文件，可以指定某一个文件夹内的页面文件才调用auto_prepend_file与auto_append_file

在需要顶部或底部加载文件的文件夹中加入.htaccess文件，内容如下：

php_value auto_prepend_file "/home/fdipzone/header.php"
php_value auto_append_file "/home/fdipzone/footer.php"

这样在指定.htaccess的文件夹内的页面文件才会加载 /home/fdipzone/header.php 与 /home/fdipzone/footer.php，其他页面文件不受影响。
使用.htaccess设置，比较灵活，不需要重启服务器，也不需要管理员权限，唯一缺点是目录中每个被读取和被解释的文件每次都要进行处理，而不是在启动时处理一次，所以性能会有所降低。
<!-- more -->
https://blog.csdn.net/fdipzone/article/details/39064001
https://www.jianshu.com/p/33a716c9a916

在 php.ini 中有两个配置参数，auto_prepend_file 和 auto_append_file，其作用相当于php代码 require 或 include，使用这两个指令包含的文件如果该文件不存在，将产生一个警告。

auto_prepend_file 表示在php程序加载应用程序前加载指定的php文件

auto_append_file 表示在php代码执行完毕后加载指定的php文件
 
在某些场合下我们可能要对所有的代码在执行前或者执行后进行统一处理，这时这2个设置项就非常有用了。例如为了实现一些自动化工作就经常用到这2个设置项，例如分析代码覆盖率，自动代码分析，自动sql分析等等，注意该指令更多的适用于测试环境调试而用。

对于Windows，其设置如下所示：

1
auto_prepend_file="c:/Program Files/Apache2.2/include/header.php"
2
auto_append_file="c:/Program Files/Apache2.2/include/footer.php"
对于UNIX，其设置如下所示：

1
auto_prepend_file="/home/username/include/header.php"
2
auto_append_file="/home/username/include/footer.php"
注意：

(1) auto_prepend_file 与 auto_append_file 只能require一个php文件，但这个php文件内可以require多个其他的php文件

(2) 配置了该参数后，所有使用该php服务器的项目都会加载相应的配置文件

当然也可以对单个目录进行不同的配置，这样做的前提是服务器允许重设其主配置文件。要给目录设定自动前加入和自动追加，需要在该目录中创建一个名为.htaccess的文 
件。这个文件需要包含如下两行代码：

apache 环境

1
php_value auto_prepend_file /home/phpernote/include/header.php
2
php_value auto_append_file /home/phpernote/include/footer.php
nginx 环境

1
fastcgi_param PHP_VALUE "auto_prepend_file=/home/phpernote/include/header.php";
2
fastcgi_param PHP_VALUE "auto_append_file=/home/phpernote/include/footer.php";
注意：其语法与配置文件php.ini中的相应选项有所不同，和行开始处的php_value一样：没有等号。许多php.ini中的配置设定也可以按这种方法进行修改。

在.htaccess中设置选项，而不是php.ini中或是在Web服务器的配置文件中进行设置，将带来极大的灵活性。可以在一台只影响你的目录的共享机器上进行。不需要重新启动服务器而且不需要管理员权限。使用.htaccess方法的一个缺点就是目录中每个被读取和被解析的文件每次都要进行处理，而不是只在启动时处理一次，所以性能会有所降低。