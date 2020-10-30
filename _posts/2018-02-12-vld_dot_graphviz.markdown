---
title: vld_dot_graphviz
layout: post
category: php
author: 夏泽民
---
<!-- more -->
PHP7.1下 vld扩展的安装使用
PHP7.1下 vld扩展的安装使用
原创 2017年03月29日 08:03:05 595
1）git clone https://github.com/derickr/vld.git

2）cd vld

3）phpize

4）./configure

5）make && make install

6）添加ext-vld.ini配置文件

7）重启fpm 

8）php -m | grep vld 查看扩展

9）php -dvld.active test.php 测试vld扩展

PHP7.1下 vld扩展的安装使用
原创 2017年03月29日 08:03:05 595
1）git clone https://github.com/derickr/vld.git

2）cd vld

3）phpize

4）./configure

5）make && make install

6）添加ext-vld.ini配置文件

7）重启fpm 

8）php -m | grep vld 查看扩展

9）php -dvld.active test.php 测试vld扩展



关于VLD扩展显示信息的一点点解释


其中：

branch analysis from position 在分析数组时使用

return found是否返还

filename 分析的文件名

function name函数名

number of ops生成的操作数

compiled vars编译期间的变量，PHP5后添加，是一个缓存优化，在PHP源码中以IS_CV标记

op list生成的中间代码的变量列表



-dvld.active输出的是VLD的默认设置，使用-dvld.verbosity可以查看更加详细的内容

包含各个中间代码的操作数等



若只想看到输出的中间代码，并不想实际执行这段代码，可以使用-dvld.execute = 0来禁用代码的执行

php -dvld.active=1 -dvld.execute=0 test.php

它还可以支持输出.dot文件

php -dvld.active=1 -dvld.save_dir='D:\tmp' -dvld.save_paths=1 -dvld.dump_paths=1 t.php 会将生成的中间代码的信息输出再D:/tmp/path.dot中



-dvld.format是否以自定义的格式输出，默认为否，是指以-dvld.col_sep指定的参数间隔

-dvld.col_sep在-dvld.format参数启用时才会有效，默认为 \t

-dvld.verbosity是否显示更加详细的信息，默认为1，其值可以是0，1，2，3 或者小于0只是比1小的效果会喝0一样，负数的效果和3的效果一样

-dvld.save_dir指定文件的输出路径，默认/tmp

-dvld.save_path指定文件输出的路径，默认0表示不输出文件

-dvld.dump_paths控制输出的内容，0或1 默认1，即输出内容


2、linux下咋安装graphviz
http://www.cnblogs.com/sld666666/archive/2010/06/25/1765510.html
2.1）CentOS 下安装 graphviz

$ sudo yum install graphviz

Install 39 Package(s)

总下载量：13 M
Installed size: 35 M
确定吗？[y/N]：y

已安装:
graphviz.i686 0:2.26.0-10.el6

完毕！

3、在Linux下如何使用

　　假设我们把上面的代码写到了一个叫做aa.gv的文本文件里面，那么我们执行如下命令就可以了：

　　$ dot -Tpng -ohehe.png aa.gv

　　这样就会在当前目录下生成一个叫做hehe.png的图片文件
　　

mac

wget http://pecl.php.net/get/vld-0.14.0.tgz
tar zxvf vld-0.14.0.tgz
cd vld-0.14.0
phpize
locate php-config
./configure --with-php-config=/usr/local/Cellar/php70/7.0.25_17/bin/php-config
 make && make install
 php -r 'phpinfo();' |grep vld
 ls /usr/local/Cellar/php70/7.0.25_17/lib/php/extensions/no-debug-non-zts-20151012/vld.so
 vi  /usr/local/etc/php/7.0/php.ini
 php -r 'phpinfo();' |grep vld
 php -dvld.active=1 -dvld.execute=0 test.php
 php -dvld.active=1 -dvld.save_paths=1 test.php
 ls /tmp
 brew install graphviz
 dot -Tpng /tmp/paths.dot -o paths.png
<img src="{{site.url}}{{site.baseurl}}/img/paths.png"/>
