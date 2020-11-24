---
title: openresty
layout: post
category: php
author: 夏泽民
---
curl -O https://openresty.org/download/openresty-1.17.8.2.tar.gz

tar -xzvf openresty-VERSION.tar.gz

curl -O https://www.openssl.org/source/openssl-1.1.1d.tar.gz

tar -zxvf openssl-1.1.1d.tar.gz

./configure --prefix=/opt/openresty \
            --with-luajit --with-openssl=/Users/didi/c/openssl-1.1.1d
            
http://openresty.org/cn/installation.html

make -j4
sudo make install
<!-- more -->
{% raw %}
~ % /opt/openresty/bin/openresty -h
nginx version: openresty/1.17.8.2
Usage: nginx [-?hvVtTq] [-s signal] [-c filename] [-p prefix] [-g directives]

Options:
  -?,-h         : this help
  -v            : show version and exit
  -V            : show version and configure options then exit
  -t            : test configuration and exit
  -T            : test configuration, dump it and exit
  -q            : suppress non-error messages during configuration testing
  -s signal     : send signal to a master process: stop, quit, reopen, reload
  -p prefix     : set prefix path (default: /opt/openresty/nginx/)
  -c filename   : set configuration file (default: conf/nginx.conf)
  -g directives : set global directives out of configuration file
{% endraw %}

ln -s /opt/openresty/bin/openresty /usr/local/bin/openresty

 ~ % openresty
nginx: [alert] could not open error log file: open() "/opt/openresty/nginx/logs/error.log" failed (13: Permission denied)
2020/11/17 14:07:45 [emerg] 45398#0: mkdir() "/opt/openresty/nginx/client_body_temp" failed (13: Permission denied)

mkdir logs/ conf/
vi conf/nginx.conf
openresty  -p `pwd`/ -c conf/nginx.conf


worker_processes  1;
error_log logs/error.log;
events {
    worker_connections 1024;
}
http {
    server {
        listen 8080;
        location / {
            default_type text/html;
            content_by_lua_block {
                ngx.say("<p>hello, world</p>")
            }
        }
    }
}

% ps aux |grep nginx
didi             46629   0.0  0.0  4279564    732 s000  R+    2:25下午   0:00.00 grep nginx
didi             46619   0.0  0.0  4311488   1056   ??  S     2:25下午   0:00.00 nginx: worker process
didi             46618   0.0  0.0  4311056    520   ??  Ss    2:25下午   0:00.00 nginx: master process openresty -p /Users/didi/www/ -c conf/nginx.conf

% curl http://localhost:8080/
<p>hello, world</p>

http://openresty.org/cn/components.html

http://openresty.org/cn/dynamic-routing-based-on-redis.html

https://openresty.org/download/agentzh-nginx-tutorials-en.pdf

https://github.com/openresty/nginx-tutorials

https://wdicc.com/intro-openresty/

http://chenxiaoyu.org/2011/10/30/nginx-modules/
http://wendal.net/338.html


Cannot find autoconf. Please check your autoconf installation and the
$PHP_AUTOCONF environment variable. Then, rerun this script.

是因为缺少了依赖包，安装即可解决
yum install autoconf

https://github.com/openresty/nginx-tutorials/blob/master/zh-cn/00-Foreword01.tut

http://openresty.org/cn/dynamic-routing-based-on-redis.html