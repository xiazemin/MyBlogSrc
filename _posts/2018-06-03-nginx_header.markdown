---
title: nginx_header
layout: post
category: php 
author: 夏泽民
---
<!-- more -->
ngx_http_headers_module模块
一. 前言
ngx_http_headers_module模块提供了两个重要的指令add_header和expires，来添加 “Expires” 和 “Cache-Control” 头字段，对响应头添加任何域字段。add_header可以用来标示请求访问到哪台服务器上，这个也可以通过nginx模块nginx-http-footer-filter研究使用来实现。expires指令用来对浏览器本地缓存的控制。
二. add_header指令
语法: add_header name value;
默认值: —
配置段: http, server, location, if in location
对响应代码为200，201，204，206，301，302，303，304，或307的响应报文头字段添加任意域。如：

add_header From jb51.net
三. expires指令
语法: expires [modified] time;
expires epoch | max | off;
默认值: expires off;
配置段: http, server, location, if in location
在对响应代码为200，201，204，206，301，302，303，304，或307头部中是否开启对“Expires”和“Cache-Control”的增加和修改操作。
可以指定一个正或负的时间值，Expires头中的时间根据目前时间和指令中指定的时间的和来获得。
epoch表示自1970年一月一日00:00:01 GMT的绝对时间，max指定Expires的值为2037年12月31日23:59:59，Cache-Control的值为10 years。
Cache-Control头的内容随预设的时间标识指定：
·设置为负数的时间值:Cache-Control: no-cache。
·设置为正数或0的时间值：Cache-Control: max-age = #，这里#的单位为秒，在指令中指定。
参数off禁止修改应答头中的"Expires"和"Cache-Control"。
实例一：对图片，flash文件在浏览器本地缓存30天
location ~ .*\.(gif|jpg|jpeg|png|bmp|swf)$
 {
      expires 30d;
 }
实例二：对js，css文件在浏览器本地缓存1小时
location ~ .*\.(js|css)$
 {
      expires 1h;
 }
ngx_headers_more模块
一. 介绍ngx_headers_more
ngx_headers_more 用于添加、设置和清除输入和输出的头信息。nginx源码没有包含该模块，需要另行添加。
该模块是ngx_http_headers_module模块的增强版，提供了更多的实用工具，比如复位或清除内置头信息，如Content-Type, Content-Length, 和Server。
可以允许你使用-s选项指定HTTP状态码，使用-t选项指定内容类型，通过more_set_headers 和 more_clear_headers 指令来修改输出头信息。如：
more_set_headers -s 404 -t 'text/html' 'X-Foo: Bar';
输入头信息也可以这么修改，如：
location /foo {
  more_set_input_headers 'Host: foo' 'User-Agent: faked';
  # now $host, $http_host, $user_agent, and
  #  $http_user_agent all have their new values.
}
-t选项也可以在more_set_input_headers和more_clear_input_headers指令中使用。
不像标准头模块，该模块的指示适用于所有的状态码，包括4xx和5xx的。 add_header只适用于200，201，204，206，301，302，303，304，或307。

二. 安装ngx_headers_more
wget 'http://nginx.org/download/nginx-1.5.8.tar.gz'
tar -xzvf nginx-1.5.8.tar.gz
cd nginx-1.5.8/
  
# Here we assume you would install you nginx under /opt/nginx/.
./configure --prefix=/opt/nginx \
  --add-module=/path/to/headers-more-nginx-module
make
make install
ngx_headers_more 包下载地址：http://github.com/agentzh/headers-more-nginx-module/tags
ngx_openresty包含该模块。
三. 指令说明
more_set_headers
语法：more_set_headers [-t <content-type list>]... [-s <status-code list>]... <new-header>...
默认值：no
配置段：http, server, location, location if
阶段：输出报头过滤器
替换（如有）或增加（如果不是所有）指定的输出头时响应状态代码与-s选项相匹配和响应的内容类型的-t选项指定的类型相匹配的。
如果没有指定-s或-t，或有一个空表值，无需匹配。因此，对于下面的指定，任何状态码和任何内容类型都讲设置。
more_set_headers  "Server: my_server";
具有相同名称的响应头总是覆盖。如果要添加头，可以使用标准的add_header指令代替。
单个指令可以设置/添加多个输出头。如：
more_set_headers 'Foo: bar' 'Baz: bah';
在单一指令中，选项可以多次出现，如：
more_set_headers -s 404 -s '500 503' 'Foo: bar';
等同于：
more_set_headers -s '404 500 503' 'Foo: bar';
新的头是下面形式之一：

Name: Value
Name:
Name
最后两个有效清除的头名称的值。Nginx的变量允许是头值，如：
set $my_var "dog";
more_set_headers "Server: $my_var";
注意：more_set_headers允许在location的if块中，但不允许在server的if块中。下面的配置就报语法错误：
# This is NOT allowed!
 server {
    if ($args ~ 'download') {
      more_set_headers 'Foo: Bar';
    }
    ...
  }
more_clear_headers
语法：more_clear_headers [-t <content-type list>]... [-s <status-code list>]... <new-header>...
默认值：no
配置段：http, server, location, location if
阶段：输出报头过滤器
清除指定的输出头。
more_clear_headers -s 404 -t 'text/plain' Foo Baz;
等同于
more_set_headers -s 404 -t 'text/plain' "Foo: " "Baz: ";
或
more_clear_headers -s 404 -t 'text/plain' Foo Baz;
等同于
more_set_headers -s 404 -t 'text/plain' "Foo: " "Baz: ";
或
more_set_headers -s 404 -t 'text/plain' Foo Baz
也可以使用通配符*，如：
more_clear_headers 'X-Hidden-*';
清除开始由“X-Hidden-”任何输出头。
more_set_input_headers
语法：more_set_input_headers [-r] [-t <content-type list>]... <new-header>...
默认值：no
配置段：http, server, location, location if
阶段： rewrite tail
非常类似more_set_headers，不同的是它工作在输入头（或请求头），它仅支持-t选项。
注意：使用-t选项的是过滤请求头的Content-Type，而不是响应头的。
more_clear_input_headers
语法：more_clear_input_headers [-t <content-type list>]... <new-header>...
默认值：no
配置段：http, server, location, location if
阶段： rewrite tail
清除指定输入头。如：
more_clear_input_headers -s 404 -t 'text/plain' Foo Baz;
等同于
more_set_input_headers -s 404 -t 'text/plain' "Foo: " "Baz: ";
或
more_clear_input_headers -s 404 -t 'text/plain' Foo Baz;
等同于
more_set_input_headers -s 404 -t 'text/plain' "Foo: " "Baz: ";
或
more_set_input_headers -s 404 -t 'text/plain' Foo Baz
四. ngx_headers_more局限性
1. 不同于标准头模块，该模块不会对下面头有效： Expires, Cache-Control, 和Last-Modified。
2. 使用此模块无法删除Connection的响应报头。唯一方法是更改src/ HTTP/ ngx_http_header_filter_module.c文件。
五. 使用ngx_headers_more
# set the Server output header
more_set_headers 'Server: my-server';
  
# set and clear output headers
location /bar {
  more_set_headers 'X-MyHeader: blah' 'X-MyHeader2: foo';
  more_set_headers -t 'text/plain text/css' 'Content-Type: text/foo';
  more_set_headers -s '400 404 500 503' -s 413 'Foo: Bar';
  more_clear_headers 'Content-Type';
  
  # your proxy_pass/memcached_pass/or any other config goes here...
}
  
# set output headers
location /type {
  more_set_headers 'Content-Type: text/plain';
  # ...
}
  
# set input headers
location /foo {
  set $my_host 'my dog';
  more_set_input_headers 'Host: $my_host';
  more_set_input_headers -t 'text/plain' 'X-Foo: bah';
  
  # now $host and $http_host have their new values...
  # ...
}
  
 # replace input header X-Foo *only* if it already exists
more_set_input_headers -r 'X-Foo: howdy';
六. 应用ngx_headers_more
修改web服务器是什么软件，什么版本，同时隐藏Centent-Type、Accept-Range、Content-Length头信息。