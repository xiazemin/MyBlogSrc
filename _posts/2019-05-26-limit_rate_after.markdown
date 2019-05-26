---
title: limit_rate_after nginx限速配置
layout: post
category: web
author: 夏泽民
---
Nginx的http核心模块ngx_http_core_module中提供limit_rate这个指令可以用于控制速度，limit_rate_after用于设置http请求传输多少字节后开始限速。
另外两个模块ngx_http_limit_conn_module和ngx_http_limit_req_module分别用于连接数和连接频率的控制。


限制向客户端传送响应数据的速度，可以用来限制客户端的下载速度。参数rate的单位是字节/秒，0为关闭限速。

nginx按连接限速，所以如果某个客户端同时开启了两个连接，那么客户端的整体速度是这条指令设置值的2倍。

nginx限速示例：

location /flv/ {
flv;
limit_rate_after 500k;     #当传输量大于此值时，超出部分将限速传送
limit_rate 50k;
}

limit_rate_after size;

默认值: limit_rate_after 0;
上下文: http, server, location, if in location

这个指令出现在版本 0.8.0。当传输量大于此值时，超出部分将限速传送，小于设置值时不限速。
<!-- more -->
可以利用$limit_rate变量设置流量限制。如果想在特定条件下限制响应传输速率，可以使用这个功能：

server {

if ($slow) {
set $limit_rate 4k;
}

…
}
此外，也可以通过“X-Accel-Limit-Rate”响应头来完成速率限制。 这种机制可以用proxy_ignore_headers指令和 fastcgi_ignore_headers指令关闭。

