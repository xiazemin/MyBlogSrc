---
title: proxy_pass
layout: post
category: php
author: 夏泽民
---
FastCGI 是一个协议，它是应用程序和 WEB 服务器连接的桥梁。Nginx 并不能直接与 PHP-FPM 通信，而是将请求通过 FastCGI 交给 PHP-FPM 处理。

location ~.php$ {
fastcgi_pass 127.0.0.1:9000;
}

Nginx 反向代理最重要的指令是 proxy_pass
location ^~/seckill_query/{
proxy_pass http://ris.filemail.gdrive:8090/;
}

Nginx 负载均衡

介绍一下 upstream 模块：

负载均衡模块用于从”upstream”指令定义的后端主机列表中选取一台主机。nginx先使用负载均衡模块找到一台主机，再使用upstream模块实现与这台主机的交互。
upstream php-upstream {
ip_hash;
server 192.168.0.1;
server 192.168.0.2;
}

建立upstream
这个是fastcgi的例子，如果是http的则把端口改下就可以了

upstream fastcgi_backend {
    server 127.0.0.1:9000;
    keepalive 60;
}
<!-- more -->
http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_pass
http://nginx.org/en/docs/http/ngx_http_fastcgi_module.html#fastcgi_pass

cgi-fcgi命令实现了 fastcgi 客户端协议，可以直接访问php-fpm实现HTTP请求。
https://github.com/FastCGI-Archives

https://fastcgi-archives.github.io/

php代码

<?php
print_r($_GET);
命令

env USER=dev FCGI_ROLE=RESPONDER SCRIPT_FILENAME=/path/php/test/get.php QUERY_STRING="arg1=1" REQUEST_METHOD=GET SCRIPT_NAME=/get.php REQUEST_URI=/get.php DOCUMENT_URI=/get.php DOCUMENT_ROOT=/path/php/test SERVER_PROTOCOL=HTTP/1.1 GATEWAY_INTERFACE=CGI/1.1 SERVER_SOFTWARE=nginx/1.8.1 REMOTE_ADDR=127.0.0.1 REMOTE_PORT=50815 SERVER_ADDR=127.0.0.1 SERVER_PORT=80 SERVER_NAME=localhost REDIRECT_STATUS=200 HTTP_HOST=localhost HTTP_CONNECTION=keep-alive HTTP_CACHE_CONTROL=max-age=0 PHP_SELF=/index.php REQUEST_TIME_FLOAT=1499780545.7094 REQUEST_TIME=1499780545 cgi-fcgi -bind -connect 127.0.0.1:9000

输出

X-Powered-By: PHP/5.6.30
Content-type: text/html; charset=UTF-8

Array
(
    [arg1] => 1
)
命令分为两部分，第一部分设置环境变量，也就是env命令，cgi-fcgi后面才是连接php-fpm的参数

关于环境变量的配置
env命令设置的环境变量对应的php $_SERVER中的值，如果搞不清楚这里的环境变量怎么设置，可以将普通HTTP请求产生的$_SERVER变量，然后稍微修改使用，或者参考nginx fastcgi_params配置文件。

实现蛋疼的POST的请求
实现POST请求，我查遍了cgi-fcgi的文档、google、stackoverflow等都没有现成的实现方法，最后我根据fcgi的c/c++的实现，还有文档中post数据来源于stdin的描述，搞清楚了怎么用cgi-fcgi实现post请求

php代码

<php?
print_r($_POST);
命令

env USER=dev FCGI_ROLE=RESPONDER SCRIPT_FILENAME=/path/php/test/post.php QUERY_STRING="" REQUEST_METHOD=POST CONTENT_TYPE="application/x-www-form-urlencoded" CONTENT_LENGTH=5 SCRIPT_NAME=/post.php REQUEST_URI=/post.php DOCUMENT_URI=/post.php DOCUMENT_ROOT=/path/php/test SERVER_PROTOCOL=HTTP/1.1 GATEWAY_INTERFACE=CGI/1.1 SERVER_SOFTWARE=nginx/1.8.1 REMOTE_ADDR=127.0.0.1 REMOTE_PORT=50815 SERVER_ADDR=127.0.0.1 SERVER_PORT=80 SERVER_NAME=localhost REDIRECT_STATUS=200 HTTP_HOST=localhost HTTP_CONNECTION=keep-alive HTTP_CACHE_CONTROL=max-age=0 PHP_SELF=/post.php REQUEST_TIME_FLOAT=1456138229.7094 REQUEST_TIME=1456138229 cgi-fcgi -bind -connect 127.0.0.1:9000 <<< "arg=1"
输出

X-Powered-By: PHP/5.6.30
Content-type: text/html; charset=UTF-8

Array
(
    [arg] => 1
)
post请求要注意CONTENT_LENGTH环境变量要和post的数据长度保持一致。

这样就可以访问任意php服务器了，可以向它们提交配置文件到local cache、刷新部分php文件的opcache等特殊操作了，当然它的玩法还不仅限于此。

