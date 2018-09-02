---
title: haproxy_exp
layout: post
category: web
author: 夏泽民
---
<!-- more -->
定义独立日志文件

[root@node1 haproxy]# vim /etc/rsyslog.conf #为其添加日志功能
# Provides UDP syslog reception
$ModLoad imudp
$UDPServerRun 514 ------>启动udp，启动端口后将作为服务器工作
# Provides TCP syslog reception
$ModLoad imtcp
$InputTCPServerRun 514 ------>启动tcp监听端口
local2.* /var/log/haproxy.log

[root@node1 haproxy]# service rsyslog restar
[root@LB haproxy]# vim haproxy.cfg
log 127.0.0.1 local2 --------->在global端中添加此行

[root@node1 haproxy]# vim /etc/rsyslog.conf #为其添加日志功能
# Provides UDP syslog reception
$ModLoad imudp
$UDPServerRun 514 ------>启动udp，启动端口后将作为服务器工作
# Provides TCP syslog reception
$ModLoad imtcp
$InputTCPServerRun 514 ------>启动tcp监听端口
local2.* /var/log/haproxy.log
 
[root@node1 haproxy]# service rsyslog restar
[root@LB haproxy]# vim haproxy.cfg
log 127.0.0.1 local2 --------->在global端中添加此行
一个最简单的http服务的配置

global
log 127.0.0.1 local2
chroot /var/lib/haproxy
pidfile /var/run/haproxy.pid
maxconn 4000
user haproxy
group haproxy
daemon
stats socket /var/lib/haproxy/stats
defaults
mode http
log global
option httplog
option dontlognull
option http-server-close
option forwardfor except 127.0.0.0/8
option redispatch
retries 3
timeout http-request 10s
timeout queue 1m
timeout connect 10s
timeout client 1m
timeout server 1m
timeout http-keep-alive 10s
timeout check 10s
maxconn 3000
frontend webser #webser为名称
option forwardfor
bind *:80
default_backend app
backend app
balance roundrobin #使拥roundrobin 算法
server app1 192.168.1.111:80 check
server app2 192.168.1.112:80 check

global
log 127.0.0.1 local2
chroot /var/lib/haproxy
pidfile /var/run/haproxy.pid
maxconn 4000
user haproxy
group haproxy
daemon
stats socket /var/lib/haproxy/stats
defaults
mode http
log global
option httplog
option dontlognull
option http-server-close
option forwardfor except 127.0.0.0/8
option redispatch
retries 3
timeout http-request 10s
timeout queue 1m
timeout connect 10s
timeout client 1m
timeout server 1m
timeout http-keep-alive 10s
timeout check 10s
maxconn 3000
frontend webser #webser为名称
option forwardfor
bind *:80
default_backend app
backend app
balance roundrobin #使拥roundrobin 算法
server app1 192.168.1.111:80 check
server app2 192.168.1.112:80 check
 

haproxy统计页面的输出机制

frontend webser
log 127.0.0.1 local3
option forwardfor
bind *:80
default_backend app
backend app
cookie node insert nocache
balance roundrobin
server app1 192.168.1.111:80 check cookie node1 intval 2 rise 1 fall 2
server app2 192.168.1.112:80 check cookie node2 intval 2 rise 1 fall 2
server backup 127.0.0.1:8010 check backup
listen statistics
bind *:8009 # 自定义监听端口
stats enable # 启用基于程序编译时默认设置的统计报告
stats auth admin:admin # 统计页面用户名和密码设置
stats uri /admin?stats # 自定义统计页面的URL，默认为/haproxy?stats
stats hide-version # 隐藏统计页面上HAProxy的版本信息
stats refresh 30s # 统计页面自动刷新时间
stats admin if TRUE #如果认证通过就做管理功能，可以管理后端的服务器
stats realm Hapadmin # 统计页面密码框上提示文本，默认为Haproxy\ Statistics

frontend webser
log 127.0.0.1 local3
option forwardfor
bind *:80
default_backend app
backend app
cookie node insert nocache
balance roundrobin
server app1 192.168.1.111:80 check cookie node1 intval 2 rise 1 fall 2
server app2 192.168.1.112:80 check cookie node2 intval 2 rise 1 fall 2
server backup 127.0.0.1:8010 check backup
listen statistics
bind *:8009 # 自定义监听端口
stats enable # 启用基于程序编译时默认设置的统计报告
stats auth admin:admin # 统计页面用户名和密码设置
stats uri /admin?stats # 自定义统计页面的URL，默认为/haproxy?stats
stats hide-version # 隐藏统计页面上HAProxy的版本信息
stats refresh 30s # 统计页面自动刷新时间
stats admin if TRUE #如果认证通过就做管理功能，可以管理后端的服务器
stats realm Hapadmin # 统计页面密码框上提示文本，默认为Haproxy\ Statistics
 

动静分离示例：

frontend webservs
bind *:80
acl url_static path_beg -i /static /images /javascript /stylesheets
acl url_static path_end -i .jpg .gif .png .css .js .html
acl url_<a href="http://www.ttlsa.com/php/" title="php"target="_blank">php</a> path_end -i .php
acl host_static hdr_beg(host) -i img. imgs. video. videos. ftp. image. download.
use_backend static if url_static or host_static
use_backend dynamic if url_php
default_backend dynamic
backend static
balance roundrobin
server node1 192.168.1.111:80 check maxconn 3000
backend dynamic
balance roundrobin
server node2 192.168.1.112:80 check maxconn 1000

frontend webservs
bind *:80
acl url_static path_beg -i /static /images /javascript /stylesheets
acl url_static path_end -i .jpg .gif .png .css .js .html
acl url_php path_end -i .php
acl host_static hdr_beg(host) -i img. imgs. video. videos. ftp. image. download.
use_backend static if url_static or host_static
use_backend dynamic if url_php
default_backend dynamic
backend static
balance roundrobin
server node1 192.168.1.111:80 check maxconn 3000
backend dynamic
balance roundrobin
server node2 192.168.1.112:80 check maxconn 1000
 

http服务器配置完整示例

#---------------------------------------------------------------------
# Global settings
#---------------------------------------------------------------------
global
# to have these messages end up in /var/log/haproxy.log you will
# need to:
#
# 1) configure syslog to accept network log events. This is done
# by adding the '-r' option to the SYSLOGD_OPTIONS in
# /etc/sysconfig/syslog
#
# 2) configure local2 events to go to the /var/log/haproxy.log
# file. A line like the following can be added to
# /etc/sysconfig/syslog
#
# local2.* /var/log/haproxy.log
#
log 127.0.0.1 local2
chroot /var/lib/haproxy
pidfile /var/run/haproxy.pid
maxconn 4000
user haproxy
group haproxy
daemon
defaults
mode http
log global
option httplog
option dontlognull
option http-server-close
option forwardfor except 127.0.0.0/8
option redispatch
retries 3
timeout http-request 10s
timeout queue 1m
timeout connect 10s
timeout client 1m
timeout server 1m
timeout http-keep-alive 10s
timeout check 10s
maxconn 30000
listen stats
mode http
bind 0.0.0.0:1080
stats enable
stats hide-version
stats uri /haproxyadmin?stats
stats realm Haproxy\ Statistics
stats auth admin:admin
stats admin if TRUE
frontend http-in
bind *:80
mode http
log global
option httpclose
option logasap #不等待响应结束就记录日志，表示提前记录日志，一般日志会记录响应时长，此不记录响应时长
option dontlognull #不记录空信息
capture request header Host len 20 #记录请求首部的前20个字符
capture request header Referer len 60 #referer跳转引用，就是上一级
default_backend servers
frontend healthcheck
bind :1099 #定义外部检测机制
mode http
option httpclose
option forwardfor
default_backend servers
backend servers
balance roundrobin
server websrv1 192.168.1.111:80 check maxconn 2000
server websrv2 192.168.1.112:80 check maxconn 2000

#---------------------------------------------------------------------
# Global settings
#---------------------------------------------------------------------
global
# to have these messages end up in /var/log/haproxy.log you will
# need to:
#
# 1) configure syslog to accept network log events. This is done
# by adding the '-r' option to the SYSLOGD_OPTIONS in
# /etc/sysconfig/syslog
#
# 2) configure local2 events to go to the /var/log/haproxy.log
# file. A line like the following can be added to
# /etc/sysconfig/syslog
#
# local2.* /var/log/haproxy.log
#
log 127.0.0.1 local2
chroot /var/lib/haproxy
pidfile /var/run/haproxy.pid
maxconn 4000
user haproxy
group haproxy
daemon
defaults
mode http
log global
option httplog
option dontlognull
option http-server-close
option forwardfor except 127.0.0.0/8
option redispatch
retries 3
timeout http-request 10s
timeout queue 1m
timeout connect 10s
timeout client 1m
timeout server 1m
timeout http-keep-alive 10s
timeout check 10s
maxconn 30000
listen stats
mode http
bind 0.0.0.0:1080
stats enable
stats hide-version
stats uri /haproxyadmin?stats
stats realm Haproxy\ Statistics
stats auth admin:admin
stats admin if TRUE
frontend http-in
bind *:80
mode http
log global
option httpclose
option logasap #不等待响应结束就记录日志，表示提前记录日志，一般日志会记录响应时长，此不记录响应时长
option dontlognull #不记录空信息
capture request header Host len 20 #记录请求首部的前20个字符
capture request header Referer len 60 #referer跳转引用，就是上一级
default_backend servers
frontend healthcheck
bind :1099 #定义外部检测机制
mode http
option httpclose
option forwardfor
default_backend servers
backend servers
balance roundrobin
server websrv1 192.168.1.111:80 check maxconn 2000
server websrv2 192.168.1.112:80 check maxconn 2000
 

负载均衡MySQL服务的配置示例

#---------------------------------------------------------------------
# Global settings
#---------------------------------------------------------------------
global
# to have these messages end up in /var/log/haproxy.log you will
# need to:
#
# 1) configure syslog to accept network log events. This is done
# by adding the '-r' option to the SYSLOGD_OPTIONS in
# /etc/sysconfig/syslog
#
# 2) configure local2 events to go to the /var/log/haproxy.log
# file. A line like the following can be added to
# /etc/sysconfig/syslog
#
# local2.* /var/log/haproxy.log
#
log 127.0.0.1 local2
chroot /var/lib/haproxy
pidfile /var/run/haproxy.pid
maxconn 4000
user haproxy
group haproxy
daemon
defaults
mode tcp
log global
option httplog
option dontlognull
retries 3
timeout http-request 10s
timeout queue 1m
timeout connect 10s
timeout client 1m
timeout server 1m
timeout http-keep-alive 10s
timeout check 10s
maxconn 600
listen stats
mode http
bind 0.0.0.0:1080
stats enable
stats hide-version
stats uri /haproxyadmin?stats
stats realm Haproxy\ Statistics
stats auth admin:admin
stats admin if TRUE
frontend mysql
bind *:3306
mode tcp
log global
default_backend mysqlservers
backend mysqlservers
balance leastconn
server dbsrv1 192.168.1.111:3306 check port 3306 intval 2 rise 1 fall 2 maxconn 300
server dbsrv2 192.168.1.112:3306 check port 3306 intval 2 rise 1 fall 2 maxconn 300
 #---------------------------------------------------------------------
# Global settings
#---------------------------------------------------------------------
global
# to have these messages end up in /var/log/haproxy.log you will
# need to:
#
# 1) configure syslog to accept network log events. This is done
# by adding the '-r' option to the SYSLOGD_OPTIONS in
# /etc/sysconfig/syslog
#
# 2) configure local2 events to go to the /var/log/haproxy.log
# file. A line like the following can be added to
# /etc/sysconfig/syslog
#
# local2.* /var/log/haproxy.log
#
log 127.0.0.1 local2
chroot /var/lib/haproxy
pidfile /var/run/haproxy.pid
maxconn 4000
user haproxy
group haproxy
daemon
defaults
mode tcp
log global
option httplog
option dontlognull
retries 3
timeout http-request 10s
timeout queue 1m
timeout connect 10s
timeout client 1m
timeout server 1m
timeout http-keep-alive 10s
timeout check 10s
maxconn 600
listen stats
mode http
bind 0.0.0.0:1080
stats enable
stats hide-version
stats uri /haproxyadmin?stats
stats realm Haproxy\ Statistics
stats auth admin:admin
stats admin if TRUE
frontend mysql
bind *:3306
mode tcp
log global
default_backend mysqlservers
backend mysqlservers
balance leastconn
server dbsrv1 192.168.1.111:3306 check port 3306 intval 2 rise 1 fall 2 maxconn 300
server dbsrv2 192.168.1.112:3306 check port 3306 intval 2 rise 1 fall 2 maxconn 300

6、haproxy启用监控页面
编辑haproxy.cfg  加上下面参数  
listen admin_stats
        stats   enable
        bind    *:9090    //监听的ip端口号
        mode    http    //开关
        option  httplog
        log     global
        maxconn 10
        stats   refresh 30s   //统计页面自动刷新时间
        stats   uri /admin    //访问的uri   ip:8080/admin
        stats   realm haproxy
        stats   auth admin:Redhat  //认证用户名和密码
        stats   hide-version   //隐藏HAProxy的版本号
        stats   admin if TRUE   //管理界面，如果认证成功了，可通过webui管理节点
保存退出后
重起service haproxy restart
然后访问 http://ip:9090/admin          用户名:admin 密码:Redhat
