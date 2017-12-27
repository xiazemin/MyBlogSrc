---
title: Django_nginx_uwsgi
layout: post
category: web
author: 夏泽民
---
<!-- more -->
安装Python包管理
easy_install 包 https://pypi.python.org/pypi/distribute
wget https://pypi.python.org/packages/source/d/distribute/distribute-0.6.49.tar.gz
tar xf distribute-0.6.49.tar.gz
cd distribute-0.6.49
python2.7 setup.py install
easy_install --version
安装 uwsgi
uwsgi:https://pypi.python.org/pypi/uWSGI

uwsgi 参数详解：http://uwsgi-docs.readthedocs.org/en/latest/Options.html

pip install uwsgi
uwsgi --version    # 查看 uwsgi 版本
测试 uwsgi 是否正常：

新建 test.py 文件，内容如下：

def application(env, start_response):
    start_response('200 OK', [('Content-Type','text/html')])
    return "Hello World"
然后在终端运行：

uwsgi --http :8001 --wsgi-file test.py

在浏览器内输入：http://127.0.0.1:8001，查看是否有"Hello World"输出，若没有输出，请检查你的安装过程。

安装 Django
pip install django
测试 django 是否正常，运行：

django-admin.py startproject demosite
cd demosite
python2.7 manage.py runserver 0.0.0.0:8002
在浏览器内输入：http://127.0.0.1:8002，检查django是否运行正常。

安装 Nginx
安装命令如下：

cd ~
wget http://nginx.org/download/nginx-1.5.6.tar.gz
tar xf nginx-1.5.6.tar.gz
cd nginx-1.5.6
./configure --prefix=/usr/local/nginx-1.5.6 \
--with-http_stub_status_module \
--with-http_gzip_static_module
make && make install
你可以阅读 Nginx 安装配置 了解更多内容。

uwsgi 配置
uwsgi支持ini、xml等多种配置方式，本文以 ini 为例， 在/ect/目录下新建uwsgi9090.ini，添加如下配置：

[uwsgi]
socket = 127.0.0.1:9090
master = true         //主进程
vhost = true          //多站模式
no-site = true        //多站模式时不设置入口模块和文件
workers = 2           //子进程数
reload-mercy = 10     
vacuum = true         //退出、重启时清理文件
max-requests = 1000   
limit-as = 512
buffer-size = 30000
pidfile = /var/run/uwsgi9090.pid    //pid文件，用于下面的脚本启动、停止该进程
daemonize = /website/uwsgi9090.log
Nginx 配置
找到nginx的安装目录（如：/usr/local/nginx/），打开conf/nginx.conf文件，修改server配置：

server {
        listen       80;
        server_name  localhost;
        
        location / {            
            include  uwsgi_params;
            uwsgi_pass  127.0.0.1:9090;              //必须和uwsgi中的设置一致
            uwsgi_param UWSGI_SCRIPT demosite.wsgi;  //入口文件，即wsgi.py相对于项目根目录的位置，“.”相当于一层目录
            uwsgi_param UWSGI_CHDIR /demosite;       //项目根目录
            index  index.html index.htm;
            client_max_body_size 35m;
        }
    }
    
WSGI是Web Server Gateway Interface的缩写。以层的角度来看，WSGI所在层的位置低于CGI。但与CGI不同的是WSGI具有很强的伸缩性且能运行于多线程或多进程的环境下，这是因为WSGI只是一份标准并没有定义如何去实现。实际上WSGI并非CGI，因为其位于web应用程序与web服务器之间，而web服务器可以是CGI，mod_python（注：现通常使用mod_wsgi代替），FastCGI或者是一个定义了WSGI标准的web服务器就像python标准库提供的独立WSGI服务器称为wsgiref。

Python Paste - WSGI底层工具集. 包括多线程, SSL和 基于Cookies, sessions等的验证(authentication)库. 可以用Paste方便地搭建自己的Web框架。
WSGI:Python Web Server Gateway Interface v1.0
它是 PEP3333中定义的（PEP3333的目标建立一个简单的普遍适用的服务器与Web框架之间的接口）
WSGI是Python应用程序或框架和Web服务器之间的一种接口