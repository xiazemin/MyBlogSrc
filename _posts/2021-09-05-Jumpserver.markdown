---
title: Jumpserver
layout: post
category: python
author: 夏泽民
---
https://gitee.com/jumpserver/jumpserver/
https://www.jumpserver.org/
https://github.com/jumpserver/jumpserver/
Jumpserver开源跳板机系统
公司服务器多了，不可避免的需要用到跳板机，如果是自己运维人员用，没有太高需求使用跳板机系统，直接利用一台具有公网IP的服务器登录然后内网SSH到目标服务器即可。
但是有以下一些情况，让我们需要部署跳板机系统了：
　　1. 网安等保评测（这点儿不用多说，做过等保的朋友应该都清楚）；
　　2. 部分开发人员因为需要登录线上服务器协助排查问题，这种情况一般是没有实时的日志收集分析系统，所以需要开发人员登录服务器跟踪分析日志；
　　3. 运维团队内部需要做审计、记录服务器操作等；

　　而个人对比使用了一些跳板机系统，个人认为jumpserver这款软件是当前最好用
　　
　　Jumpserver 使用 Python / Django 进行开发，可以管理SSH、 Telnet、 RDP、 VNC 协议资产，对于Linux和Windows服务器来说都能控制，具有美观的web管理页面，登录资产都可以在web页面操作，这样也避免开发人员等少量使用服务器的人都安装SSH 客户端软件。
<!-- more -->
1. Jumpserver 为管理后台, 管理员可以通过 Web 页面进行资产管理、用户管理、资产授权等操作, 用户可以通过 Web 页面进行资产登录, 文件管理等操作；
　　2. Coco 为 SSH Server 和 Web Terminal Server 。用户可以使用自己的账户通过 SSH 或者 Web Terminal 访问 SSH 协议和 Telnet 协议资产；
　　3. Luna 为 Web Terminal Server 前端页面, 用户使用 Web 方式登录受控服务器所需要的组件；
　　4. Guacamole 为 RDP 协议和 VNC 协议资产组件, 用户可以通过 Web Terminal 来连接 RDP 协议和 VNC 协议资产 (暂时只能通过 Web Terminal 来访问)，简单说就是受控服务器为Windows时需要安装；
　　

https://www.cnblogs.com/vilenx/p/12447196.html
https://www.jianshu.com/p/653c60987b0f
https://blog.csdn.net/dreamli199/article/details/81704996




