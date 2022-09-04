---
title: Valid_referers nginx防止盗链
layout: post
category: nginx
author: 夏泽民
---
Valid_referers  设置信任网站

None      浏览器中referer（Referer是header的一部分，当浏览器向web服务器发送请求的时候，一般会带上Referer，告诉服务器我是从哪个页面连接过来的，服务器基此可以获得一些信息用于处理）为空的情况，就直接在浏览器访问图片
Blocked  referrer不为空的情况，但是值被代理或防火墙删除了，这些值不以http://或者https:// 开头。
《1》：vim /usr/local/nginx/conf/nginx.conf  编辑文件

《2》：写入（注：写入的内容要放在缓存的上面。）

        location ~* \.(wma|wmv|asf|mp3|mmf|zip|rar|jpg|gif|png|swf|flv)$ {

          valid_referers none blocked *.source.com source.com;

          if ($invalid_referer) {

             rewrite ^/ http://www.source.com/error.html;

            }

        }
<!-- more -->
释：

第一行：wma|wmv|asf|mp3|mmf|zip|rar|jpg|gif|png|swf|flv 表示对这些后缀的文件实行防盗链

第二行：none blocked *.source.com source.com;           不区分大小写

表示referers信息中匹配none blocked *.source.com source.com（*代表任何，任何的二级域名）

If{ }里面内容的意思是，如果连接不是来自第二行指定的就强制跳转到403错误页面，当然直接返>回404也是可以的，也可以是图片。
https://blog.csdn.net/m0_54434140/article/details/122489818

1、Nginx Referer模块

nginx模块ngx_http_referer_module通常用于阻挡来源非法的域名请求。当一个请求头的Referer字段中包含一些非正确的字段，这个模块可以禁止这个请求访问站点。构造Referer的请求很容易实现，所以使用这个模块并不能100%的阻止这些请求。

2、valid_referers 指令

语法: valid_referers none | blocked | server_names | string … ;

配置段: server, location

指定合法的来源'referer', 他决定了内置变量$invalid_referer的值，如果referer头部包含在这个合法网址里面，这个变量被设置为0，否则设置为1. 需要注意的是：这里并不区分大小写的.

参数说明：

none：请求头缺少Referer字段，即空Referer
blocked：请求头Referer字段不为空（即存在Referer），但是值被代理或者防火墙删除了，这些值不以“http://”或“https://”开头，通俗点说就是允许“http://”或"https//"以外的请求。
server_names：Referer请求头白名单。
arbitrary string：任意字符串，定义服务器名称或可选的URI前缀，主机名可以使用*号开头或结尾，Referer字段中的服务器端口将被忽略掉。
regular expression：正则表达式，以“~”开头，在“http://”或"https://"之后的文本匹配。

例子:

server { 
    listen   80     ; 
    server_name  img.abc.com; 
    # 必须使用域名访问 
 
    if ($host != 'img.abc.com') { 
            return 403 ; 
    } 
 
    # 拦截非法referer ,none  和 blocked 的区别看以上说明
    valid_referers blocked www.abc.com img.abc.com ; 
    if ($invalid_referer) { 
        return 403 ; 
        #rewrite ^.*$ http://www.baidu.com/403.jpg; 
    } 
    charset utf-8; 
    location / { 
        ………… 
    ｝ 
｝
上面配置合法的Referer为 www.abc.com / img.abc.com 和 无Referer(浏览器直接访问，就没有Referer) ; 其他非法Referer请求过来时， $invalid_referer 值为1 ， 就return 403 , 或者重定向到一个403图片。

http://events.jianshu.io/p/1f3ccaf93e7c