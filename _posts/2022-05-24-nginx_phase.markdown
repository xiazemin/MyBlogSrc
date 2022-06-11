---
title: nginx_phase 处理请求的11个阶段
layout: post
category: nginx
author: 夏泽民
---
post-read
接受到完整的http头部后，读取请求内容阶段，nginx读取并解析完请求头之后就立即开始执行；

server-rewrite
在uri与location匹配之前修改请求的URI（重定向），在server块中的请求地址重写阶段；

find-config
配置查找阶段，根据请求uri匹配location表达式，这个阶段不支持nginx模块注册处理程序，而是由ngx_http_core_module模块来完成当前请求与location配置快之间的配对工作；

rewrite
location块中的请求地址重写阶段，当rewrite指令用于location中，即运行。另外，ngx_lua模块中的set_by_lua指令和rewrite_by_lua指令也在此阶段；

post-rewrite
请求地址重写提交阶段，防止递归修改uri造成死循环，（一个请求执行10次就会被nginx认定为死循环）该阶段只能由ngx_http_core_module模块实现

preaccess
访问权限检查准备阶段，http模块介入处理阶段，标准模块ngx_limit_req和ngx_limit_zone就运行在此阶段，前置可以控制访问的频率，后者限制访问的并发度

access
访问权限检查阶段，标准模块ngx_access,第三方模块nginx_auth_request以及第三方模块ngx_lua的access_by_lua 指令运行在此阶段，配置指令多是执行访问控制性质的任务，比如检查用户的访问权限，检查用户的来源IP地址是否合法；

post-access
访问权限检查提交阶段；如果请求不被允许访问nginx服务器，该阶段负责向用户返回错误响应；

try-files
配置项try_files处理阶段

如果http请求访问静态文件资源，try_files配置项可以使这个请求顺序地访问多个静态文件资源，直到某个静态文件资源符合选取条件；

content
内容产生阶段，大部分HTTP模块会介入该阶段，是所有请求处理阶段中最重要的阶段，因为这个阶段的指令通常是用来生成HTTP响应内容的；

log
日志模块处理阶段，记录日志；

以上阶段中，有些阶段是必备的，有些阶段是可选的，各个阶段可以允许多个HTTP模块同时介入，nginx会按照各个HTTP模块的ctx_index顺序执行这些模块的hadler方法。
<!-- more -->
但是ngx_http_find_config_phase,nginx_http_post_rewrite_phase,nginx_http_post_access_phase,ngx_http_try_files_phase这四个阶段是不允许HTTP模块加入自己的ngx_http_handler_py方法处理用户请求的，他们仅由HTTP框架自身实现
https://blog.csdn.net/Carmelo_a/article/details/118419198
