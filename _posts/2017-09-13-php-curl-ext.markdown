---
title: php-curl-ext
layout: post
category: php
author: 夏泽民
---
php curl  扩展安装
<!-- more -->

方法一

安装cURL

wget http://curl.haxx.se/download/curl-7.17.1.tar.gz

 tar -zxf curl-7.17.1.tar.gz

./configure --prefix=/usr/local/curl

make & make install

安装php

   只要打开开关 --with-curl=/usr/local/curl

   就可以了。

   这个扩展库还是非常棒，是fsockopen等等相关的有效的替代品。

方法二

进入安装原php的源码目录，

cd ext

cd curl

phpize

./configure --with-curl=DIR

make & make install

就会在PHPDIR/ext/curl/moudles/下生成curl.so的文件。

复制curl.so文件到extensions的配置目录，修改php.ini就好了

extension=curl.so

第一种方法试了N遍一直在失败中，于是放弃。

使用第二种方法安装，

phpize提示找不到，其实命令在/usr/local/php/bin/目标下:
{% highlight bash %}
$/usr/local/php/bin/phpize
./configure --with-curl=DIR需要指定php的配置路径，应该如下：
$./configure --with-php-config=/usr/local/php/bin/php-config --with-curl=DIR
{% endhighlight %}

注：上面的资料中错把--with-php-config写成了--with-php-php-config

然后就是编译安装：
{% highlight bash %}
$ make
$ make install
{% endhighlight %}
到这里会提示生成文件curl.so的路径： /usr/local/php/lib/php/extensions/no-debug-non-zts-20060613/

进入到这个路径下，复制curl到extension_dir目录下(本目录路径可以看phpinfo可是直接看php.int)，

修改php.ini
{% highlight bash %}
extension=curl.so
$ /usr/local/php/bin/php -m
{% endhighlight %}
如果看到有curl项表示成功。

{% highlight bash linenos %}
git clone https://github.com/xiazemin/php-src.git
cd php-src/ext/curl/
 /usr/local/bin/phpize
./configure   --with-curl= /usr/local/Cellar/php70/7.0.8/
configure: error: invalid value of canonical build
 ./configure   --with-curl=/usr/local/etc/php/7.0
 make & make install
 cp modules/curl.so  /usr/local/Cellar/php70/7.0.8/lib/php/extensions/no-debug-non-zts-20151012/curl.so
 vi /usr/local/etc/php/7.0/php.ini
 Warning: Module 'curl' already loaded in Unknown on line 0
{% endhighlight %}
ibtool: link: `xhprof.lo’ is not a valid libtool object
解决方法
用命令
make clean
然后在重新执行命令

