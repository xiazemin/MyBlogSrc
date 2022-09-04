---
title: add_header允许跨域
layout: post
category: nginx
author: 夏泽民
---
1.跨域是浏览器的问题

明确一点跨域是浏览器问题，浏览器出现跨域问题的时候，会发送两次请求，第一次请求的方法是‘Option’，第二次请求是正常的Get或者post。

2.跨域解决方法

利用nginx解决跨域，新建一个server，server中内容如下，其中proxy_pass是代理的ip地址以及端口。也就是我们最终请求的地址。也就是在第一次请求的时候，由nginx直接返回204状态，不经过最终的ip地址。（注意修改proxy_pass中的地址）
<!-- more -->
server {
        listen       8099;
        server_name  localhost;
 
        #charset koi8-r;
 
        #access_log  logs/host.access.log  main;
 
        location  / {
			root   html;
			index  index.html index.htm;
			# 配置html以文件方式打开
			# 配置html以文件方式打开
			if ($request_method = 'POST') {
				add_header 'Access-Control-Allow-Origin' '*' always;
				add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS,PUT,DELETE,OPTION';
				add_header 'Access-Control-Allow-Credentials' 'true';
				add_header 'Access-Control-Allow-Headers' 'DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,Accept,Referer,Accept-Encoding,Accept-Language,Access-Control-Request-Headers,Access-Control-Request-Method,Connection,Host,Origin,Sec-Fetch-Mode';
			}
			if ($request_method = 'GET') {
				add_header 'Access-Control-Allow-Origin' '*' always;
				add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS,PUT,DELETE,OPTION';
				add_header 'Access-Control-Allow-Credentials' 'true';
				add_header 'Access-Control-Allow-Headers' 'DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,Accept,Referer,Accept-Encoding,Accept-Language,Access-Control-Request-Headers,Access-Control-Request-Method,Connection,Host,Origin,Sec-Fetch-Mode';
			}
			if ($request_method = 'OPTIONS') {
				add_header 'Access-Control-Allow-Origin' '*' always;
				add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS,PUT,DELETE,OPTION';
				add_header 'Access-Control-Allow-Credentials' 'true';
				add_header 'Access-Control-Allow-Headers' 'DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,Accept,Referer,Accept-Encoding,Accept-Language,Access-Control-Request-Headers,Access-Control-Request-Method,Connection,Host,Origin,Sec-Fetch-Mode';
				return 204;
			}
			# 代理到ip地址端口
			proxy_pass       http://xxxx:xxx;
 
		}
 
        #error_page  404              /404.html;
 
        # redirect server error pages to the static page /50x.html
        #
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }
 
        # proxy the PHP scripts to Apache listening on 127.0.0.1:80
        #
        #location ~ \.php$ {
        #    proxy_pass   http://127.0.0.1;
        #}
 
        # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000
        #
        #location ~ \.php$ {
        #    root           html;
        #    fastcgi_pass   127.0.0.1:9000;
        #    fastcgi_index  index.php;
        #    fastcgi_param  SCRIPT_FILENAME  /scripts$fastcgi_script_name;
        #    include        fastcgi_params;
        #}
 
        # deny access to .htaccess files, if Apache's document root
        # concurs with nginx's one
        #
        #location ~ /\.ht {
        #    deny  all;
        #}
    }

3.问题汇总

(1)出现如下问题。

The 'Access-control-allow-origin' header contains multiple values '',but only one is allow'
是出现了两次的问题，修改上述server中内容。注释get和post中access-control-allow-origin。访问成功。（注意修改proxy_pass）
https://blog.csdn.net/weixin_41415235/article/details/122940569
