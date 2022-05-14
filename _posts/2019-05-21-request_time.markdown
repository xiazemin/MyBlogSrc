---
title: nginx request_time
layout: post
category: php
author: 夏泽民
---
1、request_time
官网描述：request processing time in seconds with a milliseconds resolution; time elapsed between the first bytes were read from the client and the log write after the last bytes were sent to the client 。
指的就是从接受用户请求的第一个字节到发送完响应数据的时间，即包括接收请求数据时间、程序响应时间、输出响应数据时间。
 
2、upstream_response_time
官网描述：keeps times of responses obtained from upstream servers; times are kept in seconds with a milliseconds resolution. Several response times are separated by commas and colons like addresses in the $upstream_addr variable
 
是指从Nginx向后端（php-cgi)建立连接开始到接受完数据然后关闭连接为止的时间。
 
从上面的描述可以看出，$request_time肯定大于等于$upstream_response_time，特别是使用POST方式传递参数时，因为Nginx会把request body缓存住，接受完毕后才会把数据一起发给后端。所以如果用户网络较差，或者传递数据较大时，$request_time会比$upstream_response_time大很多。
 
所以如果使用nginx的accesslog查看php程序中哪些接口比较慢的话，记得在log_format中加入$upstream_response_time。
 
根据引贴对官网描述的翻译：
upstream_response_time:从 Nginx 建立连接 到 接收完数据并关闭连接
request_time:从 接受用户请求的第一个字节 到 发送完响应数据

如果把整个过程补充起来的话 应该是：
［1用户请求］［2建立 Nginx 连接］［3发送响应］［4接收响应］［5关闭  Nginx 连接］
那么 upstream_response_time 就是 2+3+4+5 
但是 一般这里面可以认为 ［5关闭 Nginx 连接］ 的耗时接近 0
所以 upstream_response_time 实际上就是 2+3+4 
而 request_time 是 1+2+3+4
二者之间相差的就是 ［1用户请求］ 的时间

如果用户端网络状况较差 或者传递数据本身较大 
再考虑到 当使用 POST 方式传参时 Nginx 会先把 request body 缓存起来
而这些耗时都会累积到 ［1用户请求］ 头上去
这样就解释了
为什么 request_time 有可能会比 upstream_response_time 要大

因为用户端的状况通常千差万别 无法控制 
所以并不应该被纳入到测试和调优的范畴里面
更值得关注的应该是 upstream_response_time
所以在实际工作中 如果想要关心哪些请求比较慢的话 
记得要在配置文件的 log_format 中加入 $upstream_response_time 
<!-- more -->
